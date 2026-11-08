# quarter-date-calculator

Maps civil dates to fiscal quarters given a configurable start month and computes the inclusive first and last dates of each quarter.

## Usage

```go
package main

import (
	"fmt"

	qdc "example.com/quarter-date-calculator"
)

func main() {
	calc, err := qdc.NewQuarterCalculator(2) // fiscal year starts in February
	if err != nil {
		panic(err)
	}

	d, err := qdc.NewQuarterDate(2024, 3, 15)
	if err != nil {
		panic(err)
	}

	fmt.Println(calc.Quarter(d))      // 1
	fmt.Println(calc.FiscalYear(d))   // 2024

	start, end := calc.QuarterBounds(d)
	fmt.Println(start) // {2024 2 1}
	fmt.Println(end)   // {2024 4 30}
}
```

## Why this exists

Fiscal calendars rarely align with the Gregorian calendar. A company whose fiscal year starts in February has quarters that begin in February, May, August, and November, and a date in January belongs to the fiscal year that began the previous calendar year. The standard library's `time` package does not model this, and pulling in a date library for one calculation is heavier than the problem warrants.

The trade-off here is that `QuarterDate` is a plain civil date with no location. Fiscal calendars are defined in local terms, so carrying a location would invite off-by-one errors at midnight boundaries without adding any real capability. Callers who have a `time.Time` can convert with `FromTime`, which uses the time's own location.

## Edge cases

A fiscal year that starts in February means January dates belong to the previous fiscal year. `FiscalYear` returns the calendar year in which the fiscal year *starts*, so 31 January 2024 with a February start is fiscal year 2023. This is the interpretation the tests assert.

Quarter bounds are inclusive on both ends. The end date is the last day of the third month in the quarter, which is why `QuarterBounds` needs to know about leap years.
