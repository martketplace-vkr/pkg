package inbox

import (
	"github.com/martketplace-vkr/pkg/inbox/internal/locker"
	"github.com/martketplace-vkr/pkg/inbox/internal/provider/kafka"
	"github.com/martketplace-vkr/pkg/utils/duration"
)

type (
	Config struct {
		// Metrics configuration (Prometheus namespace, etc.)
		Metrics Metrics `validation:"required"`

		// MessageCleaner defines how old messages and logs are periodically cleaned up.
		MessageCleaner *MessageCleaner

		// MessagesScheduler defines event scheduling and processing parameters.
		MessagesScheduler MessagesScheduler `validation:"required"`

		// KafkaMessageReceiver cfg for kafka consumer.
		KafkaMessageReceiver *kafka.MessageReceiverConfig `validation:"required"`

		// DistributedLocker locks events by key in the current instance replicas.
		DistributedLocker *locker.Config
	}

	// Metrics holds configuration for metric collection and namespace settings.
	Metrics struct {
		// Namespace defines the Prometheus metrics namespace.
		// Default: "inbox"
		Namespace string `validate:"required" default:"inbox"`
	}

	// MessageCleaner defines configuration for periodic cleanup tasks,
	// such as archiving old inbox entries and deleting expired logs.
	MessageCleaner struct {
		// SchedulePeriod defines a cron-style schedule for running cleanup.
		SchedulePeriod string `validate:"required" default:"*/10 2-8 * * *"`
		// InboxDeletedRetentionDays defines how long records in `inbox_archive`
		// should be kept before permanent removal.
		// Default: 1440h (60 days)
		InboxDeletedRetentionDays duration.Days `validate:"required" default:"60"`

		// InboxArchiveAfterDays defines how long processed inbox messages
		// should be kept before being moved to `inbox_archive`.
		// Default: 168h (7 days)
		InboxArchiveAfterDays duration.Days `validate:"required" default:"7"`

		// EffectLogRetentionDays defines how long entries in the `effect_log` table
		// should be retained before deletion.
		// Default: 720h (30 days)
		EffectLogRetentionDays duration.Days `validate:"required" default:"30"`

		// Batch defines the number of rows to process in one cleanup iteration.
		// Default: 1000
		Batch int `validate:"required" default:"1000"`
	}

	// MessagesScheduler defines configuration for event scheduling, locking,
	// and background task processing.
	MessagesScheduler struct {
		// ScheduleInterval defines how often the scheduler checks for new events (in seconds).
		// Default: 10
		ScheduleInterval duration.Seconds `validation:"required" default:"10"`

		// EventLockTimeout defines how long event locks are held before expiring (in seconds).
		// Default: 60
		EventLockTimeout duration.Seconds `validation:"required" default:"60"`

		// EventKeyFilters defines event keys that should be skipped during processing.
		EventKeyFilters []string

		// EventListenKeys defines event keys that the scheduler should actively listen for.
		EventListenKeys []string

		// ProcessAttempts defines how many times an event processing attempt is made before giving up.
		// Default: 1
		ProcessAttempts int `validation:"required" default:"1"`

		// ScheduleWorkersCount defines the number of concurrent scheduler workers.
		// Default: 1
		ScheduleWorkersCount int `validation:"required" default:"1"`

		// BatchSize defines how many events are fetched and scheduled per iteration.
		// Default: 1
		BatchSize int `validation:"required" default:"1"`

		// EnableAutoRenewLock enables automatic periodic extension of event locks (heartbeat).
		EnableAutoRenewLock bool

		// ShutdownWaitTimeout defines how long the scheduler waits (in seconds)
		// for active jobs to complete before shutting down.
		// Default: 5
		ShutdownWaitTimeout duration.Seconds `validation:"required" default:"5"`

		// Debug defines whether we will log processing events or not
		// Default: false
		Debug bool

		// DisableSequentialEntriesSelection is a feature-flag to enable/disable new feature
		// in select query, where we consider the order of events by each key.
		// It is enabled by default, to disable set DisableSequentialEntriesSelection to true
		DisableSequentialEntriesSelection bool
	}
)
