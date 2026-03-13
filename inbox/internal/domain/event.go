package domain

import (
	"encoding/json"
	"time"

	"github.com/guregu/null"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/martketplace-vkr/pkg/inbox/utils"
	"github.com/samber/lo"
)

type Event struct {
	ID             int64              `db:"id"`
	Topic          string             `db:"topic"`
	Key            null.String        `db:"key"`
	EventType      null.String        `db:"event_type"`
	Payload        json.RawMessage    `db:"payload"`
	TraceCarrier   utils.TraceCarrier `db:"trace_carrier"`
	IdempotencyKey string             `db:"idempotency_key"`
	CreatedAt      time.Time          `db:"created_at"`
	ProcessedAt    null.Time          `db:"processed_at"`
	Error          null.String        `db:"error"`
	Attempts       int32              `db:"attempts"`
	State          string             `db:"state"`
}

type Events []Event

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
	EventTypes      pq.StringArray
	CreatedAts      pgtype.TimestampArray
	EntityIDs       pq.StringArray
	IdempotencyKeys pq.StringArray
	Errors          pq.StringArray
	Attempts        pq.Int32Array
	States          pq.StringArray
}

func NewEventBatchFromCommand(e Events) EventBatch {
	batch := EventBatch{
		Payloads:        make(pq.StringArray, len(e)),
		TraceCarriers:   make(pq.StringArray, len(e)),
		Topics:          make(pq.StringArray, len(e)),
		Keys:            make(pq.StringArray, len(e)),
		EventTypes:      make(pq.StringArray, len(e)),
		EntityIDs:       make(pq.StringArray, len(e)),
		IdempotencyKeys: make(pq.StringArray, len(e)),
		Errors:          make(pq.StringArray, len(e)),
		Attempts:        make([]int32, len(e)),
		States:          make(pq.StringArray, len(e)),
		CreatedAts:      pgtype.TimestampArray{},
	}

	for i, event := range e {
		batch.Payloads[i] = string(event.Payload)
		batch.TraceCarriers[i] = event.TraceCarrier.String()
		batch.Topics[i] = event.Topic
		batch.IdempotencyKeys[i] = event.IdempotencyKey
		batch.Errors[i] = event.Error.String
		batch.Attempts[i] = event.Attempts
		batch.Keys[i] = event.Key.String
		batch.EventTypes[i] = event.EventType.String
		batch.States[i] = event.State
	}

	batch.CreatedAts.Elements = make([]pgtype.Timestamp, len(e))
	batch.CreatedAts.Dimensions = []pgtype.ArrayDimension{{Length: int32(len(e)), LowerBound: 1}}
	batch.CreatedAts.Status = pgtype.Present

	for i, event := range e {
		batch.CreatedAts.Elements[i] = pgtype.Timestamp{
			Time:   event.CreatedAt,
			Status: pgtype.Present,
		}
	}

	var timestampArray pgtype.TimestampArray

	timestampArray.Set(e.CreatedAts())

	batch.CreatedAts = timestampArray

	return batch
}
