package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/google/uuid"
	eventtype "github.com/martketplace-vkr/pkg/eventtypes"
	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/logger/log"
	"go.opentelemetry.io/otel"
)

type groupHandler struct {
	eventsChan  chan dto.Event
	metrics     *kafkaMetrics
	debugLog    bool
	maxLogValue int
	workerID    int
}

func (*groupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (*groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (gh *groupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	if gh.debugLog {
		log.Debugf("kafka(worker=%d): assigned partition=%d topic=%s", gh.workerID, claim.Partition(), claim.Topic())
	}

	if gh.metrics != nil {
		gh.metrics.setPartitionAssigned(gh.workerID, claim.Topic(), claim.Partition())
		defer gh.metrics.unsetPartition(gh.workerID, claim.Topic(), claim.Partition())
	}

	for msg := range claim.Messages() {
		if gh.metrics != nil {
			gh.metrics.recordHeartbeat(gh.workerID)
			gh.metrics.incKafkaMessagesIn()
			gh.metrics.incPartitionMessage(msg.Topic, msg.Partition)
		}

		ctx := otel.GetTextMapPropagator().Extract(
			context.Background(),
			otelsarama.NewConsumerMessageCarrier(msg),
		)

		eventType := string(msg.Key)
		idempotencyKey := uuid.NewString()
		for _, h := range msg.Headers {
			switch string(h.Key) {
			case eventtype.EventsTypeKafkaKey:
				eventType = string(h.Value)
			case eventtype.IdempotencyKeyKafkaKey:
				idempotencyKey = string(h.Value)
			}
		}

		evt := dto.Event{
			Topic:          msg.Topic,
			Key:            string(msg.Key),
			EventType:      eventType,
			Payload:        msg.Value,
			Context:        ctx,
			IdempotencyKey: idempotencyKey,
			CreatedAt:      time.Now(),
		}

		if gh.metrics != nil {
			gh.metrics.incKafkaMessagesIn()
		}

		if gh.debugLog {
			valInfo := ""
			if gh.maxLogValue > 0 && len(msg.Value) > 0 {
				lim := gh.maxLogValue
				if lim > len(msg.Value) {
					lim = len(msg.Value)
				}
				valInfo = fmt.Sprintf(", valueSnippet=%q (len=%d)", string(msg.Value[:lim]), len(msg.Value))
			}
			if gh.debugLog {
				log.Debugf("kafka(worker=%d): new event: topic=%s partition=%d offset=%d key=%q eventType=%q idem=%q headers=%d%s",
					gh.workerID, msg.Topic, msg.Partition, msg.Offset, string(msg.Key), eventType, idempotencyKey, len(msg.Headers), valInfo)
			}
		}

		select {
		case gh.eventsChan <- evt:
		case <-session.Context().Done():
			return session.Context().Err()
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
