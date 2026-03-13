package outbox

import (
	"context"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *TestSuite) Test_messageScheduler() {
	_, err := s.outbox.CreateEvent(context.Background(), CreateEvent{
		EventType: "1",
		Payload:   []byte(`{"1": "2"}`),
		Topics:    []string{"transactional-outbox"},
		Key:       "1",
		Context:   context.Background(),
		CreatedAt: time.Now().AddDate(-1, 0, 0),
	})
	assert.Equal(s.T(), nil, err)

	s.outbox.scheduleBroadcast()

	time.Sleep(10 * time.Second)
}
