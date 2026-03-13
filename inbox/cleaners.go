package inbox

import (
	"context"
	"sync"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
)

func (i *inbox) сleanOnce(ctx context.Context) func() {
	if i.cfg.MessageCleaner == nil {
		return func() {
			log.Warn("message cleaner not configured, can't start")
		}
	}

	metricsCli.incGcHeartbeat()

	timeLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		timeLocation = time.UTC
	}

	return func() {
		wg := sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			i.runEffectLogCleaner(ctx, timeLocation)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			i.runInboxCleaner(ctx, timeLocation)
		}()

		wg.Wait()
	}
}

func (i *inbox) runEffectLogCleaner(ctx context.Context, loc *time.Location) {
	if err := i.cleanEffectLog(ctx); err != nil {
		log.Errorf("failed to clean effect log cleaner: %v", err)
	}
}

func (i *inbox) runInboxCleaner(ctx context.Context, loc *time.Location) {
	err := i.txManager.Do(ctx, func(ctx context.Context) error {
		err := i.archiveInbox(ctx)
		if err != nil {
			return err
		}

		err = i.purgeInboxDeleted(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("failed to launch inbox cleaner: %s", err.Error())
	}
}

func (i *inbox) cleanEffectLog(ctx context.Context) error {
	older := time.Now().Add(-i.cfg.MessageCleaner.EffectLogRetentionDays.Duration)

	deleted, err := i.storage.DeleteEffectLogOlderThan(ctx, older, i.cfg.MessageCleaner.Batch)
	if err != nil {
		log.Warnf("effect_log cleaner: %v", err)
		return err
	}
	if deleted == 0 {
		return nil
	}
	metricsCli.addEffectLogDeleted(int(deleted))

	if deleted < int64(i.cfg.MessageCleaner.Batch) {
		return nil
	}

	return nil
}

func (i *inbox) archiveInbox(ctx context.Context) error {
	older := time.Now().Add(-i.cfg.MessageCleaner.InboxArchiveAfterDays.Duration)

	moved, err := i.storage.ArchiveProcessedInbox(ctx, older, i.cfg.MessageCleaner.Batch)
	if err != nil {
		return err
	}
	if moved == 0 {
		return nil
	}

	metricsCli.addInboxArchived(int(moved))

	if moved < int64(i.cfg.MessageCleaner.Batch) {
		return nil
	}

	return nil
}

func (i *inbox) purgeInboxDeleted(ctx context.Context) error {
	delOlder := time.Now().Add(-i.cfg.MessageCleaner.InboxDeletedRetentionDays.Duration)

	deleted, err := i.storage.DeleteFromInboxDeletedOlderThan(ctx, delOlder, i.cfg.MessageCleaner.Batch)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return nil
	}

	metricsCli.addInboxDeletedPurged(int(deleted))

	if deleted < int64(i.cfg.MessageCleaner.Batch) {
		return nil
	}

	return nil
}
