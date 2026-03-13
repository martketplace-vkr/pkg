package collector

import (
	"context"

	"github.com/martketplace-vkr/pkg/kafkaconnector"
)

type (
	Collector interface {
		Collect(ctx context.Context)
		Cleanup(ctx context.Context)
	}
	KafkaCollector struct {
		cfg           Config
		consumerGroup kafkaconnector.ConsumerGroup
	}
)

func New(
	consumerGroup kafkaconnector.ConsumerGroup,
) Collector {
	return &KafkaCollector{
		consumerGroup: consumerGroup,
	}
}

func (kc *KafkaCollector) Collect(_ context.Context) {
	kc.consumerGroup.Consume(context.Background())
}

func (kc *KafkaCollector) Cleanup(_ context.Context) {
	kc.consumerGroup.Cleanup(context.Background())
}
