package compiler

import "github.com/connordoman/cadence/internal/scanner"

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

func (c *Compiler) Compile() ([]scanner.Token, error) {
	scanner := scanner.NewScanner(c.source)
	tokens, err := scanner.Scan()
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
