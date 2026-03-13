package kafka

import (
	"context"

	"github.com/IBM/sarama"
	eventtype "github.com/martketplace-vkr/pkg/eventtypes"
	"github.com/martketplace-vkr/pkg/outbox/dto"

	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/martketplace-vkr/pkg/tracer"
)

type kafkaBroker struct {
	producer kafkaconnector.SyncProducer
}

func NewSync(
	producer kafkaconnector.SyncProducer,
) *kafkaBroker {
	return &kafkaBroker{
		producer: producer,
	}
}

func (b *kafkaBroker) Publish(ctx context.Context, events dto.Events) (success dto.SuccessEvents, fail dto.FailedEvents, err error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	for _, event := range events {
		func(event dto.Event) {
			// todo тут пренебрегаем контекстом ивента, чтобы повесить таймаут на запрос.
			localCtx, msgSpan, _ := tracer.NewSpan(ctx, "Broker.Kafka.PublishEvents")
			defer msgSpan.End()

			partition, offset, err := b.producer.SendMessage(localCtx, kafkaconnector.Message{
				Topic: event.Topic,
				Key:   []byte(event.Key),
				Value: event.Payload,
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte(eventtype.EventsTypeKafkaKey),
						Value: []byte(event.EventType),
					},
					{
						Key:   []byte(eventtype.IdempotencyKeyKafkaKey),
						Value: []byte(event.IdempotencyKey),
					},
				},
			})
			if err != nil {
				fail = append(fail, dto.FailedEvent{
					ID:    event.ID,
					Error: err.Error(),
				})

				return
			}

			success = append(success, dto.SuccessEvent{
				ID:        event.ID,
				Offset:    offset,
				Partition: partition,
			})
		}(event)
	}

	return success, fail, nil
}

func (b *kafkaBroker) PublishSingle(ctx context.Context, event dto.Event) (partition int32, offset int64, err error) {
	ctx, span, _ := tracer.NewSpan(ctx)
	defer span.End()

	return b.producer.SendMessage(ctx, kafkaconnector.Message{
		Topic: event.Topic,
		Key:   []byte(event.Key),
		Value: event.Payload,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(eventtype.EventsTypeKafkaKey),
				Value: []byte(event.EventType),
			},
			{
				Key:   []byte(eventtype.IdempotencyKeyKafkaKey),
				Value: []byte(event.IdempotencyKey),
			},
		},
	})
}
