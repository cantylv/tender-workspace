package functions

import "time"

// FormatTime
// Returns format time. For formatting we should use this datetimetz `Mon Jan 2 15:04:05 MST 2006`.
func FormatTime(t time.Time) string {
	return t.Format("02.01.2006 15:04:05 UTC-07") // for me, it's more readable format
}
