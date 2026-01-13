# Context

Managing cancellation of long-running processes in Go.

### The Problem

Your web server kicks off long-running work (database queries, API calls). If the user cancels the request, that expensive work keeps running, wasting resources.

### Basic Server Setup

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, store.Fetch())
    }
}
```

When Fetch() takes 100ms but user cancels in 5ms, the work still completes.

### Cancellable Context Test

The test simulates user cancelling a request:

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
    store := &SpyStore{response: "hello"}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)

    // 1. Create cancellable context
    cancellingCtx, cancel := context.WithCancel(request.Context())

    // 2. Schedule cancel to be called after 5ms
    time.AfterFunc(5*time.Millisecond, cancel)

    // 3. Replace request's context
    request = request.WithContext(cancellingCtx)

    // 4. Call server
    response := httptest.NewRecorder()
    svr.ServeHTTP(response, request)

    // 5. Check store was cancelled
    if !store.cancelled {
        t.Error("store was not told to cancel")
    }
})
```

**Timeline:**

- 0ms: Create context, schedule cancel at 5ms
- 0ms: Call server, which calls Store.Fetch() (sleeps 100ms)
- 5ms: cancel() called, context is done
- 100ms: Fetch() finishes, returns data
- Test fails because server never called Cancel()

### First Approach: Manual Cancel (Not Idiomatic)

Server detects cancellation and calls Cancel():

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        data := make(chan string, 1)

        go func() {
            data <- store.Fetch()
        }()

        select {
        case d := <-data:
            fmt.Fprint(w, d)
        case <-ctx.Done():
            store.Cancel()
        }
    }
}
```

**Pattern:**

- Run work in goroutine
- Use select to race result against cancellation
- Call Cancel() if context is done

**Why this approach is bad:**

- Server needs to know how to cancel downstream processes
- Creates coupling throughout the codebase
- Each component manages its own cancellation logic

### Idiomatic Approach: Pass Context Through

Instead of manual Cancel(), pass context to Store:

```go
type Store interface {
    Fetch(ctx context.Context) (string, error)
}
```

Store becomes responsible for respecting cancellation:

```go
func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
    data := make(chan string, 1)

    // Simulate slow work that can be cancelled
    go func() {
        var result string
        for _, c := range s.response {
            select {
            case <-ctx.Done():
                return
            default:
                time.Sleep(10 * time.Millisecond)
                result += string(c)
            }
        }
        data <- result
    }()

    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case res := <-data:
        return res, nil
    }
}
```

Server becomes simpler:

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := store.Fetch(r.Context())

        if err != nil {
            return
        }

        fmt.Fprint(w, data)
    }
}
```

**Why this is better:**

- Context propagates through call chain
- Each function handles its own cancellation
- Server doesn't need to know about downstream processes
- Consistent pattern throughout Go codebase

### Context Tree and Propagation

Contexts form a tree. When parent is cancelled, all children are cancelled:

```
Root Context
       │
       ├─── cancellingCtx (derived from Root)
       │           │
       │           └─── ctx.Done() channel closes
       │                    when cancel() called
       │
       └─── All derived contexts cancel too
```

### Testing That No Response Was Written

When cancelled, shouldn't write response. Use spy for ResponseWriter:

```go
type SpyResponseWriter struct {
    written bool
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
    s.written = true
    return 0, errors.New("not implemented")
}
```

Test:

```go
response := &SpyResponseWriter{}
svr.ServeHTTP(response, request)

if response.written {
    t.Error("a response should not have been written")
}
```

### Proper Context.Value Usage: The userip Package Pattern

When you need to use context.Value, the right approach is to create a package that hides the details and provides strongly-typed access. This prevents key collisions and provides type safety.

**How Context.Value works:**

- Context provides a key-value mapping
- Keys must support equality (comparable)
- Keys and values are both `interface{}` (any)
- Values must be safe for goroutine use

**The Pattern: Unexported Key Type**

To avoid key collisions with other packages, define an unexported key type:

```go
// Package userip

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address. Its value of zero is
// arbitrary. If this package defined other context keys, they would have
// different integer values.
const userIPKey key = 0
```

