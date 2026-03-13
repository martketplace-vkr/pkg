package inbox

import (
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type (
	inboxOpts struct {
		sqlxDb      *sqlx.DB
		pgxDB       *pgxpool.Pool
		redisClient *redis.Client

		eventsProviderOpts  eventsProviderOpts
		eventsHandlerFunc   EventsHandlerFunc
		eventsSchedulerFunc EventsSchedulerFunc
		handlerMap          handlerMap
		skipEvent           skipEventFunc
		msgRetryAttempts    int
	}
	inboxOption interface {
		apply(*inboxOpts)
	}
)

type dbSqlxOption struct{ db *sqlx.DB }

func (o *dbSqlxOption) apply(opts *inboxOpts) { opts.sqlxDb = o.db }

func WithSqlxDB(db *sqlx.DB) inboxOption { return &dbSqlxOption{db: db} }

type dbPgxPoolOption struct{ db *pgxpool.Pool }

func (o *dbPgxPoolOption) apply(opts *inboxOpts) { opts.pgxDB = o.db }

func WithPgxPoolDB(db *pgxpool.Pool) inboxOption { return &dbPgxPoolOption{db: db} }

type eventsProviderOpts struct {
	provider     EventsProvider
	groupFactory func(groupID string) (sarama.ConsumerGroup, error)
}

func (o *eventsProviderOpts) apply(opts *inboxOpts) {
	opts.eventsProviderOpts = eventsProviderOpts{
		groupFactory: o.groupFactory,
	}
}

func WithCustomEventsProvider(provider EventsProvider) inboxOption {
	return &eventsProviderOpts{provider: provider}
}

func WithKafkaEventsProvider(
	groupFactory func(groupID string) (sarama.ConsumerGroup, error),
) inboxOption {
	return &eventsProviderOpts{
		groupFactory: groupFactory,
	}
}

type eventsHandlerFuncOption struct{ handlerFunc EventsHandlerFunc }

func (o *eventsHandlerFuncOption) apply(opts *inboxOpts) { opts.eventsHandlerFunc = o.handlerFunc }

func WithEventsHandlerFunc(handlerFunc EventsHandlerFunc) inboxOption {
	return &eventsHandlerFuncOption{handlerFunc: handlerFunc}
}

type eventsSchedulerFuncOption struct{ schedulerFunc EventsSchedulerFunc }

func (o *eventsSchedulerFuncOption) apply(opts *inboxOpts) {
	opts.eventsSchedulerFunc = o.schedulerFunc
}

func WithEventsSchedulerFunc(schedulerFunc EventsSchedulerFunc) inboxOption {
	return &eventsSchedulerFuncOption{schedulerFunc: schedulerFunc}
}

type handlerMapOption struct{ handlers handlerMap }

func (o *handlerMapOption) apply(opts *inboxOpts) { opts.handlerMap = o.handlers }

func WithHandlerMap(handlers handlerMap) inboxOption {
	return &handlerMapOption{handlers: handlers}
}

type skipEventOption struct{ skipEvent skipEventFunc }

func (o *skipEventOption) apply(opts *inboxOpts) { opts.skipEvent = o.skipEvent }

func WithCustomSkipEventFunc(skipEvent skipEventFunc) inboxOption {
	return &skipEventOption{skipEvent: skipEvent}
}

type redisClientOption struct {
	redisCli *redis.Client
}

func (o *redisClientOption) apply(opts *inboxOpts) {
	opts.redisClient = o.redisCli
}

func WithRedisClient(redis *redis.Client) inboxOption {
	return &redisClientOption{redisCli: redis}
}
