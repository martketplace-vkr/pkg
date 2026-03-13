package duration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeconds_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		seconds Seconds
		want    string
		wantErr bool
	}{
		{
			name:    "marshal 5 seconds",
			seconds: Seconds{Duration: 5 * time.Second},
			want:    "5",
			wantErr: false,
		},
		{
			name:    "marshal 0 seconds",
			seconds: Seconds{Duration: 0},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.seconds.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestSeconds_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Seconds
		wantErr bool
	}{
		{
			name:    "unmarshal 5 seconds",
			input:   "5",
			want:    Seconds{Duration: 5 * time.Second},
			wantErr: false,
		},
		{
			name:    "unmarshal 0 seconds",
			input:   "0",
			want:    Seconds{Duration: 0},
			wantErr: false,
		},
		{
			name:    "unmarshal invalid",
			input:   `"invalid"`,
			want:    Seconds{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Seconds
			err := got.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
