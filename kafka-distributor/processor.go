package kafka_distributor

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	"github.com/martketplace-vkr/pkg/kafka-distributor/collector"
	"github.com/martketplace-vkr/pkg/kafka-distributor/distributor"
	"github.com/martketplace-vkr/pkg/kafkaconnector"
)

const (
	cmpName = "kafka-distributor"
)

type EventProcessor struct {
	cfg         Config
	collector   collector.Collector
	distributor distributor.Distributor
}

func NewDefaultProcessor(
	cfg Config,
	kafkaClient kafkaconnector.Client,
	eventProcessorMap map[string]distributor.DistributeFunc,
) *EventProcessor {
	processedEvents := make(chan struct{}, 100)
	eventsKafka := make(chan *sarama.ConsumerMessage, 100)

	handler := kafkaconnector.NewDefaultHandler(eventsKafka)
	consumerGroup := kafkaClient.NewConsumerGroup(
		cfg.Collector.GroupID,
		cfg.Collector.Topics,
		handler,
	)

	eventCollector := collector.New(consumerGroup)
	eventDistributor := distributor.New(
		cfg.Distributor,
		eventProcessorMap,
		eventsKafka,
		processedEvents,
	)

	return &EventProcessor{
		cfg:         cfg,
		collector:   eventCollector,
		distributor: eventDistributor,
	}
}

func NewCustomProcessor(
	cfg Config,
	eventCollector collector.Collector,
	eventDistributor distributor.Distributor,
) *EventProcessor {

	return &EventProcessor{
		cfg:         cfg,
		collector:   eventCollector,
		distributor: eventDistributor,
	}
}

func (ep *EventProcessor) Start(ctx context.Context) (err error) {

	go ep.collector.Collect(
		context.Background(),
	)
	go ep.distributor.Distribute(
		context.Background(),
	)

	return nil
}

func (ep *EventProcessor) Stop(ctx context.Context) (err error) {
	ep.collector.Cleanup(ctx)
	ep.distributor.Cleanup(ctx)

	return err
}

func (ep *EventProcessor) GetStartTimeout() time.Duration {
	return ep.cfg.StartTimeout.Duration
}

func (ep *EventProcessor) GetStopTimeout() time.Duration {
	return ep.cfg.StopTimeout.Duration
}

func (ep *EventProcessor) GetShutdownDelay() time.Duration {
	return ep.cfg.ShutdownDelay.Duration
}

func (ep *EventProcessor) GetName() string {
	return cmpName
}
