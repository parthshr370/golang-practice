# Concurrency

Concurrency means having more than one thing in progress. Go handles this using **Goroutines** and **Channels**.

## The Problem: Sequential Blocking

If you check 100 websites sequentially, and each takes 1 second, your program runs for 100 seconds.

```go
for _, url := range urls {
    results[url] = wc(url) // Each call blocks the next one
}
```

## Goroutines (`go`)

A **goroutine** is a lightweight thread managed by the Go runtime.

- **Syntax**: `go doSomething()`
- **Behavior**: The main function continues immediately without waiting for `doSomething` to finish.
- **Anonymous Functions**: Often used to spawn goroutines inside loops.
  ```go
  go func() {
      // Work happens here
  }()
  ```

## Race Conditions

A **race condition** occurs when multiple goroutines try to write to the same memory location (like a map) at the same time.

- **Panic**: `fatal error: concurrent map writes`
- **Detection**: Use `go test -race` to find these bugs.

## Channels (`chan`)

Channels are pipes that connect concurrent goroutines. They allow you to **communicate** and **synchronize**.

### The Channel Operator (`<-`)

The arrow shows the direction of data flow.

1. **Send (`ch <- value`)**: Puts data **INTO** the channel.

   ```go
   resultChannel <- result{url, wc(url)} // "Mail" the result
   ```

   _Blocks until someone is ready to receive._

2. **Receive (`value := <-ch`)**: Takes data **OUT OF** the channel.
   ```go
   r := <-resultChannel // Wait for mail
   ```
   _Blocks until data is available._

## The Solution: Coordinating with Channels

To fix the race condition, we coordinate the writes.

1. **Workers (Goroutines)**: Calculate the result and send it to the channel.
2. **Manager (Main Routine)**: Receives from the channel and writes to the map safely.

```go
// 1. Start Workers
for _, url := range urls {
    go func() {
        resultChannel <- result{url, wc(url)}
    }()
}

// 2. Collect Results
for i := 0; i < len(urls); i++ {
    r := <-resultChannel // Blocks here waiting for next result
    results[r.string] = r.bool // Safe, sequential write
}
```

This pattern is often summarized as:

> **"Do not communicate by sharing memory; instead, share memory by communicating."**

---

## Analogy: The Coffee Shop ☕

**Sequential Code (No Goroutines):**
Imagine a Coffee Shop with **one employee** (The Main Thread).
1.  Customer A orders.
2.  Employee makes coffee for A (takes 2 mins).
3.  Customer B orders.
4.  Employee makes coffee for B (takes 2 mins).
*   **Result:** Customer B waits 4 minutes.

**Concurrent Code (Goroutines):**
Imagine the same shop, but the employee can summon **magic assistants** instanty.
1.  Customer A orders.
2.  Employee yells `go makeCoffee(A)!` -> A magic assistant appears and starts brewing.
3.  Employee immediately turns to Customer B (no waiting).
4.  Employee yells `go makeCoffee(B)!` -> Another assistant appears.
*   **Result:** Customer B orders immediately. Both coffees brew at the same time.

---

## Technical Deep Dive

### OS Threads vs. Goroutines
*   **OS Thread (Heavy)**: Standard threads (Java/C++) have large stacks (1-2MB). Creating thousands will crash your RAM. Context switching is slow.
*   **Goroutine (Light)**: Start with tiny stacks (2KB). Managed by the Go Runtime, not the OS. You can spawn **millions**.

### M:N Scheduling
The Go Scheduler sits between your code and the OS.
*   **M (Goroutines)**: Thousands of logical tasks.
*   **N (OS Threads)**: Only a few actual CPU threads (e.g., 8).
The Scheduler juggles them instantly. If one blocks (waiting for a website), it swaps it out for another.

### The "Main" Trap
If `main()` finishes, the program exits **immediately**, killing all other goroutines.
*   **Why?** Because `main` didn't wait.
*   **The Fix:** Use Channels or WaitGroups to synchronize termination.

---

## Visualization: The Pipeline

Without channels, everyone fights for the map (Data Race). With channels, we create a pipeline:

