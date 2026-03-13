package duration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDays_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		days    Days
		want    string
		wantErr bool
	}{
		{
			name:    "marshal 5 days",
			days:    Days{Duration: 5 * 24 * time.Hour},
			want:    "5",
			wantErr: false,
		},
		{
			name:    "marshal 0 days",
			days:    Days{Duration: 0},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.days.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}
