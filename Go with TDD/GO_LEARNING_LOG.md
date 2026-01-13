# Go Learning Log

_This file will continuously grow as I learn Go concepts through practice and testing._

## Session 1: Basic Go Structure

### What I Learned Today

#### `package main`

- Every Go file starts with a package declaration
- `main` is a special package that creates an executable program
- Only one `main` function allowed per package

#### `func main()`

- Entry point of the program
- Where execution begins
- No parameters, no return value

#### `import "fmt"`

- `fmt` = Format package
- Provides I/O functions like `Println()`
- Used for console output and formatted strings

#### `func Hello() string`

- Function declaration with return type `string`
- Returns a string value
- Can be called from other functions

#### `testing` Package

- Go's built-in testing framework
- Test files must end with `_test.go`
- Test functions must start with `Test` and take `*testing.T` parameter

#### Test Structure

```go
func TestHello(t *testing.T) {
    got := Hello()           // actual result
    want := "hello world"    // expected result

    if got != want {
        t.Errorf("got %q want %q", got, want)  // fail test
    }
}
```

#### Key Go Concepts

- **Strong typing**: Must declare return types
- **Exported names**: Capitalized names (like `Hello`) are public
- **Semicolons**: Not needed (Go inserts them automatically)
- **String formatting**: `%q` adds quotes around strings in output

## Commands Used

- `go test` - runs all tests in current directory
- `go mod` - manages Go modules

## Errors Encountered & Solutions

1. **"main redeclared"**: Multiple `main()` functions in same package
   - Solution: Only one main function per package
2. **"[no test files]"**: Test file not named correctly
   - Solution: Test files must end with `_test.go`
3. **Test failure**: Expected vs actual values don't match
   - Solution: Check function output matches test expectation

## Session 2: Testing Framework Deep Dive

### Understanding `t *testing.T` Parameter

#### Breaking down `t *testing.T`:

- **`t`**: Variable name (convention like `err` for errors)
- **`*`**: Pointer to the struct (pass-by-reference)
- **`testing`**: Go's built-in testing package
- **`T`**: Struct type that provides test functionality

#### Why pointer?

- Pass-by-reference: Changes affect original test object
- More efficient than copying entire struct
- Allows test framework to track results across function

#### Key methods available:

- `t.Error(msg)` - Mark test failed but continue
- `t.Fatal(msg)` - Mark test failed and stop immediately
- `t.Logf(format, args)` - Log information during test
- `t.Parallel()` - Mark test safe for parallel execution
- `t.Skip()` - Skip this test

#### Usage pattern:

```go
func TestHello(t *testing.T) {
    got := Hello()
    want := "Hello, world"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

### New Commands Learned

- `go test -run TestHello` - Run specific test function
- `go run .` - Run all main files in directory
- `./hello_testing` - Execute built binary
- `go test -v` - Verbose test output
- `go test -cover` - Test coverage report

### Learning Pod Setup

- Created CLAUDE.md with learning approach guidelines
- Established "Build → Test → Run → Log" workflow
- Set rule: Claude guides, user writes all code
- Learning progression: basics → intermediate → advanced Go concepts

---

## Session 3: Extending Hello Function (Language Support)

- Switch statements: Cleaner than multiple if statements
- Named return values: `func foo() (x string) { return }`
- Default case: Catches unmatched switch cases
- Public vs Private: Capital letter = exported/public

---

## 3

devlog on my golang journey

1. Completed **Pointers & Errors** and **Maps** modules from the TDD series.
2. Deep dive into **Memory Addresses and Pointers**: Understanding how Go copies values by default and using pointer receivers (`*Type`) to mutate state.
3. Explored **Custom Type Definitions**: Learned how to create domain-specific types (e.g., `type Bitcoin int`) and implement the `Stringer` interface for custom printing.
4. Studied **Sentinel Errors**: Implementing package-level error constants for immutable, reusable, and testable error handling.
5. Understanding **Maps Internals**: Learned that maps are pointers to `hmap` structures (making them feel like reference types) and the nuances of handling `nil` maps.
6. Refined usage of the **Blank Identifier (`_`)** for ignoring return values and handling the "two-value lookup" property of maps.
7. Experimented with **errcheck**: Using linters to identify unchecked errors in the test suite and implementation.

---

## 4

devlog on my golang journey

1. Completed **Dependency Injection**, **Mocking**, and **Concurrency** modules.
2. **Dependency Injection**: Learned to inject `io.Writer` (the "Universal Plug") to make code testable, decoupled, and reusable across terminals, files, and web servers.
3. **Mocking & Spying**: Implemented custom Spies to record the _order_ of operations (e.g., Sleep -> Write -> Sleep), verifying behavior without relying on real-time delays.
4. **Concurrency Internals**: Deep dive into Goroutines vs. OS Threads, the Go Scheduler (M:N model), and how `main` kills background processes.
5. **Channels**: Visualized channels as "Pipelines" that prevent Race Conditions by serializing access to shared memory ("Share memory by communicating").
6. **Go Standard Library**: Explored the internals of `fmt.Fprintf`, `bytes.Buffer` (the testing secret weapon), and variadic functions.

---

## 5

devlog on my golang journey

1. Completed the **Select** module, mastering multi-channel synchronization.
2. **The `select` Switchboard**: Learned how to wait for multiple concurrent events and act on the first one that completes (racing goroutines).
3. **Signaling with `chan struct{}`**: Used empty structs for memory-efficient signaling and understood why `close(ch)` is a powerful broadcast mechanism.
4. **Timeouts**: Implemented robust error handling using `time.After` to prevent system hangs in concurrent code.
5. **HTTP Testing**: Leveraged `net/http/httptest` to spin up mock servers, making HTTP tests fast, reliable, and independent of external networks.
6. **Channel Lifecycle**: Deep dive into the importance of `make(chan)` versus `nil` channels to avoid permanent blocking and runtime panics.

---

## 6

devlog on my golang journey

1. Completed the **Reflection** module, mastering runtime introspection.
2. **The `interface{}` "Blind Box"**: Learned that `any` (alias for `interface{}`) accepts any type but obscures its structure from the compiler.
3. **The `reflect` Package**: Used `reflect.ValueOf` to "X-ray" interfaces and inspect their `Kind`, `NumField`, and values at runtime.
4. **Recursive Walking**: Implemented a generic `walk` function that can traverse deeply nested structs, maps, slices, and arrays using recursion.
5. **Handling Pointers**: Solved the panic trap by using `val.Elem()` to dereference pointers before inspecting their fields.
6. **The "Decision Guide"**: Learned when to use Structs (Data), Interfaces (Behavior), and Reflection (Generic Libraries) - and why Reflection should be a last resort due to lack of type safety.
