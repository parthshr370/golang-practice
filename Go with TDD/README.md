# golang-practice

A personal repository for Go programming practice, learning, and experimentation. This is my workspace for improving Go skills through various exercises, tutorials, and projects.

## Contents

### Go with TDD

This folder contains a complete journey through Go using Test-Driven Development, based on the "Learn Go with Tests" methodology. Each module focuses on specific Go concepts with comprehensive notes and examples.

#### Modules Covered:

| Module          | Description               | Key Concepts                                                 |
| --------------- | ------------------------- | ------------------------------------------------------------ |
| **helloWorld**  | Introduction to Go basics | Packages, functions, testing, TDD workflow                   |
| **integers**    | Numeric operations        | TDD red-green-refactor, documentation, examples              |
| **iteration**   | Loops & performance       | For loops, benchmarks, string.Builder optimization           |
| **slice_arr**   | Data structures           | Arrays vs slices, memory layout, append behavior             |
| **structs**     | Custom types & interfaces | Structs, methods, interfaces, table-driven tests             |
| **pointers**    | Memory management         | Pointers, dereferencing, custom types, Stringer interface    |
| **maps**        | Key-value storage         | Maps, reference semantics, CRUD operations                   |
| **dependency**  | Dependency injection      | io.Writer interface, decoupling, testing patterns            |
| **mocking**     | Test doubles              | Spies, configurable sleepers, mock design                    |
| **concurrency** | Parallel programming      | Goroutines, channels, race conditions, M:N scheduling        |
| **select**      | Channel synchronization   | Select statements, timeouts, HTTP testing                    |
| **context**     | Cancellation              | Context propagation, cancellation signals, request lifecycle |
| **sync**        | Shared state              | Mutex, WaitGroup, happens-before relation                    |
| **reflection**  | Runtime introspection     | reflect package, interface{}, type safety                    |

#### Key Features:

- **Comprehensive Notes**: Each module includes detailed learning notes with Q&A sections
- **TDD Approach**: All code developed following test-driven development
- **Design Patterns**: Notes on creational patterns and Go idioms
- **Learning Log**: Progress tracking through devlogs and session notes

### Getting Started

#### Prerequisites

- Go 1.22 or later
- Basic programming knowledge

#### Running Tests

```bash
# Run all tests in Go with TDD
cd "Go with TDD"
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

#### Running Specific Modules

```bash
cd "Go with TDD/[module-name]"
go test -v
```

## Learning Resources

### Recommended Reading

- [Effective Go](https://go.dev/doc/effective_go)
- [The Go Programming Language](https://gopl.io/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)

### Useful Tools

- `go doc` - View package documentation
- `go vet` - Run static analysis
- `go fmt` - Format code
- `errcheck` - Check for unchecked errors
- `pkgsite` - Local documentation viewer

## Project Structure

```
golang-practice/
├── README.md
└── Go with TDD/
    ├── go.mod
    ├── AGENTS.md
    ├── GO_LEARNING_LOG.md
    ├── learning_notes.md
    └── [modules]/
        ├── *.go           # Implementation files
        ├── *_test.go      # Test files
        └── notes.md       # Learning notes
```

## Contribution Philosophy

This is a personal learning repository. The primary goals are:

- Practice Go programming concepts
- Experiment with different approaches
- Build a reference of idiomatic Go patterns
- Document learning progress

## License

This repository is for personal learning purposes. Code examples may reference or be inspired by various learning resources.

---

_Last updated: January 2026_
