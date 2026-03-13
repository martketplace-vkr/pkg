package outbox

import (
	"slices"
	"sync"

	"github.com/martketplace-vkr/pkg/outbox/dto"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	once = sync.Once{}
)

type metrics struct {
	namespace                string
	registerer               prometheus.Registerer
	heartbeatExtendErrors    prometheus.Counter
	lockedUntilExtendedTotal prometheus.Counter
	createdEventsTotal       prometheus.Counter
	sentEventsTotal          *prometheus.CounterVec
	deletedEventsTotal       prometheus.Counter
	gcHeartbeat              prometheus.Counter
	sendErrorsTotal          *prometheus.CounterVec
	eventProcessDuration     prometheus.Histogram
	duplicateAttemptsTotal   *prometheus.CounterVec
	redisLockConflicts       prometheus.Counter
	outboxErrorsTotal        *prometheus.CounterVec

	// ─── Cleaner metrics ─────────────────────────────────────────────────────
	sentOutboxDeleted   prometheus.Counter
	outboxArchived      prometheus.Counter
	outboxDeletedPurged prometheus.Counter

	// ─── Worker metrics ──────────────────────────────────────────────────────
	workersTotal *prometheus.GaugeVec
	workersBusy  *prometheus.GaugeVec
	workersIdle  *prometheus.GaugeVec

	mu           sync.Mutex
	totalWorkers map[string]int
	busyWorkers  map[string]int
}

func newMetrics(namespace string, registerer prometheus.Registerer) *metrics {
	latencyBuckets := []float64{
		0.001, 0.005, 0.01, 0.025, 0.05,
		0.1, 0.25, 0.5, 1, 2.5,
		5, 10, 30, 60,
	}

	return &metrics{
		namespace:    namespace,
		registerer:   registerer,
		mu:           sync.Mutex{},
		totalWorkers: make(map[string]int),
		busyWorkers:  make(map[string]int),

		createdEventsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "created_events_total",
				Help:      "Total number of events created.",
			},
		),
		outboxErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "errors_total", Help: "Total number of errors"},
			[]string{"error"},
		),
		sentEventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "sent_events_total",
				Help:      "Total number of events sent successfully.",
			},
			[]string{"event_type"},
		),
		deletedEventsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "deleted_events_total",
				Help:      "Total number of events deleted by GC.",
			},
		),
		gcHeartbeat: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "gc_heartbeat",
				Help:      "Total number of events deleted by GC",
			},
		),
		eventProcessDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "event_process_duration_seconds",
				Help:      "Time spent sending an event (seconds).",
				Buckets:   latencyBuckets,
			},
		),
		heartbeatExtendErrors: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "heartbeat_extend_errors_total", Help: "Number of errors while extending locked_until"},
		),
		lockedUntilExtendedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "heartbeat_locked_until_extended_total", Help: "Number of successful locked_until extensions"},
		),
		duplicateAttemptsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "duplicate_attempts_total",
				Help:      "Total number of duplicate handling attempts",
			},
			[]string{"event_type"},
		),
		workersTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "workers_total",
				Help:      "Total number of workers",
			},
			[]string{"type"},
		),
		workersBusy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "workers_busy",
				Help:      "Number of busy workers",
			},
			[]string{"type"},
		),
		workersIdle: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "workers_idle",
				Help:      "Number of idle workers",
			},
			[]string{"type"},
		),

		// ─── Cleaner ──────────────────────────────
		sentOutboxDeleted: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "sent_deleted_total", Help: "Total rows deleted from sent_outbox by cleaner"},
		),
		outboxArchived: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "archived_total", Help: "Total rows archived from outbox to outbox_archive"},
		),
		outboxDeletedPurged: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "archive_purged_total", Help: "Total rows purged from outbox_archive"},
		),
		redisLockConflicts: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "redis_lock_conflicts_total", Help: "Total number of redis lock conflicts"},
		),
	}
}

