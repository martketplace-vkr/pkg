package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/martketplace-vkr/pkg/inbox/dto"
)

// EffectsRecorder DB-level effect dedup for sqlx storage
// Before processing some event we register the action in database
// So we decrease the chance of double processing.
const queryInsertEffectsBatch = `
	with input as (
	  select *
	  from unnest($1::text[], $2::text[], $3::text[])
		   with ordinality as t(idempotency_key, effect, payload_text, ord)
	),
	inserted as (
	  insert into effect_log (idempotency_key, effect, payload)
	  select i.idempotency_key, i.effect, i.payload_text::jsonb
	  from input i
	  on conflict (idempotency_key, effect) do nothing
	  returning idempotency_key, effect
	)
	select i.ord, (inserted.idempotency_key is not null) as inserted
	from input i
	left join inserted using (idempotency_key, effect)
	order by i.ord;
`

func (s *pgxPoolStorage) RegisterEffectsBatch(
	ctx context.Context,
	events []dto.Event,
) (insertedIdx []int, conflictedIdx []int, err error) {
	if len(events) == 0 {
		return nil, nil, nil
	}

	idKeys := make([]string, len(events))
	payloads := make([]string, len(events))
	effects := make([]string, len(events))

	for i, ev := range events {
		idKeys[i] = ev.IdempotencyKey
		payloads[i] = string(ev.Payload)
		effects[i] = getEffectName(ev)
	}

	rows, qerr := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		queryInsertEffectsBatch,
		pq.Array(idKeys),
		pq.Array(effects),
		pq.Array(payloads),
	)
	if qerr != nil {
		return nil,
			nil,
			fmt.Errorf("RegisterEffectsBatch (pgx): query: %w", qerr)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ord      int64
			inserted bool
		)

		if scanErr := rows.Scan(
			&ord,
			&inserted,
		); scanErr != nil {
			return nil,
				nil,
				fmt.Errorf("RegisterEffectsBatch (pgx): scan: %w", scanErr)
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

func (s *sqlxStorage) RegisterEffectsBatch(
	ctx context.Context,
	events []dto.Event,
) (insertedIdx []int, conflictedIdx []int, err error) {
	if len(events) == 0 {
		return nil, nil, nil
	}

	idKeys := make([]string, len(events))
	payloads := make([]string, len(events))
	effects := make([]string, len(events))

	for i, ev := range events {
		idKeys[i] = ev.IdempotencyKey
		payloads[i] = string(ev.Payload)
		effects[i] = getEffectName(ev)
	}

	rows, qerr := s.ctxGetter.DefaultTrOrDB(ctx, s.db).QueryContext(
		ctx,
		queryInsertEffectsBatch,
		pq.Array(idKeys),
		pq.Array(effects),
		pq.Array(payloads),
	)
	if qerr != nil {
		return nil,
			nil,
			fmt.Errorf("RegisterEffectsBatch (sqlx): query: %w", qerr)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ord      int64
			inserted bool
		)

		if scanErr := rows.Scan(
			&ord,
			&inserted,
		); scanErr != nil {
			return nil,
				nil,
				fmt.Errorf("RegisterEffectsBatch (sqlx): scan: %w", scanErr)
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

func getEffectName(event dto.Event) string {
	return fmt.Sprintf("handle: %s", event.EventType)
}
