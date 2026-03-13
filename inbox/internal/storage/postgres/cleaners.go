package postgres

import (
	"context"
	"github.com/lib/pq"
	"time"
)

const (
	queryDeleteEffectLogs = `
		WITH d AS (
			SELECT ctid FROM effect_log WHERE created_at < $1 ORDER BY created_at LIMIT $2
		)
		DELETE FROM effect_log e USING d WHERE e.ctid = d.ctid;
	`

	querySelectInboxForArchive = `
		SELECT id FROM inbox 
		WHERE processed_at IS NOT NULL AND processed_at < $1
		ORDER BY processed_at 
		LIMIT $2;
	`

	queryInsertInboxToDeleted = `
		INSERT INTO inbox_archive (
			id, key, topic, payload, trace_carrier, idempotency_key, 
			created_at, processed_at, locked_until, error, attempts, 
			event_type, meta, deleted_at
		)
		SELECT 
			id, key, topic, payload, trace_carrier, idempotency_key, 
			created_at, processed_at, locked_until, error, attempts, 
			event_type, meta, now()
		FROM inbox 
		WHERE id = ANY($1);
	`

	queryDeleteFromInboxByIDs = `
		DELETE FROM inbox WHERE id = ANY($1);
	`

	queryDeleteFromInboxDeletedOlderThan = `
		WITH d AS (
			SELECT ctid FROM inbox_archive WHERE deleted_at < $1 ORDER BY deleted_at LIMIT $2
		)
		DELETE FROM inbox_archive e USING d WHERE e.ctid = d.ctid;
	`
)

func (s *sqlxStorage) DeleteEffectLogOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	res, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).ExecContext(
		ctx,
		queryDeleteEffectLogs,
		olderThan,
		limit,
	)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *pgxPoolStorage) DeleteEffectLogOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	ct, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryDeleteEffectLogs,
		olderThan,
		limit,
	)
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), err
}

func (s *sqlxStorage) ArchiveProcessedInbox(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	db := s.ctxGetter.DefaultTrOrDB(ctx, s.db)

	// 1. Select events to archive.
	var ids []int64
	if err := db.SelectContext(
		ctx,
		&ids,
		querySelectInboxForArchive,
		olderThan,
		limit,
	); err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}

	// 2. Write them to archive.
	if _, err := db.ExecContext(
		ctx,
		queryInsertInboxToDeleted,
		pq.Array(ids),
	); err != nil {
		return 0, err
	}

	// 3. Delete events from the main table.
	res, err := db.ExecContext(
		ctx,
		queryDeleteFromInboxByIDs,
		pq.Array(ids),
	)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *pgxPoolStorage) ArchiveProcessedInbox(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	db := s.ctxGetter.DefaultTrOrDB(ctx, s.db)

	rows, err := db.Query(ctx, querySelectInboxForArchive, olderThan, limit)
	if err != nil {
		return 0, err
	}
	ids := make([]int64, 0, limit)

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			rows.Close()
			return 0, err
		}
		ids = append(ids, id)
	}

	rows.Close()
	if len(ids) == 0 {
		return 0, nil
	}

	if _, err = db.Exec(ctx, queryInsertInboxToDeleted, pq.Array(ids)); err != nil {
		return 0, err
	}

	ct, err := db.Exec(ctx, queryDeleteFromInboxByIDs, pq.Array(ids))
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), nil
}

func (s *sqlxStorage) DeleteFromInboxDeletedOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	res, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).ExecContext(
		ctx,
		queryDeleteFromInboxDeletedOlderThan,
		olderThan,
		limit,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *pgxPoolStorage) DeleteFromInboxDeletedOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	ct, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryDeleteFromInboxDeletedOlderThan,
		olderThan,
		limit,
	)
	return ct.RowsAffected(), err
}
