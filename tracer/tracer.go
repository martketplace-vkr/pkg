package tracer

import (
	"go.opentelemetry.io/otel"
	//nolint:staticcheck // SA1019: deprecated Jaeger exporter
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Tracer struct {
	cfg Config
	exp *jaeger.Exporter
	tp  *tracesdk.TracerProvider
}

func New(cfg Config) (*Tracer, error) {
	options := []jaeger.CollectorEndpointOption{
		jaeger.WithEndpoint(cfg.URL),
	}

	if cfg.Auth != nil {
		options = append(options, []jaeger.CollectorEndpointOption{
			jaeger.WithPassword(cfg.Auth.Password),
			jaeger.WithUsername(cfg.Auth.Username),
		}...)
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			options...,
		),
	)
	if err != nil {
		return nil, err
	}

	var traceProviderOpts []tracesdk.TracerProviderOption

	traceProviderOpts = append(
		traceProviderOpts,
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)),
	)

	if cfg.Sampler != nil {
		traceProviderOpts = append(
			traceProviderOpts,
			tracesdk.WithSampler(
				tracesdk.TraceIDRatioBased(cfg.Sampler.Ratio),
			),
		)
	}

	tp := tracesdk.NewTracerProvider(
		traceProviderOpts...,
	)

	if cfg.IgnoreErrors {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {}))
	}

	return &Tracer{
		cfg: cfg,
		tp:  tp,
		exp: exp,
	}, nil
}

func (t *Tracer) GetTracerProvider() *tracesdk.TracerProvider {
	return t.tp
}

func (t *Tracer) GetExporter() *jaeger.Exporter {
	return t.exp
}
