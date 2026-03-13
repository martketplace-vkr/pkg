package kafka

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

var kafkaOnce = sync.Once{}

type kafkaMetrics struct {
	namespace  string
	registerer prometheus.Registerer

	kafkaMessagesIn    prometheus.Counter
	kafkaSendTimeouts  prometheus.Counter
	kafkaConsumeErrors prometheus.Counter

	consumerPartitions   *prometheus.GaugeVec
	consumerHeartbeats   *prometheus.GaugeVec
	messagesPerPartition *prometheus.CounterVec
}

func newKafkaMetrics(namespace string, registerer prometheus.Registerer) *kafkaMetrics {
	m := &kafkaMetrics{
		namespace:  namespace,
		registerer: registerer,
		kafkaMessagesIn: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "kafka_messages_in_total",
				Help:      "Total messages received from Kafka",
			},
		),
		kafkaSendTimeouts: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "kafka_send_timeouts_total",
				Help:      "Total number of send timeouts to events channel",
			},
		),
		kafkaConsumeErrors: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "kafka_consume_errors_total",
				Help:      "Total consume errors in Kafka provider",
			},
		),
		consumerPartitions: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "kafka_consumer_partitions",
				Help:      "Kafka consumer assigned partitions by worker and topic",
			},
			[]string{"worker", "topic", "partition"},
		),
		consumerHeartbeats: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "kafka_consumer_last_heartbeat_timestamp",
				Help:      "Last heartbeat timestamp for Kafka consumer workers",
			},
			[]string{"worker"},
		),
		messagesPerPartition: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "kafka_messages_per_partition_total",
				Help:      "Number of messages consumed per topic and partition",
			},
			[]string{"topic", "partition"},
		),
	}

	m.register()

	return m
}

func (k *kafkaMetrics) register() {
	kafkaOnce.Do(func() {
		k.registerer.MustRegister(
			k.kafkaMessagesIn,
			k.kafkaSendTimeouts,
			k.kafkaConsumeErrors,
			k.consumerPartitions,
			k.consumerHeartbeats,
			k.messagesPerPartition,
		)
	})
}

// Helpers
func (k *kafkaMetrics) incKafkaMessagesIn()    { k.kafkaMessagesIn.Inc() }
func (k *kafkaMetrics) incKafkaSendTimeouts()  { k.kafkaSendTimeouts.Inc() }
func (k *kafkaMetrics) incKafkaConsumeErrors() { k.kafkaConsumeErrors.Inc() }

func (k *kafkaMetrics) setPartitionAssigned(workerID int, topic string, partition int32) {
	k.consumerPartitions.WithLabelValues(fmt.Sprintf("%d", workerID), topic, fmt.Sprintf("%d", partition)).Set(1)
}

func (k *kafkaMetrics) unsetPartition(workerID int, topic string, partition int32) {
	k.consumerPartitions.WithLabelValues(fmt.Sprintf("%d", workerID), topic, fmt.Sprintf("%d", partition)).Set(0)
}

func (k *kafkaMetrics) recordHeartbeat(workerID int) {
	k.consumerHeartbeats.WithLabelValues(fmt.Sprintf("%d", workerID)).Set(float64(time.Now().Unix()))
}

func (k *kafkaMetrics) incPartitionMessage(topic string, partition int32) {
	k.messagesPerPartition.WithLabelValues(topic, fmt.Sprintf("%d", partition)).Inc()
}
