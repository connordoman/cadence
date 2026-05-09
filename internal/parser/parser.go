package parser

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/connordoman/cadence/internal/expr"
	"github.com/connordoman/cadence/internal/scanner"
)

var (
	ErrParse = errors.New("parse error")
)

type Parser struct {
	Tokens  []scanner.Token
	Current int
	Errors  []error
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) Error(err error) error {
	p.Errors = append(p.Errors, err)
	return err
}

func (p *Parser) ErrExpected(expected string, got string) error {
	return p.Error(fmt.Errorf("expected %s, got '%s'", expected, got))
}

func (p *Parser) ErrNested(msg string, err error) error {
	p.Error(fmt.Errorf("%s", msg))
	return err
}

func (p *Parser) Parse() (*expr.Expression, error) {
	expr, err := p.expression()
	if err != nil {
		fmt.Println("Parser error:")
		for i, err := range p.Errors {
			fmt.Println(i, ":", err)
		}
		return nil, ErrParse
	}
	return expr, nil
}

func (p *Parser) expression() (*expr.Expression, error) {
	if !p.matchSome(scanner.EVERY) {
		return nil, p.ErrExpected("'EVERY'", p.peek().Type.String())
	}

	expression := &expr.Expression{}

	if p.checkSome(scanner.INTEGER, scanner.DAY, scanner.WEEK, scanner.MONTH) {
		explicitInterval, selector, err := p.explicitSchedule()
		if err != nil {
			return nil, p.ErrNested("explicit schedule", err)
		}

		expression.Interval = explicitInterval
		expression.Selector = selector
	} else {
		implicitSchedule, err := p.implicitSchedule()
		if err != nil {
			return nil, p.ErrNested("implicit schedule", err)
		}
		expression.Interval = &expr.IntervalSpec{
			Count: 1,
			Unit:  expr.UnitWeek,
		}
		expression.Selector = implicitSchedule
	}

	if p.check(scanner.FROM) {
		dateRange, err := p.dateRange()
		if err != nil {
			return nil, p.ErrNested("date range", err)
		}
		expression.DateRange = dateRange
	}

	return expression, nil
}

func (p *Parser) day() (time.Weekday, error) {
	if p.matchSome(scanner.MON, scanner.TUE, scanner.WED, scanner.THU, scanner.FRI, scanner.SAT, scanner.SUN) {
		prev := p.previous()
		switch prev.Type {
		case scanner.MON:
			return time.Monday, nil
		case scanner.TUE:
			return time.Tuesday, nil
		case scanner.WED:
			return time.Wednesday, nil
		case scanner.THU:
			return time.Thursday, nil
		case scanner.FRI:
			return time.Friday, nil
		case scanner.SAT:
			return time.Saturday, nil
		case scanner.SUN:
			return time.Sunday, nil
		default:
			return -1, fmt.Errorf("invalid day: %s", prev.Type)
		}
	}

	return -1, p.ErrExpected("day", p.peek().Type.String())
}

func (p *Parser) distinction() (expr.Distinction, error) {
	if p.matchSome(scanner.FIRST, scanner.LAST) {
		prev := p.previous()
		switch prev.Type {
		case scanner.FIRST:
			return expr.First, nil
		case scanner.LAST:
			return expr.Last, nil
		default:
			return -1, fmt.Errorf("invalid distinction: %s", prev.Type)
		}
	}

	return -1, p.ErrExpected("distinction", p.peek().Type.String())
}

func (p *Parser) ordinalDay() (*expr.OrdinalDay, error) {
	distinction, err := p.distinction()
	if err != nil {
		return nil, err
	}

	day, err := p.day()
	if err != nil {
		return nil, err
	}

	return &expr.OrdinalDay{
		Distinction: distinction,
		Day:         day,
	}, nil
}

func (p *Parser) dayList() (*expr.DayListSelector, error) {
	days := []time.Weekday{}

	first, err := p.day()
	if err != nil {
		return nil, err
	}

	days = append(days, first)

	for p.matchSome(scanner.COMMA) {
		day, err := p.day()
		if err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return &expr.DayListSelector{
		Days: days,
	}, nil
}

func (p *Parser) ordinalDayList() (*expr.OrdinalDayListSelector, error) {
	items := []expr.OrdinalDay{}

	first, err := p.ordinalDay()
	if err != nil {
		return nil, err
	}

	items = append(items, *first)

	for p.matchSome(scanner.COMMA) {
		item, err := p.ordinalDay()
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	return &expr.OrdinalDayListSelector{
		Items: items,
	}, nil
}

func (p *Parser) selector() (expr.Selector, error) {
	if p.matchSome(scanner.WEEKDAY) {
		return &expr.WeekdaysSelector{}, nil
	}

	if p.check(scanner.FIRST) || p.check(scanner.LAST) {
		item, err := p.ordinalDayList()
		if err != nil {
			return nil, p.ErrNested("ordinal day list", err)
		}

		return item, nil
	}

	return p.dayList()
}

func (p *Parser) unit() (expr.Unit, error) {
	if p.matchSome(scanner.DAY, scanner.WEEK, scanner.MONTH) {
		prev := p.previous()
		switch prev.Type {
		case scanner.DAY:
			return expr.UnitDay, nil
		case scanner.WEEK:
			return expr.UnitWeek, nil
		case scanner.MONTH:
			return expr.UnitMonth, nil
		default:
			return -1, fmt.Errorf("invalid unit: %s", prev.Type)
		}
	}

	return -1, p.ErrExpected("unit", p.peek().Type.String())
}

func (p *Parser) intervalSpec() (*expr.IntervalSpec, error) {
	count := 1
	if p.matchSome(scanner.INTEGER) {
		count = p.previous().Literal.(int)
	}

	unit, err := p.unit()
	if err != nil {
		return nil, err
	}

	return &expr.IntervalSpec{
		Count: count,
		Unit:  unit,
	}, nil
}

func (p *Parser) explicitSchedule() (*expr.IntervalSpec, expr.Selector, error) {
	interval, err := p.intervalSpec()
	if err != nil {
		return nil, nil, err
	}

	var selector expr.Selector

	if p.matchSome(scanner.ON) {
		resolvedSelector, err := p.selector()
		if err != nil {
			return nil, nil, err
		}

		selector = resolvedSelector
	}

	return interval, selector, nil
}

func (p *Parser) implicitSchedule() (expr.Selector, error) {
	selector, err := p.selector()
	if err != nil {
		return nil, p.ErrNested("selector", err)
	}

	return selector, nil
}

func (p *Parser) dateRange() (*expr.DateRange, error) {
	if !p.matchSome(scanner.FROM) {
		return nil, p.ErrExpected("'FROM'", p.peek().Type.String())
	}

	fromDate, err := p.date()
	if err != nil {
		return nil, p.ErrNested("from date", err)
	}

	var toDate *time.Time

	if p.matchSome(scanner.TO) {
		date, err := p.date()
		if err != nil {
			return nil, p.ErrNested("to date", err)
		}
		toDate = &date
	}

	return &expr.DateRange{
		From: fromDate,
		To:   toDate,
	}, nil
}

func (p *Parser) date() (time.Time, error) {
	if !p.matchSome(scanner.DATE) {
		return time.Time{}, p.ErrExpected("date", p.peek().Type.String())
	}

	return p.previous().Literal.(time.Time), nil
}

func (p *Parser) matchSome(tokenTypes ...scanner.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) checkSome(tokenTypes ...scanner.TokenType) bool {
	return slices.ContainsFunc(tokenTypes, p.check)
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() scanner.Token {
	return p.Tokens[p.Current-1]
}
