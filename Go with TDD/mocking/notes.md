We know we want our Countdown function to write data somewhere and io.Writer is the de-facto way of capturing that as an interface in Go.

- In **main**, we will send to `os.Stdout` so our users see the countdown printed to the terminal.
- In **test**, we will send to `bytes.Buffer` so our tests can capture what data is being generated.

We're using `fmt.Fprint` which takes an `io.Writer` (like `*bytes.Buffer`) and sends a string to it.

---

## Pointers vs Values in DI

When you see `buffer := &bytes.Buffer{}`, you are dealing with a core Go concept: The difference between a Value and a Pointer.

### 1. What `&bytes.Buffer{}` does:

- `bytes.Buffer{}`: This creates a new, empty Buffer **value** on the stack.
- `&`: This is the "address of" operator. It finds the memory location where that buffer lives.
- **Result**: The variable `buffer` doesn't hold the "bucket" itself; it holds a **pointer** (the memory address) to that bucket.

### 2. Why do we do this?

There are two main reasons we use the address (`&`) here:

**A. Satisfying the Interface**
If you look at the `bytes.Buffer` source code, the `Write` method is defined like this:
`func (b *Buffer) Write(p []byte) (n int, err error)`
Notice the `*Buffer` receiver. This means only a **pointer** to a Buffer can do the writing. If you passed the value instead of the address, Go might complain that it doesn't satisfy the `io.Writer` interface.

**B. Mutability (Efficiency)**
A `bytes.Buffer` can grow quite large.

- If you pass the **Value**, Go makes a full copy of the entire buffer every time you call a function. If your buffer has 1MB of text, you just wasted 1MB of memory making a copy.
- If you pass the **Pointer** (`&`), you are only passing a tiny 8-byte address. The function then reaches back to the original memory location to add more text.

### Summary

When you see `buffer := &bytes.Buffer{}:`

- **Value of buffer**: `0x1234abcd` (a memory address).
- **What's at that address**: The actual "bucket" holding your bytes.

In your tests, this is standard practice because we want our `Countdown` function to write into the **original** buffer so we can read it later using `buffer.String()`.

---

## io.Writer vs bytes.Buffer

- **`io.Writer` is the Interface**: It is just a set of rules (the "contract"). It says: "I don't care who you are, as long as you have a `Write` method."
- **`bytes.Buffer` is the Implementation**: It is the actual struct that follows those rules. It uses an internal byte array (`[]byte`) to store the data you give it.

When you use `io.Writer` in your function signature, you are allowing your code to work with _anything_ that follows the rules—whether it's the terminal, a file, or an in-memory `bytes.Buffer`.

The tests still pass and the software works as intended but we have some problems:

    Our tests take 3 seconds to run.

        Every forward-thinking post about software development emphasises the importance of quick feedback loops.

        Slow tests ruin developer productivity.

        Imagine if the requirements get more sophisticated warranting more tests. Are we happy with 3s added to the test run for every new test of Countdown?

    We have not tested an important property of our function.

We have a dependency on Sleeping which we need to extract so we can then control it in our tests.

If we can mock time.Sleep we can use dependency injection to use it instead of a "real" time.Sleep and then we can spy on the calls to make assertions on them.

Idiomatic is just another way of saying conventional. It is just a generally accepted way of doing things in the language. Kind of overused, like epic.
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//
// - Slice: The size specifies the length. The capacity of the slice is
// equal to its length. A second integer argument may be provided to
// specify a different capacity; it must be no smaller than the
// length. For example, make([]int, 0, 10) allocates an underlying array
// of size 10 and returns a slice of length 0 and capacity 10 that is
// backed by this underlying array.
// - Map: An empty map is allocated with enough space to hold the
// specified number of elements. The size may be omitted, in which case
// a small starting size is allocated.
// - Channel: The channel's buffer is initialized with the specified
// buffer capacity. If zero, or the size is omitted, the channel is
// unbuffered.

