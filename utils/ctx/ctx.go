package ctx

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func LightEqual(ctx context.Context) interface{} {
	return mock.MatchedBy(func(ctx context.Context) bool {
		return true
	})
}