```text
                  (Concurrency happens here)
                          |
[Worker 1] ---> [ Result ] \
                            \
[Worker 2] ---> [ Result ] --+--> [ CHANNEL PIPE ] 
                            /         |
[Worker 3] ---> [ Result ] /          |
                          |           | (Serialization happens here)
                          |           v
                                      |
                                 [ Main Routine ] <--- The "Maintainer"
                                      |
                                      v
                                  [ MAP (Safe) ]
```

## Q&A Section

### How to rewrite a traditional for-loop using modern Go syntax (Go 1.22+)?

**Question:**
How to rewrite this loop in modern go terms?
```go
for i := 0; i < wantedCount; i++ {
    go func() {
        counter.Inc()
        wg.Done()
    }()
}
```

**Answer:**
Go 1.22 introduced **range over integers**. You can now iterate a fixed number of times without manually managing the index variable:

```go
for range wantedCount {
    go func() {
        counter.Inc()
        wg.Done()
    }()
}
```

If you need the index, you can write `for i := range wantedCount`. This is cleaner and less error-prone than the traditional C-style `for` loop.

### Where is the goroutine applied in a concurrent test?

**Question:**
Where is the goroutine applied here?
```go
for i := 0; i < wantedCount; i++ {
    go func() {
        counter.Inc()
        wg.Done()
    }()
}
```

**Answer:**
The goroutine is launched right here in the loop with the `go` keyword:

```go
go func() {        // <-- the 'go' keyword launches a new goroutine
    counter.Inc()
    wg.Done()
}()
```

**How it works:**
- The `go` keyword immediately spawns a new goroutine to run the anonymous function
- It's not automatic - without `go`, each iteration would run sequentially
- Parallel execution - all 1000 goroutines start running almost simultaneously
- Main continues - the main goroutine doesn't wait for spawned goroutines to finish; it immediately moves to the next loop iteration

Without `go`:
```go
for i := 0; i < wantedCount; i++ {
    func() {           // No 'go' keyword here
        counter.Inc()
        wg.Done()
    }()
}
```
This would run **sequentially** - each iteration completes before the next one starts, effectively running on a single thread.

### Why does the anonymous function in the goroutine need no parameters?

**Question:**
This function seems lame - it has no params, just `counter` and `wg` inside:
```go
go func() {
    counter.Inc()
    wg.Done()
}()
```

**Answer:**
No parameters needed because of **closures**. The anonymous function automatically captures `counter` and `wg` from the outer scope. This is called a closure - the function "remembers" variables from where it was defined:

```go
counter := Counter{}          // outer scope
wg := sync.WaitGroup{}        // outer scope

for range 1000 {
    go func() {
        counter.Inc()  // captures counter from outer scope
        wg.Done()      // captures wg from outer scope
    }()
}
```

Each goroutine gets its own closure that references the same `counter` and `wg` variables from the outer scope.

### Do goroutines "slack" or wait for the loop to finish?

**Question:**
This for loop creates a new goroutine each time. Will the goroutines slack or wait for the loop to finish?

**Answer:**
Goroutines don't "slack" - they run immediately. Each iteration creates a goroutine that starts running right away, not at the end of the loop:

```go
Iteration 1: Launch goroutine A (starts running NOW)
Iteration 2: Launch goroutine B (starts running NOW)
Iteration 3: Launch goroutine C (starts running NOW)
...
Iteration 1000: Launch goroutine Z (starts running NOW)

Loop ends → wg.Wait() blocks until all 1000 goroutines finish
```

The goroutines don't wait for the loop to finish. They're actively executing on available CPU cores in parallel. The Go scheduler efficiently juggles them across your CPU cores.

If goroutines had no work and "slacked," they'd just finish quickly and be done. In this test, each goroutine immediately does `counter.Inc()` and `wg.Done()`, so they're doing real work, not idling.

### Why should you NOT embed sync.Mutex in your struct?

**Question:**
I've seen examples where sync.Mutex is embedded into the struct:
```go
type Counter struct {
    sync.Mutex
    value int
}
```

This looks elegant because you can call `c.Lock()` directly. The author says this is "bad and wrong" because exposing Lock and Unlock methods is confusing and harmful. What does "public API security" mean here?

