package outbox

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/martketplace-vkr/pkg/outbox/dto"

	"github.com/martketplace-vkr/pkg/tracer"
)

// CreateEvent is a request to create a new event.
type CreateEvent struct {
	// EventType is an identifier of entity (for example order_id)
	// in order to group events by entity.
	EventType string
	// Key is a Broker key for message routing.
	Key string
	// Payload is a message body.
	Payload []byte
	// CreatedAt is a time when event was created.
	CreatedAt time.Time
	// Topic is a message topic.
	Topics []string
	// Context is a context for tracing (ctx will be used if not present).
	Context context.Context
	// LockedUntil - lock time
	LockedUntil time.Time
}

const jitterMax = 200 * time.Millisecond

func (o outbox) CreateEvent(ctx context.Context, cmd CreateEvent) ([]int64, error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	// if we want events to be sent immediately.
	if o.cfg.MessagesScheduler.ImmediateSendTimeout != nil && o.cfg.MessagesScheduler.ImmediateSendTimeout.Duration > 0 {
		return o.trySendEventImmediately(ctx, cmd)
	}

	return o.createEvent(ctx, cmd)
}

func (o *outbox) createEvent(ctx context.Context, cmd CreateEvent) ([]int64, error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	events := make([]dto.Event, 0, len(cmd.Topics))

	for _, topic := range cmd.Topics {
		events = append(events, prepareEvent(ctx, cmd, topic))
	}

	ids, err := o.storage.BatchCreateOutboxEvents(ctx, events)
	if err != nil {
		return nil, err
	}

	metricsCli.observeEventsCreateResult(ids)

	return ids, nil
}

func (o *outbox) trySendEventImmediately(ctx context.Context, cmd CreateEvent) ([]int64, error) {
	events := make([]dto.Event, 0, len(cmd.Topics))
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

	for _, topic := range cmd.Topics {
		e := prepareEvent(ctx, cmd, topic)
		e.LockedUntil = nextAttemptAt

		events = append(events, e)
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

func prepareEvent(ctx context.Context, cmd CreateEvent, topic string) dto.Event {
	event := dto.Event{
		Key:            cmd.Key,
		EventType:      cmd.EventType,
		IdempotencyKey: uuid.NewString(),
		Topic:          topic,
		Payload:        cmd.Payload,
		CreatedAt:      cmd.CreatedAt,
		LockedUntil:    cmd.LockedUntil,
		State:          dto.StateNew,
	}

	if cmd.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	if cmd.Context == nil {
		cmd.Context = ctx
	}
	if event.State == "" {
		event.State = dto.StateNew
	}

	return event
}
