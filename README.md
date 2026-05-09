# Cadence

Cadence is an open-source DSL for generating ISO 8601 dates based on a simple, SQL-like grammar.

## Installation

```sh
go get github.com/connordoman/cadence
```

## Usage

```go
expression := "EVERY TUE FROM 2026-01-01 TO 2026-02-01"

results, err := cadence.CompileAsJSON(expression)
if err != nil {
  panic()
}

fmt.Println(results)

// ["2026-01-06","2026-01-13","2026-01-20","2026-01-27"]
```

## Syntax

Cadence uses an SQL-adjacent syntax for "querying" the calendar.

```
EVERY (<interval> | <selector>) [date range]
```

The most basic expression in Cadence is `EVERY DAY`. The default date range is the from the current date to one year later. If you were to execute `EVERY DAY` on January 1, 2026, you would get 365 days from `2026-01-01` to `2026-12-31`.

### Intervals

An _interval_ is how often the specified selector should repeat. An interval is optional if a selector is specified, but is required otherwise. Cadence supports 3 types of intervals:

- Days
  - `DAY`, `DAYS`
- Weeks
  - `WEEK`, `WEEKS`
- Months
  - `MONTH`, `MONTHS`

By default, an interval is 1 unit. This is called an _implicit interval_. An interval can also be _explicit_, where the length of the interval is specified:

```
EVERY 2 DAYS
```

When using an explicit `WEEKS` interval, the offset is determined by the `FROM` date. For example, January 1, 2026 is a Thursday, so when executing `EVERY 2 WEEKS FROM 2026-01-01` the result will be:

- Thu Jan 1
- Fri Jan 2
- Sat Jan 3
- Sun Jan 4
- Mon Jan 12
- Tue Jan 13
- etc.

> [!IMPORTANT]
> If no selector is specified, the default is every day in the interval, equivalent to `EVERY <interval> ON MON, TUE, WED, THU, FRI, SAT, SUN`

### Selectors

A _selector_ is used to specify some subset of dates in an interval. A selector is optional if an interval is specified, but is required otherwise. There are 3-ish selectors in Cadence:

- Day List

  - `TUE, WED, THU`

- Ordinal Day List

  - `FIRST SAT, LAST FRI`
  - May only be used with a `MONTH` interval

- Weekdays (special case of Day List)
  - `WEEKDAY` (equivalent to `MON, TUE, WED, THU, FRI`)

> [!IMPORTANT]
> If no interval is specified, the default is `1 DAY`, except in the case of ordinal day lists, in which case the default is `1 MONTH`.

### Date Range

The _date range_ is used to limit the span of the output. Date ranges are optional. A date range can be _closed_ or _open_:

- Open
  - `FROM <date>`
  - Begins at `<date>`, and continues for one year
- Closed
  - `FROM <date 1> TO <date 2>`
  - Range ends at `<date 2>`, **exclusive**

> [!IMPORTANT]
> If no date range is specified, the default is `FROM <now> TO <now + 1 year>`.

## Examples

_Assume all examples without a specified date range were executed on `2026-02-14`._

```
EVERY DAY
```

Result: `2026-02-14, 2026-02-15, 2026-02-16, ..., 2027-02-12, 2027-02-13`