**Why unexported?**

- Other packages can't use the same type (they don't have access to `key`)
- Only `userip` package can set/get this specific value
- Prevents accidental collisions

**FromRequest extracts userIP from request:**

```go
func FromRequest(req *http.Request) (net.IP, error) {
    ip, _, err := net.SplitHostPort(req.RemoteAddr)
    if err != nil {
        return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
    }
    return net.ParseIP(ip), nil
}
```

**NewContext creates a context with userIP:**

```go
func NewContext(ctx context.Context, userIP net.IP) context.Context {
    return context.WithValue(ctx, userIPKey, userIP)
}
```

**FromContext extracts userIP from context:**

```go
func FromContext(ctx context.Context) (net.IP, bool) {
    // ctx.Value returns nil if ctx has no value for the key;
    // the net.IP type assertion returns ok=false for nil.
    userIP, ok := ctx.Value(userIPKey).(net.IP)
    return userIP, ok
}
```

**Why type assertion?**

`ctx.Value(key)` returns `interface{}` (any). The type assertion `.(net.IP)` converts it to the concrete type, with `ok` indicating success.

**Usage example:**

```go
// Handler stores user IP in context
func handler(w http.ResponseWriter, r *http.Request) {
    userIP, err := userip.FromRequest(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx := userip.NewContext(r.Context(), userIP)
    // Pass ctx to downstream functions...
}

// Downstream function retrieves user IP
func processRequest(ctx context.Context) {
    userIP, ok := userip.FromContext(ctx)
    if !ok {
        log.Println("User IP not available in context")
        return
    }
    log.Printf("Processing request from %s", userIP)
}
```

### Context.Value Warning

Don't use context.Value for function inputs:

```go
// BAD
func DoSomething(ctx context.Context) {
    userID := ctx.Value("userID")  // No type safety
}

// GOOD
func DoSomething(ctx context.Context, userID string) {  // Typed
}
```

### Key Takeaways

- Context manages cancellation across goroutines
- Pass context as first argument through call chain
- Use `select { case <-ctx.Done(): }` pattern to check cancellation
- Each function respects its own context
- Don't use `context.Value` for business logic inputs

### Q&A Section

### What does the test failure mean and what's happening in the t.Run?

**Question:**
I'm getting this test failure:

```
--- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled_ (0.10s)
    context_test.go:29: store was not told to cancel
```

What's happening in the test and why is it failing?

**Answer:**
The test creates a request that gets cancelled after 5ms, but the server doesn't detect the cancellation or call Cancel().

**Here's what happens step by step:**

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
    data := "hello , world"
    store := &SpyStore{response: data}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)

    // Step 1: Create cancellable context
    cancellingCtx, cancel := context.WithCancel(request.Context())

    // Step 2: Schedule cancel to be called after 5ms
    time.AfterFunc(5*time.Millisecond, cancel)

    // Step 3: Replace request's context
    request = request.WithContext(cancellingCtx)

    response := httptest.NewRecorder()

    // Step 4: Call server
    svr.ServeHTTP(response, request)

    // Step 5: Check if Cancel() was called
    if !store.cancelled {
        t.Error("store was not told to cancel")
    }
})
```

**Timeline of events:**

```
Time 0ms:    Create request, create cancellableCtx
Time 0ms:    Schedule cancel() to run at 5ms
Time 0ms+     Call svr.ServeHTTP()

             Inside Server():
               Calls store.Fetch()
               Fetch() starts sleeping for 100ms...

Time 5ms:    cancel() is called automatically!
             cancellingCtx.Done() channel closes

             BUT: Server doesn't check the context!
             It's still waiting for Fetch() to finish...

Time 100ms:  Fetch() finally finishes sleeping
             Returns "hello , world"

             Server writes response

Test ends:   Checks store.cancelled
             Result: false (Cancel was never called)
             TEST FAILS!
```

**The issue:** Your Server function just calls `store.Fetch()` and ignores the context:

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(w, store.Fetch())  // Ignores r.Context()
    }
}
```

