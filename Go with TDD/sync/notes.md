# Sync

Notes on concurrent shared state using sync.Mutex.

## Overview

This module builds on concurrency basics and focuses on protecting shared data when multiple goroutines access it. We build a thread-safe counter that can be incremented concurrently without race conditions.

## The Problem

When 1000 goroutines try to increment the same counter value, you get a race condition:

```go
for range 1000 {
    go func() {
        counter.value++  // Data race!
    }()
}
```

Multiple goroutines read the same value, compute new value, then write. Some writes get overwritten.

## The Solution: sync.Mutex

A Mutex is a mutual exclusion lock. Only one goroutine can hold it at a time.

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

**How it works:**

- `Lock()` blocks until the mutex is available
- `defer Unlock()` ensures unlock even if function panics
- While one goroutine holds the lock, others wait

## WaitGroup Pattern

For coordinating goroutine completion:

```go
var wg sync.WaitGroup
wg.Add(1000)  // Tell it how many goroutines to expect

for range 1000 {
    go func() {
        counter.Inc()
        wg.Done()  // Mark this goroutine as done
    }()
}

wg.Wait()  // Block until all goroutines call Done()
```

## Implementation Details

### Test Helper Pattern

Create reusable assertion functions:

```go
func assertCounter(t testing.TB, got *Counter, want int) {
    t.Helper()
    if got.Value() != want {
        t.Errorf("got %d, want %d", got.Value(), want)
    }
}
```

Use `testing.TB` interface to work with both `*testing.T` and `*testing.B`.

### defer Statement

`defer` postpones a function call until the surrounding function returns.

```go
func Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()  // Executes when Inc() returns
    c.value++
}

// Even if this panics, Unlock() still runs
```

The deferred call's arguments are evaluated immediately, but the call itself waits for return.

### Don't Embed sync.Mutex

Bad:

```go
type Counter struct {
    sync.Mutex  // Lock/Unlock become public!
    value int
}
```

Good:

```go
type Counter struct {
    mu sync.Mutex  // Private field
    value int
}
```

Embedding exposes Lock/Unlock to the public API. Users might call them directly, breaking your safety guarantees and creating coupling.

## Q&A Section

### How do mutexes relate to goroutines and the "happens before" relationship?

Mutexes are NOT tied to goroutines. They're tied to data. The "happens before" relationship guarantees ordering of operations.

Without synchronization, operation order is undefined (compiler can reorder, runtime can delay). With mutex Lock/Unlock, a happens-before relationship is established: one goroutine's Unlock happens before another goroutine's Lock, so writes are visible.

### Why doesn't Go have reentrant locks?

In other languages, calling Lock() twice in the same goroutine just increments a counter and succeeds. In Go, calling Lock() twice deadlocks because you're waiting for yourself to unlock.

Go deliberately avoids reentrant locks to force explicit, clear locking design. If Lock() is nested, you need to know WHO is managing the lock and WHEN.

### What does "copies lock value" error mean?

This error means you're passing a struct containing a sync.Mutex by value. When you copy the struct, you copy the mutex too, creating two independent locks. Two goroutines each think they have "the" lock, but they have different locks.

Always pass structs containing mutexes as pointers:

```go
func assertCounter(t testing.TB, counter *Counter, want int)
```

### How does go vet detect lock copies with noCopy?

The sync package embeds a noCopy field with Lock() and Unlock() methods. go vet scans for structs with Lock/Unlock methods and warns when they're assigned, passed by value, or copied - because you should only pass pointers to such structs.

The field takes 0 memory but still prevents copying:

```go
type noCopy struct{}
func (*noCopy) Lock() {}
func (*noCopy) Unlock() {}
```

### Why does the anonymous function in the goroutine need no parameters?

Closures. The function captures variables from the outer scope automatically:

```go
counter := Counter{}
wg := WaitGroup{}

for range 1000 {
    go func() {
        counter.Inc()  // Captures counter
        wg.Done()      // Captures wg
    }()
}
```

Each goroutine gets its own closure referencing the same variables.

### Do goroutines "slack" or wait for the loop to finish?

No. Each goroutine starts running immediately when created. The loop doesn't wait - it launches all goroutines almost instantly, then moves to the next iteration. The goroutines execute in parallel on available CPU cores.

That's why you need WaitGroup to wait for all of them to finish before asserting the result.

### How to rewrite a traditional for-loop using modern Go syntax?

Use range over integers (Go 1.22+):

```go
// Old way
for i := 0; i < wantedCount; i++ {
    go func() { counter.Inc(); wg.Done() }()
}

// New way
for range wantedCount {
    go func() { counter.Inc(); wg.Done() }()
}
```

Cleaner, no manual index management.

## Key Takeaways

- Mutexes guard shared data, not goroutines
- Pass structs with mutexes by pointer only
- WaitGroup coordinates goroutine completion
- defer ensures cleanup even on panic
- Don't embed sync.Mutex - keeps locks private
- Testing concurrent code requires careful synchronization
