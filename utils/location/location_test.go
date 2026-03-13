package location

import (
	"testing"
	"time"
)

func TestLocation_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Location *time.Location
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test UnmarshalJSON",
			fields: fields{
				Location: time.UTC,
			},
			args: args{
				data: []byte(`"UTC"`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Location{
				Location: tt.fields.Location,
			}
			if err := l.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
