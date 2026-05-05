package compiler

import "fmt"

type TokenType int

const (
	EOF TokenType = iota

	// Literals
	INTEGER
	DATE

	// Single-character tokens
	COMMA

	// Keywords
	EVERY
	WEEK
	MONTH
	ON
	MON
	TUE
	WED
	THU
	FRI
	SAT
	SUN
	WEEKDAY
	FROM
	TO
	FIRST
	LAST
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case INTEGER:
		return "INTEGER"
	case DATE:
		return "DATE"
	case COMMA:
		return "COMMA"
	case EVERY:
		return "EVERY"
	case WEEK:
		return "WEEK"
	case MONTH:
		return "MONTH"
	case ON:
		return "ON"
	case MON:
		return "MON"
	case TUE:
		return "TUE"
	case WED:
		return "WED"
	case THU:
		return "THU"
	case FRI:
		return "FRI"
	case SAT:
		return "SAT"
	case SUN:
		return "SUN"
	case WEEKDAY:
		return "WEEKDAY"
	case FROM:
		return "FROM"
	case TO:
		return "TO"
	case FIRST:
		return "FIRST"
	case LAST:
		return "LAST"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s '%s' %v", t.Type, t.Lexeme, t.Literal)
}
