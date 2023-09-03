package entities

import "fmt"

const (
	AVG_WPM = 200
)

// Returns the estimated read time of the article in minutes
func ReadTime(html *string, path string) string {
	var minutes, words uint32
	for _, c := range *html {
		if c == ' ' {
			words++
		}
	}
	minutes = words / AVG_WPM
	if minutes == 0 {
		return "less than a minute"
	}
	return fmt.Sprintf("%d minute%s", minutes, func() string {
		if minutes > 1 {
			return "s"
		}
		return ""
	}()) // tenaries are cool in go ^_^
}