**Answer:**
"Public API security" is about protecting your code from being misused by others. When you embed `sync.Mutex`, you make its public methods part of your struct's public interface.

**What embedding does:**
```go
type Counter struct {
    sync.Mutex  // Embedding makes all Mutex methods part of Counter
    value int
}

// Now anyone can do this:
counter.Lock()
counter.Unlock()
```

**Why it's bad:**

1. **Unwanted public methods** - `Lock()` and `Unlock()` become public. Any package can call them.

2. **Coupling** - Other code starts depending on your internal locking mechanism. If you change your lock implementation later (e.g., to RWMutex), their code breaks.

3. **Dangerous misuse** - Users mess with locks manually:
   ```go
   counter.Lock()
   counter.Inc()
   // Oops, forgot Unlock() - deadlock!
   ```

4. **Cannot fix bugs safely** - Once `Lock()` is public, you're stuck maintaining it forever for backwards compatibility.

**Better approach:**
```go
type Counter struct {
    mu sync.Mutex    // Private field
    value int
}

// Only safe methods are public
func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

Now users can ONLY use your safe public methods. They can't mess up the locking because it's encapsulated inside your struct.

**The mental model:**
- Design for the future - assume someone will use your code in unexpected ways
- Minimize surface area - only expose what you must
- Internal implementation stays internal - locking is a detail, not an interface

### What does "copies lock value" error mean?

**Question:**
I'm getting this error from go vet:
```
call of assertCounter copies lock value: gopractice/sync.Counter contains sync.Mutex
```

**Answer:**
You're passing `Counter` by value (copying the whole struct), but it contains a `sync.Mutex`. Mutexes cannot be copied safely.

**What's happening:**

```go
type Counter struct {
    mu    sync.Mutex    // This mutex CANNOT be copied
    value int
}

// You have this test function:
func assertCounter(t testing.TB, got Counter, want int) {  // <-- BUG: got is a copy
    t.Helper()
    if got.Value() != want {
        t.Errorf("got %d, want %d", got.Value(), want)
    }
}
```

**Why copying a mutex is dangerous:**

```
ORIGINAL Counter                    COPIED Counter
┌────────────────────┐             ┌────────────────────┐
│ mu  [LOCKED]       │  COPY  ──►  │ mu  [LOCKED]       │
│ value: 5           │             │ value: 5           │
└────────────────────┘             └────────────────────┘
      ▲                                  ▲
      │                                  │
   goroutine A                       goroutine B
   has original mutex               has copy mutex

PROBLEM: Two goroutines think they have the SAME lock, but they have
different mutex instances. No synchronization! Race conditions!
```

**Visual explanation:**

```
Safe: Pass by pointer (reference)
┌────────────────┐
│   Counter {}   │
│  ┌──────────┐  │
│  │   Mutex  │  │
│  │   value  │  │
│  └──────────┘  │
└────────────────┘
       ▲
       │ Pointer (same object)   
       │
━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━
       │
       ├──► Goroutine 1 sees THE SAME mutex
       │
       └──► Goroutine 2 sees THE SAME mutex

Result: Only one goroutine can hold the lock at a time. Safe!


UNSAFE: Pass by value (copy)
┌────────────────┐         Copy created         ┌────────────────┐
│   Counter {}   │ ────────────────────────►  │   Counter {}   │
│  ┌──────────┐  │                           │  ┌──────────┐  │
│  │ Mutex A  │  │                           │  │ Mutex B  │  │
│  │   value  │  │                           │  │   value  │  │
│  └──────────┘  │                           │  └──────────┘  │
└────────────────┘                           └────────────────┘
      ▲                                           ▲
      │                                           │
   goroutine 1                                goroutine 2

