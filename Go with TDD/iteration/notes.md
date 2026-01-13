# Chapter 3: Iteration - Loops & Benchmarks

## Go Testing Conventions

- Test files must have `_test.go` suffix
- Test files should be in same directory as code
- Run all tests recursively: `go test -v ./...`

## Modern For Loops (Go 1.22+)

Go recently added support for range over integers, making loops cleaner.

### Traditional Loop

```go
for i := 0; i < repeatCount; i++ {
    repeated += character
}
```

### Modern Range Loop

```go
// "for range integer"
for range repeatCount {
    repeated += character
}
```

**Note:** If you don't use the loop variable (like `i`), you can omit it.

## Benchmarking in Go

Benchmarks are first-class citizens in Go's testing package.

### Structure

```go
func BenchmarkRepeat(b *testing.B) {
    // Setup code (not timed)

    for b.Loop() {
        // Code to measure
        Repeat("a")
    }

    // Cleanup code (not timed)
}
```

### b.Loop() (Go 1.24+)

- `b.Loop()` automatically manages the timer
- Resets timer after setup
- Stops timer before cleanup
- Prevents compiler from optimizing away the loop body
- Runs benchmark function once per measurement (unlike old `b.N` style)

### Running Benchmarks

```bash
go test -bench=.
```

**Output Analysis:**

```text
BenchmarkRepeat-8   39271040   29.37 ns/op
```

- **39271040**: Number of iterations run
- **29.37 ns/op**: Average time per operation (nanoseconds)

## String Performance Optimization

Strings in Go are **immutable**.

- `+` operator creates a new string every time (copying memory)
- This gets slow with many concatenations

### The Solution: `strings.Builder`

`strings.Builder` minimizes memory copying by buffering writes.

```go
func Repeat(character string) string {
    var repeated strings.Builder
    for range repeatCount {
        repeated.WriteString(character)
    }
    return repeated.String()
}
```

**Performance Impact:**

- Significant reduction in ns/op (nanoseconds per operation)
- Fewer memory allocations
- Much faster for heavy string operations

## Key Takeaways

1. **For is the only loop**: Go has no `while` or `do-while`
2. **Braces required**: `{ }` are always mandatory
3. **Short declaration**: `:=` declares and initializes
4. **Benchmarking**: Built-in tool to measure performance
5. **Optimization**: Use `strings.Builder` for building strings in loops

There's a new way to create a slice. make allows you to create a slice with a starting capacity of the len of the numbersToSum we need to work through.

append helps us rather create a new slice and then play with it -

However, you can use the append function which takes a slice and a new value, then returns a new slice with all the items in it.

func SumAll(numbersToSum ...[]int) []int {
var sums []int
for \_, numbers := range numbersToSum {
sums = append(sums, Sum(numbers))
}

    return sums

}

In this implementation, we are worrying less about capacity. We start with an empty slice sums and append to it the result of Sum as we work through the varargs.

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
func make(t Type, size ...IntegerType) Type

append manages memory allocation and capacity automatically, which simplifies your code.
With make (Manual):

- You must know the exact size upfront.
- You must manage indices (sums[i] = ...).
- If you calculate the wrong size, you'll panic (index out of range) or waste memory.
  With append (Automatic):
- You don't need to know the final size.
- You don't manage indices.
- Go automatically grows the underlying array if needed.
- Code is cleaner and less error-prone.
  Performance Note:
  make is technically slightly faster if you know the exact size, because append might need to resize the array and copy data multiple times as it grows. But for most readable, idiomatic code, append is preferred unless you're optimizing a very hot loop.

# what is reflect.DeepEqual

DeepEqual reports whether x and y are “deeply equal,” defined as follows. Two values of identical type are deeply equal if one of the following cases applies. Values of distinct types are never deeply equal.
Array values are deeply equal when their corresponding elements are deeply equal.
Struct values are deeply equal if their corresponding fields, both exported and unexported, are deeply equal.
Func values are deeply equal if both are nil; otherwise they are not deeply equal.
Interface values are deeply equal if they hold deeply equal concrete values.
Map values are deeply equal when all of the following are true: they are both nil or both non-nil, they have the same length, and either they are the same map object or their corresponding keys (matched using Go equality) map to deeply equal values.
Pointer values are deeply equal if they are equal using Go’s == operator or if they point to deeply equal values.
Slice values are deeply equal when all of the following are true: they are both nil or both non-nil, they have the same length, and either they point to the same initial entry of the same underlying array (that is, &x[0] == &y[0]) or their corresponding elements (up to length) are deeply equal. Note that a non-nil empty slice and a nil slice (for example, []byte{} and []byte(nil)) are not deeply equal.
Other values - numbers, bools, strings, and channels - are deeply equal if they are equal using Go’s == operator.
In general DeepEqual is a recursive relaxation of Go’s == operator. However, this idea is impossible to implement without some inconsistency. Specifically, it is possible for a value to be unequal to itself, either because it is of func type (uncomparable in general) or because it is a floating-point NaN value (not equal to itself in floating-point comparison), or because it is an array, struct, or interface containing such a value. On the other hand, pointer values are always equal to themselves, even if they point at or contain such problematic values, because they compare equal using Go’s == operator, and that is a sufficient condition to be deeply equal, regardless of content. DeepEqual has been defined so that the same short-cut applies to slices and maps: if x and y are the same slice or the same map, they are deeply equal regardless of content.
As DeepEqual traverses the data values it may find a cycle. The second and subsequent times that DeepEqual compares two pointer values that have been compared before, it treats the values as equal rather than examining the values to which they point. This ensures that DeepEqual terminates.
