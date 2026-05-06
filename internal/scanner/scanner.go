package scanner

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUnexpectedEOF  = errors.New("unexpected EOF")
	ErrUnexpectedChar = errors.New("unexpected character")
	ErrInvalidInteger = errors.New("invalid integer")
	ErrInvalidDate    = errors.New("invalid date")
	ErrUnknownKeyword = errors.New("unknown keyword")
)

var keywords = map[string]TokenType{
	"EVERY":    EVERY,
	"DAY":      DAY,
	"DAYS":     DAY,
	"WEEK":     WEEK,
	"WEEKS":    WEEK,
	"MONTH":    MONTH,
	"MONTHS":   MONTH,
	"ON":       ON,
	"MON":      MON,
	"TUE":      TUE,
	"WED":      WED,
	"THU":      THU,
	"FRI":      FRI,
	"SAT":      SAT,
	"SUN":      SUN,
	"WEEKDAY":  WEEKDAY,
	"WEEKDAYS": WEEKDAY,
	"FROM":     FROM,
	"TO":       TO,
	"FIRST":    FIRST,
	"LAST":     LAST,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) report(err error, format string, args ...any) error {
	return fmt.Errorf("line %d: %w: %s", s.line, err, fmt.Sprintf(format, args...))
}

func (s *Scanner) Scan() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
	char, err := s.advance()
	if err != nil {
		return err
	}

	switch char {
	case ',':
		s.emptyToken(COMMA)
	// case '-':
	// 	s.emptyToken(HYPHEN)
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	default:
		if s.isDigit(char) {
			if s.isDigit(s.peek()) && s.peekNext() == '-' {
				return s.date()
			} else {
				return s.integer()
			}
		} else if s.isAlpha(char) {
			return s.keyword()
		} else {
			return ErrUnexpectedChar
		}
	}

	return nil
}

func (s *Scanner) keyword() error {
	for s.isAlpha(s.peek()) {
		s.advance()
	}

	text := strings.ToUpper(s.source[s.start:s.current])
	if keyword, ok := keywords[text]; ok {
		s.emptyToken(keyword)
		return nil
	}

	return s.report(ErrUnknownKeyword, "'%s'", text)
}

func (s *Scanner) integer() error {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	literal, err := strconv.Atoi(s.source[s.start:s.current])
	if err != nil {
		return s.report(ErrInvalidInteger, "%s: %v", s.source[s.start:s.current], err)
	}
	s.addToken(INTEGER, literal)
	return nil
}

func (s *Scanner) date() error {
	// consume day
	s.advance() // second digit of day
	s.advance() // hyphen

	// consume month
	monthFirstDigit := s.peek()
	monthSecondDigit := s.peekNext()
	if !s.isDigit(monthFirstDigit) || !s.isDigit(monthSecondDigit) {
		return s.report(ErrInvalidDate, "invalid month: '%s%s'", string(monthFirstDigit), string(monthSecondDigit))
	}
	s.advance() // first digit of month
	s.advance() // second digit of month

	if ok, got := s.match('-'); !ok {
		return s.report(ErrInvalidDate, "expected '-', got '%s'", string(got))
	}

	// consume year
	for range 4 {
		if !s.isDigit(s.peek()) {
			return s.report(ErrInvalidDate, "invalid year: expected digit")
		}
		s.advance()
	}

	date, err := time.Parse("02-01-2006", s.source[s.start:s.current])
	if err != nil {
		return s.report(ErrInvalidDate, "invalid date: %v", err)
	}

	s.addToken(DATE, date)
	return nil
}

func (s *Scanner) match(expected rune) (bool, rune) {
	if s.isAtEnd() {
		return false, '\000'
	}
	if rune(s.source[s.current]) != expected {
		return false, rune(s.source[s.current])
	}
	s.current++
	return true, rune(s.source[s.current])
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	return s.peekAhead(1)
}

func (s *Scanner) peekAhead(offset int) rune {
	if s.current+offset >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+offset])
}

func (s *Scanner) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func (s *Scanner) isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) advance() (rune, error) {
	s.current++
	if s.isAtEnd() {
		return 0, ErrUnexpectedEOF
	}

	return rune(s.source[s.current-1]), nil
}

func (s *Scanner) emptyToken(tokenType TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}
