package templates

import "fmt"

// Takes a number of seconds and returns a string formatted as HH:MM
func FormatDuration(n uint32) string {
	return fmt.Sprintf("%2d:%02d", n/3600, n/60%60)
}
