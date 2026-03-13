package kafkaconnector

import (
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/IBM/sarama"

	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	Client interface {
		NewConsumerGroup(groupID string, topics []string, eventsHandler sarama.ConsumerGroupHandler) *kafkaConsumerGroup
		NewSyncProducer() *syncProducer
		NewAsyncProducer() *asyncProducer

		NewSaramaConsumerGroup(groupID string) (sarama.ConsumerGroup, error)
	}

	ClientConfig struct {
		Brokers            []string `validate:"required"`
		SASL               *SASL
		InsecureSkipVerify bool
		Producer           *ProducerConfig
		Consumer           *ConsumerConfig
	}

	kafkaClient struct {
		sarama.Client
		brokers []string
	}
)

func NewClient(cfg ClientConfig) *kafkaClient {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	if cfg.SASL.Username != "" && cfg.SASL.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.SASL.Username
		config.Net.SASL.Password = cfg.SASL.Password
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: sha512.New}
		}
	}

	if cfg.SASL.CaPath != nil {
		caCert, err := os.ReadFile(*cfg.SASL.CaPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{
			RootCAs:            caCertPool,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		}
	}

	if cfg.Producer != nil {
		// The total number of times to retry sending a message (default 3).
		config.Producer.Retry.Max = cfg.Producer.RetryMax
		// How long to wait for the cluster to settle between retries (default 100ms).
		config.Producer.Retry.Backoff = time.Millisecond * 250
		// idempotent syncProducer has a unique syncProducer ID and uses sequence IDs for each message,
		// allowing the broker to ensure, on a per-partition basis, that it is committing ordered messages with no duplication.
		config.Producer.Idempotent = cfg.Producer.Idempotent.Mode
		if cfg.Producer.Idempotent.Mode {
			config.Producer.Retry.Max = cfg.Producer.Idempotent.RetryMax
			config.Net.MaxOpenRequests = cfg.Producer.Idempotent.MaxOpenRequests
			config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.Producer.RequireAcks)
		}
		//  Successfully delivered messages will be returned on the Successes channel
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}

	if cfg.Consumer != nil {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
		if cfg.Consumer.OffsetNewest {
			config.Consumer.Offsets.Initial = sarama.OffsetNewest
		}

		switch cfg.Consumer.Assignor {
		case "sticky":
			config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
		case "round-robin":
			config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		case "range":
			config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
		default:
			log.Panicf("Unrecognized collector group partition assignor: %s", cfg.Consumer.Assignor)
		}

		config.Consumer.Return.Errors = true

		config.Consumer.Offsets.AutoCommit.Enable = cfg.Consumer.AutoCommit
		config.Consumer.Offsets.AutoCommit.Interval = cfg.Consumer.AutoCommitInterval.Duration
	}

	saramaCli, err := sarama.NewClient(cfg.Brokers, config)
	if err != nil {
		log.Fatalf("failed to create kafka client: %e", err)
	}

	return &kafkaClient{
		brokers: cfg.Brokers,
		Client:  saramaCli,
	}
}

func (k kafkaClient) NewConsumerGroup(groupID string, topics []string, eventsHandler sarama.ConsumerGroupHandler) *kafkaConsumerGroup {
	return newConsumerGroup(k.brokers, groupID, topics, k.Config(), eventsHandler)
}

func (k kafkaClient) NewSyncProducer() *syncProducer {
	producer, err := newSyncProducer(k.brokers, k.Config())
	if err != nil {
		log.Fatalf("failed to crete sync producer: %s", err.Error())
	}

	return producer
}

func (k kafkaClient) NewAsyncProducer() *asyncProducer {
	producer, err := newAsyncProducer(k.brokers, k.Config())
	if err != nil {
		log.Fatalf("failed to crete sync producer: %s", err.Error())
	}

	return producer
}

func (k kafkaClient) NewSaramaConsumerGroup(groupID string) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(k.brokers, groupID, k.Config())
}
