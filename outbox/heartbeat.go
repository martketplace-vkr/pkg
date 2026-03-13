package outbox

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/outbox/dto"
)

func (o *outbox) startHeartbeat(ctx context.Context, events dto.Events, extendBy time.Duration) (stop func()) {
	if extendBy <= 0 || len(events) == 0 {
		return func() {}
	}

	ids := events.IDs()
	interval := extendBy / 2
	if interval < time.Second {
		interval = time.Second
	}
	ctx, cancel := context.WithCancel(ctx)
	t := time.NewTicker(interval)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := o.storage.ExtendLocksByIDs(ctx, ids, extendBy); err != nil {
					log.Errorf("failed to extend locks: %v", err)
					metricsCli.incHeartbeatExtendErrors()
				} else {
					metricsCli.incLockedUntilExtended()
				}
			}
		}
	}()

	return cancel
}
