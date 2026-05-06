package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/connordoman/cadence/internal/compiler"
	"github.com/connordoman/cadence/internal/expr"
)

func main() {
	if len(os.Args) < 2 {
		repl()
		return
	}

	input := os.Args[1]

	comp := compiler.NewCompiler(input)
	expression, err := comp.Compile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(expr.Print(expression))
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
		expression, err := comp.Compile()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(expr.Print(expression))
	}
	fmt.Println("Bye!")
}
