package postgres

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/outbox/dto"
	"github.com/martketplace-vkr/pkg/outbox/internal/domain"
	"github.com/martketplace-vkr/pkg/outbox/utils"

	"github.com/guregu/null"
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
		ID:             event.ID,
		Topic:          event.Topic,
		Key:            null.NewString(event.Key, event.Key != ""),
		Payload:        event.Payload,
		IdempotencyKey: event.IdempotencyKey,
		CreatedAt:      event.CreatedAt,
		EventType:      event.EventType,
		SentAt:         null.TimeFrom(event.SentAt),
		Attempts:       event.Attempts,
		LastError:      null.StringFrom(event.LastError),
		LockedUntil:    &event.LockedUntil,
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
		Key:            event.Key.String,
		EventType:      event.EventType,
		IdempotencyKey: event.IdempotencyKey,
		Topic:          event.Topic,
		Payload:        event.Payload,
		CreatedAt:      event.CreatedAt,
		SentAt:         event.SentAt.Time,
		Attempts:       event.Attempts,
		LastError:      event.LastError.String,
		State:          event.State,
		Context:        ctx,
	}
}

func eventsToDto(events domain.Events) dto.Events {
	res := make(dto.Events, len(events))
	for i, event := range events {
		res[i] = eventToDto(event)
	}

	return res
}
