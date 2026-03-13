package clock

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *clock
	}{
		{
			name: "success",
			want: &clock{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clock_Now_Fail(t *testing.T) {
	c := NewClock()

	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "success",
			want: c.Now(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// emulate delay.
			time.Sleep(1 * time.Second)

			if got := c.Now(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clock_Now_MockedSuccess(t *testing.T) {
	c := NewMockedClock()

	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "success",
			want: c.Now(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// emulate delay.
			time.Sleep(1 * time.Second)

			if got := c.Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}
