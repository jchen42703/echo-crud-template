package templates

import "time"

func DayFromNow() time.Time {
	return time.Now().Add(24 * time.Hour)
}

// For redis expiration
func DayInSeconds() string {
	return "86400"
}
