package main

import (
	"fmt"

	"github.com/connordoman/cadence/internal/parser"
)

func main() {
	expr := &parser.Expression{

		// Interval: &parser.IntervalSpec{
		// 	Count: 1,
		// 	Unit:  parser.UnitMonth,
		// },
		Interval: nil,
		Selector: &parser.OrdinalDayListSelector{
			Items: []parser.OrdinalDay{
				{Distinction: parser.First, Day: parser.Monday},
			},
		},
	}

	fmt.Println(parser.Print(expr))
}
