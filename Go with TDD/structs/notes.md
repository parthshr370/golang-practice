# Structs, Methods & Interfaces

## What You'll Learn

- **Structs**: How to declare your own data types to bundle related data together.
- **Methods**: How to add functionality to your types using method receivers.
- **Interfaces**: How to define abstract behaviors to create decoupled and flexible code (polymorphism).
- **Table Driven Tests**: A pattern for writing clear, extensible test suites.

## Concepts

### Structs
A struct is a named collection of fields where you can store data.
```go
type Rectangle struct {
    Width  float64
    Height float64
}
```

### Methods
A method is a function with a receiver. It binds a function to a specific type.
```go
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```
In this example, `(r Rectangle)` is the receiver.

### Interfaces
Interfaces allow you to define functions that can be used by different types. Interface resolution is implicit in Go; if a type has the methods declared in the interface, it satisfies the interface.
```go
type Shape interface {
    Area() float64
}
```

## Table Driven Tests
Table driven tests are useful when you want to build a list of test cases that can be tested in the same manner.

```go
func TestArea(t *testing.T) {
    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }{
        {name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
        {name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
        {name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
    }

    for _, tt := range areaTests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.shape.Area()
            if got != tt.hasArea {
                t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
            }
        })
    }
}
```

## Q&A

**Q: Why use `%g` instead of `%f` in format strings?**
A: `%g` prints a more precise decimal number in the error message, which is helpful for debugging subtle floating-point differences.

**Q: What is the `%#v` format string?**
A: `%#v` prints out the struct with the values in its fields (Go syntax representation), so you can see exactly what properties are being tested.

**Q: Why use `t.Run`?**
A: `t.Run` allows you to run specific sub-tests and provides clearer output on failures, showing exactly which case in the table failed.

**Q: Can you overload functions in Go (e.g., `Area(Circle)` and `Area(Rectangle)`)?**
A: No, Go does not support function overloading. You must use different function names or, better yet, methods on types.

## Appendix: Strings, Bytes, and Runes

Go handles text differently than many other languages. A `string` is not just a list of characters; it's a **read-only slice of bytes**.

### The Core Difference

*   **Byte (`byte`)**: A basic unit of storage (uint8). It can store ASCII characters like 'A' (65).
*   **Rune (`rune`)**: A Unicode Code Point (int32). It represents a single "character" concept, which might take multiple bytes to store (e.g., '⌘').
*   **String**: A sequence of bytes that *usually* represents UTF-8 encoded text.

### Example: "a⌘"

If you have the string `"a⌘"`:
1.  **Bytes**: It takes **4 bytes**.
    *   'a' = 1 byte (ASCII)
    *   '⌘' = 3 bytes (UTF-8 encoding)
    *   `len("a⌘")` returns **4**.
2.  **Runes**: It has **2 runes** (characters).
    *   `utf8.RuneCountInString("a⌘")` returns **2**.

### Iterating over Strings

*   **Standard Loop (`i++`)**: Iterates over **bytes**.
    ```go
    for i := 0; i < len(str); i++ {
        // accesses bytes, not necessarily full characters
    }
    ```
*   **Range Loop**: automatically decodes **runes**.
    ```go
    for i, r := range str {
        // r is a rune (the character)
        // i is the starting byte index of that rune
    }
    ```

**Summary**: When you need to process text character-by-character (like reversing a string or counting letters), think **Runes**. When you need to store or transmit data, think **Bytes** (Strings).
