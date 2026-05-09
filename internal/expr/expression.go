package expr

import "time"

type Expression struct {
	Interval  *IntervalSpec
	Selector  Selector
	DateRange *DateRange
}

func (e *Expression) Accept(v Visitor) (any, error) {
	return v.VisitExpression(e)
}

type Unit int

const (
	UnitDay Unit = iota
	UnitWeek
	UnitMonth
)

func (u Unit) String() string {
	switch u {
	case UnitDay:
		return "DAY"
	case UnitWeek:
		return "WEEK"
	case UnitMonth:
		return "MONTH"
	}
	return "UNKNOWN"
}

/* Interval Spec (INTEGER unit) */

type IntervalSpec struct {
	Count int
	Unit  Unit
}

func (i *IntervalSpec) Accept(v Visitor) (any, error) {
	return v.VisitInterval(i)
}

/* Selector (WEEKDAYS, day_list, ordinal_day_list) */

type Selector interface {
	Node
	isSelector()
}

/* Weekdays (WEEKDAYS) */

type WeekdaysSelector struct{}

func (WeekdaysSelector) isSelector() {}

func (w *WeekdaysSelector) Accept(v Visitor) (any, error) {
	return v.VisitWeekdays(w)
}

/* Days (MON, TUE, ..., SUN) */

// type Day int

// const (
// 	Monday Day = iota
// 	Tuesday
// 	Wednesday
// 	Thursday
// 	Friday
// 	Saturday
// 	Sunday
// )

// func (d Day) String() string {
// 	switch d {
// 	case Monday:
// 		return "MON"
// 	case Tuesday:
// 		return "TUE"
// 	case Wednesday:
// 		return "WED"
// 	case Thursday:
// 		return "THU"
// 	case Friday:
// 		return "FRI"
// 	case Saturday:
// 		return "SAT"
// 	case Sunday:
// 		return "SUN"
// 	}
// 	return "UNKNOWN"
// }

/* Day List ([MON, TUE, ..., SUN]) */

type DayListSelector struct {
	Days []time.Weekday
}

func (DayListSelector) isSelector() {}

func (d *DayListSelector) Accept(v Visitor) (any, error) {
	return v.VisitDayList(d)
}

/* Ordinal Days (FIRST, LAST) */

type Distinction int

const (
	First Distinction = iota
	Last
)

func (d Distinction) String() string {
	switch d {
	case First:
		return "FIRST"
	case Last:
		return "LAST"
	}
	return "UNKNOWN"
}

type OrdinalDay struct {
	Distinction Distinction
	Day         time.Weekday
}

type OrdinalDayListSelector struct {
	Items []OrdinalDay
}

func (OrdinalDayListSelector) isSelector() {}

func (o *OrdinalDayListSelector) Accept(v Visitor) (any, error) {
	return v.VisitOrdinalDayList(o)
}

/* Date Range (FROM date, TO date) */

type DateRangeInclusivity int

const (
	DateRangeInclusive DateRangeInclusivity = iota
	DateRangeExclusive
)

func (d DateRangeInclusivity) String() string {
	switch d {
	case DateRangeInclusive:
		return "INCLUSIVE"
	case DateRangeExclusive:
		return "EXCLUSIVE"
	}
	return "UNKNOWN"
}

type DateRangeEnd struct {
	Date        time.Time
	Inclusivity DateRangeInclusivity
}

type DateRange struct {
	From DateRangeEnd
	To   *DateRangeEnd
}

func (d *DateRange) Accept(v Visitor) (any, error) {
	return v.VisitDateRange(d)
}
