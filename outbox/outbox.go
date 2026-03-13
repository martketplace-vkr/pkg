package outbox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/outbox/internal/broker/kafka"
	"github.com/martketplace-vkr/pkg/outbox/internal/locker/redis"
	"github.com/martketplace-vkr/pkg/outbox/internal/storage/postgres"
	"github.com/martketplace-vkr/pkg/utils/scheduler"
	"github.com/robfig/cron/v3"

	pgxManager "github.com/avito-tech/go-transaction-manager/pgxv5"
	txManager "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsCli *metrics
)

type (
	Outbox interface {
		Run(ctx context.Context) error
		Shutdown(ctx context.Context) error

		// CreateEvent publishes event to outbox
		CreateEvent(ctx context.Context, cmd CreateEvent) ([]int64, error)
	}
	outbox struct {
		cfg       Config
		storage   Storage
		broker    Broker
		txManager trm.Manager
		cron      *cron.Cron
		cancel    context.CancelFunc

		distributedLocker DistributedLocker

		workersWG sync.WaitGroup

		signalChan chan struct{}
	}
)

func New(
	cfg Config,
	storage Storage,
	broker Broker,
	txManager trm.Manager,
	registry prometheus.Registerer,
	cron *cron.Cron,
) *outbox {
	metricsCli = newMetrics(cfg.Metrics.Namespace, registry)

	return &outbox{
		cfg:        cfg,
		storage:    storage,
		broker:     broker,
		txManager:  txManager,
		signalChan: make(chan struct{}, cfg.MessagesScheduler.ScheduleWorkersCount),
		cron:       cron,
	}
}

func NewDefaultWithOptions(
	cfg Config,
	opts ...outboxOption,
) (*outbox, error) {
	options := &outboxOpts{}
	for _, opt := range opts {
		opt.apply(options)
	}

	var (
		storage        Storage
		producer       kafkaconnector.SyncProducer
		defaultManager trm.TrFactory
	)

	switch {
	case options.sqlxDb != nil:
		storage = postgres.NewSqlxStorage(
			options.sqlxDb,
			txManager.DefaultCtxGetter,
			cfg.MessagesScheduler.EventLockTimeout.Duration,
			cfg.MessagesScheduler.DisableSequentialEntriesSelection,
		)

		defaultManager = txManager.NewDefaultFactory(options.sqlxDb)
	case options.pgxDB != nil:
		storage = postgres.NewPgxPoolStorage(
			options.pgxDB,
			pgxManager.DefaultCtxGetter,
			cfg.MessagesScheduler.EventLockTimeout.Duration,
			cfg.MessagesScheduler.DisableSequentialEntriesSelection,
		)

		defaultManager = pgxManager.NewDefaultFactory(options.pgxDB)
	default:
		return nil, errors.New(
			"at least WithPgxPoolDB or WithSqlxDB option must be passed in arguments",
		)
	}

	if options.producer != nil {
		producer = options.producer
	} else {
		return nil, errors.New(
			"WithKafkaProducer option must be passed in arguments",
		)
	}

	metricsCli = newMetrics(cfg.Metrics.Namespace, prometheus.DefaultRegisterer)

	ob := outbox{
		cfg:        cfg,
		storage:    storage,
		broker:     kafka.NewSync(producer),
		txManager:  manager.Must(defaultManager),
		cron:       cron.New(),
		signalChan: make(chan struct{}, cfg.MessagesScheduler.ScheduleWorkersCount),
	}

	if options.customBroker != nil {
		ob.broker = options.customBroker
	}

	if cfg.DistributedLocker != nil {
		// only check for Redis for now, but could be more in the future.
		if options.redis == nil {
			return nil, errors.New(
				"you need to pass WithRedisClient option to enable redis distributed locker",
			)
		}

		if cfg.DistributedLocker.Disabled {
			log.Info("enabled redis distributed locker in mocked mode")
		} else {
			log.Info("enabled redis distributed locker")
		}

		ob.distributedLocker = redis.NewRedisDistributedLocker(*cfg.DistributedLocker, options.redis)
	}

	return &ob, nil
}

func (o outbox) Run(ctx context.Context) (err error) {
	metricsCli.register()

	runCtx, cancel := context.WithCancel(ctx)
	o.cancel = cancel

	// we need this to make sure that only 1 worker is currently scheduling.
	jobChain := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger))

	if o.cfg.MessageCleaner != nil {
		_, err := o.cron.AddJob(
			o.cfg.MessageCleaner.SchedulePeriod,
			jobChain.Then(cron.FuncJob(o.cleanOnce(ctx))),
		)
		if err != nil {
			return fmt.Errorf("failed to add deleteExpiredEvents job: %v", err)
		}
	}

	o.startEventsSchedulers(runCtx)

	_, err = o.cron.AddJob(
		scheduler.Period(o.cfg.MessagesScheduler.ScheduleTime.Duration),
		jobChain.Then(cron.FuncJob(o.broadcastWorkersSignal)),
	)
	if err != nil {
		return fmt.Errorf("failed to add scheduleBroadcast job: %v", err)
	}

	o.cron.Start()

	return nil
}

func (o outbox) Shutdown(_ context.Context) error {
	cronCtx := o.cron.Stop()

	waitTimeout := o.cfg.MessagesScheduler.ShutdownWaitTimeout.Duration

	// w8 for all workers to process with timeout.
	select {
	case <-cronCtx.Done():
		// ок
	case <-time.After(waitTimeout):
		return fmt.Errorf("timeout waiting for cron jobs to finish")
	}

	if o.cancel != nil {
		o.cancel()
	}

	done := make(chan struct{})
	go func() {
		o.workersWG.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(waitTimeout):
		return fmt.Errorf("timeout waiting for background workers to finish")
	}
}
