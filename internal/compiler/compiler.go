package compiler

type Compiler struct {
	source string
	tokens []Token
}

func NewCompiler(source string) *Compiler {
	return &Compiler{
		source: source,
		tokens: []Token{},
	}
}

func (c *Compiler) Compile() ([]Token, error) {
	scanner := NewScanner(c.source)
	tokens, err := scanner.Scan()
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