func (m *metrics) register() {
	collectors := []prometheus.Collector{
		m.createdEventsTotal,
		m.sentEventsTotal,
		m.deletedEventsTotal,
		m.eventProcessDuration,
		m.duplicateAttemptsTotal,
		m.workersTotal,
		m.workersBusy,
		m.workersIdle,
		m.sentOutboxDeleted,
		m.outboxArchived,
		m.outboxDeletedPurged,
		m.lockedUntilExtendedTotal,
		m.heartbeatExtendErrors,
		m.outboxErrorsTotal,
		m.redisLockConflicts,
	}
	once.Do(func() {
		m.registerer.MustRegister(
			collectors...,
		)
	})
}

const (
	workerTypeConsumer  = "consumer"
	workerTypeScheduler = "scheduler"
)

func (m *metrics) setWorkersTotal(workerType string, total int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalWorkers[workerType] = total
	m.workersTotal.WithLabelValues(workerType).Set(float64(total))

	busy := m.busyWorkers[workerType]
	idle := total - busy
	if idle < 0 {
		idle = 0
	}
	m.workersBusy.WithLabelValues(workerType).Set(float64(busy))
	m.workersIdle.WithLabelValues(workerType).Set(float64(idle))
}

func (m *metrics) incWorkerBusy(workerType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.busyWorkers[workerType]++
	busy := m.busyWorkers[workerType]
	total := m.totalWorkers[workerType]
	idle := total - busy
	if idle < 0 {
		idle = 0
	}

	m.workersBusy.WithLabelValues(workerType).Set(float64(busy))
	m.workersIdle.WithLabelValues(workerType).Set(float64(idle))
}

func (m *metrics) decWorkerBusy(workerType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.busyWorkers[workerType] > 0 {
		m.busyWorkers[workerType]--
	}
	busy := m.busyWorkers[workerType]
	total := m.totalWorkers[workerType]
	idle := total - busy
	if idle < 0 {
		idle = 0
	}

	m.workersBusy.WithLabelValues(workerType).Set(float64(busy))
	m.workersIdle.WithLabelValues(workerType).Set(float64(idle))
}
func (m *metrics) incHeartbeatExtendErrors() { m.heartbeatExtendErrors.Inc() }
func (m *metrics) incLockedUntilExtended()   { m.lockedUntilExtendedTotal.Inc() }
func (m *metrics) addCreatedEventsCount(val int) {
	if val > 0 {
		m.createdEventsTotal.Add(float64(val))
	}
}

func (m *metrics) addSentEventsCount(eventType string, val int) {
	if val > 0 {
		m.sentEventsTotal.WithLabelValues(eventType).Add(float64(val))
	}
}

func (m *metrics) addDeletedEventsTotal(val int) {
	if val > 0 {
		m.deletedEventsTotal.Add(float64(val))
	}
}

func (m *metrics) observeEventSendTime(v float64) {
	m.eventProcessDuration.Observe(v)
}

func (m *metrics) observeEventsPublishResult(events dto.Events, sentIds []int64) {
	for _, event := range events {
		if slices.Contains(sentIds, event.ID) {
			m.addSentEventsCount(event.EventType, 1)
		}
	}
}

func (m *metrics) observeEventsCreateResult(events []int64) {
	m.addCreatedEventsCount(len(events))
}

func (m *metrics) incDuplicateAttemptsTotal(eventType string) {
	m.duplicateAttemptsTotal.WithLabelValues(eventType).Inc()
}

func (m *metrics) incGcHeartbeat() {
	m.gcHeartbeat.Inc()
}

// ─── Cleaner ──────────────────────────────
func (m *metrics) addOutboxArchived(val int) {
	m.outboxArchived.Add(float64(val))
}

func (m *metrics) addOutboxDeletedPurged(val int) {
	m.outboxDeletedPurged.Add(float64(val))
}

func (m *metrics) addSentOutboxDeleted(val int) {
	m.sentOutboxDeleted.Add(float64(val))
}

func (m *metrics) incErrorsTotal(err string) {
	m.outboxErrorsTotal.WithLabelValues(err).Inc()
}

func (m *metrics) incRedisLocksConflicts(conflictedCount int) {
	m.redisLockConflicts.Add(float64(conflictedCount))
}
