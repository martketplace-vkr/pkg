package distributor

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/martketplace-vkr/pkg/utils/duration"
	"github.com/stretchr/testify/suite"
)

type DistributorSuite struct {
	suite.Suite
	distributor       *KafkaDistributor
	producer          kafkaconnector.SyncProducer
	ctx               context.Context
	cancel            context.CancelFunc
	consumer          kafkaconnector.ConsumerGroup
	processedMessages chan struct{}
	consumedMessages  chan *sarama.ConsumerMessage
	workerCount       int
}

func TestKafkaDistributor(t *testing.T) {
	suite.Run(t, new(DistributorSuite))
}

func (suite *DistributorSuite) SetupSuite() {
	client := kafkaconnector.NewClient(kafkaconnector.ClientConfig{
		Brokers: []string{"localhost:9092"},
		SASL: &kafkaconnector.SASL{
			Username: "",
			Password: "",
			CaPath:   nil,
		},
		InsecureSkipVerify: false,
		Producer: &kafkaconnector.ProducerConfig{
			ReadTimeout:  duration.Seconds{},
			WriteTimeout: duration.Seconds{},
			RequireAcks:  1,
			MaxAttempts:  1,
			Compression:  1,
			RetryMax:     4,
			Idempotent: struct {
				Mode            bool
				MaxOpenRequests int
				RetryMax        int
			}{false, 1, 1},
		},
		Consumer: &kafkaconnector.ConsumerConfig{
			Addresses:          []string{"localhost:9092"},
			Assignor:           "round-robin",
			OffsetNewest:       true,
			AutoCommit:         true,
			AutoCommitInterval: duration.Seconds{Duration: 5},
		},
	})
	suite.producer = client.NewSyncProducer()

	suite.workerCount = 3

	cfg := Config{suite.workerCount}
	eventMethodMap := map[string]DistributeFunc{
		"key-1": func(ctx context.Context, event *EventMsg) (err error) {
			suite.Equal("val-1", string(event.Value))
			return nil
		},
		"key-2": func(ctx context.Context, event *EventMsg) (err error) {
			suite.Equal("val-2", string(event.Value))
			return nil
		},
		"key-3": func(ctx context.Context, event *EventMsg) (err error) {
			suite.Equal("val-3", string(event.Value))
			return nil
		},
	}
	suite.processedMessages = make(chan struct{})
	suite.consumedMessages = make(chan *sarama.ConsumerMessage)
	suite.distributor = New(cfg, eventMethodMap, suite.consumedMessages, suite.processedMessages)

	defaultHandler := kafkaconnector.NewDefaultHandler(suite.consumedMessages)
	suite.consumer = client.NewConsumerGroup("test", []string{"test"}, defaultHandler)
}

func (suite *DistributorSuite) TearDownSuite() {
	if suite.distributor != nil {
		suite.distributor.Cleanup(suite.ctx)
	}
}

func (suite *DistributorSuite) SetupTest() {
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
}

func (suite *DistributorSuite) TearDownTest() {
	suite.cancel()
}

func (suite *DistributorSuite) TestDistribution() {
	go func() {
		suite.consumer.Consume(suite.ctx)
	}()

	go suite.distributor.Distribute(suite.ctx)

	time.Sleep(4 * time.Second)

	messages := []*kafkaconnector.Message{
		{Topic: "test", Key: []byte("key-1"), Value: []byte("val-1")},
		{Topic: "test", Key: []byte("key-2"), Value: []byte("val-2")},
		{Topic: "test", Key: []byte("key-3"), Value: []byte("val-3")},
	}

	for _, msg := range messages {
		partition, offset, err := suite.producer.SendMessage(suite.ctx, *msg)
		if err == nil {
			suite.T().Logf("Successfully sent message: %s to partition %d at offset %d",
				string(msg.Value), partition, offset)
		}
		suite.Require().NoError(err, "Failed to send message after retries")
	}
}
