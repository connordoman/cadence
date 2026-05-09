package interpreter

import (
	"errors"
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

		FirstDays: nil,
		LastDays:  nil,

		MonthCount: 0,
		WeekCount:  0,
		DayCount:   1,

		From: time.Now(),
		To:   nil,
	}
}

func (i *Interpreter) evaluate(node expr.Node) (any, error) {
	return node.Accept(i)
}

func (i *Interpreter) evaluateOrdinalMonths() ([]time.Time, error) {
	if i.MonthCount == 0 {
		return nil, errors.New("ordinal day lists are only supported with month intervals")
	}

	results := []time.Time{}

	for month := i.From; month.Compare(*i.To) < 0; month = month.AddDate(0, 1, 0) {
		if i.FirstDays != nil {
			for firstDay, has := range *i.FirstDays {
				if !has {
					continue
				}

				firstDayOfMonth := firstWeekdayOfMonth(month.Year(), month.Month(), firstDay, month.Location())

				if firstDayOfMonth.Compare(i.From) < 0 {
					continue
				}
				results = append(results, firstDayOfMonth)
			}
		}

		if i.LastDays != nil {
			for lastDay, has := range *i.LastDays {
				if !has {
					continue
				}

				lastDayOfMonth := lastWeekdayOfMonth(month.Year(), month.Month(), lastDay, month.Location())

				if lastDayOfMonth.Compare(*i.To) > 0 {
					continue
				}

				results = append(results, lastDayOfMonth)
			}
		}
	}

	return results, nil
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

	if i.FirstDays != nil || i.LastDays != nil {
		return i.evaluateOrdinalMonths()
	}

	_, weekOffset := i.From.ISOWeek()

	for date := i.From; date.Compare(*i.To) < 0; date = date.AddDate(0, 0, i.DayCount) {

		if i.MonthCount > 0 {
			if date.Month()%time.Month(i.MonthCount) != 0 {
				date = date.AddDate(0, 1, 0)
			}
		}

		_, week := date.ISOWeek()
		if i.WeekCount > 0 && (week-weekOffset)%i.WeekCount != 0 {
			continue
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

	return nil, nil
}

func (i *Interpreter) VisitDateRange(dateRange *expr.DateRange) (any, error) {
	i.From = dateRange.From
	i.To = dateRange.To
	return nil, nil
}
