package outbox

import (
	"context"
	"sync"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
)

func (o *outbox) cleanOnce(ctx context.Context) func() {
	if o.cfg.MessageCleaner == nil {
		return func() {
			log.Warn("outbox cleaner not configured, can't start")
		}
	}

	metricsCli.incGcHeartbeat()

	timeLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		timeLocation = time.UTC
	}

	return func() {
		var wg sync.WaitGroup

		// 1) чистим outbox_sending (аналог effect_log cleaner)
		wg.Add(1)
		go func() {
			defer wg.Done()
			o.runOutboxSendingCleaner(ctx, timeLocation)
		}()

		// 2) архивируем + чистим архив outbox
		wg.Add(1)
		go func() {
			defer wg.Done()
			o.runOutboxCleaner(ctx, timeLocation)
		}()

		wg.Wait()
	}
}

func (o *outbox) runOutboxSendingCleaner(ctx context.Context, _ *time.Location) {
	if o.cfg.MessageCleaner.OutboxSentRetentionDays.Duration <= 0 {
		return
	}
	older := time.Now().Add(-o.cfg.MessageCleaner.OutboxSentRetentionDays.Duration)

	deleted, err := o.storage.DeleteOutboxSentOlderThan(ctx, older, o.cfg.MessageCleaner.Batch)
	if err != nil {
		log.Warnf("outbox_sent cleaner: %v", err)
		return
	}
	if deleted > 0 {
		metricsCli.addSentOutboxDeleted(int(deleted))
	}
}

func (o *outbox) runOutboxCleaner(ctx context.Context, _ *time.Location) {
	// архивирование делаем в транзакции
	if err := o.txManager.Do(ctx, func(ctx context.Context) error {
		if o.cfg.MessageCleaner.OutboxArchiveAfterDays.Duration > 0 {
			older := time.Now().Add(-o.cfg.MessageCleaner.OutboxArchiveAfterDays.Duration)
			moved, err := o.storage.ArchiveSentOutbox(ctx, older, o.cfg.MessageCleaner.Batch)
			if err != nil {
				return err
			}
			if moved > 0 {
				metricsCli.addOutboxArchived(int(moved))
			}
		}

		if o.cfg.MessageCleaner.OutboxDeletedRetentionDays.Duration > 0 {
			delOlder := time.Now().Add(-o.cfg.MessageCleaner.OutboxDeletedRetentionDays.Duration)
			deleted, err := o.storage.DeleteFromOutboxArchiveOlderThan(ctx, delOlder, o.cfg.MessageCleaner.Batch)
			if err != nil {
				return err
			}
			if deleted > 0 {
				metricsCli.addOutboxDeletedPurged(int(deleted))
			}
		}

		return nil
	}); err != nil {
		log.Errorf("failed to launch outbox cleaner: %s", err.Error())
	}
}
