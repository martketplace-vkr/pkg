package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/martketplace-vkr/pkg/outbox/dto"
	"github.com/martketplace-vkr/pkg/tracer"
)

const queryInsertOutboxSendingBatch = `
	with input as (
		select *
		from unnest(
			$1::text[],  -- idempotency_key
			$2::text[],  -- topic
			$3::text[],  -- payload_text
			$4::text[]   -- event_type
		) with ordinality as t(idempotency_key, topic, payload_text, event_type, ord)
	),
	inserted as (
		insert into outbox_sending (idempotency_key, topic, payload, event_type)
		select i.idempotency_key, i.topic, i.payload_text::jsonb, i.event_type
		from input i
		on conflict (idempotency_key, topic, event_type) do nothing
		returning idempotency_key, topic, event_type
	)
	select i.ord, (inserted.idempotency_key is not null) as inserted
	from input i
	left join inserted
	  on inserted.idempotency_key = i.idempotency_key
	 and inserted.topic           = i.topic
	 and inserted.event_type      = i.event_type
	order by i.ord;
`

func (s *pgxPoolStorage) RegisterOutboxSendingBatch(
	ctx context.Context,
	events []dto.Event,
) (insertedIdx []int, conflictedIdx []int, err error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	if len(events) == 0 {
		return nil, nil, nil
	}

	keys := make([]string, len(events))
	topics := make([]string, len(events))
	payloads := make([]string, len(events))
	types := make([]string, len(events))

	for i, ev := range events {
		keys[i] = ev.IdempotencyKey
		topics[i] = ev.Topic
		payloads[i] = string(ev.Payload) // передаём text[], в sql каст к jsonb
		types[i] = ev.EventType
	}

	rows, qerr := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		queryInsertOutboxSendingBatch,
		pq.Array(keys),
		pq.Array(topics),
		pq.Array(payloads),
		pq.Array(types),
	)
	if qerr != nil {
		return nil, nil, fmt.Errorf("RegisterOutboxSendingBatch (pgx): query: %w", qerr)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ord      int64
			inserted bool
		)
		if err := rows.Scan(&ord, &inserted); err != nil {
			return nil, nil, fmt.Errorf("RegisterOutboxSendingBatch (pgx): scan: %w", err)
		}
		idx := int(ord - 1)
		if inserted {
			insertedIdx = append(insertedIdx, idx)
		} else {
			conflictedIdx = append(conflictedIdx, idx)
		}
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}

	return insertedIdx, conflictedIdx, nil
}

func (s *sqlxStorage) RegisterOutboxSendingBatch(
	ctx context.Context,
	events []dto.Event,
) (insertedIdx []int, conflictedIdx []int, err error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	if len(events) == 0 {
		return nil, nil, nil
	}

	keys := make([]string, len(events))
	topics := make([]string, len(events))
	payloads := make([]string, len(events))
	types := make([]string, len(events))

	for i, ev := range events {
		keys[i] = ev.IdempotencyKey
		topics[i] = ev.Topic
		payloads[i] = string(ev.Payload)
		types[i] = ev.EventType
	}

	rows, qerr := s.ctxGetter.DefaultTrOrDB(ctx, s.db).QueryContext(
		ctx,
		queryInsertOutboxSendingBatch,
		pq.Array(keys),
		pq.Array(topics),
		pq.Array(payloads),
		pq.Array(types),
	)
	if qerr != nil {
		return nil, nil, fmt.Errorf("RegisterOutboxSendingBatch (sqlx): query: %w", qerr)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ord      int64
			inserted bool
		)
		if err := rows.Scan(&ord, &inserted); err != nil {
			return nil, nil, fmt.Errorf("RegisterOutboxSendingBatch (sqlx): scan: %w", err)
		}
		idx := int(ord - 1)
		if inserted {
			insertedIdx = append(insertedIdx, idx)
		} else {
			conflictedIdx = append(conflictedIdx, idx)
		}
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}

	return insertedIdx, conflictedIdx, nil
}
