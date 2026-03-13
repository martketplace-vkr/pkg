package inbox

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"

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
	db.SetMaxOpenConns(60)
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
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	cg, err := sarama.NewConsumerGroup([]string{"localhost:29092"}, "inbox_group", config)
	if err != nil {
		return nil, err
	}

	return cg, nil
}

func InitKafkaProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 200 * time.Millisecond
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Compression = sarama.CompressionSnappy
	config.ClientID = "my-kafka-producer"
	brokers := []string{"localhost:29092"}

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
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
