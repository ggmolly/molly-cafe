package entities

import (
	"os"
	"syscall"
	"time"
)

// Returns the creation date of the article (format: Monday, 4th September 2023)
func CreatedAt(html *string, path string) string {
	stat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	st := stat.Sys().(*syscall.Stat_t)
	return timeFormat(time.Unix(st.Ctim.Sec, st.Ctim.Nsec))
}
