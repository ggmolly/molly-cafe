package entities

import (
	"fmt"
	"time"
)

func ordinal(n int) string {
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// Returns any date in the format: Monday, 4th September 2023
func timeFormat(t time.Time) string {
	dayOfWeek := t.Weekday().String()
	dayOfMonth := t.Day()
	month := t.Month().String()
	year := t.Year()
	return fmt.Sprintf("%s, %d%s %s %d", dayOfWeek, dayOfMonth, ordinal(dayOfMonth), month, year)
}
