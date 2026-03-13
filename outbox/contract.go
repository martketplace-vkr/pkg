package outbox

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/outbox/dto"
)

type (
	Storage interface {
		EventsStorage
		ArchiveStorage
	}

	EventsStorage interface {
		// CreateOutboxEvent creates a new event in the outbox table.
		CreateOutboxEvent(ctx context.Context, event dto.Event) (int64, error)
		// BatchCreateOutboxEvents creates multiple events in the outbox table.
		BatchCreateOutboxEvents(ctx context.Context, events dto.Events) ([]int64, error)
		// FetchUnprocessedEvents fetches unprocessed events from the outbox table.
		FetchUnprocessedEvents(ctx context.Context, batchSize, maxAttempts int) (dto.Events, error)
		// MarkEventsAsProcessed marks events as processed in the outbox table.
		MarkEventsAsProcessed(ctx context.Context, ids dto.SuccessEvents) error
		// UnlockEvents unlocks events in the outbox table.
		UnlockEvents(ctx context.Context, events dto.FailedEvents) error
		UnlockSkippedEvents(ctx context.Context, events []int64) error
		// RegisterOutboxSendingBatch writes outbox events to database.
		RegisterOutboxSendingBatch(
			ctx context.Context,
			events []dto.Event,
		) (insertedIdx []int, conflictedIdx []int, err error)
		ExtendLocksByIDs(
			ctx context.Context,
			ids []int64,
			extendBy time.Duration,
		) error
	}

	ArchiveStorage interface {
		ArchiveSentOutbox(ctx context.Context, olderThan time.Time, limit int) (int64, error)
		DeleteFromOutboxArchiveOlderThan(ctx context.Context, olderThan time.Time, limit int) (int64, error)
		DeleteOutboxSentOlderThan(ctx context.Context, olderThan time.Time, limit int) (int64, error)
	}

	Broker interface {
		Publish(ctx context.Context, events dto.Events) (success dto.SuccessEvents, fail dto.FailedEvents, err error)
		PublishSingle(ctx context.Context, event dto.Event) (partition int32, offset int64, err error)
	}
	DistributedLocker interface {
		Lock(ctx context.Context, events dto.Events) (cleanEvents dto.Events, dirtyEvents []int64, err error)
		Unlock(ctx context.Context, events dto.Events) error
	}
)
