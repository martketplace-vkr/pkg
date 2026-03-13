package duration

import (
	"fmt"
	"time"
)

const (
	SecondsInMinute = 60
	SecondsInHour   = 60 * SecondsInMinute
	SecondsInDay    = 24 * SecondsInHour
	SecondsInMonth  = 30 * SecondsInDay
	SecondsInYear   = 365 * SecondsInDay
)

// Format returns a string representation of the duration in the format "XhYmZs".
func Format(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	totalSeconds %= SecondsInYear
	totalSeconds %= SecondsInMonth
	totalSeconds %= SecondsInDay
	hours := totalSeconds / SecondsInHour
	totalSeconds %= SecondsInHour
	minutes := totalSeconds / SecondsInMinute
	seconds := totalSeconds % SecondsInMinute

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
}
