package templates

import "time"

func TimeToUnix(d time.Time) int64 {
	return d.Unix()
}
