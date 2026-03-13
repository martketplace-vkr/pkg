package duration

import (
	"database/sql/driver"
	"errors"
	"math"
	"time"
)

// Duration is a time.Duration alias that supports the driver.Valuer and
// sql.Scanner interfaces.
type Duration time.Duration

// ErrTooBig is returned by interval.Duration and Duration.Scan if the
// interval would overflow a time.Duration.
var ErrTooBig = errors.New("interval overflows time.Duration")

// Duration converts an interval into a time.Duration with the same
// semantics as `EXTRACT(EPOCH from <interval>)` in PostgreSQL.
func (i interval) Duration() (time.Duration, error) {
	dur := int64(i.Years())

	if dur > math.MaxInt64/nsPerYr || dur < math.MinInt64/nsPerYr {
		return 0, ErrTooBig
	}
	dur *= hrsPerYr
	dur += int64(i.hrs)

	if dur > math.MaxInt64/int64(time.Hour) || dur < math.MinInt64/int64(time.Hour) {
		return 0, ErrTooBig
	}
	dur *= int64(time.Hour)

	us := i.Microseconds() * int64(time.Microsecond)
	if dur > 0 {
		if math.MaxInt64-dur < us {
			return 0, ErrTooBig
		}
	} else {
		if math.MinInt64-dur > us {
			return 0, ErrTooBig
		}
	}
	dur += us

	return time.Duration(dur), nil
}

// Scan implements sql.Scanner.
func (d *Duration) Scan(src interface{}) error {
	ival := interval{}
	err := (&ival).Scan(src)
	if err != nil {
		return err
	}

	result, err := ival.Duration()
	if err != nil {
		return err
	}

	*d = Duration(result)
	return nil
}

// Value implements driver.Valuer.
func (d Duration) Value() (driver.Value, error) {
	return time.Duration(d), nil
}
