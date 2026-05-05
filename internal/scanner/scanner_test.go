package scanner_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/connordoman/cadence/internal/scanner"
)

var (
	EOFToken = scanner.Token{Type: scanner.EOF, Lexeme: "", Literal: nil, Line: 1}
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []scanner.Token
		expectedError  error
	}{
		{
			name:  "simple integer",
			input: "123",
			expectedTokens: []scanner.Token{
				{Type: scanner.INTEGER, Lexeme: "123", Literal: 123, Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
		{
			name:  "simple date",
			input: "01-01-2026",
			expectedTokens: []scanner.Token{
				{Type: scanner.DATE, Lexeme: "01-01-2026", Literal: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
		{
			name:  "digit then date",
			input: "123 01-01-2026",
			expectedTokens: []scanner.Token{
				{Type: scanner.INTEGER, Lexeme: "123", Literal: 123, Line: 1},
				{Type: scanner.DATE, Lexeme: "01-01-2026", Literal: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
		{
			name:           "impossible date",
			input:          "32-01-2026",
			expectedTokens: nil,
			expectedError:  scanner.ErrInvalidDate,
		},
		{
			name:           "unknown keyword",
			input:          "EVERYDAY",
			expectedTokens: nil,
			expectedError:  scanner.ErrUnknownKeyword,
		},
		{
			name:  "keyword every",
			input: "EVERY",
			expectedTokens: []scanner.Token{
				{Type: scanner.EVERY, Lexeme: "EVERY", Literal: nil, Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
		{
			name:  "keyword week",
			input: "WEEK",
			expectedTokens: []scanner.Token{
				{Type: scanner.WEEK, Lexeme: "WEEK", Literal: nil, Line: 1},
				EOFToken,
			},
		},
		{
			name:  "keyword plurals",
			input: "WEEKS MONTHS WEEKDAYS",
			expectedTokens: []scanner.Token{
				{Type: scanner.WEEK, Lexeme: "WEEKS", Literal: nil, Line: 1},
				{Type: scanner.MONTH, Lexeme: "MONTHS", Literal: nil, Line: 1},
				{Type: scanner.WEEKDAY, Lexeme: "WEEKDAYS", Literal: nil, Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
		{
			name:  "keyword case sensitivity",
			input: "every week",
			expectedTokens: []scanner.Token{
				{Type: scanner.EVERY, Lexeme: "every", Literal: nil, Line: 1},
				{Type: scanner.WEEK, Lexeme: "week", Literal: nil, Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
		{
			name:  "comma-separated days",
			input: "EVERY MON, TUE, WED",
			expectedTokens: []scanner.Token{
				{Type: scanner.EVERY, Lexeme: "EVERY", Literal: nil, Line: 1},
				{Type: scanner.MON, Lexeme: "MON", Literal: nil, Line: 1},
				{Type: scanner.COMMA, Lexeme: ",", Literal: nil, Line: 1},
				{Type: scanner.TUE, Lexeme: "TUE", Literal: nil, Line: 1},
				{Type: scanner.COMMA, Lexeme: ",", Literal: nil, Line: 1},
				{Type: scanner.WED, Lexeme: "WED", Literal: nil, Line: 1},
				EOFToken,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanner := scanner.NewScanner(test.input)
			tokens, err := scanner.Scan()
			if err != nil {
				if test.expectedError != nil {
					if !errors.Is(err, test.expectedError) {
						t.Errorf("expected error %v, got %v", test.expectedError, err)
					}
				} else {
					t.Errorf("expected no error, got %v", err)
				}
			}
			if test.expectedTokens == nil {
				if tokens != nil {
					t.Errorf("expected no tokens, got %v", tokens)
				}
			} else {
				if tokens == nil {
					t.Errorf("expected tokens, got nil")
				} else {
					if !reflect.DeepEqual(tokens, test.expectedTokens) {
						t.Errorf("expected %v, got %v", test.expectedTokens, tokens)
					}
				}
			}
		})
	}
}
