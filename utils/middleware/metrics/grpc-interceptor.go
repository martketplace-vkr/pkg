package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func defaultBuckets() []float64 {
	return []float64{
		0.000000001,
		0.000000002,
		0.000000005,
		0.00000001,
		0.00000002,
		0.00000005,
		0.0000001,
		0.0000002,
		0.0000005,
		0.000001,
		0.000002,
		0.000005,
		0.00001,
		0.00002,
		0.00005,
		0.0001,
		0.0002,
		0.0005,
		0.001,
		0.002,
		0.005,
		0.01,
		0.02,
		0.05,
		0.1,
		0.2,
		0.5,
		1.0,
		2.0,
		5.0,
		10.0,
		15.0,
		20.0,
		30.0,
	}
}

type GRPCMetricsInterceptor struct {
	errorCount      *prometheus.CounterVec
	responseTime    *prometheus.SummaryVec
	countRequest    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	inFlight        *prometheus.GaugeVec
}

func NewGRPCInterceptor(registry prometheus.Registerer, percentiles Objectives, buckets []float64) *GRPCMetricsInterceptor {
	if registry == nil {
		registry = prometheus.DefaultRegisterer
	}

	factory := promauto.With(registry)

	if percentiles == nil {
		percentiles = Objectives{
			0.5: 0.01, 0.9: 0.005, 0.95: 0.002, 0.98: 0.001, 0.99: 0.001, 0.999: 0.0001, 0.9999: 0.00001,
		}
	}

	responseTime := factory.NewSummaryVec(
		//nolint:exhaustruct
		prometheus.SummaryOpts{
			Namespace: "grpc", Name: "response_time", Help: "Время ответа", Objectives: percentiles,
		},
		[]string{"handler"},
	)

	errorCount := factory.NewCounterVec(
		//nolint:exhaustruct
		prometheus.CounterOpts{Namespace: "grpc", Name: "error_count_total", Help: "Количество ошибок"},
		[]string{"handler", "code"},
	)

	countRequest := factory.NewCounterVec(
		//nolint:exhaustruct
		prometheus.CounterOpts{Namespace: "grpc", Name: "request_count_total", Help: "Количество запросов"},
		[]string{"handler"},
	)

	if len(buckets) == 0 {
		buckets = defaultBuckets()
	}

	requestDuration := factory.NewHistogramVec(
		prometheus.HistogramOpts{ //nolint:exhaustruct
			Namespace: "grpc",
			Name:      "request_duration",
			Help:      "Гистограмма входящих grpc запросов",
			Buckets:   buckets,
		},
		[]string{"handler"},
	)

	inFlight := factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "grpc", Name: "inflight_requests", Help: "Current in-flight gRPC requests",
		},
		[]string{"handler"},
	)

	return &GRPCMetricsInterceptor{
		errorCount:      errorCount,
		responseTime:    responseTime,
		countRequest:    countRequest,
		requestDuration: requestDuration,
		inFlight:        inFlight,
	}
}

func (c *GRPCMetricsInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		interface{}, error,
	) {
		c.inFlight.WithLabelValues(info.FullMethod).Inc()
		defer c.inFlight.WithLabelValues(info.FullMethod).Dec()

		start := time.Now()

		c.countRequest.WithLabelValues(info.FullMethod).Inc()

		result, err := handler(ctx, req)

		c.requestDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		if err == nil {
			c.responseTime.WithLabelValues(info.FullMethod).Observe(float64(time.Since(start).Milliseconds()))
		} else if e, ok := status.FromError(err); ok {
			c.errorCount.With(prometheus.Labels{"handler": info.FullMethod, "code": e.Code().String()}).Inc()
		}

		return result, err
	}
}

func (c *GRPCMetricsInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, streamServer grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		c.countRequest.WithLabelValues(info.FullMethod).Inc()

		err := handler(srv, streamServer)

		if err == nil {
			c.responseTime.WithLabelValues(info.FullMethod).Observe(float64(time.Since(start).Milliseconds()))
		} else if e, ok := status.FromError(err); ok {
			c.errorCount.With(prometheus.Labels{"handler": info.FullMethod, "code": e.Code().String()}).Inc()
		}

		return err
	}
}
