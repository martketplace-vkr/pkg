package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func (s *TestSuite) TestDefaultOptions_SqlxDB() {
	outbox, err := NewDefaultWithOptions(
		s.cfg,
		WithSqlxDB(s.dbSqlx),
		WithKafkaProducer(s.kafka.NewSyncProducer()),
	)

	s.Require().Nil(err)
	s.Require().NotNil(outbox)

	id, err := outbox.CreateEvent(context.Background(), CreateEvent{
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

func (s *TestSuite) TestDefaultOptions_PgxPoolDB() {
	outbox, err := NewDefaultWithOptions(
		s.cfg,
		WithPgxPoolDB(s.dbPgx),
		WithKafkaProducer(s.kafka.NewSyncProducer()),
	)

	s.Require().Nil(err)
	s.Require().NotNil(outbox)

	id, err := outbox.CreateEvent(context.Background(), CreateEvent{
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

func (s *TestSuite) TestDefaultOptions_Error() {
	outbox, err := NewDefaultWithOptions(
		s.cfg,
		WithKafkaProducer(s.kafka.NewSyncProducer()),
	)

	s.Require().Nil(outbox)
	s.Require().NotNil(err)
}
