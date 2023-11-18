package templates

import "time"

// Takes a unix timestamp and returns a string formatted as "HH:MM"
func FormatTimestamp(n int64) string {
	return time.Unix(n, 0).Format("15:04")
}
