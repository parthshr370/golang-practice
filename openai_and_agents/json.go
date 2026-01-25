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
	m := Message{"Alice", "Hello", 1294706395881547000}

	b, err := json.Marshal(m)

	// Using %s to print bytes as string, %v for raw bytes
	fmt.Printf("Raw bytes: %v\nError: %v\n", b, err)
	fmt.Println("JSON String:", string(b))
}
