package outbox

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/outbox/dto"
	"github.com/martketplace-vkr/pkg/tracer"
)

// CreateBatchEvents is a request to create a new events batch.
type CreateBatchEvents []CreateEvent

func (o outbox) BatchCreateEvents(ctx context.Context, cmd CreateBatchEvents) ([]int64, error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	// if we want events to be sent immediately.
	if o.cfg.MessagesScheduler.ImmediateSendTimeout != nil && o.cfg.MessagesScheduler.ImmediateSendTimeout.Duration > 0 {
		return o.trySendBatchEventsImmediately(ctx, cmd)
	}

	return o.createEventsBatch(ctx, cmd)
}

func (o *outbox) createEventsBatch(ctx context.Context, cmd CreateBatchEvents) ([]int64, error) {
	totalEvents := 0
	for _, command := range cmd {
		totalEvents += len(command.Topics)
	}

	events := make([]dto.Event, 0, totalEvents)

	for _, command := range cmd {
		for _, topic := range command.Topics {
			events = append(events, prepareEvent(ctx, command, topic))
		}
	}

	ids, err := o.storage.BatchCreateOutboxEvents(ctx, events)
	if err != nil {
		return nil, err
	}

	metricsCli.observeEventsCreateResult(ids)

	return ids, nil
}

func (o outbox) trySendBatchEventsImmediately(ctx context.Context, cmd CreateBatchEvents) ([]int64, error) {
	events := make([]dto.Event, 0)
	now := time.Now().UTC()
	sendTimeout := o.cfg.MessagesScheduler.ImmediateSendTimeout.Duration

	jitter := time.Duration(0)
	if sendTimeout > 0 {
		n := time.Now().UnixNano()
		if n < 0 {
			n = -n
		}
		jitter = time.Duration(n % int64(jitterMax))
	}

	var nextAttemptAt time.Time
	if sendTimeout > 0 {
		nextAttemptAt = now.Add(sendTimeout + jitter)
	} else {
		nextAttemptAt = now
	}

	for _, command := range cmd {
		for _, topic := range command.Topics {
			e := prepareEvent(ctx, command, topic)
			e.LockedUntil = nextAttemptAt

			events = append(events, e)
		}
	}

	ids, err := o.storage.BatchCreateOutboxEvents(ctx, events)
	if err != nil {
		return nil, err
	}
	metricsCli.observeEventsCreateResult(ids)

	// if no immediate process orders
	if sendTimeout == 0 || len(ids) == 0 {
		return ids, nil
	}

	sent := make(dto.SuccessEvents, 0, len(ids))
	failedIDs := make([]int64, 0)

	for i, id := range ids {
		ev := events[i]
		ev.ID = id

		sendCtx, cancel := context.WithTimeout(ctx, sendTimeout)
		partition, offset, pubErr := o.broker.PublishSingle(sendCtx, ev)
		cancel()

		if pubErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		sent = append(sent, dto.SuccessEvent{
			ID:        id,
			Offset:    offset,
			Partition: partition,
			SentAt:    time.Now().UTC(),
		})

	}

	if len(sent) > 0 {
		if err := o.storage.MarkEventsAsProcessed(ctx, sent); err != nil {
		}
	}

	return ids, nil
}
