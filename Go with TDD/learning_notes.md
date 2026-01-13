# Go Design Patterns Learning Notes

Based on "Creational Design Patterns in Golang" and current codebase context.

## Core Concept: Encapsulated Creation
The main goal of Creational Patterns is to stop instantiating structs directly (e.g., `&MyStruct{}`) throughout the codebase. This couples your code to the specific implementation details of that struct. Instead, use patterns to centralize *how* objects are created.

## 1. Factory Pattern
**Concept:** A function that creates and returns an object. It abstracts the struct initialization.

**Current Codebase Example:**
In `sync/sync.go`, the `NewCounter` function is a Factory:
```go
func NewCounter() *Counter {
    return &Counter{}
}
```
If you ever need to change how `Counter` is initialized (e.g., setting a default value), you only change `NewCounter`, not every place that uses it.

## 2. Abstract Factory
**Concept:** A factory of factories. It groups related object creations together. Useful when you need to ensure that a family of objects (e.g., "Windows UI Widgets" vs "Mac UI Widgets") are used together without mixing them.

## 3. Singleton Pattern
**Concept:** Ensures a class has only one instance and provides a global point of access to it.

**Go Implementation:**
Uses `sync.Once` to ensure thread-safe, one-time initialization.
```go
var once sync.Once
var instance *Database

func GetDatabase() *Database {
    once.Do(func() {
        instance = &Database{} // Executed only once
    })
    return instance
}
```

---

## 4. Builder Pattern (Deep Dive)

**The Problem: The "Telescoping Constructor"**
Imagine you have a complex `Server` struct with many configuration options. Most are optional.

Without a builder, your constructor might look like this:
```go
// What do these 'nil's and 'true's mean? Hard to read.
srv := NewServer("localhost", 8080, 30, true, nil, false)
```
If you want to add a new option, you have to break every single function call in your code.

**The Solution: The Builder**
The Builder pattern separates the construction of a complex object from its representation. It allows you to build the object step-by-step.

**Analogy: Subway Sandwich**
You don't order a sandwich by shouting a fixed list of 20 ingredients at once. You build it step-by-step:
1. "Start with Italian Bread"
2. "Add Turkey"
3. "Add Cheese"
4. "Toast it"
5. "Finish/Wrap it"

**Go Example:**

```go
package main

import "time"

// 1. The Complex Object
type Server struct {
    Host    string
    Port    int
    Timeout time.Duration
    UseTLS  bool
}

// 2. The Builder
type ServerBuilder struct {
    server Server
}

func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{
        // Set defaults here
        server: Server{
            Host: "localhost",
            Port: 80,
            Timeout: 30 * time.Second,
        },
    }
}

// 3. Chainable Methods (The "Steps")

func (b *ServerBuilder) SetPort(port int) *ServerBuilder {
    b.server.Port = port
    return b // Return the builder itself to allow chaining
}

func (b *ServerBuilder) WithTLS() *ServerBuilder {
    b.server.UseTLS = true
    return b
}

func (b *ServerBuilder) SetHost(host string) *ServerBuilder {
    b.server.Host = host
    return b
}

// 4. The Finalizer
func (b *ServerBuilder) Build() Server {
    return b.server
}

// Usage
func main() {
    // Clean, readable, and order doesn't matter
    myServer := NewServerBuilder().
        SetHost("192.168.1.1").
        SetPort(8080).
        WithTLS().
        Build()
}
```

**Key Benefits:**
1.  **Readability:** `SetPort(8080)` is clearer than passing `8080` as the 2nd argument.
2.  **Flexibility:** You can skip parameters you don't care about (defaults are used).
3.  **Immutability:** often the `Build()` method returns the final struct, and the builder is discarded.
