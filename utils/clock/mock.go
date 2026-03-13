package clock

import "time"

var _ = Clock(mockedClock{})

type mockedClock struct{}

func NewMockedClock() *mockedClock {
	return &mockedClock{}
}

func (c mockedClock) Now() time.Time {
	// hitman's birthday.
	return time.Date(1964, time.September, 5, 0, 0, 0, 0, time.UTC)
}

func MockSingletonWithMock() {
	Time = NewMockedClock()
}
