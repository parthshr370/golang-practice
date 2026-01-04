# Go Learning Log

## Chapter 1: Hello, World

### Date: 2025-12-26

---

## Core Concepts Learned

### 1. Go Module System

- **Initialize**: `go mod init example.com/hello_testing`
- Module name must have a dot (e.g., `.com`) for tool compatibility
- Creates `go.mod` file with module name and Go version
- Required for Go 1.16+ when running `go test` or `go build`

### 2. Basic Program Structure

```go
package main  // Executable programs must be in 'main' package

import "fmt"  // Import packages for external functionality

func main() {  // Entry point for executable programs
    fmt.Println("Hello, world")
}
```

**Key Points:**
- `package main` is required for executable programs
- `func main()` is the entry point
- Packages group related Go code together

### 3. Functions

```go
func Hello() string {
    return "Hello, world"
}

func Hello(name string, language string) string {
    // function body
    return englishHelloPrefix + name
}
```

**Syntax:**
- `func` keyword defines functions
- Return type specified after function name/parameters
- Multiple parameters of same type: `(name, language string)`
- Named return values: `func greetingPrefix(language string) (prefix string)`

### 4. Constants

```go
const englishHelloPrefix = "Hello, "

// Grouped constants
const (
    spanish            = "Spanish"
    french             = "French"
    englishHelloPrefix = "Hello, "
    spanishHelloPrefix = "Hola, "
    frenchHelloPrefix  = "Bonjour, "
)
```

**Benefits:**
- Capture meaning of values
- Improve code quality
- Potential performance improvements
- Group related constants for readability

### 5. Control Flow

#### If Statements
```go
if name == "" {
    name = "World"
}
```

#### Switch Statements
```go
switch language {
case spanish:
    prefix = spanishHelloPrefix
case french:
    prefix = frenchHelloPrefix
default:
    prefix = englishHelloPrefix
}
```

**Use switch when:**
- Multiple if statements checking same value
- Better readability
- Easier to extend

### 6. Testing in Go

#### Test File Structure
```go
package main  // Same package as code being tested

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Chris")
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

#### Test Requirements
1. File must end with `_test.go`
2. Test function must start with `Test`
3. Takes single argument `t *testing.T`
4. Import `"testing"` package

#### Subtests
```go
t.Run("saying hello to people", func(t *testing.T) {
    got := Hello("Chris")
    want := "Hello, Chris"
    assertCorrectMessage(t, got, want)
})

t.Run("empty string defaults to 'world'", func(t *testing.T) {
    got := Hello("")
    want := "Hello, World"
    assertCorrectMessage(t, got, want)
})
```

**Benefits:**
- Group tests around scenarios
- Share setup code
- Better test organization

#### Helper Functions
```go
func assertCorrectMessage(t testing.TB, got, want string) {
    t.Helper()  // Marks this as a helper
    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

**Key Points:**
- `testing.TB` interface works with both `*testing.T` and `*testing.B`
- `t.Helper()` ensures failure reports correct line number
- Reduces test code duplication
- Accepts interface for flexibility (works with benchmarks)

### 7. Visibility in Go

- **Public**: Starts with capital letter (e.g., `Hello`, `Println`)
- **Private**: Starts with lowercase letter (e.g., `greetingPrefix`)
- Applies to functions, types, and package-level variables

### 8. Named Return Values

```go
func greetingPrefix(language string) (prefix string) {
    // 'prefix' variable automatically created
    switch language {
    case french:
        prefix = frenchHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return  // Returns 'prefix' without specifying
}
```

**Benefits:**
- Creates variable automatically (zero value for type)
- Makes function intent clearer in documentation
- Simplifies return for single-value returns

### 9. Variable Declaration

```go
varName := value  // Short declaration (infers type)
got := Hello()
want := "Hello, World"
```

**Format Verbs:**
- `%q` - Quote-wrapped string (great for tests)
- `%v` - Default format
- `%+v` - Struct with field names
- `%#v` - Go syntax representation
- `%T` - Type representation

### 10. TDD Cycle

1. **Write a test** - Capture requirements first
2. **Make compiler pass** - Fix compilation errors
3. **Run test, see it fail** - Verify error message is meaningful
4. **Write enough code to make test pass** - Minimum implementation
5. **Refactor** - Improve code with test safety

**Why this matters:**
- Ensures relevant tests
- Helps design good software
- Fast tests = flow state
- Failing tests show clear error messages

---

## Commands Used

```bash
go run first.go           # Compile and run
go test                   # Run tests
go test -v                # Verbose test output
go test -run TestHello    # Run specific test
go mod init <name>        # Initialize module
go fmt                    # Format code
go vet                    # Static analysis
go build                  # Build executable
```

---

## My Implementation Notes

### Function Name: Mellow
- Created my own greeting function instead of `Hello`
- Tests use subtests for different scenarios
- Helper function named `asserCorrectMessage` (intentional typo to test)

### Supported Languages
- English: "Hey"
- Spanish: "Hola"
- French: "Bonjour"
- Defaults to English when language not recognized

### Current Code Structure
- `first.go` - Main implementation with Mellow function
- `hello_test.go` - Test suite with subtests
- Uses constants for language prefixes
- Private `greetingPrefix` function for language handling
- Switch statement for language selection

---

## Key Takeaways

1. **Go's testing is built-in** - No external frameworks needed
2. **The compiler is your friend** - Listen to error messages
3. **Tests are specifications** - They describe what code should do
4. **Refactor tests too** - Keep test code clean with helpers
5. **TDD feedback loop** - Write failing test then make pass then refactor
6. **Constants over magic strings** - Improve readability and maintainability
7. **Switch over multiple ifs** - Better for checking same value repeatedly

---

## Next Steps

- **Chapter 2**: Integers - Learn about numeric types and operations
- Keep practicing TDD cycle
- Explore Go's standard library documentation with `go doc` and `pkgsite`
- Consider adding more language support to Mellow function

---

## Common Gotchas

1. **Module initialization** - Must run `go mod init` before testing in Go 1.16+
2. **Helper functions** - Need `t.Helper()` for correct error reporting
3. **Zero values** - Named return values get zero value (`""` for strings)
4. **Visibility** - Capitalization determines public/private
5. **Package main** - Required for executable programs

---

## Resources

- Official Go documentation: https://pkg.go.dev/
- `go doc fmt` - View package docs offline
- `go install golang.org/x/pkgsite/cmd/pkgsite@latest` - Local docs server
