package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func (s *TestSuite) Test_CreateEvent() {
	id, err := s.outbox.CreateEvent(context.Background(), CreateEvent{
		EventType: "1",
		Payload:   json.RawMessage(`{"1": "2"}`),
		Topics:    []string{"transactional-outbox"},
		Key:       uuid.NewString(),
		Context:   context.Background(),
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	})
	s.Require().Nil(err)
	s.Require().NotNil(id)
}
