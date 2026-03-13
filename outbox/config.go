package outbox

import (
	"github.com/martketplace-vkr/pkg/outbox/internal/locker"
	"github.com/martketplace-vkr/pkg/utils/duration"
)

type (
	Config struct {
		Metrics           Metrics `validation:"required"`
		MessageCleaner    *MessageCleaner
		MessagesScheduler MessagesScheduler `validation:"required"`

		// DistributedLocker locks events by key in the current instance replicas.
		DistributedLocker *locker.Config
	}
	Metrics struct {
		Namespace string `validate:"required" default:"outbox_events"`
	}

	// MessageCleaner defines configuration for periodic cleanup tasks,
	// such as archiving old inbox entries and deleting expired logs.
	MessageCleaner struct {
		// every hour from 2 to 8 am.
		SchedulePeriod string `default:"0 2-8 * * *"`

		// OutboxDeletedRetentionDays defines how long records in `outbox_archive`
		// should be kept before permanent removal.
		// Default: 1440h (60 days)
		OutboxDeletedRetentionDays duration.Days `validate:"required" default:"60"`

		// OutboxArchiveAfterDays defines how long processed inbox messages
		// should be kept before being moved to `outbox_archive`.
		// Default: 168h (7 days)
		OutboxArchiveAfterDays duration.Days `validate:"required" default:"7"`

		// OutboxSentRetentionDays defines how long entries in the `outbox_sending` table
		// should be retained before deletion.
		// Default: 720h (30 days)А
		OutboxSentRetentionDays duration.Days `validate:"required" default:"30"`

		// Batch defines the number of rows to process in one cleanup iteration.
		// Default: 1000
		Batch int `validate:"required" default:"1000"`
	}
	MessagesScheduler struct {
		BatchSize            int              `validation:"required" default:"1"`
		ScheduleTime         duration.Seconds `validation:"required" default:"10"`
		EventLockTimeout     duration.Seconds `validation:"required" default:"10"`
		ShutdownWaitTimeout  duration.Seconds `validation:"required" default:"10"`
		ScheduleWorkersCount int              `validation:"required" default:"1"`
		MaxProcessAttempts   int              `validation:"required" default:"1"`
		ImmediateSendTimeout *duration.Seconds
		// EnableAutoRenewLock enables automatic periodic extension of event locks (heartbeat).
		EnableAutoRenewLock bool
		// DisableSequentialEntriesSelection is a feature-flag to enable/disable new feature
		// in select query, where we consider the order of events by each key.
		// It is enabled by default, to disable set DisableSequentialEntriesSelection to true
		DisableSequentialEntriesSelection bool
	}
)
