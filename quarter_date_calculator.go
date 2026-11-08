// Package quarter_date_calculator maps dates to fiscal quarters given a
// configurable start month and computes quarter boundaries.
//
// The package deliberately avoids time zones: a QuarterDate is a civil date
// with no location. Fiscal calendars are usually defined in local terms, and
// attaching a location would force callers to reason about offsets that the
// fiscal calendar itself does not care about.
package quarter_date_calculator

import (
	"errors"
	"time"
)

// QuarterDate is a civil date: a year, a month (1-12), and a day (1-31).
// It is not tied to any time zone or location.
type QuarterDate struct {
	Year  int
	Month int
	Day   int
}

// NewQuarterDate returns a QuarterDate for the given civil year, month, and day.
// It validates that the date exists in the proleptic Gregorian calendar.
func NewQuarterDate(year, month, day int) (QuarterDate, error) {
	if month < 1 || month > 12 {
		return QuarterDate{}, errors.New("month must be between 1 and 12")
	}
	if day < 1 || day > daysInMonth(year, month) {
		return QuarterDate{}, errors.New("day is out of range for the month")
	}
	return QuarterDate{Year: year, Month: month, Day: day}, nil
}

// QuarterCalculator maps dates to fiscal quarters using a configurable start
// month. StartMonth is the first month of the fiscal year, so a start month of
// 2 means the fiscal year begins in February.
type QuarterCalculator struct {
	StartMonth int
}

// NewQuarterCalculator returns a calculator whose fiscal year begins in the
// given month. The month must be between 1 and 12.
func NewQuarterCalculator(startMonth int) (*QuarterCalculator, error) {
	if startMonth < 1 || startMonth > 12 {
		return nil, errors.New("start month must be between 1 and 12")
	}
	return &QuarterCalculator{StartMonth: startMonth}, nil
}

// Quarter returns the one-based fiscal quarter number (1-4) that contains the
// given date. The fiscal year is divided into four quarters of three months
// each, starting at StartMonth.
func (c *QuarterCalculator) Quarter(d QuarterDate) int {
	offset := (d.Month - c.StartMonth + 12) % 12
	return offset/3 + 1
}

// FiscalYear returns the fiscal year label for the given date. The fiscal year
// is named after the calendar year in which it starts. A date in January with a
// February start month belongs to the fiscal year that began the previous
// calendar year.
func (c *QuarterCalculator) FiscalYear(d QuarterDate) int {
	if d.Month >= c.StartMonth {
		return d.Year
	}
	return d.Year - 1
}

// QuarterBounds returns the inclusive first date and inclusive last date of the
// fiscal quarter containing d. The quarter is determined by Quarter and
// FiscalYear, and the bounds are the first and last days of the three calendar
// months that make up the quarter.
func (c *QuarterCalculator) QuarterBounds(d QuarterDate) (start, end QuarterDate) {
	q := c.Quarter(d)
	fy := c.FiscalYear(d)

	startMonth := c.StartMonth + (q-1)*3
	startYear := fy
	if startMonth > 12 {
		startMonth -= 12
		startYear++
	}

	endMonth := startMonth + 2
	endYear := startYear
	if endMonth > 12 {
		endMonth -= 12
		endYear++
	}

	start = QuarterDate{Year: startYear, Month: startMonth, Day: 1}
	end = QuarterDate{Year: endYear, Month: endMonth, Day: daysInMonth(endYear, endMonth)}
	return start, end
}

// daysInMonth returns the number of days in the given month of the given year.
// It is used instead of time.Date so that the package does not depend on the
// time package's location handling for a purely civil calculation.
func daysInMonth(year, month int) int {
	switch month {
	case 2:
		if isLeapYear(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

// isLeapYear reports whether year is a leap year in the proleptic Gregorian
// calendar. The rule is: divisible by 4, except centuries not divisible by 400.
func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

// FromTime converts a time.Time to a QuarterDate using its calendar date in the
// time's own location. This is the only place the package touches time.Time;
// callers who need a specific location should call t.In(loc) first.
func FromTime(t time.Time) QuarterDate {
	y, m, d := t.Date()
	return QuarterDate{Year: y, Month: int(m), Day: d}
}
