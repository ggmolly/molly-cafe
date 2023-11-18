package templates

import "fmt"

// Takes a number of seconds and returns a string formatted as MM:SS
func FormatDuration(n uint32) string {
	return fmt.Sprintf("%2d:%02d", n/60, n%60)
}
