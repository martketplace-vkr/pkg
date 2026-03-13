package duration

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

// interval can represent the full range of PostgreSQL's interval type.
type interval struct {
	// the top bit is the sign for the microseconds,
	// bottom 29 are the signed year.
	yrs uint32

	hrs int32

	// it takes 33 bits (ouch) to fit microseconds-per-hour with sign,
	// but we have extra space in 'yrs' so the top bit there is the
	// sign for these microseconds
	us uint32
}

// New creates an interval.
func New(years, days, hours, minutes, seconds, microseconds int) interval {
	if years > maxYear || years < minYear || int64(hours) > maxHour || hours < minHour {
		panic("interval outside range")
	}

	microsecs := int64(microseconds) + int64(seconds)*usPerSec + int64(minutes)*usPerMin
	hrs := int64(hours) + int64(days)*24 + microsecs/usPerHr
	yrs := int64(years)

	microsecs %= usPerHr

	if yrs < 0 {
		yrs = (-yrs) | yrSignBit
	}
	if microsecs < 0 {
		yrs |= usSignBit
		microsecs *= -1
	}

	return interval{
		yrs: uint32(yrs),
		hrs: int32(hrs),
		us:  uint32(microsecs),
	}
}

// Years returns the number of years in the interval.
func (i interval) Years() int32 {
	years := int32(i.yrs & (yrSignBit - 1))
	if i.yrs&yrSignBit != 0 {
		years *= -1
	}
	return years
}

// Hours returns the number of hours in the interval.
func (i interval) Hours() int32 {
	return i.hrs
}

// Microseconds returns the number of microseconds in the interval,
// up to the number of microseconds in an hour.
func (i interval) Microseconds() int64 {
	us := int64(i.us)
	if i.yrs&usSignBit != 0 {
		us *= -1
	}
	return us
}

// Scan implements sql.Scanner.
func (i *interval) Scan(src interface{}) error {
	var s string
	switch x := src.(type) {
	case string:
		s = x
	case []byte:
		s = string(x)
	default:
		return errors.New(
			"pqinterval: converting driver.Value type %T (%q) to string: invalid syntax",
		)
	}

	result, err := parse(s)
	if err != nil {
		return err
	}

	*i = result
	return nil
}

// Value implements driver.Valuer.
func (i interval) Value() (driver.Value, error) {
	var years, months, days, hours, mins, secs, msecs, usecs int64
	years = int64(i.Years())
	hours = int64(i.Hours())
	usecs = int64(i.Microseconds())
	days, hours = divmod(hours, 24)
	mins, usecs = divmod(usecs, usPerMin)
	secs, usecs = divmod(usecs, usPerSec)
	msecs, usecs = divmod(usecs, 1000)
	return formatInput(years, months, days, hours, mins, secs, msecs, usecs), nil
}

// formatValue produces a string in the format that postgres expects for interval input.
// (https://www.postgresql.org/docs/current/static/datatype-datetime.html#DATATYPE-INTERVAL-INPUT)
func formatInput(years, months, days, hours, mins, secs, msecs, usecs int64) string {
	pieces := make([]string, 0, 8)
	if years != 0 {
		pieces = append(pieces, fmt.Sprintf("%d years", years))
	}
	if months != 0 {
		pieces = append(pieces, fmt.Sprintf("%d months", months))
	}
	if days != 0 {
		pieces = append(pieces, fmt.Sprintf("%d days", days))
	}
	if hours != 0 {
		pieces = append(pieces, fmt.Sprintf("%d hours", hours))
	}
	if mins != 0 {
		pieces = append(pieces, fmt.Sprintf("%d minutes", mins))
	}
	if secs != 0 {
		pieces = append(pieces, fmt.Sprintf("%d seconds", secs))
	}
	if msecs != 0 {
		pieces = append(pieces, fmt.Sprintf("%d milliseconds", msecs))
	}
	if usecs != 0 || len(pieces) == 0 {
		pieces = append(pieces, fmt.Sprintf("%d microseconds", usecs))
	}
	return strings.Join(pieces, " ")
}

func divmod(num int64, denom int64) (int64, int64) {
	return num / denom, num % denom
}

const (
	// the year range allowed in PostgreSQL intervals.
	maxYear = 0xaaaaaaa
	minYear = -0xaaaaaaa

	maxHour = 1 << 31
	minHour = -1 << 31

	yrSignBit = 0x10000000
	usSignBit = 0x80000000

	usPerSec = 1000000
	usPerMin = usPerSec * 60
	usPerHr  = usPerMin * 60

	// assumptions embedded in PostgreSQL's EXTRACT(EPOCH FROM <interval>)
	daysPerYr  = 365.25
	daysPerMon = 30

	hrsPerYr = daysPerYr * 24
	nsPerYr  = int64(hrsPerYr * time.Hour)
)