Result: Goroutine 1 locks Mutex A, goroutine 2 locks Mutex B.
They THINK they're synchronized, but they're not! Two lock instances!
```

**The fix: Pass by pointer**

```go
func assertCounter(t testing.TB, counter *Counter, want int) {  // Use pointer
    t.Helper()
    if counter.Value() != want {  // Dereference the pointer
        t.Errorf("got %d, want %d", counter.Value(), want)
    }
}
```

And update the test:

```go
t.Run("test name", func(t *testing.T) {
    counter := Counter{}
    // ... set up counter ...
    assertCounter(t, &counter, wantedCount)  // Pass address with &
})
```

**How go vet catches this:**
Go's `sync` package uses a special `noCopy` struct insideMutex. When you run `go vet`, it checks if you're copying structs that contain `noCopy`. If you are, it warns you.

### How do mutexes relate to goroutines and the "happens before" relationship?

**Question:**
I read this about mutexes: "mutex locking doesn't lock the goroutine. The rule is that only one goroutine can hold a mutex at a time." What does "happens before" mean and how do mutexes establish ordering?

**Answer:**
Mutexes are NOT tied to goroutines. They're tied to data. The "happens before" relationship is a concurrency primitive that guarantees ordering of operations.

**Mutexes are tied to data, not goroutines:**

```
Mutex = Gatekeeper for shared data
Not = Lock attached to a specific goroutine

SHARED DATA
┌────────────────────────┐
│     []string slice     │
│                        │
│  Gate:               │
│  ┌──────────────────┐│
│  │   sync.Mutex     ││  ← This mutex guards the slice
│  └──────────────────┘│
└────────────────────────┘
         ▲
         │
         ├──► Goroutine A: "Knock knock, can I enter?"
         │    (calls Lock(), gets in, work, Unlock())
         │
         └──► Goroutine B: "Wait, gate is closed."
              (blocks on Lock() until A leaves)

The mutex doesn't care WHICH goroutine wants it.
It only cares that ONLY ONE can hold it at a time.
```

**The "happens before" relationship:**

This is a formal concept from the Go Memory Model. It defines which operations are guaranteed to be visible to which other operations.

```
Without synchronization (NO happens-before):

Goroutine A                          Goroutine B
───────────                          ───────────
x = 42                               print(x)
  │
  └─────► ??? ◄─────┐

Result: UNDEFINED! B might print 0, 42, or anything.
The compiler can reorder. The runtime can delay.
No guarantees because A and B have no relationship.


 With synchronization (happens-before with mutex):

Goroutine A                          Goroutine B
───────────                          ───────────
mu.Lock()                            mu.Lock()
  │                                    │
  ▼                                    ▼
x = 42                                print(x)
  │                                    │
  ▼                                    ▼
mu.Unlock()                           mu.Unlock()

This establishes a happens-before relationship:
A's Unlock() happens-before B's Lock()
Therefore: A's write (x = 42) happens-before B's read (print(x))
Result: B is GUARANTEED to see 42
```

**Key rules about "happens before":**

1. **No relationship = undefined order**
   ```go
   // Goroutine A                // Goroutine B
   x = 1                          y = 2
   // Is x=1 visible to B?        // Is y=2 visible to A?
   // NO GUARANTEE!               // NO GUARANTEE!
   ```

2. **Mutex Lock/Unlock creates happens-before**
   ```go
   // Goroutine A                // Goroutine B
   mu.Lock()                     mu.Lock()
   x = 1                         print(x)  // GUARANTEED to see x=1
   mu.Unlock()                   mu.Unlock()
   ```

3. **Channel operations create happens-before**
   ```go
   // Goroutine A                // Goroutine B
   ch <- value                   v := <-ch
   // Send happens-before receive
   ```

**Critical insight from the Reddit discussion:**

> "Mutex locking doesn't lock the goroutine. The rule is that only one goroutine can hold a mutex at a time. Crucially, from the perspective of other goroutines, all work done by the locking goroutine before it holds the mutex happens before all work while it's holding the mutex, which happens before all work done by any other goroutine that later acquires the mutex."

**Visual timeline:**

```
Goroutine A:                               Goroutine B:

  work1()                                   work3()
     │                                          │
     ▼                                          ▼
────────────────────────────────────────────────────────
     │                                          │
    NO happens-before relationship!
    Can't guarantee order between work1 and work3
