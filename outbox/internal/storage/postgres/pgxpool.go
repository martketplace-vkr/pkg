package postgres

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/outbox/dto"

	"github.com/martketplace-vkr/pkg/outbox/internal/domain"

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
	queryFetchEvents string
}

func NewPgxPoolStorage(
	db *pgxpool.Pool,
	ctxGetter *txManager.CtxGetter,
	eventLockTimeout time.Duration,
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
		queryFetchEvents: query,
	}
}

func (s *pgxPoolStorage) CreateOutboxEvent(ctx context.Context, dtoEvent dto.Event) (int64, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var id int64

	event := eventFromDto(dtoEvent)

	err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).QueryRow(
		ctx,
		queryCreateEvent,
		event.Topic,
		event.Payload,
		event.TraceCarrier.String(),
		event.Key,
		event.CreatedAt,
		event.IdempotencyKey,
		event.Attempts,
		event.EventType,
		event.LastError,
		event.LockedUntil,
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

func (s *pgxPoolStorage) BatchCreateOutboxEvents(
	ctx context.Context,
	events dto.Events,
) ([]int64, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var ids []int64

	batch := domain.NewEventBatchFromCommand(eventsFromDto(events))

	rows, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		queryInsertBatch,
		batch.Topics,
		batch.Payloads,
		batch.TraceCarriers,
		batch.Keys,
		batch.CreatedAts,
		batch.EventTypes,
		batch.IdempotencyKeys,
		batch.Attempts,
		batch.LastErrors,
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

func (s *pgxPoolStorage) FetchUnprocessedEvents(ctx context.Context, batchSize, maxAttempts int) (dto.Events, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var events domain.Events

	rows, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Query(
		ctx,
		s.queryFetchEvents,
		time.Now().Add(s.eventLockTimeout),
		batchSize,
		maxAttempts,
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
			&event.SentAt,
			&event.Key,
			&event.IdempotencyKey,
			&event.Attempts,
			&event.LastError,
			&event.EventType,
		); err != nil {
			return nil, tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.Scan()",
				caller,
			)
		}
		events = append(events, event)
	}

	return eventsToDto(events), err
}

func (s *pgxPoolStorage) MarkEventsAsProcessed(ctx context.Context, events dto.SuccessEvents) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	batch := &pgx.Batch{}
	for _, event := range events {
		batch.Queue(queryMarkEventAsPublished, event.ID, event.Partition, event.Offset)
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
				events[i].ID,
			)
		}
	}

	return nil
}

func (s *pgxPoolStorage) UnlockEvents(ctx context.Context, events dto.FailedEvents) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	batch := &pgx.Batch{}

	for _, event := range events {
		batch.Queue(queryUnlockEvents, event.ID, event.Error)
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
