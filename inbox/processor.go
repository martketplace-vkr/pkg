package inbox

import (
	"context"
	"slices"
	"sync"

	"github.com/avast/retry-go/v4"
	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/logger/log"
)

func (i inbox) processIncomingEvents(ctx context.Context) {
	events := i.eventsProvider.GetEventsChan()
	workers := 1
	if i.cfg.KafkaMessageReceiver != nil {
		workers = i.cfg.KafkaMessageReceiver.ConsumersCount
	}

	metricsCli.setWorkersTotal(workerTypeConsumer, workers)

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			i.process(ctx, events)
		}()
	}

	wg.Wait()
}

func (i inbox) process(ctx context.Context, events chan dto.Event) {
	for {
		select {
		case <-i.stopChan:
			log.Info("stopping process of incoming inbox events")
			return
		case <-ctx.Done():
			log.Info("stopping process of incoming inbox events")
			return
		case event := <-events:
			metricsCli.incWorkerBusy(workerTypeConsumer)
			metricsCli.addReceivedEventsTotal(event.EventType, 1)

			if i.skipEvent(event) {
				metricsCli.incSkippedEventsTotal(event.EventType)
				metricsCli.decWorkerBusy(workerTypeConsumer)
				continue
			}

			eventCtx := event.Context
			if eventCtx == nil {
				eventCtx = ctx
			}

			err := retry.Do(
				func() error {
					_, writeErr := i.CreateEvent(eventCtx, event)
					return writeErr
				},
				retry.Context(ctx),
				retry.Attempts(uint(i.saveMsgRetryAttempts)),
				retry.OnRetry(func(attempt uint, err error) {
					log.Infof("retry inbox write cause [err: %s], attempt: %d", err.Error(), attempt)
				}),
			)
			if err != nil {
				log.Errorf("faield to retry event after %d attempts", i.saveMsgRetryAttempts)
			}

			metricsCli.decWorkerBusy(workerTypeConsumer)
		}
	}
}

func (i inbox) defaultSkipEvent(event dto.Event) bool {
	if i.cfg.MessagesScheduler.EventListenKeys != nil {
		if slices.Contains(
			i.cfg.MessagesScheduler.EventListenKeys,
			event.EventType,
		) {
			return false
		} else {
			return true
		}
	} else if i.cfg.MessagesScheduler.EventKeyFilters != nil {
		if slices.Contains(
			i.cfg.MessagesScheduler.EventKeyFilters,
			event.EventType,
		) {
			return true
		} else {
			return false
		}
	}

	return false
}
