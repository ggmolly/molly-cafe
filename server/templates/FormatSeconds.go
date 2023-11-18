package templates

import "fmt"

// Takes a number of seconds and returns a string formatted as HHhMM
func FormatSeconds(n int) string {
	return fmt.Sprintf("%2dh%02d", n/3600, n/60%60)
}