The server needs to check `r.Context()` and call `store.Cancel()` when the context is done.

### When to use pointer receivers vs value receivers and pass pointers vs values?

**Question:**
I see `func (s *SpyStore)` in method declarations and `&SpyStore{data}` when creating. When do I use `*` vs `&`?

**Answer:**

**Type declaration (parameter or variable type):**

```go
func (s *SpyStore) Fetch() string
//       ^ This means: "s is a pointer to SpyScore"
```

**Creating and passing:**

```go
svr := Server(&SpyScore{data})
//             ^ This means: "create SpyScore, then get its address"
```

**When to use pointer receiver `*Type`:**

```go
func (s *SpyStore) Cancel() {
    s.cancelled = true  // ← WRITING to field
}
```

Use when methods need to read or write fields. Especially when you need to change the struct's state.

**When to use value receiver `Type`:**

```go
func (s SpyStore) Fetch() string {
    return s.response  // ← Only READING
}
```

Use when the method only reads data and doesn't need to modify state.

**Why we pass `&SpyScore{data}`:**

```go
store := &SpyScore{response: "hello"}
//            Create struct: SpyScore{response: "hello"}
//            & takes its address
```

We pass pointer because:

1. SpyScore methods use pointer receivers `(s *SpyScore)`
2. We might want to read the `cancelled` field later in tests
3. Avoids copying the struct for efficiency with larger structs

**If you passed by value:**

```go
store := SpyScore{response: "hello"}
svr := Server(store)

if store.cancelled {  // This would always be false!
    // Server got a COPY, not the original
}
```

**Decision tree:**

- Method needs to modify struct? Use pointer receiver `(s *Type)`
- Need to read struct state after passing? Pass as pointer `&Type{}`

### Why do we cancel context and what are we testing?

**Question:**
Why are we cancelling the context in the test? What even are we trying to verify?

**Answer:**

**Real-world scenario:**

```
User opens web browser
  ↓
User clicks a link that takes 10 seconds to load
  ↓
User gets impatient and clicks again (hits refresh)
  ↓
Browser cancels the first request
```

**Without context cancellation:**

```
Server receives first request
  ↓
Starts slow database query (10 seconds)
  ↓
Browser cancels request (user hit refresh)
  ↓
Server KEEPS RUNNING the query!
  ↓
10 seconds later: Query finishes, wasting resources
```

**With context cancellation:**

```
Server receives first request
  ↓
Starts slow database query
  ↓
Browser cancels request
  ↓
Server detects cancellation via context
  ↓
Stops database query immediately
  ↓
Saves time and resources
```

**What the test is verifying:**

```go
// Simulate: User cancels request after 5ms
time.AfterFunc(5*time.Millisecond, cancel)

// Server should detect this and tell store to stop
// We verify: Was the store actually told to cancel?
if !store.cancelled {
    t.Error("store was not told to cancel")  // ← This is what we're checking
}
```

**Visual explanation:**

```
WITHOUT cancellation (BAD):
┌──────────────────────────────────┐
│ Request arrives                  │
│                                  │
┌──────────────────┐              │
│ Query: 100ms ────┼── Still runs │
└──────────────────┘              │
     ↑                            ↑
  Cancelled                     Still
  at 5ms                      Running!
     │                            │
     └────────────────────────────┘
         Wasted time/resources


WITH cancellation (GOOD):
┌──────────────────────────────────┐
│ Request arrives                  │
│                                  │
┌──────────────────┐              │
│ Query: 100ms ────┼── Stopped immediately │
└──────────────────┘              │
     ↑                            ↑
  Cancelled                     Query
  at 5ms                      Aborted
     │                            │
     └────────────────────────────┘
         Saved time/resources
```

**Why test this matters:**

- Without context, servers waste resources on cancelled requests
- With context, servers respond faster and scale better
- The test ensures your cancellation logic actually works

**The test in plain English:**
"Hey server, I'm going to start a request then cancel it after 5ms. After it's cancelled, I expect you to have told the store to stop working. If the store kept running, you failed."
func (s \*SpyStore) Fetch() string
// ^ Means: "s is a pointer to SpyStore"

