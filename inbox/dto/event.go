package dto

import (
	"context"
	"encoding/json"
	"time"

	"github.com/guregu/null"
	"github.com/samber/lo"
)

const (
	StateNew       = "new"
	StateProcessed = "processing"
	StateDone      = "done"
	StateError     = "error"
)

type (
	Event struct {
		ID             int64
		IdempotencyKey string
		Topic          string
		Key            string
		EventType      string
		Payload        json.RawMessage
		CreatedAt      time.Time
		ProcessedAt    null.Time
		Context        context.Context
		Error          string
		Attempts       int
		State          string // new, pending, done
	}
	SuccessEvent struct {
		ID        int64 `db:"id"`
		Offset    int64 `db:"offset"`
		Partition int32 `db:"partition"`
	}
	Events        []Event
	SuccessEvents []int64
	FailedEvents  []FailedEvent
	FailedEvent   struct {
		ID    int64
		Error string
	}
)

func (e Events) IDs() []int64 {
	return lo.Map(e, func(event Event, _ int) int64 {
		return event.ID
	})
}

func (e FailedEvents) IDs() []int64 {
	return lo.Map(e, func(event FailedEvent, _ int) int64 {
		return event.ID
	})
}

func (e Events) IdempotencyKeys() []string {
	return lo.Map(e, func(event Event, _ int) string {
		return event.IdempotencyKey
	})
}

type (
	ProcessMode string
)

const (
	ProcessModeBatch          ProcessMode = "batch"
	ProcessModeSingleWIthLock ProcessMode = "single_with_lock"
)
