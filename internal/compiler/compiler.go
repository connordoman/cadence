package compiler

import (
	"encoding/json"
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

func Compile(source string) ([]time.Time, error) {
	comp := NewCompiler(source)
	results, err := comp.Compile(false)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func CompileAsStringSlice(source string) ([]string, error) {
	results, err := Compile(source)
	if err != nil {
		return nil, err
	}
	stringified := []string{}
	for _, result := range results {
		stringified = append(stringified, result.Format(time.DateOnly))
	}
	return stringified, nil
}

func CompileAsJSON(source string) (string, error) {
	results, err := CompileAsStringSlice(source)
	if err != nil {
		return "", err
	}
	json, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
