package outbox

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/redis/go-redis/v9"
)

type (
	outboxOpts struct {
		sqlxDb       *sqlx.DB
		pgxDB        *pgxpool.Pool
		producer     kafkaconnector.SyncProducer
		redis        *redis.Client
		customBroker Broker
	}
	outboxOption interface {
		apply(*outboxOpts)
	}
)

type dbSqlxOption struct {
	db *sqlx.DB
}

func (o *dbSqlxOption) apply(opts *outboxOpts) {
	opts.sqlxDb = o.db
}

func WithSqlxDB(db *sqlx.DB) outboxOption {
	return &dbSqlxOption{db: db}
}

type dbPgxPoolOption struct {
	db *pgxpool.Pool
}

func (o *dbPgxPoolOption) apply(opts *outboxOpts) {
	opts.pgxDB = o.db
}

func WithPgxPoolDB(db *pgxpool.Pool) outboxOption {
	return &dbPgxPoolOption{db: db}
}

type producerOption struct {
	producer kafkaconnector.SyncProducer
}

func (o *producerOption) apply(opts *outboxOpts) {
	opts.producer = o.producer
}

func WithKafkaProducer(producer kafkaconnector.SyncProducer) outboxOption {
	return &producerOption{producer: producer}
}

type idempotencyProtection struct {
	rdClient *redis.Client
}

func (o *idempotencyProtection) apply(opts *outboxOpts) {
	opts.redis = o.rdClient
}

func WithIdempotencyProtection(redisClient *redis.Client) outboxOption {
	return &idempotencyProtection{rdClient: redisClient}
}

type (
	customBroker struct {
		customPublisher Broker
	}
)

func (o *customBroker) apply(opts *outboxOpts) {
	opts.customBroker = o.customPublisher
}

func WithCustomEventPublisher(customPublisher Broker) outboxOption {
	return &customBroker{
		customPublisher: customPublisher,
	}
}

type redisClientOption struct {
	redisCli *redis.Client
}

func (o *redisClientOption) apply(opts *outboxOpts) {
	opts.redis = o.redisCli
}

func WithRedisClient(redis *redis.Client) outboxOption {
	return &redisClientOption{redisCli: redis}
}
