package uid

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	mockUID = "989c9308-b919-4395-9764-9e7d137474e0"
)

var _ = Generator(MockUID{})

type MockUID struct {
	uid string
}

func NewMockUID() *MockUID {
	return &MockUID{
		uid: mockUID,
	}
}

func NewCustomUID(uid string) *MockUID {
	return &MockUID{
		uid: uid,
	}
}

func (uid MockUID) NewString() string {
	return uid.uid
}

func (uid MockUID) NewUUID() uuid.UUID {
	return uuid.MustParse(uid.uid)
}

func (uid MockUID) NewUUIDWithVersion(_ version) (uuid.UUID, error) {
	return uuid.Parse(uid.uid)
}

func (uid MockUID) ToPgType(_ uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: uuid.MustParse(uid.uid),
		Valid: true,
	}
}
func ReplaceSingletonWithMock() {
	UID = NewMockUID()
}
