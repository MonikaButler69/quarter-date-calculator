package quarter_date_calculator

import (
	"testing"
	"time"
)

func TestNewQuarterCalculatorRejectsInvalidStartMonth(t *testing.T) {
	for _, m := range []int{0, 13, -1} {
		if _, err := NewQuarterCalculator(m); err == nil {
			t.Fatalf("NewQuarterCalculator(%d) expected error, got nil", m)
		}
	}
}

func TestNewQuarterCalculatorAcceptsAllMonths(t *testing.T) {
	for m := 1; m <= 12; m++ {
		c, err := NewQuarterCalculator(m)
		if err != nil {
			t.Fatalf("NewQuarterCalculator(%d) unexpected error: %v", m, err)
		}
		if c.StartMonth != m {
			t.Fatalf("StartMonth = %d, want %d", c.StartMonth, m)
		}
	}
}

func TestQuarterWithJanuaryStart(t *testing.T) {
	c, _ := NewQuarterCalculator(1)
	cases := []struct {
		month int
		want  int
	}{
		{1, 1}, {2, 1}, {3, 1},
		{4, 2}, {5, 2}, {6, 2},
		{7, 3}, {8, 3}, {9, 3},
		{10, 4}, {11, 4}, {12, 4},
	}
	for _, tc := range cases {
		d := QuarterDate{Year: 2024, Month: tc.month, Day: 15}
		if got := c.Quarter(d); got != tc.want {
			t.Errorf("Quarter(%v) = %d, want %d", d, got, tc.want)
		}
	}
}

func TestQuarterWithFebruaryStart(t *testing.T) {
	c, _ := NewQuarterCalculator(2)
	cases := []struct {
		month int
		want  int
	}{
		{2, 1}, {3, 1}, {4, 1},
		{5, 2}, {6, 2}, {7, 2},
		{8, 3}, {9, 3}, {10, 3},
		{11, 4}, {12, 4}, {1, 4},
	}
	for _, tc := range cases {
		d := QuarterDate{Year: 2024, Month: tc.month, Day: 15}
		if got := c.Quarter(d); got != tc.want {
			t.Errorf("Quarter(%v) = %d, want %d", d, got, tc.want)
		}
	}
}

func TestQuarterWithDecemberStart(t *testing.T) {
	c, _ := NewQuarterCalculator(12)
	cases := []struct {
		month int
		want  int
	}{
		{12, 1}, {1, 1}, {2, 1},
		{3, 2}, {4, 2}, {5, 2},
		{6, 3}, {7, 3}, {8, 3},
		{9, 4}, {10, 4}, {11, 4},
	}
	for _, tc := range cases {
		d := QuarterDate{Year: 2024, Month: tc.month, Day: 15}
		if got := c.Quarter(d); got != tc.want {
			t.Errorf("Quarter(%v) = %d, want %d", d, got, tc.want)
		}
	}
}

func TestFiscalYearWithFebruaryStart(t *testing.T) {
	c, _ := NewQuarterCalculator(2)
	cases := []struct {
		date QuarterDate
		want int
	}{
		{QuarterDate{2024, 2, 1}, 2024},
		{QuarterDate{2024, 1, 31}, 2023},
		{QuarterDate{2024, 12, 31}, 2024},
		{QuarterDate{2023, 1, 1}, 2022},
	}
	for _, tc := range cases {
		if got := c.FiscalYear(tc.date); got != tc.want {
			t.Errorf("FiscalYear(%v) = %d, want %d", tc.date, got, tc.want)
		}
	}
}

func TestQuarterBoundsWithFebruaryStart(t *testing.T) {
	c, _ := NewQuarterCalculator(2)
	d := QuarterDate{Year: 2024, Month: 3, Day: 15}
	start, end := c.QuarterBounds(d)
	wantStart := QuarterDate{Year: 2024, Month: 2, Day: 1}
	wantEnd := QuarterDate{Year: 2024, Month: 4, Day: 30}
	if start != wantStart {
		t.Errorf("start = %v, want %v", start, wantStart)
	}
	if end != wantEnd {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestQuarterBoundsAcrossYearBoundary(t *testing.T) {
	c, _ := NewQuarterCalculator(11)
	d := QuarterDate{Year: 2024, Month: 1, Day: 10}
	start, end := c.QuarterBounds(d)
	wantStart := QuarterDate{Year: 2023, Month: 11, Day: 1}
	wantEnd := QuarterDate{Year: 2024, Month: 1, Day: 31}
	if start != wantStart {
		t.Errorf("start = %v, want %v", start, wantStart)
	}
	if end != wantEnd {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestQuarterBoundsLeapYearFebruary(t *testing.T) {
	c, _ := NewQuarterCalculator(1)
	d := QuarterDate{Year: 2024, Month: 2, Day: 15}
	start, end := c.QuarterBounds(d)
	wantStart := QuarterDate{Year: 2024, Month: 1, Day: 1}
	wantEnd := QuarterDate{Year: 2024, Month: 3, Day: 31}
	if start != wantStart {
		t.Errorf("start = %v, want %v", start, wantStart)
	}
	if end != wantEnd {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestQuarterBoundsNonLeapYearFebruary(t *testing.T) {
	c, _ := NewQuarterCalculator(1)
	d := QuarterDate{Year: 2023, Month: 2, Day: 15}
	start, end := c.QuarterBounds(d)
	wantStart := QuarterDate{Year: 2023, Month: 1, Day: 1}
	wantEnd := QuarterDate{Year: 2023, Month: 3, Day: 31}
	if start != wantStart {
		t.Errorf("start = %v, want %v", start, wantStart)
	}
	if end != wantEnd {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestNewQuarterDateValidation(t *testing.T) {
	if _, err := NewQuarterDate(2024, 2, 29); err != nil {
		t.Fatalf("leap day should be valid: %v", err)
	}
	if _, err := NewQuarterDate(2023, 2, 29); err == nil {
		t.Fatal("Feb 29 in non-leap year should be invalid")
	}
	if _, err := NewQuarterDate(2024, 4, 31); err == nil {
		t.Fatal("April 31 should be invalid")
	}
	if _, err := NewQuarterDate(2024, 13, 1); err == nil {
		t.Fatal("month 13 should be invalid")
	}
	if _, err := NewQuarterDate(2024, 1, 0); err == nil {
		t.Fatal("day 0 should be invalid")
	}
}

func TestFromTime(t *testing.T) {
	tm := time.Date(2024, time.March, 15, 12, 0, 0, 0, time.UTC)
	got := FromTime(tm)
	want := QuarterDate{Year: 2024, Month: 3, Day: 15}
	if got != want {
		t.Errorf("FromTime(%v) = %v, want %v", tm, got, want)
	}
}
