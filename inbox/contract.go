package inbox

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/inbox/dto"
)

type (
	Storage interface {
		EventsStorage
		ArchiveStorage
		EffectStorage
		LeaseExtender
	}

	LeaseExtender interface {
		// ExtendLocksByIDs updates lock timeout if event are processed longer, than we expected.
		ExtendLocksByIDs(ctx context.Context, ids []int64, extendBy time.Duration) error
	}
	EventsStorage interface {
		// CreateInboxEvent creates a new event in the inbox table.
		CreateInboxEvent(ctx context.Context, event dto.Event) (int64, error)
		// BatchCreateInboxEvents creates multiple events in the inbox table.
		BatchCreateInboxEvents(ctx context.Context, events dto.Events) ([]int64, error)
		// FetchUnprocessedEvents fetches unprocessed events from the inbox table.
		FetchUnprocessedEvents(ctx context.Context, batchSize int) (dto.Events, error)
		// MarkEventsAsProcessed marks events as processed in the inbox table.
		MarkEventsAsProcessed(ctx context.Context, ids dto.SuccessEvents) error
		// FetchExpiredEvents fetches expired events from the inbox table.
		FetchExpiredEvents(ctx context.Context, TTL, batchSize int) (dto.Events, error)
		// DeleteEventsByIDs deletes events by their IDs from the inbox table.
		DeleteEventsByIDs(ctx context.Context, ids []int64) error
		// UnlockFailedEvents unlocks failed events in the inbox table.
		UnlockFailedEvents(ctx context.Context, events dto.FailedEvents) error
		// UnlockSkippedEvents unlocks skipped events in the inbox table.
		UnlockSkippedEvents(ctx context.Context, events []int64) error
	}
	ArchiveStorage interface {
		// DeleteEffectLogOlderThan deletes expired effect logs.
		DeleteEffectLogOlderThan(ctx context.Context, olderThan time.Time, limit int) (int64, error)
		// ArchiveProcessedInbox archives old inbox events.
		ArchiveProcessedInbox(ctx context.Context, olderThan time.Time, limit int) (int64, error)
		// DeleteFromInboxDeletedOlderThan deletes archived inbox events.
		DeleteFromInboxDeletedOlderThan(ctx context.Context, olderThan time.Time, limit int) (int64, error)
	}
	EffectStorage interface {
		// RegisterEffectsBatch registers events in effect_log table.
		RegisterEffectsBatch(
			ctx context.Context,
			events []dto.Event,
		) (insertedIdx []int, conflictedIdx []int, err error)
	}
	// EventsProvider provides a channel with incoming events, that inbox must save
	EventsProvider interface {
		GetEventsChan() chan dto.Event
		Consume(ctx context.Context)
	}

	DistributedLocker interface {
		Lock(
			ctx context.Context,
			events dto.Events,
		) (cleanEvents dto.Events, dirtyEvents []int64, err error)
		Unlock(
			ctx context.Context,
			events dto.Events,
		) (err error)
	}
	DistributedLockerStorage interface {
		Clear(ctx context.Context, key string) error
		checkAndSet(ctx context.Context, key string) error
	}
	// EventsSchedulerFunc is a func that launches inbox worker with some period
	EventsSchedulerFunc func()
	// EventsHandlerFunc is a function that actually processes incoming events
	// if you want to specify your own logic of events process u may specify custom function
	// in options.
	EventsHandlerFunc func(ctx context.Context, events dto.Events) (
		successEvents dto.SuccessEvents,
		failedEventIDs dto.FailedEvents,
		err error,
	)
	skipEventFunc func(event dto.Event) bool
)
