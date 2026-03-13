package postgres

import (
	"context"
	"github.com/lib/pq"
	"time"
)

const (
	querySelectOutboxForArchive = `
		SELECT id FROM outbox
		WHERE sent_at IS NOT NULL AND sent_at < $1
		ORDER BY sent_at
		LIMIT $2;
	`

	queryInsertOutboxToArchive = `
		INSERT INTO outbox_archive (
		    id, key, topic, payload, trace_carrier, idempotency_key,
		    created_at, sent_at, locked_until, error, attempts, event_type, meta, deleted_at
		)
		SELECT
		    id, key, topic, payload, trace_carrier, idempotency_key,
		    created_at, sent_at, locked_until, last_error, attempts, event_type, meta, now()
		FROM outbox
		WHERE id = ANY($1);
	`

	queryDeleteFromOutboxByIDs = `
		DELETE FROM outbox WHERE id = ANY($1);
	`

	queryDeleteFromOutboxArchiveOlderThan = `
		WITH d AS (
			SELECT ctid FROM outbox_archive WHERE deleted_at < $1 ORDER BY deleted_at LIMIT $2
		)
		DELETE FROM outbox_archive e USING d WHERE e.ctid = d.ctid;
	`

	queryDeleteOutboxSent = `
		WITH d AS (
			SELECT ctid FROM outbox_sending WHERE created_at < $1 ORDER BY created_at LIMIT $2
		)
		DELETE FROM outbox_sending e USING d WHERE e.ctid = d.ctid;
	`
)

func (s *sqlxStorage) ArchiveSentOutbox(
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
		querySelectOutboxForArchive,
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
		queryInsertOutboxToArchive,
		pq.Array(ids),
	); err != nil {
		return 0, err
	}

	// 3. Delete events from the main table.
	res, err := db.ExecContext(
		ctx,
		queryDeleteFromOutboxByIDs,
		pq.Array(ids),
	)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *pgxPoolStorage) ArchiveSentOutbox(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	db := s.ctxGetter.DefaultTrOrDB(ctx, s.db)

	rows, err := db.Query(ctx, querySelectOutboxForArchive, olderThan, limit)
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

	if _, err = db.Exec(ctx, queryInsertOutboxToArchive, pq.Array(ids)); err != nil {
		return 0, err
	}

	ct, err := db.Exec(ctx, queryDeleteFromOutboxByIDs, pq.Array(ids))
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), nil
}

func (s *sqlxStorage) DeleteFromOutboxArchiveOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	res, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).ExecContext(
		ctx,
		queryDeleteFromOutboxArchiveOlderThan,
		olderThan,
		limit,
	)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *pgxPoolStorage) DeleteFromOutboxArchiveOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	ct, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryDeleteFromOutboxArchiveOlderThan,
		olderThan,
		limit,
	)
	return ct.RowsAffected(), err
}

func (s *sqlxStorage) DeleteOutboxSentOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	res, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).ExecContext(
		ctx,
		queryDeleteOutboxSent,
		olderThan,
		limit,
	)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *pgxPoolStorage) DeleteOutboxSentOlderThan(
	ctx context.Context,
	olderThan time.Time,
	limit int,
) (int64, error) {
	ct, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryDeleteOutboxSent,
		olderThan,
		limit,
	)
	if err != nil {
		return 0, err
	}

	return ct.RowsAffected(), err
}
