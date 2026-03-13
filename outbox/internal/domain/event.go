package domain

import (
	"encoding/json"
	"time"

	"github.com/martketplace-vkr/pkg/outbox/utils"

	"github.com/guregu/null"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type Event struct {
	ID             int64              `db:"id"`
	Topic          string             `db:"topic"`
	Key            null.String        `db:"key"`
	Payload        json.RawMessage    `db:"payload"`
	TraceCarrier   utils.TraceCarrier `db:"trace_carrier"`
	IdempotencyKey string             `db:"idempotency_key"`
	CreatedAt      time.Time          `db:"created_at"`
	EventType      string             `db:"event_type"`
	SentAt         null.Time          `db:"sent_at"`
	Attempts       int64              `db:"attempts"`
	LastError      null.String        `db:"last_error"`
	LockedUntil    *time.Time         `db:"locked_until"`
	State          string             `db:"state"`
}

type Events []Event

func (e Events) IDs() []int64 {
	return lo.Map(e, func(event Event, _ int) int64 {
		return event.ID
	})
}

func (e Events) CreatedAts() []time.Time {
	return lo.Map(e, func(event Event, _ int) time.Time {
		return event.CreatedAt
	})
}

type EventBatch struct {
	Payloads        pq.StringArray
	TraceCarriers   pq.StringArray
	Topics          pq.StringArray
	Keys            pq.StringArray
	CreatedAts      pgtype.TimestampArray
	IdempotencyKeys pq.StringArray
	Attempts        pq.Int64Array
	LastErrors      pq.StringArray
	EventTypes      pq.StringArray
	States          pq.StringArray
}

func NewEventBatchFromCommand(e Events) EventBatch {
	batch := EventBatch{
		Payloads:        make(pq.StringArray, len(e)),
		TraceCarriers:   make(pq.StringArray, len(e)),
		Topics:          make(pq.StringArray, len(e)),
		Keys:            make(pq.StringArray, len(e)),
		EventTypes:      make(pq.StringArray, len(e)),
		IdempotencyKeys: make(pq.StringArray, len(e)),
		Attempts:        make(pq.Int64Array, len(e)),
		LastErrors:      make(pq.StringArray, len(e)),
		States:          make(pq.StringArray, len(e)),
	}

	for i, event := range e {
		batch.Payloads[i] = string(event.Payload)
		batch.TraceCarriers[i] = event.TraceCarrier.String()
		batch.Topics[i] = event.Topic
		batch.EventTypes[i] = event.EventType

		if event.Key.Valid {
			batch.Keys[i] = event.Key.String
		}

		batch.IdempotencyKeys[i] = event.IdempotencyKey
		batch.Attempts[i] = event.Attempts
		batch.LastErrors[i] = event.LastError.String
		batch.States[i] = event.State
	}

	var timestampArray pgtype.TimestampArray

	err := timestampArray.Set(e.CreatedAts())
	if err != nil {
		log.Err(err)
	}

	batch.CreatedAts = timestampArray

	return batch
}

type (
	SuccessEvent struct {
		ID        int64 `db:"id"`
		Offset    int64 `db:"offset"`
		Partition int32 `db:"partition"`
	}

	EventResult struct {
		ID    int64 `db:"id"`
		Error error `db:"error"`
	}

	SuccessEvents []SuccessEvent
)
