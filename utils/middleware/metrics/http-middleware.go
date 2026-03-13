package metrics

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "http"
)

type (
	HttpMetricsConfig struct {
		IgnoredHTTPPaths []string
	}
	HTTPMetricsInterceptor struct {
		cfg             HttpMetricsConfig
		requestsTotal   *prometheus.CounterVec
		requestDuration *prometheus.HistogramVec
		requestInFlight *prometheus.GaugeVec
	}
)

func NewHTTPMetrics(registry prometheus.Registerer, buckets []float64) *HTTPMetricsInterceptor {
	if registry == nil {
		registry = prometheus.DefaultRegisterer
	}

	m := &HTTPMetricsInterceptor{}
	factory := promauto.With(registry)

	m.requestsTotal = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "", "requests_total"),
		},
		[]string{"status_code", "path"},
	)

	if buckets == nil {
		buckets = defaultBuckets()
	}

	m.requestDuration = factory.NewHistogramVec(prometheus.HistogramOpts{
		Name:    prometheus.BuildFQName(namespace, "", "request_duration"),
		Buckets: buckets,
	},
		[]string{"status_code", "path"},
	)

	m.requestInFlight = factory.NewGaugeVec(prometheus.GaugeOpts{
		Name: prometheus.BuildFQName(namespace, "", "request_in_flight"),
	},
		[]string{"path"},
	)

	return m
}

func (it HTTPMetricsInterceptor) MetricsMiddleware(c *fiber.Ctx) (err error) {
	start := time.Now()
	path := clone(c.OriginalURL())

	it.IncRequestInFlight(path)
	defer it.DecRequestInFlight(path)

	err = c.Next()
	status := c.Response().StatusCode()

	statusCode := strconv.Itoa(status)
	it.IncRequestTotal(statusCode, path)

	elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
	it.ObserveRequestDuration(statusCode, path, elapsed)

	return err
}

func (it HTTPMetricsInterceptor) IncRequestInFlight(path string) {
	if contains(it.cfg.IgnoredHTTPPaths, path) {
		return
	}

	it.requestInFlight.WithLabelValues(path).Inc()
}

func (it HTTPMetricsInterceptor) DecRequestInFlight(path string) {
	if contains(it.cfg.IgnoredHTTPPaths, path) {
		return
	}

	it.requestInFlight.WithLabelValues(path).Dec()
}

func (it HTTPMetricsInterceptor) IncRequestTotal(statusCode, path string) {
	if contains(it.cfg.IgnoredHTTPPaths, path) {
		return
	}

	it.requestsTotal.WithLabelValues(statusCode, path).Inc()
}

func (it HTTPMetricsInterceptor) ObserveRequestDuration(statusCode, path string, elapsed float64) {
	if contains(it.cfg.IgnoredHTTPPaths, path) {
		return
	}

	it.requestDuration.WithLabelValues(statusCode, path).Observe(elapsed)
}

func clone(s string) string {
	b := make([]byte, len(s))
	copy(b, s)

	return *(*string)(unsafe.Pointer(&b))
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(a, e) {
			return true
		}
	}

	return false
}
