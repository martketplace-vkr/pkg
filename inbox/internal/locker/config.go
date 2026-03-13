package locker

import "time"

type Config struct {
	EventLockTTL time.Duration `validate:"required"`
	// InstanceID defines the service, that locks that row
	// e.g. crypto-orchestrator
	InstanceID string `validate:"required"`
	Disabled   bool
}
