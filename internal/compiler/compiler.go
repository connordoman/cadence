package compiler

import (
	"fmt"
	"time"

	"github.com/connordoman/cadence/internal/expr"
	"github.com/connordoman/cadence/internal/interpreter"
	"github.com/connordoman/cadence/internal/parser"
	"github.com/connordoman/cadence/internal/scanner"
)

type Compiler struct {
	source string
	tokens []scanner.Token
}

func NewCompiler(source string) *Compiler {
	return &Compiler{
		source: source,
		tokens: []scanner.Token{},
	}
}

func (c *Compiler) Compile(verbose bool) ([]time.Time, error) {
	s := scanner.NewScanner(c.source)
	tokens, err := s.Scan()
	if err != nil {
		return nil, err
	}

	p := parser.NewParser(tokens)
	expression, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Println(expr.Print(expression))
	}

	i := interpreter.NewInterpreter()
	results, err := i.Interpret(expression)
	if err != nil {
		return nil, err
	}

	return results, nil
}
