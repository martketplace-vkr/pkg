package outbox

import (
	"context"

	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/outbox/dto"
	"github.com/martketplace-vkr/pkg/tracer"
	"github.com/prometheus/client_golang/prometheus"
)

func (o outbox) scheduleBroadcast() {
	ctx, span, _ := tracer.NewSpan(context.Background())
	defer span.End()

	var (
		cleanEvents dto.Events
		dirtyEvents dto.FailedEvents
		err         error
	)
	defer func() {
		if err != nil {
			log.Errorf("outbox schedule broadcast err: %v", err)
		}
	}()

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metricsCli.observeEventSendTime(v)
	}))

	err = o.txManager.Do(ctx, func(ctx context.Context) error {
		events, err := o.storage.FetchUnprocessedEvents(
			ctx,
			o.cfg.MessagesScheduler.BatchSize,
			o.cfg.MessagesScheduler.MaxProcessAttempts,
		)
		if err != nil {
			return err
		}
		cleanEvents = events

		return nil
	})
	if err != nil || cleanEvents == nil || len(cleanEvents) == 0 {
		return
	}

	if o.distributedLocker != nil {
		// 1. Lock events in redis.
		clean, skippedEvents, err := o.distributedLocker.Lock(ctx, cleanEvents)
		if err != nil {
			metricsCli.incErrorsTotal(err.Error())
			return
		}

		defer func() {
			err = o.distributedLocker.Unlock(ctx, clean)
			if err != nil {
				metricsCli.incErrorsTotal(err.Error())
				return
			}
		}()

		// if we found skippedEvents, then we need to unlock them (locked_until = false)
		if len(skippedEvents) != 0 {
			metricsCli.incRedisLocksConflicts(len(skippedEvents))

			err = o.storage.UnlockSkippedEvents(ctx, skippedEvents)
			if err != nil {
				metricsCli.incErrorsTotal(err.Error())
				return
			}
		}

		// if no events to handle then skip
		if len(clean) == 0 {
			return
		}

		cleanEvents = clean
	}

	// фиксируем обработку idempotency_key + event_type + topic
	inserted, conflictedIdx, err := o.storage.RegisterOutboxSendingBatch(ctx, cleanEvents)
	if err != nil {
		log.Errorf("register effects: %v", err)
		return
	}
	// если добавили ивентов меньше, чем хотели, значит есть конфликты, фильтруем их.
	if len(inserted) != len(cleanEvents) {
		cleanEvents, dirtyEvents = filterEvents(cleanEvents, inserted, conflictedIdx)
	}

	if o.cfg.MessagesScheduler.EnableAutoRenewLock {
		lease := o.cfg.MessagesScheduler.EventLockTimeout.Duration
		stop := o.startHeartbeat(ctx, cleanEvents, lease)
		defer stop()
	}

	processedEvents, failedEvents, err := o.broker.Publish(ctx, cleanEvents)
	if err != nil {
		log.Errorf("failed to process events: %v", err)
	}

	failedEvents = append(failedEvents, dirtyEvents...)

	err = o.storage.MarkEventsAsProcessed(ctx, processedEvents)
	if err != nil {
		return
	}

	if len(failedEvents) > 0 {
		err = o.storage.UnlockEvents(ctx, failedEvents)
		if err != nil {
			return
		}
	}

	metricsCli.observeEventsPublishResult(cleanEvents, processedEvents.IDs())
	timer.ObserveDuration()

	if err != nil {
		log.ErrorSentry(err, tracer.GetTraceID(span))
	}

	return
}

func (o *outbox) startEventsSchedulers(ctx context.Context) {
	metricsCli.setWorkersTotal(workerTypeScheduler, o.cfg.MessagesScheduler.ScheduleWorkersCount)

	for w := 0; w < o.cfg.MessagesScheduler.ScheduleWorkersCount; w++ {
		o.workersWG.Add(1)
		go func(workerID int) {
			defer o.workersWG.Done()
			for {
				select {
				case <-ctx.Done():
					log.Info("shutting down schedulers")
					return
				case <-o.signalChan:
					metricsCli.incWorkerBusy(workerTypeScheduler)
					o.scheduleBroadcast()
					metricsCli.decWorkerBusy(workerTypeScheduler)
				}
			}
		}(w)
	}
}

func (o *outbox) broadcastWorkersSignal() {
	for w := 0; w < o.cfg.MessagesScheduler.ScheduleWorkersCount; w++ {
		select {
		case o.signalChan <- struct{}{}:
		default:
			// do nothing if signal is full
		}
	}
}

func filterEvents(clean dto.Events, inserted, conflicted []int) (cleaned []dto.Event, failed dto.FailedEvents) {
	cleaned = filterByIdx(clean, inserted)
	dirtyEvents := filterByIdx(clean, conflicted)

	failed = make(dto.FailedEvents, len(dirtyEvents))
	for i, ev := range dirtyEvents {
		metricsCli.incDuplicateAttemptsTotal(ev.EventType)

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
