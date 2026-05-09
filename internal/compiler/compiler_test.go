package compiler_test

import (
	"testing"
	"time"

	"github.com/connordoman/cadence/internal/compiler"
)

func buildTimeArray(dates ...string) []time.Time {
	times := make([]time.Time, len(dates))
	for i, date := range dates {
		times[i], _ = time.Parse("2006-01-02", date)
	}
	return times
}

func TestCompile(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedResults []time.Time
	}{
		{
			name:            "simple every monday",
			input:           "EVERY MON FROM 2026-01-01 TO 2026-02-01",
			expectedResults: buildTimeArray("2026-01-05", "2026-01-12", "2026-01-19", "2026-01-26"),
		},
		{
			name:            "simple every week",
			input:           "EVERY MON, TUE FROM 2026-01-01 TO 2026-02-01",
			expectedResults: buildTimeArray("2026-01-05", "2026-01-06", "2026-01-12", "2026-01-13", "2026-01-19", "2026-01-20", "2026-01-26", "2026-01-27"),
		},
		{
			name:  "simple every weekday",
			input: "EVERY WEEKDAY FROM 2026-01-01 TO 2026-02-01",
			expectedResults: buildTimeArray(
				"2026-01-01", "2026-01-02",
				"2026-01-05", "2026-01-06", "2026-01-07", "2026-01-08", "2026-01-09",
				"2026-01-12", "2026-01-13", "2026-01-14", "2026-01-15", "2026-01-16",
				"2026-01-19", "2026-01-20", "2026-01-21", "2026-01-22", "2026-01-23",
				"2026-01-26", "2026-01-27", "2026-01-28", "2026-01-29", "2026-01-30",
			),
		},
		{
			name:  "simple every day",
			input: "EVERY DAY FROM 2026-01-01 TO 2026-02-01",
			expectedResults: buildTimeArray(
				"2026-01-01", "2026-01-02", "2026-01-03", "2026-01-04",
				"2026-01-05", "2026-01-06", "2026-01-07", "2026-01-08", "2026-01-09", "2026-01-10", "2026-01-11",
				"2026-01-12", "2026-01-13", "2026-01-14", "2026-01-15", "2026-01-16", "2026-01-17", "2026-01-18",
				"2026-01-19", "2026-01-20", "2026-01-21", "2026-01-22", "2026-01-23", "2026-01-24", "2026-01-25",
				"2026-01-26", "2026-01-27", "2026-01-28", "2026-01-29", "2026-01-30", "2026-01-31",
			),
		},
		{
			name:  "simple every 2 days",
			input: "EVERY 2 DAYS FROM 2026-01-01 TO 2026-02-01",
			expectedResults: buildTimeArray(
				"2026-01-01", "2026-01-03", "2026-01-05", "2026-01-07", "2026-01-09", "2026-01-11",
				"2026-01-13", "2026-01-15", "2026-01-17", "2026-01-19", "2026-01-21", "2026-01-23",
				"2026-01-25", "2026-01-27", "2026-01-29", "2026-01-31",
			),
		},
		{
			name:  "simple every 2 weeks",
			input: "EVERY 2 WEEKS FROM 2026-01-01 TO 2026-02-01",
			expectedResults: buildTimeArray(
				"2026-01-01", "2026-01-02", "2026-01-03", "2026-01-04",
				"2026-01-12", "2026-01-13", "2026-01-14", "2026-01-15", "2026-01-16", "2026-01-17", "2026-01-18",
				"2026-01-26", "2026-01-27", "2026-01-28", "2026-01-29", "2026-01-30", "2026-01-31",
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			compiler := compiler.NewCompiler(test.input)
			results, err := compiler.Compile(false)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}

			if len(results) != len(test.expectedResults) {
				t.Errorf("expected %d results, got %d", len(test.expectedResults), len(results))
				return
			}

			for i, result := range results {
				if !result.Equal(test.expectedResults[i]) {
					t.Errorf("expected %v, got %v", test.expectedResults[i], result)
				}
			}
		})
	}
}
