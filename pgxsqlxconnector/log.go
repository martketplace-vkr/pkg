package pgxsqlxconnector

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
)

type logHook struct {
	cfg *Config
}

func NewLogHook(cfg *Config) *logHook {
	return &logHook{cfg: cfg}
}

func (h *logHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "begin", time.Now()), nil
}

func (h *logHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin, ok := ctx.Value("begin").(time.Time)
	if !ok {
		return ctx, nil
	}

	go func(begin time.Time) {
		log.DatabaseQuery(h.cfg.DBName, query, time.Since(begin))
		if h.cfg.HookFunc != nil {
			h.cfg.HookFunc(query, time.Since(begin).Milliseconds(), h.cfg.AppName)
		}
	}(begin)

	return ctx, nil
}
