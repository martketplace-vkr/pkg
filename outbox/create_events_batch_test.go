package outbox

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

func (s *TestSuite) Test_BatchCreateEvent() {
	ids, err := s.outbox.BatchCreateEvents(context.Background(), CreateBatchEvents{
		CreateEvent{
			EventType: "2",
			Payload:   json.RawMessage(`{"3": "4"}`),
			Topics:    []string{"transactional-outbox"},
			Key:       uuid.NewString(),
			Context:   context.Background(),
			CreatedAt: time.Now().AddDate(-1, 0, 0),
		},
		CreateEvent{
			EventType: "3",
			Payload:   json.RawMessage(`{"5": "6"}`),
			Topics:    []string{"transactional-outbox"},
			Key:       uuid.NewString(),
			Context:   context.Background(),
		},
	})
	s.Require().Nil(err)
	s.Require().NotNil(ids)
}
