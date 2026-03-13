package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPMetrics(t *testing.T) {
	buckets := []float64{0.1, 0.2, 0.5, 1.0}
	m := NewHTTPMetrics(nil, buckets)

	assert.NotNil(t, m.requestsTotal)
	assert.NotNil(t, m.requestDuration)
	assert.NotNil(t, m.requestInFlight)
}

func TestHTTPMetricsInterceptor_IncRequestInFlight(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewHTTPMetrics(registry, nil)
	m.IncRequestInFlight("/test")

	assert.Equal(t, 1, testutil.CollectAndCount(m.requestInFlight))
}

func TestHTTPMetricsInterceptor_IncRequestTotal(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewHTTPMetrics(registry, nil)
	m.IncRequestTotal("200", "/test")

	assert.Equal(t, 1, testutil.CollectAndCount(m.requestsTotal))
}

func TestHTTPMetricsInterceptor_ObserveRequestDuration(t *testing.T) {
	registry := prometheus.NewRegistry()
	m := NewHTTPMetrics(registry, nil)
	m.ObserveRequestDuration("200", "/test", 0.5)

	assert.Equal(t, 1, testutil.CollectAndCount(m.requestDuration))
}
