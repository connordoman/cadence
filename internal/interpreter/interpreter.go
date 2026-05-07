package interpreter

import (
	"errors"
	"log"
	"time"

	"github.com/connordoman/cadence/internal/expr"
)

type Weekdays map[time.Weekday]bool

type Interpreter struct {
	Weekdays Weekdays

	FirstDays *Weekdays
	LastDays  *Weekdays

	MonthCount int
	WeekCount  int
	DayCount   int

	From time.Time
	To   *time.Time
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		Weekdays: Weekdays{},

		FirstDays: &Weekdays{},
		LastDays:  &Weekdays{},

		MonthCount: 0,
		WeekCount:  0,
		DayCount:   1,

		From: time.Now(),
		To:   nil,
	}
}

func (i *Interpreter) copyFirstDays() Weekdays {
	if i.FirstDays == nil {
		return nil
	}

	firstDays := Weekdays{}
	for day := range *i.FirstDays {
		firstDays[day] = true
	}

	return firstDays
}

func (i *Interpreter) copyLastDays() Weekdays {
	if i.LastDays == nil {
		return nil
	}

	lastDays := Weekdays{}
	for day := range *i.LastDays {
		lastDays[day] = true
	}

	return lastDays
}

func (i *Interpreter) evaluate(node expr.Node) (any, error) {
	return node.Accept(i)
}

func (i *Interpreter) Interpret(expression *expr.Expression) ([]time.Time, error) {
	result, err := i.evaluate(expression)
	if err != nil {
		return nil, err
	}
	return result.([]time.Time), nil
}

func (i *Interpreter) VisitExpression(expression *expr.Expression) (any, error) {
	var results []time.Time

	if expression.Interval != nil {
		_, err := i.evaluate(expression.Interval)
		if err != nil {
			return nil, err
		}
	}

	if expression.Selector != nil {
		_, err := i.evaluate(expression.Selector)
		if err != nil {
			return nil, err
		}
	}

	if expression.DateRange != nil {
		i.evaluate(expression.DateRange)
	}

	if i.To == nil {
		to := i.From.AddDate(1, 0, 0)
		i.To = &to
	}

	firstDays := i.copyFirstDays()
	lastDays := i.copyLastDays()

	lastMonth := i.From.Month()

	for date := i.From; date.Compare(*i.To) < 0; date = date.AddDate(0, 0, i.DayCount) {
		if date.Month() != lastMonth {
			lastMonth = date.Month()
			firstDays = i.copyFirstDays()
			lastDays = i.copyLastDays()
		}

		if i.MonthCount > 0 {
			if date.Month()%time.Month(i.MonthCount) != 0 {
				date = date.AddDate(0, 1, 0)
			}
		}

		_, week := date.ISOWeek()
		if i.WeekCount > 0 && week%i.WeekCount != 0 {
			continue
		}

		if len(firstDays) > 0 {
			if has, ok := firstDays[date.Weekday()]; !ok || (has && date.Day() > 7) {
				continue
			} else {
				delete(firstDays, date.Weekday())
			}
		}

		if len(lastDays) > 0 {
			if has, ok := lastDays[date.Weekday()]; !ok || (has && date.AddDate(0, 0, 7).Day() >= date.Day()) {
				continue
			} else {
				delete(lastDays, date.Weekday())
			}
		}

		if len(i.Weekdays) > 0 {
			if _, ok := i.Weekdays[date.Weekday()]; !ok {
				continue
			}
		}
		results = append(results, date)
	}

	return results, nil
}

func (i *Interpreter) VisitInterval(interval *expr.IntervalSpec) (any, error) {
	switch interval.Unit {
	case expr.UnitDay:
		i.DayCount = interval.Count
	case expr.UnitWeek:
		i.WeekCount = interval.Count
	case expr.UnitMonth:
		i.MonthCount = interval.Count
	default:
		i.DayCount = 1
	}

	return nil, nil
}

func (i *Interpreter) VisitWeekdays(weekdays *expr.WeekdaysSelector) (any, error) {
	i.Weekdays = Weekdays{
		time.Monday:    true,
		time.Tuesday:   true,
		time.Wednesday: true,
		time.Thursday:  true,
		time.Friday:    true,
	}
	return nil, nil
}

func (i *Interpreter) VisitDayList(dayList *expr.DayListSelector) (any, error) {
	for _, day := range dayList.Days {
		i.Weekdays[day] = true
	}
	return nil, nil
}

func (i *Interpreter) VisitOrdinalDayList(ordinalDayList *expr.OrdinalDayListSelector) (any, error) {
	if i.MonthCount == 0 {
		return nil, errors.New("ordinal day list must be used with a month interval")
	}

	for _, item := range ordinalDayList.Items {
		switch item.Distinction {
		case expr.First:
			if i.FirstDays == nil {
				i.FirstDays = &Weekdays{}
			}
			(*i.FirstDays)[item.Day] = true
		case expr.Last:
			if i.LastDays == nil {
				i.LastDays = &Weekdays{}
			}
			(*i.LastDays)[item.Day] = true
		}
	}

	log.Println("first days:", i.FirstDays)
	log.Println("last days:", i.LastDays)
	return nil, nil
}

func (i *Interpreter) VisitDateRange(dateRange *expr.DateRange) (any, error) {
	i.From = dateRange.From
	i.To = dateRange.To
	return nil, nil
}
