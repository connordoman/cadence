package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/connordoman/cadence/internal/compiler"
	"github.com/connordoman/cadence/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		repl()
		return
	}

	input := os.Args[1]

	comp := compiler.NewCompiler(input)
	expr, err := comp.Compile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(parser.Print(expr))
}

func repl() {
	fmt.Println("Cadence REPL (type 'exit' to quit)")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		}
		comp := compiler.NewCompiler(line)
		expr, err := comp.Compile()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(parser.Print(expr))
	}
	fmt.Println("Bye!")
}
