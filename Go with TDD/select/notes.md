# Select & Concurrency Patterns

The goal of this chapter is to race two URLs and return the one that responds faster.

## 1. The Naive Approach (Sequential)

We could simply measure the time for `A`, then measure `B`, and compare.

```go
startA := time.Now()
http.Get(a)
durationA := time.Since(startA)
// ... repeat for B
```

**Problem:** This is blocking. If A takes 5s and B takes 1s, we wait 6s total. We want to wait only 1s (the fastest time).

## 2. The `select` Statement

`select` is a switch statement for channels. It blocks until **one** of its cases is ready.

```go
select {
case <-chanA:
    // Run this if A finishes first
case <-chanB:
    // Run this if B finishes first
}
```

## 3. The `ping` Function (The Runner)

We need a function that starts a background process and returns a channel signal.

```go
func ping(url string) chan struct{} {
    ch := make(chan struct{}) // 1. Create channel
    go func() {               // 2. Start runner
        http.Get(url)
        close(ch)             // 3. Signal completion
    }()
    return ch                 // 4. Return channel immediately
}
```

- **`chan struct{}`**: The smallest data type (0 bytes). Used purely for signaling.
- **`close(ch)`**: Closing a channel sends an immediate signal to receivers.

## 4. How the Race Works (Step-by-Step)

```go
case <-ping(a):
```

1. **Parallel Start**: `ping(a)` runs and spawns a goroutine. `ping(b)` does the same. Both are running.
2. **The Hook**: `<-` connects the `select` statement to the channels returned by `ping`.
3. **The Wait**: `select` pauses the main thread.
4. **The Trigger**: The first URL to return calls `close(ch)`.
5. **The Win**: `select` sees that channel is ready, executes that case, and returns the winner.

## 5. Timeouts with `time.After`

What if both servers hang? We add a timeout case.
`time.After(duration)` returns a channel that sends a signal after the time passes.

```go
select {
case <-ping(a):
    return a, nil
case <-ping(b):
    return b, nil
case <-time.After(10 * time.Second):
    return "", fmt.Errorf("timeout")
}
```

This acts as a third racer. If the timer finishes before A or B, the error case wins.

## 6. `httptest` (Testing HTTP)

Instead of calling real websites (slow, flaky), we spawn fake local servers.

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}))
// Use server.URL in your test
```

This keeps tests fast, offline, and reliable.

ping

We have defined a function ping which creates a chan struct{} and returns it.

In our case, we don't care what type is sent to the channel, we just want to signal we are done and closing the channel works perfectly!

Why struct{} and not another type like a bool? Well, a chan struct{} is the smallest data type available from a memory perspective so we get no allocation versus a bool. Since we are closing and not sending anything on the chan, why allocate anything?

select is this referee that waits for multiple channels at once and the first to send a value wins the code underneath.

Case is just us putting and creating channels.

## 7. Always use `make` for channels

In Go, you must use `make` to initialize a channel.

```go
ch := make(chan struct{})
```

If you just declare it without `make`:
```go
var ch chan struct{} // This is a NIL channel
```
- **Blocking**: Sending to or receiving from a `nil` channel blocks **forever**.
- **Panic**: Closing a `nil` channel causes a panic (crash).

Using `make` ensures the channel's internal data structures are allocated and ready for communication.

---

## Summary
- **`select`**: Helps you synchronise processes by waiting on multiple channels.
- **`time.After`**: A handy pattern to prevent systems from blocking forever.
- **`httptest`**: Allows you to create controllable, reliable test servers using the standard `net/http` interfaces.
- **`chan struct{}`**: The memory-efficient way to signal events without passing data.
- **`defer`**: Useful for cleaning up resources (like `server.Close()`) at the end of a function.

