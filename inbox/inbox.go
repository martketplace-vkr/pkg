package inbox

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/inbox/internal/locker/redis"
	"github.com/martketplace-vkr/pkg/inbox/internal/storage/postgres"
	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/utils/scheduler"

	pgxManager "github.com/avito-tech/go-transaction-manager/pgxv5"
	txManager "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

var (
	metricsCli *metrics
)

type (
	Inbox interface {
		Run(ctx context.Context) error
		Shutdown(ctx context.Context) error
		CreateEvent(ctx context.Context, event dto.Event) (int64, error)
	}
	inbox struct {
		cfg  *Config
		opts *inboxOpts

		// storage stores all the incoming events
		storage Storage
		// eventsProvider provides a channel with incoming events, that inbox must save
		eventsProvider EventsProvider
		txManager      trm.Manager
		// eventsSchedulerFunc is a func that launches inbox worker with some period
		eventsSchedulerFunc EventsSchedulerFunc
		// eventsHandlerFunc is a function that actually processes incoming events
		// if you want to specify your own logic of events process
		// u may pass a custom function  in options.
		eventsHandlerFunc EventsHandlerFunc
		// saveMsgRetryAttempts specifies incoming event save attempts
		saveMsgRetryAttempts int
		// may be nil if custom handler func were provided
		// otherwise it will be filled with passed by user value
		// so there is 2 ways to use it
		// 1) you pass custom eventsHandlerFunc so ure decide how to process events without map
		// 2) you pass the handler map, so it will use default events processor func to process
		handlerMap handlerMap
		skipEvent  skipEventFunc
		stopChan   chan struct{}
		signalChan chan struct{}

		cron *cron.Cron

		cancel    context.CancelFunc
		workersWG sync.WaitGroup

		locker DistributedLocker
	}
)

func New(
	cfg Config,
	storage Storage,
	eventsProvider EventsProvider,
	txManager trm.Manager,
	registry prometheus.Registerer,
	cron *cron.Cron,
	eventsHandlerFunc EventsHandlerFunc,
	eventsSchedulerFunc EventsSchedulerFunc,
) *inbox {
	metricsCli = newMetrics(cfg.Metrics.Namespace, registry)

	i := &inbox{
		cfg:                 &cfg,
		storage:             storage,
		eventsProvider:      eventsProvider,
		txManager:           txManager,
		cron:                cron,
		eventsSchedulerFunc: eventsSchedulerFunc,
		eventsHandlerFunc:   eventsHandlerFunc,
		stopChan:            make(chan struct{}),
	}

	if eventsHandlerFunc == nil {
		i.eventsHandlerFunc = i.defaultEventsHandler
	}
	if eventsSchedulerFunc == nil {
		i.eventsSchedulerFunc = i.scheduleInboxFunc
	}

	return i
}