────────────────────────────────────────────────────────
     │                                          │
     ▼                                          ▼
  mu.Lock()        ◄───── happens-before ──────►  mu.Lock() waits
     │                                          │
     ▼                                          ▼
  work2()                                   (blocked)
     │                                          │
     ▼                                          │
────────────────────────── happens-before ──────┘
     │
  mu.Unlock()  ◄─── happens-before ──────────────►  mu.Lock() succeeds
     │                                          │
     ▼                                          ▼
  work1()                                 work2()  ← now sees work2's changes
```

### Why doesn't Go have reentrant locks?

**Question:**
I learned Go doesn't have reentrant mutexes. If I do `mtx.Lock()` twice in a row, it hangs forever. Other languages allow this. Why?

**Answer:**
Reentrant (recursive) locks automatically detect if the "owner" already holds the lock and let them acquire it again. Go deliberately doesn't support this because it encourages bad design.

**What reentrant locking does (in other languages like Java):**

```
In Java (has reentrant locks):

function nestedAcquisition() {
    lock.lock()      // First Lock: SUCCESS (now you own it)
    doSomething()
    innerFunction()  // This also calls lock.lock()
    lock.unlock()    // Unlock: releases the lock
}

function innerFunction() {
    lock.lock()      // Second Lock: SUCCESS (same owner, just increments counter)
    doInnerWork()
    lock.unlock()    // Unlock: decrements counter
}

This works because the lock knows "this thread already owns me"
```

**What happens in Go (no reentrant locks):**

```
In Go (NO reentrant locks):

func nestedAcquisition() {
    mu.Lock()        // First Lock: SUCCESS
    doSomething()
    innerFunction()  // This also calls mu.Lock()
    mu.Unlock()
}

func innerFunction() {
    mu.Lock()        // Second Lock: DEADLOCK!
                    // YOU ARE ALREADY THE OWNER
                    // Waiting for yourself to unlock
                    // But you can't unlock because you're stuck here
    doInnerWork()
    mu.Unlock()
}

DEADLOCK: Goroutine blocked waiting for itself
```

**Deadlock visualization:**

```
Goroutine A calls Lock():
┌────────────────────┐
│     Mutex          │
│   [UNLOCKED]       │
└────────────────────┘
       ▲
       │ "I want this lock"
       │
  Goroutine A acquires lock
       │
       ▼
┌────────────────────┐
│     Mutex          │
│   [LOCKED]         │  ← Owner: Goroutine A
│   lockCount: 1     │
└────────────────────┘

Goroutine A calls Lock() AGAIN (in nested function):
       ▲
       │ "I want this lock"
       │
       │ Mutex checks: "Who holds me?"
       │            ⟶ "Goroutine A holds me"
       │            ⟶ "Is this reentrant?"
       │            ⟶ "NO! Wait for unlock!"
       │
       │ But Goroutine A is the one asking!
       │ A is waiting for A to unlock.
       │ A can't unlock because A is blocked.
       │
       ▼
    DEADLOCK!
```

**Why Go deliberately avoids reentrant locks:**

```
Problem with reentrant locking: It hides design problems.

Reentrant locks let you write:
  function() {
      lock()
      helper() {
          lock()    // Reentrant! Works!
      }
      unlock()
  }

But who actually knows helper() needs locking?
Is it the caller's job or helper's job?
Can helper() be called without the outer function?
Does helper() really need locking at all?

The reentrant lock papers over the fact that you
don't know who should be managing the lock.
```

**Go's philosophy: Be explicit**

```go
// Better design: Separate concerns

// Option 1: Caller manages the lock
func outerFunction() {
    mu.Lock()
    defer mu.Unlock()
    innerFunction()  // No locking needed, just does work
}

// Option 2: Helper function takes a locked context
func outerFunction() {
    mu.Lock()
    defer mu.Unlock()
    innerFunctionWithLockedContext()
}

// Option 3: Split into public/private methods
func (c *MyStruct) DoSomething() {
    mu.Lock()
    defer mu.Unlock()
    c.doSomethingUnsafe()  // Private: assumes lock held
}
```

Go forces you to design your code so the locking responsibility is clear and explicit. No hidden magic.

### How does go vet detect lock copies with noCopy?

**Question:**
I found that sync.WaitGroup uses a noCopy struct to warn when someone accidentally copies a lock. How does this trick work?

**Answer:**
Go's `go vet` tool checks if you're copying structs that have a special `Lock()` method. The `sync` package adds a `noCopy` field that triggers this check.

**How the trick works:**

```go
// In the sync package (implementation detail):

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// Now use it in your struct:

