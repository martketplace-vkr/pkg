package postgres

import (
	"context"
	"time"

	"github.com/guregu/null"
	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/inbox/internal/domain"
	"github.com/martketplace-vkr/pkg/inbox/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func eventFromDto(event dto.Event) domain.Event {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	if event.Context == nil {
		event.Context = context.Background()
	}

	resEvent := domain.Event{
		Topic:          event.Topic,
		Key:            null.NewString(event.Key, event.Key != ""),
		EventType:      null.NewString(event.EventType, event.EventType != ""),
		Payload:        event.Payload,
		IdempotencyKey: event.IdempotencyKey,
		CreatedAt:      event.CreatedAt,
		ProcessedAt:    event.ProcessedAt,
		State:          event.State,
	}

	carrier := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(event.Context, carrier)
	resEvent.TraceCarrier = utils.NewTraceCarrierFromTraceCarrier(carrier)

	return resEvent
}

func eventsFromDto(events dto.Events) domain.Events {
	res := make(domain.Events, len(events))
	for i, event := range events {
		res[i] = eventFromDto(event)
	}

	return res
}

func eventToDto(event domain.Event) dto.Event {
	ctx := context.Background()
	otel.GetTextMapPropagator().Extract(ctx, event.TraceCarrier)

	return dto.Event{
		ID:             event.ID,
		IdempotencyKey: event.IdempotencyKey,
		Topic:          event.Topic,
		Key:            event.Key.String,
		EventType:      event.EventType.String,
		Payload:        event.Payload,
		CreatedAt:      event.CreatedAt,
		ProcessedAt:    event.ProcessedAt,
		Context:        ctx,
		State:          event.State,
	}
}

func eventsToDto(events domain.Events) dto.Events {
	res := make(dto.Events, len(events))
	for i, event := range events {
		res[i] = eventToDto(event)
	}

	return res
}
