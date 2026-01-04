# Chapter 2: Integers - TDD & Documentation

## Project Structure

```
GO Practice/
├── go.mod
├── integers/
│   ├── adder_test.go
│   └── second.go
└── helloWorld/
    ├── first.go
    └── hello_test.go
```

**Important Rule:** Go source files can only have one package per directory

## Test Naming Convention

```go
// WRONG - lowercase t
func testadder(t *testing.T) {}

// CORRECT - uppercase T
func TestAdder(t *testing.T) {}
```

The first letter must be an uppercase `T`. The testing tool will ignore functions that don't follow this convention.

## TDD Workflow

### 1. Write Test First

```go
func TestAdder(t *testing.T) {
    sum := Add(2, 2)
    expected := 4

    if sum != expected {
        t.Errorf("expected '%d' but got '%d'", expected, sum)
    }
}
```

### 2. Run Test - Compilation Error

```bash
$ go test
./adder_test.go:6:9: undefined: Add
```

### 3. Write Minimal Code to Compile

```go
func Add(x, y int) int {
    return 0
}
```

### 4. Run Test - Failing Test

```bash
$ go test
adder_test.go:10: expected '4' but got '0'
```

### 5. Write Code to Make Pass

```go
func Add(x, y int) int {
    return x + y
}
```

### 6. Run Test - Passing

```bash
$ go test
PASS
```

## Multiple Arguments of Same Type

```go
// Long form
func Add(x int, y int) int {}

// Short form (preferred)
func Add(x, y int) int {}
```

## Format Strings

- `%d` - for integers
- `%q` - for strings (quoted)

```go
t.Errorf("expected '%d' but got '%d'", expected, sum)
```

## Named Return Values

Use named return values when the meaning isn't clear from context.

```go
// Named - useful when purpose isn't obvious
func Calculate() (result int) {
    // ...
    return  // returns result automatically
}

// Not needed - Add is obvious
func Add(x, y int) int {
    return x + y
}
```

## Documentation Comments

Comments appear in `go doc`:

```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
    return x + y
}
```

## Testable Examples

Example functions are compiled and tested with your test suite:

```go
import "fmt"

func ExampleAdd() {
    sum := Add(1, 5)
    fmt.Println(sum)
    // Output: 6
}
```

**Why use examples?**
- README examples get outdated
- Example functions are compiled and validated
- If code changes, build fails - you know to update docs

**Example vs Test:**
- `ExampleAdd()` - appears in documentation
- `TestAdder()` - regular test

**The `// Output:` comment is required** for the example to run as a test.

```bash
$ go test -v
=== RUN   TestAdder
--- PASS: TestAdder (0.00s)
=== RUN   ExampleAdd
--- PASS: ExampleAdd (0.00s)
```

## pkgsite - Documentation Viewer

View documentation in web browser:

```bash
# Install pkgsite
go install golang.org/x/pkgsite/cmd/pkgsite@latest

# Run in project root
cd "/home/parthshr370/Downloads/GO Practice"
pkgsite -open .
# Opens http://localhost:8080
```

Shows:
- All packages
- Functions with documentation
- Example code nicely formatted

## pkg.go.dev - Public Documentation

Public site for Go documentation:

`https://pkg.go.dev/github.com/quii/learn-go-with-tests`

Publish your code to GitHub to get public docs.

## Key Concepts

### TDD Red-Green-Refactor

1. **Red** - Write failing test
2. **Green** - Write minimum code to pass
3. **Refactor** - Improve code

### Packages

- `package main` - executable programs
- `package integers` - library code (can be imported)
- Package name should match directory name
- All lowercase, no punctuation

### Testing Tools

- `go test` - run tests
- `go test -v` - verbose output
- Test files end with `_test.go`
- Test functions start with `Test`

---

## Structs, Interfaces, and Pointers Review

### Struct = Data + Methods

```go
type Car struct {
    brand string
    speed int
}

func (c Car) Drive() {
    fmt.Println("Vroom!")
}
```

### Interface = Contract

```go
type Vehicle interface {
    Drive()
}

// Car implements Vehicle automatically (has Drive method)
```

### Why `*testing.T` Needs Pointer?

```go
type T struct {
    failed bool
    errors []string
}

// Needs pointer to modify internal state
func (t *testing.T) Errorf(msg string) {
    t.failed = true
    t.errors = append(t.errors, msg)
}
```

### Why `testing.TB` Doesn't Need Pointer?

```go
// Interface is just a method list, no data
type TB interface {
    Errorf(msg string)
    Helper()
}

// *testing.T implements TB automatically
```

### Key Insight

- `*testing.T` - pointer to struct (modifies state)
- `testing.TB` - interface type (just contract)
- `*testing.T` can be passed where `testing.TB` is expected

### Simple Analogy

- **Struct** = Actual car with data
- **Interface** = "Vehicle" category
- **Pointer `*`** = Keys to modify the car

---

## Summary

- TDD workflow: test then fail then fix then pass
- Multiple args of same type: `(x, y int)`
- Documentation comments with `//`
- Example functions for tested documentation
- `pkgsite` for local documentation viewer
- `pkg.go.dev` for public documentation
