package distributor

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"

	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	Distributor interface {
		Distribute(ctx context.Context)
		Cleanup(ctx context.Context)
	}
	DistributeFunc func(ctx context.Context, event *EventMsg) (err error)
)

type KafkaDistributor struct {
	cfg               Config
	cancel            context.CancelFunc
	consumedMessages  chan *sarama.ConsumerMessage
	processedMessages chan struct{}
	eventMethodMap    map[string]DistributeFunc
}

func New(
	cfg Config,
	eventMethodMap map[string]DistributeFunc,
	consumedMessages chan *sarama.ConsumerMessage,
	processedMessages chan struct{},
) *KafkaDistributor {
	return &KafkaDistributor{
		cfg:               cfg,
		eventMethodMap:    eventMethodMap,
		consumedMessages:  consumedMessages,
		processedMessages: processedMessages,
	}
}

func (d *KafkaDistributor) Distribute(
	_ context.Context,
) {
	wg := &sync.WaitGroup{}
	wg.Add(d.cfg.WorkersCount)

	for i := 0; i < d.cfg.WorkersCount; i++ {
		go func(consumeChan <-chan *sarama.ConsumerMessage, eventsHandler map[string]DistributeFunc) {
			defer wg.Done()
			for event := range consumeChan {
				event := event
				log.Infof("consumed message on topic: %s | offset: %d", event.Topic, event.Offset)

				method, ok := eventsHandler[string(event.Key)]
				if !ok {
					log.Infof("Unknown event type: %v | consumed message: %s | topic: %s",
						string(event.Key),
						string(event.Value),
						event.Topic,
					)

					continue
				}

				ctx := otel.GetTextMapPropagator().Extract(context.Background(),
					otelsarama.NewConsumerMessageCarrier(event))

				err := method(ctx, &EventMsg{event})
				if err != nil {
					log.Error(err)
				}
			}
		}(d.consumedMessages, d.eventMethodMap)
	}
	wg.Wait()
}

func (d *KafkaDistributor) Cleanup(_ context.Context) {
	close(d.consumedMessages)
	close(d.processedMessages)
}
