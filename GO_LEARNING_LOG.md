# Go Learning Log

*This file will continuously grow as I learn Go concepts through practice and testing.*

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

### Current State
- `Mellow(name string)` returns greetings
- Empty string defaults to "World"
- Uses `englishHelloPrefix` constant

### Next Steps (Learn Go with Tests - Chapter 1)

**1. Add language parameter**
```go
func Mellow(name, language string) string { ... }
```

**2. Add constants for languages**
```go
const (
    englishHelloPrefix = "Hey"
    spanishHelloPrefix = "Hola, "
    frenchHelloPrefix  = "Bonjour, "
)
```

**3. Use switch for language selection**
```go
prefix := englishHelloPrefix
switch language {
case "Spanish":
    prefix = spanishHelloPrefix
case "French":
    prefix = frenchHelloPrefix
}
return prefix + name
```

**4. Extract to helper function**
```go
func greetingPrefix(language string) (prefix string) {
    switch language {
    // ...
    default:
        prefix = englishHelloPrefix
    }
    return
}
```

### New Concepts
- **Switch statements**: Cleaner than multiple if statements
- **Named return values**: `func foo() (x string) { return }`
- **Default case**: Catches unmatched switch cases
- **Public vs Private**: Capital letter = exported/public

### TDD Cycle Reminder
1. Write failing test
2. Make compiler pass
3. See test fail with clear error
4. Write minimal code to pass
5. Refactor

---

*Next: Chapter 2 - Integers*