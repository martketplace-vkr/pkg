package kafkaconnector

import (
	"context"
	"errors"
	
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type (
	EventsHandler struct {
		eventsKafka chan *sarama.ConsumerMessage
	}

	CustomEventsHandler struct {
		processEvent func(context.Context, *sarama.ConsumerMessage) error
	}
)

func NewDefaultHandler(
	eventsKafka chan *sarama.ConsumerMessage,
) *EventsHandler {
	return &EventsHandler{eventsKafka: eventsKafka}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (kc *EventsHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (kc *EventsHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	session.Commit()

	return nil
}

// ConsumeClaim must start a consumerGroup loop of ConsumerGroupClaim's Messages().
func (kc *EventsHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	for {
		select {
		case message := <-claim.Messages():
			session.MarkMessage(message, "")

			kc.eventsKafka <- message
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewCustomHandler(processEvent func(context.Context, *sarama.ConsumerMessage) error) *CustomEventsHandler {
	return &CustomEventsHandler{processEvent: processEvent}
}

func (kc *CustomEventsHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (kc *CustomEventsHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	session.Commit()

	return nil
}

func (kc *CustomEventsHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	for message := range claim.Messages() {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(),
			otelsarama.NewConsumerMessageCarrier(message))

		if err := kc.processEvent(ctx, message); err != nil {
			tr := otel.GetTracerProvider().Tracer("error handler")
			_, span := tr.Start(ctx, message.Topic)
			span.RecordError(errors.New(err.Error()))
			span.SetAttributes(attribute.String("value", string(message.Value)))
			span.SetStatus(codes.Error, err.Error())
			span.End()

			session.ResetOffset(message.Topic, message.Partition, message.Offset, "")

			return err
		}

		session.MarkMessage(message, "")
	}

	return err
}
