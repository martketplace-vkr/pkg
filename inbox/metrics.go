package inbox

import (
	"slices"
	"sync"

	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	once = sync.Once{}
)

const (
	workerTypeConsumer  = "consumer"
	workerTypeScheduler = "scheduler"
)

// metrics contains Prometheus metrics for inbox events, cleaners, heartbeat, and workers.
type metrics struct {
	namespace  string
	registerer prometheus.Registerer

	// ─── Event-related metrics ────────────────────────────────────────────────
	inboxErrorsTotal        *prometheus.CounterVec
	createdEventsTotal      *prometheus.CounterVec
	skippedEventsTotal      *prometheus.CounterVec
	deletedEventsTotal      prometheus.Counter
	receivedEventsTotal     *prometheus.CounterVec
	processedEventsTotal    *prometheus.CounterVec
	duplicateAttemptsTotal  *prometheus.CounterVec
	redisLockConflicts      prometheus.Counter
	eventProcessDurationSec *prometheus.HistogramVec
	gcHeartbeat             prometheus.Counter

	// ─── Cleaner metrics ─────────────────────────────────────────────────────
	effectLogDeleted   prometheus.Counter
	inboxArchived      prometheus.Counter
	inboxDeletedPurged prometheus.Counter

	// ─── Heartbeat metrics ───────────────────────────────────────────────────
	heartbeatExtendErrors    prometheus.Counter
	lockedUntilExtendedTotal prometheus.Counter

	// ─── Worker metrics ──────────────────────────────────────────────────────
	workersTotal *prometheus.GaugeVec
	workersBusy  *prometheus.GaugeVec
	workersIdle  *prometheus.GaugeVec

	mu           sync.Mutex
	totalWorkers map[string]int
	busyWorkers  map[string]int
}

// newMetrics creates and initializes all non-Kafka metrics for inbox.
func newMetrics(namespace string, registerer prometheus.Registerer) *metrics {
	latencyBuckets := []float64{
		0.001, 0.005, 0.01, 0.025, 0.05,
		0.1, 0.25, 0.5, 1, 2.5,
		5, 10, 30, 60,
	}

	return &metrics{
		namespace:  namespace,
		registerer: registerer,

		// ─── Base event metrics ──────────────────────────────
		receivedEventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "received_events_total", Help: "Total number of events received"},
			[]string{"event_type"},
		),
		createdEventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "created_events_total", Help: "Total number of successfully created events"},
			[]string{"event_type"},
		),
		processedEventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "processed_events_total", Help: "Total number of events processed"},
			[]string{"event_type"},
		),
		redisLockConflicts: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "redis_lock_conflicts_total", Help: "Total number of redis lock conflicts"},
		),
		deletedEventsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "deleted_events_total", Help: "Total number of events deleted by GC"},
		),
		gcHeartbeat: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "gc_heartbeat_total", Help: "Total GC heartbeat executions"},
		),
		inboxErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "errors_total", Help: "Total number of errors"},
			[]string{"error"},
		),
		skippedEventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "skipped_events_total", Help: "Total number of skipped events"},
			[]string{"event_type"},
		),
		duplicateAttemptsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Namespace: namespace, Name: "duplicate_attempts_total", Help: "Total number of duplicate handling attempts"},
			[]string{"event_type"},
		),
		eventProcessDurationSec: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Namespace: namespace, Name: "event_process_duration_seconds", Help: "Time spent processing an event", Buckets: latencyBuckets},
			[]string{"event_type"},
		),

		// ─── Cleaner ──────────────────────────────
		effectLogDeleted: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "effect_log_deleted_total", Help: "Total rows deleted from effect_log by cleaner"},
		),
		inboxArchived: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "archived_total", Help: "Total rows archived from inbox to inbox_archive"},
		),
		inboxDeletedPurged: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "archive_purged_total", Help: "Total rows purged from inbox_archive"},
		),

		// ─── Heartbeat ──────────────────────────────
		heartbeatExtendErrors: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "heartbeat_extend_errors_total", Help: "Number of errors while extending locked_until"},
		),
		lockedUntilExtendedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Namespace: namespace, Name: "heartbeat_locked_until_extended_total", Help: "Number of successful locked_until extensions"},
		),

		// ─── Workers ──────────────────────────────
		workersTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Name: "workers_total", Help: "Total number of workers"},
			[]string{"type"},
		),
		workersBusy: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Name: "workers_busy", Help: "Number of busy workers"},
			[]string{"type"},
		),
		workersIdle: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Namespace: namespace, Name: "workers_idle", Help: "Number of idle workers"},
			[]string{"type"},
		),

		totalWorkers: make(map[string]int),
		busyWorkers:  make(map[string]int),
	}
}

