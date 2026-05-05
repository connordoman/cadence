package compiler_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/connordoman/cadence/internal/compiler"
)

var (
	EOFToken = compiler.Token{Type: compiler.EOF, Lexeme: "", Literal: nil, Line: 1}
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []compiler.Token
	}{
		{
			name:  "simple integer",
			input: "123",
			expectedTokens: []compiler.Token{
				{Type: compiler.INTEGER, Lexeme: "123", Literal: 123, Line: 1},
				EOFToken,
			},
		},
		{
			name:  "simple date",
			input: "01-01-2026",
			expectedTokens: []compiler.Token{
				{Type: compiler.DATE, Lexeme: "01-01-2026", Literal: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Line: 1},
				EOFToken,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := compiler.NewScanner(test.input)
			tokens, err := scanner.Scan()
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !reflect.DeepEqual(tokens, test.expectedTokens) {
				t.Errorf("expected %v, got %v", test.expectedTokens, tokens)
			}
		})
	}
}