func reflect.DeepEqual(x any, y any) bool
DeepEqual reports whether x and y are “deeply equal,” defined as follows. Two values of identical type are deeply equal if one of the following cases applies. Values of distinct types are never deeply equal.
Array values are deeply equal when their corresponding elements are deeply equal.
Struct values are deeply equal if their corresponding fields, both exported and unexported, are deeply equal.
Func values are deeply equal if both are nil; otherwise they are not deeply equal.
Interface values are deeply equal if they hold deeply equal concrete values.
Map values are deeply equal when all of the following are true: they are both nil or both non-nil, they have the same length, and either they are the same map object or their corresponding keys (matched using Go equality) map to deeply equal values.
Pointer values are deeply equal if they are equal using Go’s == operator or if they point to deeply equal values.
Slice values are deeply equal when all of the following are true: they are both nil or both non-nil, they have the same length, and either they point to the same initial entry of the same underlying array (that is, &x[0] == &y[0]) or their corresponding elements (up to length) are deeply equal. Note that a non-nil empty slice and a nil slice (for example, []byte{} and []byte(nil)) are not deeply equal.

---

## Mocking and Spying: The Camera Analogy

We have two different dependencies (Writing and Sleeping) and we want to record all of their operations into one list. So we'll create one spy for them both.

### The Security Camera (`SpyCountdownOperations`)
Imagine a camera pointing at your code. It has a tape (`Calls []string`) inside it.

Your `Countdown` function needs two things:
1. **A Printer** (`io.Writer`) to show the numbers.
2. **A Sleeper** (`Sleeper`) to pause between them.

We want **one camera** to record both activities.

### How the Spy Works
The spy has two "lenses" (methods):

1. **`Write()`**: When the code calls this, the spy stamps the word **"write"** onto its tape.
2. **`Sleep()`**: When the code calls this, the spy stamps the word **"sleep"** onto its tape.

Because it's the **same tape**, the words appear in the exact order they happened.

```go
type SpyCountdownOperations struct {
    Calls []string
}

func (s *SpyCountdownOperations) Sleep() {
    s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
    s.Calls = append(s.Calls, write)
    return
}
```

### Why do we need this?
We want to prove that the code sleeps **between** the prints.

If the test shows:
`["write", "sleep", "write", "sleep", "write", "sleep", "write"]`
-> The code is correct!

If the test shows:
`["sleep", "write", "sleep", "write", "sleep", "write", "write"]`
-> The code is broken! It slept before printing the first number.

By checking the **sequence** in `s.Calls`, we are testing the **logic and timing** of the program, not just the output.

---

## Configurable Sleepers

We can make our `Sleeper` smarter by allowing it to configure how long to sleep.

### The ConfigurableSleeper
This struct holds a `duration` (how long) and a `sleep` function (how to do it).

```go
type ConfigurableSleeper struct {
    duration time.Duration
    sleep    func(time.Duration)
}
```

**Why a function for sleeping?**
- In the **real app**, we pass `time.Sleep` (the standard library function).
- In the **test**, we pass a `SpyTime` function that just records how long it *would* have slept, so we can verify the configuration is correct.

### The Final Refactor
Now we can use our configurable sleeper in `main`:

```go
func main() {
    sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
    Countdown(os.Stdout, sleeper)
}
```

---

## When to Mock (and when not to)

Mocking is a tool. It's not evil, but it can be overused.

### When Mocking is Good:
- **Testing Side Effects**: How do you test that code paused for 1 second? Mock the timer.
- **Testing Failure**: How do you test what happens if the database crashes? Mock the database.
- **Speed**: You don't want to wait 10 seconds for every test to run.

### Signs of Bad Design:
If you find yourself creating 5 mocks just to test one simple function, it's usually a sign that:
1. Your function is doing too much (too many dependencies).
2. Your dependencies are too fine-grained (consolidate them).
3. You are testing implementation details instead of useful behavior.

**Rule of Thumb**: If a test uses more than 3 mocks, it's time to rethink the design.