```

**Creating and passing:**
```

store := &SpyStore{data}
// ^ Takes address of the struct

````

**When to use pointer receiver `*Type`:**
```go
func (s *SpyStore) Cancel() {
    s.cancelled = true  // WRITING to field
}
````

Use when methods need to read/write or modify the struct's state.

**When to use value receiver `Type`:**

```go
func (s SpyStore) Fetch() string {
    return s.response  // Only READING
}
```

Use when methods only read data and don't modify state.

**Why we pass `&SpyScore{data}`:**

```go
store := &SpyScore{response: "hello"}
// Create SpyScore and take address
```

We pass pointer because:

1. Methods might need pointer receivers
2. We want to read `cancelled` field later in tests
3. Avoids copying larger structs for efficiency

**If you pass by value:**

```go
store := SpyScore{response: "hello"}
Server(store)

if store.cancelled {  // Always false!
    // Server got a COPY, not the original
}
```

type ResponseWriter interface {
// Header returns the header map that will be sent by
// [ResponseWriter.WriteHeader]. The [Header] map also is the mechanism with which
// [Handler] implementations can set HTTP trailers.
//
// Changing the header map after a call to [ResponseWriter.WriteHeader] (or
// [ResponseWriter.Write]) has no effect unless the HTTP status code was of the
// 1xx class or the modified headers are trailers.
//
// There are two ways to set Trailers. The preferred way is to
// predeclare in the headers which trailers you will later
// send by setting the "Trailer" header to the names of the
// trailer keys which will come later. In this case, those
// keys of the Header map are treated as if they were
// trailers. See the example. The second way, for trailer
// keys not known to the [Handler] until after the first [ResponseWriter.Write],
// is to prefix the [Header] map keys with the [TrailerPrefix]
// constant value.
//
// To suppress automatic response headers (such as "Date"), set
// their value to nil.
Header() Header

    // Write writes the data to the connection as part of an HTTP reply.
    //
    // If [ResponseWriter.WriteHeader] has not yet been called, Write calls
    // WriteHeader(http.StatusOK) before writing the data. If the Header
    // does not contain a Content-Type line, Write adds a Content-Type set
    // to the result of passing the initial 512 bytes of written data to
    // [DetectContentType]. Additionally, if the total size of all written
    // data is under a few KB and there are no Flush calls, the
    // Content-Length header is added automatically.
    //
    // Depending on the HTTP protocol version and the client, calling
    // Write or WriteHeader may prevent future reads on the
    // Request.Body. For HTTP/1.x requests, handlers should read any
    // needed request body data before writing the response. Once the
    // headers have been flushed (due to either an explicit Flusher.Flush
    // call or writing enough data to trigger a flush), the request body
    // may be unavailable. For HTTP/2 requests, the Go HTTP server permits
    // handlers to continue to read the request body while concurrently
    // writing the response. However, such behavior may not be supported
    // by all HTTP/2 clients. Handlers should read before writing if
    // possible to maximize compatibility.
    Write([]byte) (int, error)

    // WriteHeader sends an HTTP response header with the provided
    // status code.
    //
    // If WriteHeader is not called explicitly, the first call to Write
    // will trigger an implicit WriteHeader(http.StatusOK).
    // Thus explicit calls to WriteHeader are mainly used to
    // send error codes or 1xx informational responses.
    //
    // The provided code must be a valid HTTP 1xx-5xx status code.
    // Any number of 1xx headers may be written, followed by at most
    // one 2xx-5xx header. 1xx headers are sent immediately, but 2xx-5xx
    // headers may be buffered. Use the Flusher interface to send
    // buffered data. The header map is cleared when 2xx-5xx headers are
    // sent, but not with 1xx headers.
    //
    // The server will automatically send a 100 (Continue) header
    // on the first read from the request body if the request has
    // an "Expect: 100-continue" header.
    WriteHeader(statusCode int)

}
A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.
A ResponseWriter may not be used after [Handler.ServeHTTP] has returned.
