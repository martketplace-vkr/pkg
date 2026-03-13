package postgres

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/outbox/dto"

	"github.com/martketplace-vkr/pkg/outbox/internal/domain"

	txManager "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/martketplace-vkr/pkg/tracer"
)

type sqlxStorage struct {
	db        *sqlx.DB
	ctxGetter *txManager.CtxGetter

	eventLockTimeout time.Duration
	queryFetchEvents string
}

func NewSqlxStorage(
	db *sqlx.DB,
	ctxGetter *txManager.CtxGetter,
	eventLockTimeout time.Duration,
	disableSequentialEntriesSelection bool,
) *sqlxStorage {
	query := queryFetchUnprocessedEventsWithOrder
	if disableSequentialEntriesSelection {
		query = queryFetchUnprocessedEventsUnordered
	}

	return &sqlxStorage{
		db:               db,
		ctxGetter:        ctxGetter,
		eventLockTimeout: eventLockTimeout,
		queryFetchEvents: query,
	}
}

func (s *sqlxStorage) CreateOutboxEvent(ctx context.Context, eventDto dto.Event) (int64, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var id int64

	event := eventFromDto(eventDto)

	err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).GetContext(
		ctx,
		&id,
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
	)
	if err != nil {
		return id, tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.GetContext(event: %+v)",
			caller,
			event,
		)
	}

	return id, err
}

func (s *sqlxStorage) BatchCreateOutboxEvents(
	ctx context.Context,
	events dto.Events,
) ([]int64, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var ids []int64

	batch := domain.NewEventBatchFromCommand(eventsFromDto(events))

	err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).SelectContext(
		ctx,
		&ids,
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
			"%s.SelectContext(events: %+v)",
			caller,
			events,
		)
	}

	return ids, err
}

func (s *sqlxStorage) FetchUnprocessedEvents(ctx context.Context, batchSize, maxAttempts int) (dto.Events, error) {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	var events domain.Events

	err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).SelectContext(
		ctx,
		&events,
		s.queryFetchEvents,
		time.Now().Add(s.eventLockTimeout),
		batchSize,
		maxAttempts,
	)
	if err != nil {
		return nil, tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.SelectContext(batchSize: %d)",
			caller,
			batchSize,
		)
	}

	return eventsToDto(events), err
}

func (s *sqlxStorage) MarkEventsAsProcessed(ctx context.Context, events dto.SuccessEvents) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	stmt, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).PreparexContext(ctx, queryMarkEventAsPublished)
	if err != nil {
		return tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.PrepareContext()",
			caller,
		)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err = stmt.ExecContext(
			ctx,
			event.ID,
			event.Partition,
			event.Offset,
		)
		if err != nil {
			return tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.ExecContext(req: %s)",
				caller,
				event.ID,
			)
		}
	}

	return err
}

func (s *sqlxStorage) UnlockEvents(ctx context.Context, events dto.FailedEvents) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	stmt, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).Prepare(queryUnlockEvents)
	if err != nil {
		return tracer.SpanSetErrWrapf(
			span,
			err,
			"%s.Prepare(queryUnlockEvents) failed",
			caller,
		)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err := stmt.ExecContext(ctx, event.ID, event.Error)
		if err != nil {
			return tracer.SpanSetErrWrapf(
				span,
				err,
				"%s.ExecContext(event: %+v) failed",
				caller,
				event,
			)
		}
	}

	return nil
}

func (s *sqlxStorage) UnlockSkippedEvents(ctx context.Context, events []int64) error {
	ctx, span, caller := tracer.NewSpan(ctx)
	defer span.End()

	_, err := s.ctxGetter.DefaultTrOrDB(ctx, s.db).ExecContext(
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
