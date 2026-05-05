package compiler_test

import (
	"reflect"
	"testing"

	"github.com/connordoman/cadence/internal/compiler"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []compiler.Token
	}{
		{
			name:  "simple every monday",
			input: "EVERY MON",
			expectedTokens: []compiler.Token{
				{Type: compiler.EVERY, Lexeme: "EVERY", Literal: nil, Line: 1},
				{Type: compiler.MON, Lexeme: "MON", Literal: nil, Line: 1},
			},
		},
		{
			name:  "simple every week",
			input: "EVERY MON, TUE",
			expectedTokens: []compiler.Token{
				{Type: compiler.EVERY, Lexeme: "EVERY", Literal: nil, Line: 1},
				{Type: compiler.MON, Lexeme: "MON", Literal: nil, Line: 1},
				{Type: compiler.COMMA, Lexeme: ",", Literal: nil, Line: 1},
				{Type: compiler.TUE, Lexeme: "TUE", Literal: nil, Line: 1},
			},
		},
		{
			name:  "simple every weekday",
			input: "EVERY WEEKDAY",
			expectedTokens: []compiler.Token{
				{Type: compiler.EVERY, Lexeme: "EVERY", Literal: nil, Line: 1},
				{Type: compiler.WEEKDAY, Lexeme: "WEEKDAY", Literal: nil, Line: 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			compiler := compiler.NewCompiler(test.input)
			tokens, err := compiler.Compile()
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}
			if !reflect.DeepEqual(tokens, test.expectedTokens) {
				t.Errorf("expected %v, got %v", test.expectedTokens, tokens)
			}
		})
	}
}
