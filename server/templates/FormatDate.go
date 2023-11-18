package templates

import (
	"html/template"
	"time"
)

func FormatDate(d time.Time) template.HTML {
	suffix := "th"
	switch d.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return template.HTML(d.Format("Monday, 2" + suffix + " January 2006"))
}
