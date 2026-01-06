# Maps

Maps are like key-value (KV) stores or dictionaries. They allow you to store items by a key and look them up quickly.

## Declaring a Map

Declaring a map is similar to an array, but it starts with the `map` keyword and requires two types:
- **Key type**: Written inside `[]`. It must be a **comparable** type (so Go can check if keys are equal).
- **Value type**: Goes right after the `[]`. This can be any type, even another map!

```go
dictionary := map[string]string{"test": "this is just a test"}
```

## Custom Types with Maps

You can create a custom type around a map to add methods to it:

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }
    return definition, nil
}
```

## Map Properties

### Reference-like Behavior
A map value is a pointer to a `runtime.hmap` structure. This means when you pass a map to a function, you are copying the pointer, not the underlying data. You can modify the map inside a function without passing its address (no need for `*Dictionary`).

### Nil Maps
- Reading from a `nil` map is fine (returns empty/zero value).
- **Writing to a `nil` map causes a runtime panic.**
- Always initialize maps using a literal `{}` or `make()`:
  ```go
  m := Dictionary{}
  // OR
  m := make(map[string]string)
  ```

### Two-Value Lookup
Map lookups can return two values:
1. The **value** itself.
2. A **boolean** (`ok`) indicating if the key was found.

```go
definition, ok := d[word]
```

## CRUD Operations

### Create/Add
```go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

### Read/Search
(See Custom Types section above)

### Update
Updating a value uses the same syntax as adding. To make it safe, check if the word exists first.

### Delete
Go has a built-in `delete` function for maps. It takes the map and the key, and returns nothing.

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

---

## Q&A

**Q: What is this `_ err := ...` thing?**
A: The underscore (`_`) is the **Blank Identifier**. It's used when a function returns multiple values but you don't need one of them. Since Go doesn't allow unused variables, the `_` tells the compiler: "I know there's a value here, but throw it away."

**Q: Why do I get `(no value) used as value` for `Delete`?**
A: This happens when your test expects a return value (like `err := d.Delete(word)`) but your function is defined to return nothing. To fix it, either change the test to not capture the value, or update the function to return an `error`.

**Q: Why use constants for errors instead of `errors.New` in every function?**
A: `errors.New` creates a **unique** error every time it's called. If you return `errors.New("not found")` and your test compares it against *another* `errors.New("not found")`, they will not be equal. Using constants (Sentinel Errors) provides a single source of truth for equality checks.

---

## Summary

- Maps store data in key-value pairs.
- Keys must be comparable.
- Maps are pointers to a structure; you don't usually need pointer receivers to modify them.
- Always initialize maps to avoid panics on write.
- Use the blank identifier `_` to ignore unwanted return values.
- Use sentinel errors (constants) for reliable error checking in tests.
