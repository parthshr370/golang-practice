# Reflection

Reflection is the ability of a program to examine its own structure, particularly through types. It's a form of metaprogramming that allows you to inspect and manipulate objects at runtime.

## The `interface{}` (or `any`) Gateway

In Go, `interface{}` is the **Empty Interface**. It defines zero methods. Because it has no requirements, **every single type in Go satisfies it.**

When you pass a variable into a function taking `interface{}`, you are putting it into a **locked, opaque box.** The compiler can no longer see what's inside.

### Reflection: The "X-Ray"
Reflection is the tool that allows you to look through that blind box.
- **`reflect.ValueOf(x)`**: Returns a `Value` object—an "X-ray" of your variable.
- **`val.Kind()`**: Identifies the category (Struct, String, Pointer, etc.).
- **`val.Interface()`**: The reverse operation; turns a reflected value back into an `interface{}`.

---

## When to use what? (The Decision Guide)

| Tool | Focus | When to use |
| :--- | :--- | :--- |
| **Struct** | **Data/State** | To group related variables together (e.g., `User`, `Point`). |
| **Type** | **Meaning** | To give primitives domain-specific names (e.g., `type Bitcoin int`). |
| **Interface** | **Behavior** | To define a contract of what a thing **does** (e.g., `io.Writer`). |
| **Any (`interface{}`)** | **Generic Box** | Only when you literally cannot know the type at compile time. |

**The Go Proverb:** "Accept interfaces, return structs."

---

## Technical Patterns in this Chapter

### 1. Handling Structs and Fields
To find string fields inside a struct:
```go
val := reflect.ValueOf(x)
for i := 0; i < val.NumField(); i++ {
    field := val.Field(i)
    if field.Kind() == reflect.String {
        fn(field.String())
    }
}
```

### 2. The Pointer Trap
Pointers don't have fields; they have addresses. You must follow the pointer to the data using `Elem()`.
```go
if val.Kind() == reflect.Pointer {
    val = val.Elem() // Follow the pointer to the underlying struct
}
```

### 3. Recursive Inspection (The Recursive Switch)
The most robust way to handle any type is a `switch` on the `Kind()`. If a field is a struct, slice, map, etc., we call `walk` recursively.
```go
func walk(x interface{}, fn func(string)) {
    val := getValue(x) // Handles Pointers

    walkValue := func(value reflect.Value) {
        walk(value.Interface(), fn)
    }

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        for i := 0; i < val.NumField(); i++ {
            walkValue(val.Field(i))
        }
    case reflect.Slice, reflect.Array:
        for i := 0; i < val.Len(); i++ {
            walkValue(val.Index(i))
        }
    // ... handles Map, Chan, Func ...
    }
}
```

### 4. Handling Collections
| Type | Inspection Method |
| :--- | :--- |
| **Slice / Array** | `val.Len()` and `val.Index(i)` |
| **Map** | `val.MapKeys()` and `val.MapIndex(key)` |
| **Channel** | `val.Recv()` (Pulls data until closed) |
| **Function** | `val.Call(nil)` (Executes and inspects return values) |

---

## Final Lesson
Reflection code is messy and dangerous. 
- **Assignment Mismatch**: If you call `Field()` on something that isn't a struct, Go **Panics**.
- **Invisibility**: You lose all the benefits of the compiler checking your types.

As the chapter says: **"Now that you know about reflection, do your best to avoid using it."** Use it for library level code (serializers, ORMs) but avoid it in your business logic.

---

## Summary: The Reflection Warnings

1.  **No Type Safety**: Errors that usually happen at compile time will now cause **Panics** at runtime.
2.  **Performance Hit**: Inspecting types at runtime is significantly slower than using concrete types.
3.  **Complexity**: Code becomes harder to read and maintain.

**Rule of Thumb:** Only use reflection if you absolutely must (e.g., writing a JSON serializer or a DB driver). Otherwise, stick to Interfaces and Concrete Types.
