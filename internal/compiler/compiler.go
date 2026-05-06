package compiler

import (
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

func (c *Compiler) Compile() (*parser.Expression, error) {
	scanner := scanner.NewScanner(c.source)
	tokens, err := scanner.Scan()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return expr, nil
}
