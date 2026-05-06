package parser

type Node interface {
	Accept(v Visitor) (any, error)
}

type Visitor interface {
	VisitExpression(*Expression) (any, error)
	VisitInterval(*IntervalSpec) (any, error)
	VisitWeekdays(*WeekdaysSelector) (any, error)
	VisitDayList(*DayListSelector) (any, error)
	VisitOrdinalDayList(*OrdinalDayListSelector) (any, error)
	VisitDateRange(*DateRange) (any, error)
}
