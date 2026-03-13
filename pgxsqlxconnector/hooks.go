package pgxsqlxconnector

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

type logData struct {
	Query     string    `bson:"query"`
	Took      int64     `bson:"took"`
	CreatedAt time.Time `bson:"createdAt"`
}

type CompositeHook struct {
	Hooks []sqlhooks.Hooks
}

func NewCompositeHook(hooks ...sqlhooks.Hooks) sqlhooks.Hooks {
	return &CompositeHook{
		Hooks: hooks,
	}
}

func (ch *CompositeHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	var err error
	for _, h := range ch.Hooks {
		ctx, err = h.Before(ctx, query, args...)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func (ch *CompositeHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	var err error
	for _, h := range ch.Hooks {
		ctx, err = h.After(ctx, query, args...)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
