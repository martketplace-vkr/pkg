package inbox

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/prometheus/client_golang/prometheus"
)

type handlerMap map[string]func(context.Context, dto.Event) error

func (i *inbox) startEventsSchedulers(ctx context.Context) {
	metricsCli.setWorkersTotal(workerTypeScheduler, i.cfg.MessagesScheduler.ScheduleWorkersCount)

	for w := 0; w < i.cfg.MessagesScheduler.ScheduleWorkersCount; w++ {
		i.workersWG.Add(1)
		go func(workerID int) {
			defer i.workersWG.Done()
			for {
				select {
				case <-ctx.Done():
					log.Info("shutting down schedulers")
					return
				case <-i.signalChan:
					metricsCli.incWorkerBusy(workerTypeScheduler) // занят
					i.eventsSchedulerFunc()
					metricsCli.decWorkerBusy(workerTypeScheduler) // свободен
				}
			}
		}(w)
	}
}

func (i *inbox) broadcastWorkersSignal() {
	for w := 0; w < i.cfg.MessagesScheduler.ScheduleWorkersCount; w++ {
		select {
		case i.signalChan <- struct{}{}:
		default:
			// do nothing if signal is full
		}
	}
}

func (i *inbox) scheduleInboxFunc() {
	var (
		clean      dto.Events
		dirty      dto.FailedEvents
		err        error
		ctx        = context.Background()
		scheduleID = uuid.NewString()
	)
	defer func() {
		if err != nil {
			metricsCli.incErrorsTotal(err.Error())
		}
	}()

	// claim 1
	if err = i.txManager.Do(ctx, func(ctx context.Context) error {
		events, err := i.storage.FetchUnprocessedEvents(ctx, i.cfg.MessagesScheduler.BatchSize)
		if err != nil {
			log.Errorf("fetch unprocessed events: %v", err)

			return err
		}
		clean = events

		return nil
	}); err != nil || len(clean) == 0 {
		return
	}
	if i.cfg.MessagesScheduler.Debug {
		log.Debugf("scheduling %d events", len(clean))
	}

	if i.locker != nil {
		// 1. Lock events in redis.
		cleanEvents, skippedEvents, err := i.locker.Lock(ctx, clean)
		if err != nil {
			metricsCli.incErrorsTotal(err.Error())
			return
		}
		defer func() {
			err = i.locker.Unlock(ctx, cleanEvents)
			if err != nil {
				metricsCli.incErrorsTotal(err.Error())
				return
			}
		}()
		// if we found skippedEvents, then we need to unlock them (locked_until = false)
		if len(skippedEvents) != 0 {
			metricsCli.incRedisLocksConflicts(len(skippedEvents))

			err = i.storage.UnlockSkippedEvents(ctx, skippedEvents)
			if err != nil {
				metricsCli.incErrorsTotal(err.Error())
				return
			}
		}

		// if no events to handle then skip
		if len(cleanEvents) == 0 {
			return
		}

		clean = cleanEvents
	}

	// фиксируем обработку пары idempotency_key + event_type
	inserted, conflictedIdx, err := i.storage.RegisterEffectsBatch(ctx, clean)
	if err != nil {
		log.Errorf("register effects: %v", err)
		return
	}
	if len(inserted) != len(clean) {
		if i.cfg.MessagesScheduler.Debug {
			log.Debugf(
				"scheduleID: %s, found dirty events, cleaning... overall events: %d, inserted: %d, dirty: %d",
				scheduleID,
				len(clean),
				len(inserted),
				len(dirty),
			)
		}

		clean, dirty = filterEvents(clean, inserted, conflictedIdx)
	}

	// heartbeat while processing
	if i.cfg.MessagesScheduler.EnableAutoRenewLock {
		lease := i.cfg.MessagesScheduler.EventLockTimeout.Duration
		stop := i.startHeartbeat(ctx, clean, lease)
		defer stop()
	}

	if i.cfg.MessagesScheduler.Debug {
		log.Debugf(
			"scheduleID: %s, starting processing events... overall events: %d, dirty: %d",
			scheduleID,
			len(clean),
			len(dirty),
		)
	}
	processed, failed, perr := i.eventsHandlerFunc(ctx, clean)
	if perr != nil {
		err = perr
	}
	failed = append(failed, dirty...)

	if i.cfg.MessagesScheduler.Debug {
		log.Debugf(
			"scheduleID: %s, event processing result: processed evnts: %d, failed + dirty: %d",
			scheduleID,
			len(processed),
			len(failed),
		)
	}

	err = i.storage.MarkEventsAsProcessed(ctx, processed)
	if err != nil {
		log.Errorf("failed to mark events as processed: %v", err)
		return
	}
	if len(failed) > 0 {
		err = i.storage.UnlockFailedEvents(ctx, failed)
		if err != nil {
			log.Errorf("failed to unlock failed events: %v", err)
		}
	}

	if i.cfg.MessagesScheduler.Debug {
		log.Debugf(
			"scheduleID: %s, finished events processing: processed evnts: %d, failed + dirty: %d",
			scheduleID,
			len(processed),
			len(failed),
		)
	}

	metricsCli.observeEventsProcessResult(clean, processed)
}

func (i inbox) defaultEventsHandler(
	ctx context.Context,
	events dto.Events,
) (successEvents dto.SuccessEvents, failedEventsIDs dto.FailedEvents, err error) {
	for _, event := range events {
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metricsCli.observeEventProcessDuration(event.EventType, v)
		}))

		eventCtx := event.Context
		if event.Context == nil {
			eventCtx = ctx
		}

		handler, ok := i.handlerMap[event.EventType]
		if !ok {
			failedEventsIDs = append(failedEventsIDs, dto.FailedEvent{
				ID: event.ID,
				Error: fmt.Sprintf(
					"no handler for event with key %s, skipping...",
					event.EventType,
				),
			})
		} else {
			if err = handler(eventCtx, event); err != nil {
				failedEventsIDs = append(failedEventsIDs, dto.FailedEvent{
					ID:    event.ID,
					Error: err.Error(),
				})
			} else {
				successEvents = append(successEvents, event.ID)
			}
		}

		timer.ObserveDuration()
	}

	return successEvents, failedEventsIDs, err
}

func filterEvents(clean dto.Events, inserted, conflicted []int) (cleaned []dto.Event, failed dto.FailedEvents) {
	cleaned = filterByIdx(clean, inserted)
	dirtyEvents := filterByIdx(clean, conflicted)

	failed = make(dto.FailedEvents, len(dirtyEvents))
	for i, ev := range dirtyEvents {
		if metricsCli != nil {
			metricsCli.incDuplicateAttemptsTotal(ev.EventType)
		}

		failed[i] = dto.FailedEvent{
			ID:    ev.ID,
			Error: "duplicate process attempt",
		}
	}

	return cleaned, failed
}

func filterByIdx[T any](in []T, keepIdx []int) []T {
	if in == nil || keepIdx == nil || len(in) == 0 || len(keepIdx) == 0 {
		return []T{}
	}

	out := make([]T, 0, len(keepIdx))
	for _, i := range keepIdx {
		out = append(out, in[i])
	}

	return out
}
