package main

import (
	"fmt"
	"os"

	"github.com/connordoman/cadence/internal/compiler"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cadence <input>")
		os.Exit(1)
	}

	input := os.Args[1]

	comp := compiler.NewCompiler(input)
	tokens, err := comp.Compile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
}
