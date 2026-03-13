package outbox

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/martketplace-vkr/pkg/utils/duration"
	"github.com/redis/go-redis/v9"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func InitPostgres() (*sqlx.DB, error) {
	dsn := "postgres://root:dev@localhost:5432/postgres?sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

func InitRedis() (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis",
	})

	err := cli.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func InitPool() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("postgres://root:dev@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, err
	}
	config.MaxConns = 10
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func InitKafkaCg() (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	cg, err := sarama.NewConsumerGroup([]string{"localhost:29092"}, "inbox_group", config)
	if err != nil {
		return nil, err
	}

	return cg, nil
}

func InitKafka() (kafkaconnector.Client, error) {
	client := kafkaconnector.NewClient(kafkaconnector.ClientConfig{
		Brokers:            []string{"localhost:29092"},
		SASL:               &kafkaconnector.SASL{},
		InsecureSkipVerify: false,
		Producer: &kafkaconnector.ProducerConfig{
			ReadTimeout: duration.Seconds{
				Duration: 30 * time.Second,
			},
			WriteTimeout: duration.Seconds{
				Duration: 30 * time.Second,
			},
			RequireAcks: 1,
			MaxAttempts: 3,
			Compression: 0,
			RetryMax:    5,
			Idempotent: struct {
				Mode            bool
				MaxOpenRequests int
				RetryMax        int
			}{
				Mode:            false,
				MaxOpenRequests: 5,
				RetryMax:        3,
			},
		},
		Consumer: &kafkaconnector.ConsumerConfig{
			Addresses:    []string{"localhost:29092"},
			Assignor:     "round-robin",
			OffsetNewest: true,
			AutoCommit:   true,
			AutoCommitInterval: duration.Seconds{
				Duration: 5 * time.Second,
			},
		},
	})

	return client, nil
}

func InitJaeger() (*jaeger.Exporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint("http://localhost:16686"),
		),
	)
}

func InitTraceProvider(exp *jaeger.Exporter) *tracesdk.TracerProvider {
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("InboxTest"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)

	return tp
}
