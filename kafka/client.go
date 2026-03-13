package kafka

import (
	"crypto/tls"
	"time"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/logger/log"
)

const (
	autocommitInterval = 500 * time.Millisecond
)

var (
	SaramaPanicHandler func(i interface{})
)

type (
	Client interface {
		// Consumer You must call Close() on sarama.Consumer
		Consumer(options ...ConfigOption) sarama.Consumer
		// ConsumerGroup You must call ConsumerGroup.Close()
		ConsumerGroup(groupID string, options ...ConfigOption) ConsumerGroup
		// SyncProducer You must call SyncProducer.Close()
		SyncProducer(topic string, options ...ConfigOption) SyncProducer
		// AsyncProducer You must call AsyncProducer.Close()
		AsyncProducer(topic string, options ...ConfigOption) AsyncProducer
	}

	ConfigOption = func(config *sarama.Config)

	ClientImpl struct {
		sarama.Client
		brokers []string
	}
)

func NewClient(cfg Config, options ...ConfigOption) *ClientImpl {
	panicHandler := SaramaPanicHandler
	if panicHandler == nil {
		panicHandler = func(i interface{}) {
			log.Errorf("sarama panic: %v", i)
		}
	}

	if sarama.PanicHandler == nil {
		sarama.PanicHandler = panicHandler
	}

	saramaCfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		log.Fatalf("failed to parse kafka version %s: %e", cfg.Version, err)
	}

	saramaCfg.Version = version

	saramaCfg.Consumer.Offsets.AutoCommit.Enable = true
	saramaCfg.Consumer.Offsets.AutoCommit.Interval = autocommitInterval

	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	if cfg.OffsetNewest {
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	if cfg.SASL != nil {
		sasl := cfg.SASL
		saramaCfg.Net.SASL.Enable = true
		saramaCfg.Net.SASL.Handshake = true
		saramaCfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		saramaCfg.Net.SASL.User = sasl.User
		saramaCfg.Net.SASL.Password = sasl.Password
		saramaCfg.Net.TLS.Enable = true
		saramaCfg.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	for _, option := range options {
		option(saramaCfg)
	}

	saramaCli, err := sarama.NewClient(cfg.Brokers, saramaCfg)
	if err != nil {
		log.Fatalf("failed to create kafka client: %e", err)
	}

	return &ClientImpl{
		brokers: cfg.Brokers,
		Client:  saramaCli,
	}
}

func (c *ClientImpl) Consumer(options ...ConfigOption) sarama.Consumer {
	consumer, err := sarama.NewConsumer(c.brokers, withOptions(*c.Client.Config(), options...))
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %e", err)
	}

	return consumer
}

func (c *ClientImpl) ConsumerGroup(groupID string, options ...ConfigOption) ConsumerGroup {
	return NewConsumerGroup(c.brokers, withOptions(*c.Client.Config(), options...), groupID)
}

func (c *ClientImpl) SyncProducer(topic string, options ...ConfigOption) SyncProducer {
	syncProducer, err := NewSyncProducer(c.brokers, withOptions(*c.Client.Config(), options...), topic)
	if err != nil {
		log.Fatalf("failed to create kafka sync producer: %e", err)
	}

	return syncProducer
}

func (c *ClientImpl) AsyncProducer(topic string, options ...ConfigOption) AsyncProducer {
	asyncProducer, err := NewAsyncProducer(c.brokers, withOptions(*c.Client.Config(), options...), topic)
	if err != nil {
		log.Fatalf("failed to create kafka async producer: %e", err)
	}

	return asyncProducer
}

func withOptions(config sarama.Config, options ...ConfigOption) *sarama.Config {
	for _, option := range options {
		option(&config)
	}

	return &config
}
