package entities

import (
	"os"
)

// Returns the last modification date of the article (format: Monday, 4th September 2023)
func UpdatedAt(html *string, path string) string {
	stat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return timeFormat(stat.ModTime())
}
