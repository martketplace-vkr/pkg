package kafka_distributor

import (
	"github.com/martketplace-vkr/pkg/build/components"
	"github.com/martketplace-vkr/pkg/kafka-distributor/collector"
	"github.com/martketplace-vkr/pkg/kafka-distributor/distributor"
)

type Config struct {
	components.ComponentConfig
	Collector   collector.Config
	Distributor distributor.Config
}