type MyType struct {
    _ noCopy      // This field has Lock() method!
    data   int
}
```

**Why this triggers go vet:**

```
go vet has a special rule:
"Any struct type that has a Lock() or Unlock() method
 must be passed by pointer, not copied"

go vet sees:
    type MyType struct {
        _ *noCopy    // Has Lock() method!
        data int
    }

Then checks your code:
    var x MyType
    y := x      // ← COPYING! go vet complains!

Why? Because go vet also scans sync.Mutex, sync.WaitGroup,
sync.RWMutex - they all have Lock() methods and similar
hidden noCopy fields.
```

**Visual demonstration:**

```
Normal struct (can be copied):

┌────────────────────┐
│ type Point struct  │
│     X, Y float64   │
└────────────────────┘
       │
       ├──► a := Point{1, 2}
       │    b := a         // COPYING: OK
       │    fmt.Println(b) // Works fine
       │
       └──► No Lock() method, so go vet doesn't care


Struct with noCopy (warns on copy):

┌────────────────────────────┐
│ type SafeCounter struct    │
│    _ noCopy               │ ← Hidden field
│    value int              │
└────────────────────────────┘
       │
       ├──► x := SafeCounter{value: 5}
       │    y := x                  // COPYING: go vet warns!
       │    // go vet: "call of ... copies lock value"
       │
       └──► Why? Because noCopy has Lock() method!
            go vet: "This struct looks like it has a lock, 
                     don't copy it!"


How it works internally:

┌──────────────────────────────────┐
│ go vet scan:                    │
│                                  │
│ 1. Find struct type definitions  │
│ 2. Check: Does this struct have  │
│    Lock() or Unlock() methods?   │
│ 3. If YES:                      │
│    - Search for assignment: b = a│
│    - Search for function param:  │
│      func foo(x MyStruct)        │ ← Pass by value!
│    - Find: YELL at the user     │
└──────────────────────────────────┘


The secret: embedding or pointer

type MyType struct {
    _ noCopy    // ← Go vet sees: MyType has Lock() because
               // noCopy has Lock()
}

// Even though _ is a zero-sized field (takes 0 memory),
// the Lock() method is still considered part of MyType
```

**Zero-sized field trick:**

```
type MyType struct {
    _ noCopy     // ← Zero-sized (0 bytes!)
    value int     ← Takes 8 bytes
}

Memory layout:
┌────────────────┐
│   value (8B)   │
└────────────────┘
Only 8 bytes total! The _ field pads to 0 bytes.

But functionally:
myType.Lock()    // ← valid! calls noCopy.Lock()
myType.Unlock()  // ← valid! calls noCopy.Unlock()

Go is clever: zero field but methods still available
                AND go vet still checks it!
```

**Bonus tip: Preventing map keys**

```go
// You can also use zero-sized fields to prevent
// your struct from being used as a map key:

type Email struct {
    username string
    domain   string
    _ [0]map[any]any  // ← Zero-sized map
}

// This is now:
// - NOT comparable (can't be map key)
// - Can't use ==
// - Takes 0 extra memory

email1 := Email{"bob", "example.com"}
email2 := Email{"bob", "example.com"}

// This would panic:
// email1 == email2  // INVALID!

// This won't compile:
// m := map[Email]string{}  // INVALID key type!
```

**Zero-sized field comparison:**

```
Email struct layout:
┌────────────────────────────────┐
│       username (string)        │ ← 16 bytes
│       domain   (string)        │ ← 16 bytes
│       _ [0]map[any]any         │ ← 0 bytes!
└────────────────────────────────┘

Total: 32 bytes (the map is fake, zero size)

But:
email1 == email2     // ← Compile error!
// Why? Maps are uncomparable. Any struct with
// an uncomparable field is also uncomparable
```

