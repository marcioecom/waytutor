package quota

import "time"

// Period is the UTC month that usage counts against.
func Period(now time.Time) string {
	return now.UTC().Format("2006-01")
}
