package uid

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *generator
	}{
		{
			name: "success",
			want: &generator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generator_NewString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: uuid.NewString(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator{}
			if got := g.NewString(); len(got) != len(tt.want) {
				t.Errorf("NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generator_NewUUID(t *testing.T) {
	tests := []struct {
		name string
		want uuid.UUID
	}{
		{
			name: "success",
			want: uuid.New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator{}
			if got := g.NewUUID(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generator_NewUUIDWithVersion(t *testing.T) {
	type args struct {
		version version
	}
	uidV6, _ := uuid.NewV6()
	uidV7, _ := uuid.NewV7()

	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name:    "success_v4",
			args:    args{version: v4},
			want:    uuid.New(),
			wantErr: false,
		},
		{
			name:    "success_v6",
			args:    args{version: v6},
			want:    uidV6,
			wantErr: false,
		},
		{
			name:    "success_v7",
			args:    args{version: v7},
			want:    uidV7,
			wantErr: false,
		},
		{
			name:    "invalid_version",
			args:    args{version: 999},
			want:    uuid.New(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generator{}
			got, err := g.NewUUIDWithVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUUIDWithVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewUUIDWithVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockedUIDOnEquality(t *testing.T) {
	mock := NewMockUID()

	assert.Equal(t, mock.NewString(), mock.NewString())
}
