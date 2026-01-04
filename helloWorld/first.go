package main

import (
	"fmt"
)

const (
	// we can init a constant to use this and better code quality
	englishHelloPrefix = `Hey`
	// doing same for langauges
	spanishHelloPrefix = "Hola"
	spanish            = "Spanish"
	french             = "French"
	frenchHelloPrefix  = "Bonjour"
)

func Mellow(name string, language string) string {
	if name == "" {
		name = "World" // creating a case where the World is the suffix when nothing given
	}

	return greetingPrefix(language) + name
}

func main() {

	fmt.Println(Mellow("Carol", "English"))
}

func greetingPrefix(language string) (prefix string) {

	switch language {
	case spanish:
		prefix = spanishHelloPrefix

	case french:
		prefix = frenchHelloPrefix

	default:
		prefix = englishHelloPrefix

	}
	return
}
