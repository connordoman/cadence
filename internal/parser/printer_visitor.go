package parser

import (
	"fmt"
	"strings"
)

type Printer struct {
	Errors []error
}

func (p *Printer) Error(err error) {
	p.Errors = append(p.Errors, err)
}

func Print(e *Expression) string {
	p := &Printer{}
	result, err := e.Accept(p)
	if err != nil {
		p.Errors = append(p.Errors, err)
	}
	return result.(string)
}

// Expression

func (p *Printer) VisitExpression(e *Expression) (any, error) {
	var parts []string

	parts = append(parts, "EVERY")

	if e.Interval != nil {
		interval, err := e.Interval.Accept(p)
		if err != nil {
			p.Error(err)
		}
		parts = append(parts, interval.(string))

		if e.Selector != nil {
			selector, err := e.Selector.Accept(p)
			if err != nil {
				p.Error(err)
			}
			parts = append(parts, selector.(string))
		}
	} else if e.Selector != nil {
		selector, err := e.Selector.Accept(p)
		if err != nil {
			p.Error(err)
		}
		parts = append(parts, selector.(string))
	}

	if e.DateRange != nil {
		dateRange, err := e.DateRange.Accept(p)
		if err != nil {
			p.Error(err)
		}
		parts = append(parts, dateRange.(string))
	}

	return strings.Join(parts, " "), nil
}

// Interval

func (p *Printer) VisitInterval(i *IntervalSpec) (any, error) {
	count := i.Count
	if count == 0 {
		count = 1
	}
	unit := i.Unit.String()

	if count == 1 {
		return unit, nil
	}

	return fmt.Sprintf("%d %s", count, unit), nil
}

// Selectors

func (p *Printer) VisitWeekdays(w *WeekdaysSelector) (any, error) {
	return "WEEKDAYS", nil
}

func (p *Printer) VisitDayList(d *DayListSelector) (any, error) {
	var parts []string
	for _, day := range d.Days {
		parts = append(parts, day.String())
	}
	return strings.Join(parts, ", "), nil
}

func (p *Printer) VisitOrdinalDayList(o *OrdinalDayListSelector) (any, error) {
	var parts []string
	for _, item := range o.Items {
		parts = append(parts, fmt.Sprintf("%s %s", item.Distinction.String(), item.Day.String()))
	}
	return strings.Join(parts, ", "), nil
}

// Date Range

func (p *Printer) VisitDateRange(d *DateRange) (any, error) {
	var result strings.Builder
	from := d.From.Format("02-01-2006")
	fmt.Fprintf(&result, "FROM %s", from)
	to := ""
	if d.To != nil {
		to = d.To.Format("02-01-2006")
		fmt.Fprintf(&result, " TO %s", to)
	}
	return result.String(), nil
}
