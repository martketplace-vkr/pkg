package clock

import (
	"sync"
	"time"
)

var (
	timeOnce = new(sync.Once)
	Time     Clock
)

func init() {
	timeOnce.Do(func() {
		Time = NewClock()
	})
}

func Now() time.Time {
	return Time.Now()
}

type Clock interface {
	Now() time.Time
}

var _ = Clock(clock{})

type clock struct{}

func NewClock() *clock {
	return &clock{}
}

func (c clock) Now() time.Time {
	return time.Now()
}