func NewDefaultWithOptions(
	cfg Config,
	opts ...inboxOption,
) (*inbox, error) {
	options := &inboxOpts{}
	for _, opt := range opts {
		opt.apply(options)
	}

	var (
		msgRetryAttempts = 5
		eventsProvider   EventsProvider
		storage          Storage
		defaultManager   trm.TrFactory
	)

	switch {
	case options.sqlxDb != nil:
		if cfg.MessagesScheduler.DisableSequentialEntriesSelection {
			log.Info("disabling sequential entries selection")
		}

		storage = postgres.NewSqlxStorage(
			options.sqlxDb,
			txManager.DefaultCtxGetter,
			cfg.MessagesScheduler.EventLockTimeout.Duration,
			cfg.MessagesScheduler.ProcessAttempts,
			cfg.MessagesScheduler.DisableSequentialEntriesSelection,
		)

		defaultManager = txManager.NewDefaultFactory(options.sqlxDb)
	case options.pgxDB != nil:
		if cfg.MessagesScheduler.DisableSequentialEntriesSelection {
			log.Info("disabling sequential entries selection")
		}

		storage = postgres.NewPgxPoolStorage(
			options.pgxDB,
			pgxManager.DefaultCtxGetter,
			cfg.MessagesScheduler.EventLockTimeout.Duration,
			cfg.MessagesScheduler.ProcessAttempts,
			cfg.MessagesScheduler.DisableSequentialEntriesSelection,
		)

		defaultManager = pgxManager.NewDefaultFactory(options.pgxDB)
	default:
		return nil, errors.New(
			"either WithPgxPoolDB or WithSqlxDB option must be passed in arguments",
		)
	}

	switch {
	case options.eventsProviderOpts.provider != nil:
		eventsProvider = options.eventsProviderOpts.provider
	case options.eventsProviderOpts.groupFactory != nil && cfg.KafkaMessageReceiver != nil:
		eventsProvider = NewKafkaProvider(
			*cfg.KafkaMessageReceiver,
			options.eventsProviderOpts.groupFactory,
		)
	default:
		return nil, errors.New(
			"either WithEventsProvider or WithKafkaProvider must be passed in arguments and KafkaMessageReceiver != nil",
		)
	}

	if cfg.MessagesScheduler.EventListenKeys != nil &&
		cfg.MessagesScheduler.EventKeyFilters != nil {
		return nil, errors.New(
			"you must specify only EventListenKeys or EventKeyFilters in config",
		)
	}

	if options.msgRetryAttempts != 0 {
		msgRetryAttempts = options.msgRetryAttempts
	}

	cronWorker := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)
	metricsCli = newMetrics(cfg.Metrics.Namespace, prometheus.DefaultRegisterer)

	i := &inbox{
		cfg:                  &cfg,
		opts:                 options,
		storage:              storage,
		eventsProvider:       eventsProvider,
		txManager:            manager.Must(defaultManager),
		cron:                 cronWorker,
		saveMsgRetryAttempts: msgRetryAttempts,
		stopChan:             make(chan struct{}),
		signalChan:           make(chan struct{}, cfg.MessagesScheduler.ScheduleWorkersCount),
	}

	if options.eventsHandlerFunc == nil {
		if options.handlerMap == nil {
			return nil, errors.New(
				"if you want to use default event handler, then you need to pass WithEventsHandlerMap in arguments",
			)
		}

		i.handlerMap = options.handlerMap
		i.eventsHandlerFunc = i.defaultEventsHandler
	} else {
		i.eventsHandlerFunc = options.eventsHandlerFunc
	}

	if options.eventsSchedulerFunc == nil {
		i.eventsSchedulerFunc = i.scheduleInboxFunc
	} else {
		i.eventsSchedulerFunc = options.eventsSchedulerFunc
	}

	if options.skipEvent == nil {
		i.skipEvent = i.defaultSkipEvent
	} else {
		i.skipEvent = options.skipEvent
	}

	if cfg.DistributedLocker != nil {
		// only check for Redis for now, but could be more in the future.
		if options.redisClient == nil {
			return nil, errors.New(
				"you need to pass WithRedisClient option to enable redis distributed locker",
			)
		}

		if cfg.DistributedLocker.Disabled {
			log.Info("enabled redis distributed locker in mocked mode")
		} else {
			log.Info("enabled redis distributed locker")
		}

		i.locker = redis.NewRedisDistributedLocker(*cfg.DistributedLocker, options.redisClient)
	}

	return i, nil
}

func (i inbox) Run(ctx context.Context) (err error) {
	metricsCli.register()

	runCtx, cancel := context.WithCancel(ctx)
	i.cancel = cancel

	// we need this to make sure that only 1 worker is currently scheduling.
	jobChain := cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger))

	if i.cfg.MessageCleaner != nil {
		_, err = i.cron.AddJob(
			i.cfg.MessageCleaner.SchedulePeriod,
			jobChain.Then(cron.FuncJob(i.сleanOnce(ctx))),
		)
		if err != nil {
			return fmt.Errorf("failed to add deleteExpiredEvents job: %v", err)
		}
	}

	i.startEventsSchedulers(ctx)

	_, err = i.cron.AddJob(
		scheduler.Period(i.cfg.MessagesScheduler.ScheduleInterval.Duration),
		jobChain.Then(cron.FuncJob(i.broadcastWorkersSignal)),
	)
	if err != nil {
		return fmt.Errorf("failed to add scheduleBroadcast job: %v", err)
	}

	i.workersWG.Add(2)
	go func() {
		defer i.workersWG.Done()
		i.processIncomingEvents(runCtx)
	}()
	go func() {
		defer i.workersWG.Done()
		i.eventsProvider.Consume(runCtx)
	}()

	i.cron.Start()

	return nil
}

func (i *inbox) Shutdown(ctx context.Context) error {
	cronCtx := i.cron.Stop()

	waitTimeout := i.cfg.MessagesScheduler.ShutdownWaitTimeout.Duration
	timeoutCtx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()

	select {
	case <-cronCtx.Done():
		// ok
	case <-timeoutCtx.Done():
		return fmt.Errorf("timeout waiting for cron jobs to finish: %w", timeoutCtx.Err())
	}

	if i.cancel != nil {
		i.cancel()
	}

	done := make(chan struct{})
	go func() {
		i.workersWG.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		return fmt.Errorf("timeout waiting for background workers to finish: %w", timeoutCtx.Err())
	}
}
