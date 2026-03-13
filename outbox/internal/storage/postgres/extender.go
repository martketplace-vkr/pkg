package postgres

import (
	"context"
	"github.com/lib/pq"
	"time"
)

const (
	maxLockExtensionInterval = "1 hour"

	queryExtendLockByID = `
		UPDATE outbox
		SET locked_until = LEAST(now() + $2::interval, now() + interval '` + maxLockExtensionInterval + `')
		WHERE id = ANY($1) AND processed_at IS NULL;
	`
)

func (s *sqlxStorage) ExtendLocksByIDs(
	ctx context.Context,
	ids []int64,
	extendBy time.Duration,
) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).ExecContext(
		ctx,
		queryExtendLockByID,
		pq.Array(ids),
		extendBy.String(),
	)
	return err
}

func (s *pgxPoolStorage) ExtendLocksByIDs(
	ctx context.Context,
	ids []int64,
	extendBy time.Duration,
) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryExtendLockByID,
		pq.Array(ids),
		extendBy.String(),
	)
	return err
}
