package uid

import (
	"github.com/jackc/pgx/v5/pgtype"
	"sync"

	"github.com/google/uuid"
)

var (
	uidOnce = new(sync.Once)
	UID     Generator
)

func init() {
	uidOnce.Do(func() {
		UID = New()
	})
}

func NewString() string {
	return UID.NewString()
}

func NewUUID() uuid.UUID {
	return UID.NewUUID()
}

func NewUUIDWithVersion(version version) (uuid.UUID, error) {
	return UID.NewUUIDWithVersion(version)
}

func ToPgType(uid uuid.UUID) pgtype.UUID {
	return UID.ToPgType(uid)
}

type (
	Generator interface {
		NewString() string
		NewUUID() uuid.UUID
		NewUUIDWithVersion(version version) (uuid.UUID, error)
		ToPgType(uid uuid.UUID) pgtype.UUID
	}
)

var _ = Generator(&generator{})

type generator struct{}

func New() *generator {
	return &generator{}
}

func (g generator) NewString() string {
	return uuid.NewString()
}

func (g generator) NewUUID() uuid.UUID {
	return uuid.New()
}

func (g generator) NewUUIDWithVersion(version version) (uuid.UUID, error) {
	switch version {
	case v4:
		return uuid.New(), nil
	case v6:
		return uuid.NewV6()
	case v7:
		return uuid.NewV7()
	default:
		return uuid.New(), nil
	}
}

func (g generator) ToPgType(uid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: uid,
		Valid: true,
	}
}
