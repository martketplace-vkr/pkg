package postgres

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/inbox/internal/domain"

	txManager "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/martketplace-vkr/pkg/tracer"
)

type pgxPoolStorage struct {
	db        *pgxpool.Pool
	ctxGetter *txManager.CtxGetter

	eventLockTimeout time.Duration
	maxAttempts      int
	queryFetchEvents string
}

func NewPgxPoolStorage(
	db *pgxpool.Pool,
	ctxGetter *txManager.CtxGetter,
	eventLockTimeout time.Duration,
	maxAttempts int,
	disableSequentialEntriesSelection bool,
) *pgxPoolStorage {
	query := queryFetchUnprocessedEventsWithOrder
	if disableSequentialEntriesSelection {
		query = queryFetchUnprocessedEventsUnordered
	}

	return &pgxPoolStorage{
		db:               db,
		ctxGetter:        ctxGetter,
		eventLockTimeout: eventLockTimeout,
		maxAttempts:      maxAttempts,
		queryFetchEvents: query,
	}
}

func (s *pgxPoolStorage) CreateInboxEvent(ctx context.Context, dtoEvent dto.Event) (int64, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var (
		id    int64
		event = eventFromDto(dtoEvent)
	)

	err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).QueryRow(
		ctx,
		queryCreateEvent,
		event.Topic,
		event.Payload,
		event.TraceCarrier.String(),
		event.Key,
		event.EventType,
		event.CreatedAt,
		event.IdempotencyKey,
		event.Error,
		event.State,
	).Scan(&id)
	if err != nil {
		return id, tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.QueryRow(event: %+v)",
			caller,
			event,
		)
	}

	return id, err
}

func (s *pgxPoolStorage) BatchCreateInboxEvents(
	ctx context.Context,
	dtoEvents dto.Events,
) ([]int64, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var (
		ids    []int64
		events = eventsFromDto(dtoEvents)
	)

	batch := domain.NewEventBatchFromCommand(events)

	rows, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		queryInsertBatch,
		batch.Topics,
		batch.Payloads,
		batch.TraceCarriers,
		batch.Keys,
		batch.EventTypes,
		batch.CreatedAts,
		batch.IdempotencyKeys,
		batch.Errors,
		batch.Attempts,
		batch.States,
	)
	if err != nil {
		return ids, tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.Query(events: %+v)",
			caller,
			events,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return ids, tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.Scan()",
				caller,
			)
		}
		ids = append(ids, id)
	}

	return ids, err
}

func (s *pgxPoolStorage) FetchUnprocessedEvents(ctx context.Context, batchSize int) (dto.Events, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var events dto.Events

	rows, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		s.queryFetchEvents,
		time.Now().Add(s.eventLockTimeout),
		s.maxAttempts,
		batchSize,
	)
	if err != nil {
		return nil, tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.Query(batchSize: %d)",
			caller,
			batchSize,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(
			&event.ID,
			&event.Topic,
			&event.Payload,
			&event.TraceCarrier,
			&event.CreatedAt,
			&event.ProcessedAt,
			&event.Key,
			&event.EventType,
			&event.IdempotencyKey,
		); err != nil {
			return nil, tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.Scan()",
				caller,
			)
		}

		events = append(events, eventToDto(event))
	}

	return events, err
}

func (s *pgxPoolStorage) MarkEventsAsProcessed(ctx context.Context, events dto.SuccessEvents) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	batch := &pgx.Batch{}
	for _, event := range events {
		batch.Queue(queryMarkEventAsPublished, event, dto.StateDone)
	}

	br := s.ctxGetter.DefaultTrOrDB(ctx, s.db).SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(events); i++ {
		_, err := br.Exec()
		if err != nil {
			return tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.Exec(eventID: %d)",
				caller,
				events[i],
			)
		}
	}

	return nil
}

func (s *pgxPoolStorage) FetchExpiredEvents(ctx context.Context, TTL, batchSize int) (dto.Events, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var events dto.Events

	rows, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		queryFetchExpiredEvents,
		TTL,
		batchSize,
	)
	if err != nil {
		return nil, tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.Query(EventLockTTL: %d, batchSize: %d)",
			caller,
			TTL,
			batchSize,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(
			&event.ID,
			&event.Topic,
			&event.Payload,
			&event.TraceCarrier,
			&event.CreatedAt,
			&event.ProcessedAt,
			&event.Key,
			&event.EventType,
			&event.IdempotencyKey,
		); err != nil {
			return nil, tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.Scan()",
				caller,
			)
		}
		events = append(events, eventToDto(event))
	}

	return events, err
}

func (s *pgxPoolStorage) DeleteEventsByIDs(ctx context.Context, ids []int64) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	_, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryDeleteEventsByIDs,
		pq.Int64Array(ids),
	)
	if err != nil {
		return tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.Exec(ids: %+v)",
			caller,
			ids,
		)
	}

	return err
}

func (s *pgxPoolStorage) UnlockFailedEvents(ctx context.Context, events dto.FailedEvents) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	batch := &pgx.Batch{}

	for _, event := range events {
		batch.Queue(queryMarkEventsAsFailed, event.ID, event.Error)
	}

	results := s.ctxGetter.DefaultTrOrDB(ctx, s.db).SendBatch(ctx, batch)
	defer results.Close()

	for range events {
		_, err := results.Exec()
		if err != nil {
			return tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.SendBatch failed",
				caller,
			)
		}
	}

	return nil
}

func (s *pgxPoolStorage) UnlockSkippedEvents(ctx context.Context, events []int64) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	_, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Exec(
		ctx,
		queryMarkEventsAsSkipped,
		pq.Array(events),
	)
	if err != nil {
		return tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.SendBatch failed",
			caller,
		)
	}

	return nil
}
