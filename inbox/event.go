package inbox

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/tracer"
	"github.com/samber/lo"
)

// CreateEvent is a request to create a new event.
type CreateEvent struct {
	// Key is a Broker key for message routing.
	Key string
	//
	EventType string
	// Payload is a message body.
	Payload []byte
	// CreatedAt is a time when event was created.
	CreatedAt time.Time
	// Topic is a message topic.
	Topic string
	// Context is a context for tracing (ctx will be used if not present).
	Context context.Context
	// IdempotenceKey is a key for identification and control of repeated events
	IdempotencyKey string
}

func (i inbox) CreateEvent(ctx context.Context, event dto.Event) (int64, error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	id, err := i.storage.CreateInboxEvent(ctx, prepareEvent(ctx, event))
	if err != nil {
		return 0, err
	}

	metricsCli.addCreatedEventsTotal(event.EventType, 1)

	return id, nil
}

// CreateBatchEvents is a request to create a new events batch.
type CreateBatchEvents []CreateEvent

func (i inbox) BatchCreateEvents(ctx context.Context, events dto.Events) ([]int64, error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	ids, err := i.storage.BatchCreateInboxEvents(ctx, lo.Map(events, func(event dto.Event, _ int) dto.Event {
		return prepareEvent(ctx, event)
	}))
	if err != nil {
		return nil, err
	}

	metricsCli.observeEventsProcessResult(events, ids)

	return ids, err
}

func prepareEvent(ctx context.Context, event dto.Event) dto.Event {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	if event.Context == nil {
		event.Context = ctx
	}
	if event.IdempotencyKey == "" {
		event.IdempotencyKey = uuid.New().String()
	}
	if event.State == "" {
		event.State = dto.StateNew
	}

	return event
}
