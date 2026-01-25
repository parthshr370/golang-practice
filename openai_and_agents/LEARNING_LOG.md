# OpenAI & Agents Learning Log

_Notes from learning Go with OpenAI and Agents concepts._

## Session 1: JSON Marshaling & Interfaces

### What I Learned Today

#### `interface{}` - The Empty Interface

- `interface{}` means "any type" - it's Go's way of saying "I'll accept anything"
- Every Go type satisfies `interface{}` (even basic types like `int`, `string`, `struct{}`)
- At compile time, no type checking - compiler trusts the code is correct
- At runtime, you can use reflection to inspect the actual type

```go
func describe(v interface{}) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}

describe(42)        // Type: int, Value: 42
describe("hello")   // Type: string, Value: hello
describe(struct{X int}{10}) // Type: main.Struct, Value: {10}
```

#### `json.Marshal(v interface{}) ([]byte, error)`

- Converts any Go value to JSON format
- Returns `[]byte` (raw bytes) and `error`
- Works with: structs, maps, slices, arrays, primitives
- Does NOT work with: channels, functions, complex numbers

```go
type Message struct {
    Name string
    Body string
    Time int64
}

m := Message{"Alice", "Hello", 1294706395881547000}
b, _ := json.Marshal(m)
// b = []byte(`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`)
```

#### Structs & JSON Field Mapping

- By default, struct field names become JSON keys
- Use tags to customize: `json:"custom_name"`
- Unexported fields (lowercase) are ignored

```go
type Person struct {
    Name string `json:"name"`           // exported, tagged
    age  int    `json:"-"`              // unexported, ignored
    City string                        // exported, default name "City"
}
```

#### Common JSON Marshal Options

```go
// Omit empty fields
type Data struct {
    Name string `json:"name,omitempty"`  // absent if empty
}

// Ignore field
type Data struct {
    Secret string `json:"-"`             // never appears
}

// Inline/flatten
type Inner struct {
    X, Y int
}
```

### Commands Used

- `go run .` - Run the program
- `go build .` - Compile to binary
- `go fmt .` - Format code

### Key Concepts

| Term | Meaning |
|------|---------|
| Marshal | Convert Go value → JSON bytes |
| Unmarshal | Convert JSON bytes → Go value |
| Serialization | Converting data for storage/transmission |
| Deserialization | Reconstructing data from serialized form |

### Errors Encountered & Solutions

1. **"undefined: json.const"** - `json` package doesn't have `const()`
   - Solution: Use `json.Marshal()` function
2. **"non-declaration statement outside function body"** - Can't assign outside functions
   - Solution: All assignments must be inside a function like `main()`
3. **"cannot use field (string) in struct as type string in map literal"** - Type mismatch
   - Solution: Ensure struct field types match values

### Practice Code

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Message struct {
    Name string
    Body string
    Time int64
}

func main() {
    m := Message{
        Name: "Alice",
        Body: "Hello!",
        Time: 1294706395881547000,
    }

    // Marshal to JSON
    b, err := json.Marshal(m)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(string(b))
    // Output: {"Name":"Alice","Body":"Hello!","Time":1294706395881547000}
}
```

---

## Session 2: Working with JSON in Practice

### JSON to Go Mapping

| JSON Type | Go Type |
|-----------|---------|
| `null` | `nil` |
| `true/false` | `bool` |
| `number` | `float64` (unmarshal), any numeric type (marshal) |
| `string` | `string` |
| `array` | `[]any` or `[]T` |
| `object` | `map[string]any` or `struct{}` |

### Unmarshal (JSON → Go)

```go
var m Message
err := json.Unmarshal(b, &m)
if err != nil {
    log.Fatal(err)
}
fmt.Println(m.Name)  // Access unmarshaled data
```

### Pretty Print (Indent)

```go
b, _ := json.MarshalIndent(m, "", "  ")
fmt.Println(string(b))
```

Output:
```json
{
  "Name": "Alice",
  "Body": "Hello!",
  "Time": 1294706395881547000
}
```

---

## Session 3: Deep Dive - What IS Marshal?

### Marshal is NOT Built-in

`json.Marshal` is just a function from Go's standard library - NOT magic built into the language. Someone wrote it!

```go
// This is real Go code that exists in the standard library
func Marshal(v interface{}) ([]byte, error) {
    // Look at the type of v
    // If it's a struct, read each field
    // Convert field names to strings
    // Convert values to JSON format
    // Return bytes
}
```

### The Simple Flow

```
Go Struct          Marshal          JSON Bytes          HTTP/Send/Store
┌──────────┐    ┌──────────┐    ┌──────────────┐    ┌──────────────┐
│ Person { │───►│ Marshal  │───►│ []byte{      │───►│ Send over    │
│   Name   │    │          │    │   123, 34... │    │ network/file │
│   Age    │    │          │    │ }            │    │              │
└──────────┘    └──────────┘    └──────────────┘    └──────────────┘
                                              ▲
                                              │
                                         Unmarshal
                                         (reverse)
```

### `v interface{}` Explained

- `interface{}` means "accept ANY type"
- Think of it as a "blind box" - the function doesn't know what's inside
- The function just says: "I don't care what you give me, I'll convert it to JSON"

```go
// These all work:
json.Marshal("hello")           // string
json.Marshal(42)                // int
json.Marshal(map[string]int{"a": 1})  // map
json.Marshal(Person{Name: "A"}) // struct
```

### Bytes vs String

**Important:** The bytes returned BY Marshal ARE the JSON - not something before/after.

```go
b, _ := json.Marshal(Person{Name: "Alice"})
// b = []byte{123, 34, 78, 97, 109, 101, 34, 58, 34, 65, 108, 105, 99, 101, 34, 125}
//    = {"Name":"Alice"} (just in raw bytes)

fmt.Println(string(b))  // Convert to string for human reading
```

### Why Bytes?

| Reason | Explanation |
|--------|-------------|
| Efficiency | Bytes are what networks/files use natively |
| Flexibility | Can convert to string (`string(b)`) or keep as bytes |
| Binary data | Some JSON is base64-encoded (binary-safe) |

### Marshal vs Unmarshal

| Function | Direction | Example |
|----------|-----------|---------|
| `Marshal` | Go → JSON | `json.Marshal(person)` returns `[]byte` |
| `Unmarshal` | JSON → Go | `json.Unmarshal(bytes, &person)` returns struct |

### Printf vs Println

**Wrong:**
```go
fmt.Println("name: %s, age: %d", name, age)  // Ignores format verbs!
```

**Correct:**
```go
fmt.Printf("name: %s, age: %d\n", name, age)  // Processes format verbs
```

| Function | Format Verbs | Use Case |
|----------|--------------|----------|
| `fmt.Print` | No | Simple output |
| `fmt.Printf` | Yes | Formatted output |
| `fmt.Println` | No | Values with space separator + newline |

### Common Format Verbs

| Verb | Type |
|------|------|
| `%s` | string |
| `%d` | int/decimal |
| `%f` | float |
| `%v` | any value (default format) |
| `%T` | type of value |
| `%#v` | Go-syntax format |

```go
b, err := json.Marshal(person)
fmt.Printf("bytes: %s, err: %v\n", b, err)  // Show JSON as readable string
fmt.Printf("bytes: %v, err: %v\n", b, err)  // Show raw bytes
```

---

## Next Steps

- Learn about `json.Decoder` vs `json.Unmarshal`
- Explore streaming JSON (large files)
- Study `encoding/json` performance alternatives (e.g., `ffjson`, `easyjson`)
