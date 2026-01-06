# Pointers & Errors

## Why Pointers?

Go copies values when you pass them to functions/methods. So when we change the value of the balance inside the code, we are working on a copy of what came from the test. The balance in the test is unchanged.

We can fix this with pointers. Pointers let us point to some values and then let us change them. Rather than taking a copy of the whole Wallet, we instead take a pointer to that wallet so that we can change the original values within it.

Crazy - we dont need to explicitly declare something is a pointer and it automatically takes that notation on its own to make it work together.

---

## Dereferencing

Now you might wonder, why did they pass? We didn't dereference the pointer in the function, like so:

```go
func (w *Wallet) Balance() int {
    return (*w).balance
}
```

We seemingly addressed the object directly. In fact, the code above using `(*w)` is absolutely valid. However, the makers of Go deemed this notation cumbersome, so the language permits us to write `w.balance`, without an explicit dereference. These pointers to structs even have their own name: **struct pointers** and they are automatically dereferenced.

**Dereferencing** means "I am now talking about the value held at the place this pointer refers to, not the address that it holds". You are reading, changing, or writing to the resource held at the memory address, but not making any change to the memory address itself that the pointer is holding.

Think of it like a URL analogy - when you fetch an image or webpage from a webserver using a browser, you are dereferencing. You tell your browser: fetch me the data held at this URL/address.

---

## Creating New Types

Go lets you create new types from existing ones.

The syntax is: `type MyName OriginalType`

```go
type Bitcoin int
```

This creates a new type named `Bitcoin` that underlyingly stores an `int`. It's distinct from `int` - you can't mix `Bitcoin` and `int` directly without casting (e.g., `Bitcoin(10)`).

**Why do this?**
- You can define methods on `Bitcoin`, which you can't do on a standard `int`
- Adds domain-specific meaning to values
- Can let you implement interfaces

---

## The Stringer Interface

The `fmt` package checks if a type has a `String() string` method. If it does, it calls that method to get the text representation when using `%s` format string.

```go
func (b Bitcoin) String() string {
    return fmt.Sprintf("%d BTC", b)
}
```

Now when you use `t.Errorf("got %s", got)`, Go sees `got` is of type `Bitcoin`, finds the `String()` method and calls it. Instead of printing just `10`, it prints `10 BTC`.

---

## Error Handling in Go

### The `error` Type

Go has an `error` data type with its own shenanigans. Errors are the way to signify failure when calling a function/method.

```go
func errors.New(text string) error
```

`New` returns an error that formats as the given text. Each call to `New` returns a distinct error value even if the text is identical.

### Sentinel Errors

In Go, errors are values, so we can refactor them into a variable and have a single source of truth.

**The "Brittle" Way (String Checking):**
```go
if err.Error() != "cannot withdraw, insufficient funds" { ... }
```
If you change the text in the code but forget the test, it fails.

**The "Right" Way (Sentinel Errors):**

1. **Define it** at the top of your file (package level):
    ```go
    var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")
    ```

2. **Use it** in your implementation:
    ```go
    if amount > w.balance {
        return ErrInsufficientFunds
    }
    ```

3. **Check it** in your test:
    ```go
    if got != ErrInsufficientFunds {
        t.Errorf("got %q, want %q", got, ErrInsufficientFunds)
    }
    ```

This makes the error a single source of truth.

### The errors Package

Package `errors` implements functions to manipulate errors.

- The `New` function creates errors whose only content is a text message
- An error `e` wraps another error if `e`'s type has one of the methods:
    - `Unwrap() error`
    - `Unwrap() []error`
- Easy way to create wrapped errors: `fmt.Errorf("... %w ...", ..., err, ...)`

**`errors.Is` vs simple equality:**
```go
// Preferable:
if errors.Is(err, fs.ErrExist)

// Less flexible:
if err == fs.ErrExist
```
The former will succeed if `err` wraps `fs.ErrExist`.

**`errors.As` for type assertions:**
```go
var perr *fs.PathError
if errors.As(err, &perr) {
    fmt.Println(perr.Path)
}
```
This will succeed if `err` wraps an `*fs.PathError`.

---

## Useful Tool: errcheck

A really cool package for catching unchecked errors:

```bash
go install github.com/kisielk/errcheck@latest
```

Then, inside the directory with your code run:
```bash
errcheck .
```

You might get something like:
```
wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))
```

This tells us we have not checked the error being returned on that line.

---

## nil

- Pointers can be `nil`
- When a function returns a pointer to something, you need to make sure you check if it's `nil` or you might raise a runtime exception - the compiler won't help you here
- Useful for when you want to describe a value that could be missing

---

## Summary

### Pointers
- Go copies values when you pass them to functions/methods
- If you're writing a function that needs to mutate state, you'll need it to take a pointer to the thing you want to change
- Sometimes you won't want your system to make a copy of something (very large data structures, database connection pools, etc.) - in which case you need to pass a reference

### Errors
- Errors are the way to signify failure when calling a function/method
- Checking for a string in an error results in a flaky test - use sentinel errors instead
- Don't just check errors, handle them gracefully

### Create New Types
- Useful for adding domain-specific meaning to values
- Can let you implement interfaces

Pointers and errors are a big part of writing Go that you need to get comfortable with. Thankfully the compiler will usually help you out if you do something wrong - just take your time and read the error.

---

## Q&A

**Q: What is happening with `type Bitcoin int` and the `Stringer` interface?**

This is a powerful feature in Go that allows you to create domain-specific types from basic types and customize their behavior.

1. **Creating the Type**: `type Bitcoin int` creates a new type that stores an `int` but is distinct from it. You can define methods on `Bitcoin` which you can't do on a standard `int`.

2. **The Stringer Interface**: The `fmt` package checks if a type has a `String() string` method. If it does, it uses that for text representation.

3. **How it works**: When you use `t.Errorf("got %s", got)` and `got` is a `Bitcoin`, Go finds and calls the `String()` method, printing `10 BTC` instead of just `10`.

**Q: How do we write specific error messages in tests?**

Use **Sentinel Errors** - define the error as a package-level variable and check if the returned error is that variable, rather than checking string content. This creates a single source of truth and prevents brittle tests.
