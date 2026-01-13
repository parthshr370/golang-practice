# Chapter 4: Arrays and Slices

## Arrays vs Slices

### Arrays
- Fixed capacity defined at declaration
- Size is part of the type: `[4]int` and `[5]int` are different types
- Less common in Go compared to slices

```go
// Two ways to initialize
numbers := [5]int{1, 2, 3, 4, 5}
numbers := [...]int{1, 2, 3, 4, 5} // Compiler counts elements
```

### Slices
- Variable capacity (can grow)
- Reference type (points to an underlying array)
- Format: `[]T` (no size specified)

```go
mySlice := []int{1, 2, 3}
```

## Internal Structure of Slices

Think of a slice as a header struct with 3 fields:
1. **Pointer**: Points to the first element of the underlying array
2. **Length (`len`)**: Number of elements currently in the slice
3. **Capacity (`cap`)**: Number of elements the underlying array can hold before resizing

```
+---------+      +---+---+---+---+---+
| Pointer | ---> | 1 | 2 | 3 |   |   |
+---------+      +---+---+---+---+---+
| Length  | = 3    0   1   2   3   4
+---------+
| Capacity| = 5
+---------+
```

## How Append Works (Memory & Capacity)

When you use `append`, Go behaves differently depending on available capacity.

### Scenario A: Full Capacity (`len == cap`)

```go
sliceA := []int{1, 2, 3} // len=3, cap=3
sliceB := append(sliceA, 4)
```

**Before Append (sliceA):**
```
[1][2][3]
 ^
 | sliceA (len=3, cap=3)
```

**After Append (sliceB):**
1. Capacity is full.
2. Go allocates NEW array (usually 2x size).
3. Copies elements [1, 2, 3].
4. Appends [4].

```
[1][2][3]         [1][2][3][4][ ][ ]
 ^                 ^
 | sliceA          | sliceB (len=4, cap=6)
```
**Result**: `sliceA` and `sliceB` are completely independent.

### Scenario B: Available Capacity (`cap > len`)

```go
sliceD := make([]int, 0, 4) // cap=4
sliceD = append(sliceD, 1, 2, 3) // len=3, cap=4
```

**State of sliceD:**
```
[1][2][3][ ]
 ^
 | sliceD (len=3, cap=4)
```

**Forking the slice:**
```go
sliceE := append(sliceD, 4) // writes to index 3
sliceF := append(sliceD, 5) // OVERWRITES index 3
```

1. **sliceE**: Uses existing capacity. Writes `4`.
   ```
   [1][2][3][4]
    ^        ^
    | sliceD | sliceE (len=4)
   ```

2. **sliceF**: Also uses existing capacity. Writes `5`, overwriting `4`.
   ```
   [1][2][3][5]
    ^        ^
    | sliceD | sliceF (len=4)
             | sliceE (also sees 5 now!)
   ```

**Result**: `sliceE` and `sliceF` share memory. `sliceE` thinks the last element is 4, but `sliceF` changed it to 5.

### The "Gotcha"
Because they share the underlying array, modifying an element in `sliceF` will be visible in `sliceE` and `sliceD` if they overlap in memory range.

**Rule of Thumb:** Never assume `append` will copy or reuse. Always assign back to the same variable `a = append(a, x)` unless you explicitly want to branch/fork the slice (in which case, be very careful about shared memory).

## Slicing Slices

You can create a new slice from an existing one:
```go
// Syntax: slice[low:high]
// low: inclusive, high: exclusive
newSlice := slice[1:4]
```

- If you omit low: `slice[:high]` (from start)
- If you omit high: `slice[low:]` (to end)
- `slice[:]` (copy of whole slice)

**Important:** The new slice points to the **same underlying array**. Changing one affects the other!

## Modern For Loops & Range

To iterate over collections:

```go
// Idiomatic
for _, number := range numbers {
    sum += number
}
```
- **`_` (Blank Identifier)**: Used to ignore the index return value
- **`range`**: Returns (index, value)

## Variadic Functions

Functions that take a variable number of arguments:

```go
func SumAll(numbersToSum ...[]int) {
    // numbersToSum is now a slice of slices: [][]int
}
```
You can pass multiple arguments individually or explode a slice with `...`:
```go
nums := []int{1, 2, 3}
sum(nums...)
```

## Creating & Growing Slices

### `make`
Creates a slice with specific length/capacity:
```go
// length 5, capacity 5
slice := make([]int, 5) 
```

### `append`
Adds elements to a slice, handling memory allocation automatically:
```go
var sums []int
sums = append(sums, newValue)
```
- Preferred over manual index management
- Handles resizing the underlying array if capacity is exceeded
- Returns a new slice reference (must assign back)

## Testing Techniques

### 1. Test Coverage
Run coverage analysis to find untested code paths:
```bash
go test -cover
```

### 2. Helper Functions inside Tests
You can define functions inside test functions:
- Access to local scope
- Reduces API surface (not visible outside test)
- Adds type safety (compared to strict dynamic checking)

```go
checkSums := func(t testing.TB, got, want []int) {
    t.Helper()
    if !slices.Equal(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

### 3. Comparing Slices
- You cannot use `==` with slices (except `nil`)
- Use `slices.Equal` (Go 1.21+) for simple equality
- Use `reflect.DeepEqual` for checking complex structures (but lose type safety)

## Key Takeaways

1. **Prefer Slices**: Use slices (`[]int`) over arrays (`[5]int`) for flexibility
2. **Refactor**: Start with working code, then improve. We refactored `Sum` from array to slice without changing logic.
3. **Compiler Errors**: Friend, not foe. They catch type mismatches early.
4. **Runtime Errors**: Avoid them. "Index out of range" is a common runtime panic with arrays/slices.
5. **Coverage**: Use `go test -cover` to ensure confidence, not just metrics.

## Common Pitfalls

- **Empty Slices**: Slicing an empty slice `[][1:]` causes a runtime panic. Always check length first!
- **Type Mismatch**: `[5]int` != `[]int`. You can't pass an array to a function expecting a slice.
- **Nil Slices**: `nil` is a valid slice with 0 length. `append` works on `nil` slices.