// register registers all inbox (non-Kafka) metrics.
func (m *metrics) register() {
	once.Do(func() {
		m.registerer.MustRegister(
			m.receivedEventsTotal,
			m.createdEventsTotal,
			m.processedEventsTotal,
			m.deletedEventsTotal,
			m.inboxErrorsTotal,
			m.skippedEventsTotal,
			m.duplicateAttemptsTotal,
			m.eventProcessDurationSec,
			m.gcHeartbeat,
			m.effectLogDeleted,
			m.inboxArchived,
			m.inboxDeletedPurged,
			m.heartbeatExtendErrors,
			m.lockedUntilExtendedTotal,
			m.workersTotal,
			m.workersBusy,
			m.workersIdle,
			m.redisLockConflicts,
		)
	})
}

// Cleaner metric helpers
func (m *metrics) addEffectLogDeleted(val int)   { m.effectLogDeleted.Add(float64(val)) }
func (m *metrics) addInboxArchived(val int)      { m.inboxArchived.Add(float64(val)) }
func (m *metrics) addInboxDeletedPurged(val int) { m.inboxDeletedPurged.Add(float64(val)) }
func (m *metrics) incHeartbeatExtendErrors()     { m.heartbeatExtendErrors.Inc() }
func (m *metrics) incRedisLocksConflicts(conflictedCount int) {
	m.redisLockConflicts.Add(float64(conflictedCount))
}
func (m *metrics) incLockedUntilExtended()        { m.lockedUntilExtendedTotal.Inc() }
func (m *metrics) incGcHeartbeat()                { m.gcHeartbeat.Inc() }
func (m *metrics) incErrorsTotal(err string)      { m.inboxErrorsTotal.WithLabelValues(err).Inc() }
func (m *metrics) incSkippedEventsTotal(t string) { m.skippedEventsTotal.WithLabelValues(t).Inc() }
func (m *metrics) incDuplicateAttemptsTotal(t string) {
	m.duplicateAttemptsTotal.WithLabelValues(t).Inc()
}
func (m *metrics) observeEventProcessDuration(t string, v float64) {
	m.eventProcessDurationSec.WithLabelValues(t).Observe(v)
}
func (m *metrics) observeEventsProcessResult(events dto.Events, sentIds []int64) {
	for _, e := range events {
		if slices.Contains(sentIds, e.ID) {
			m.addProcessedEventsTotal(e.EventType, 1)
		}
	}
}

// Common counters
func (m *metrics) addReceivedEventsTotal(t string, v int) {
	m.receivedEventsTotal.WithLabelValues(t).Add(float64(v))
}
func (m *metrics) addCreatedEventsTotal(t string, v int) {
	m.createdEventsTotal.WithLabelValues(t).Add(float64(v))
}
func (m *metrics) addProcessedEventsTotal(t string, v int) {
	m.processedEventsTotal.WithLabelValues(t).Add(float64(v))
}
func (m *metrics) addDeletedEventsTotal(v int) { m.deletedEventsTotal.Add(float64(v)) }

// Worker helpers
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
