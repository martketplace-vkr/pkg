package dto

import (
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"time"
)

const (
	StateNew        = "new"
	StateProcessing = "processing"
	StateDone       = "done"
	StateError      = "done"
)

type Event struct {
	ID             int64
	Key            string
	EventType      string
	IdempotencyKey string
	Topic          string
	Payload        json.RawMessage
	CreatedAt      time.Time
	SentAt         time.Time
	Attempts       int64
	LastError      string
	Context        context.Context
	LockedUntil    time.Time
	Partition      int32
	Offset         int64
	State          string // new, pending, done
}

type Events []Event

func (e Events) IDs() []int64 {
	return lo.Map(e, func(event Event, _ int) int64 {
		return event.ID
	})
}

type (
	FailedEvent struct {
		ID    int64
		Error string
	}
	FailedEvents []FailedEvent

	SuccessEvent struct {
		ID        int64
		Offset    int64
		Partition int32
		SentAt    time.Time
	}
	EventResult struct {
		ID    int64
		Error error
	}

	SuccessEvents []SuccessEvent
)

func (e FailedEvents) IDs() []int64 {
	return lo.Map(e, func(event FailedEvent, _ int) int64 {
		return event.ID
	})
}

func (e SuccessEvents) IDs() []int64 {
	return lo.Map(e, func(event SuccessEvent, _ int) int64 {
		return event.ID
	})
}
