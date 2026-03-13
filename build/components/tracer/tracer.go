package tracer

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	cmpName = "tracer"
)

type Tracer struct {
	cfg    Config
	tracer *tracer.Tracer
}

func New(cfg Config) (*Tracer, error) {
	tr, err := tracer.New(cfg.Config)
	if err != nil {
		return nil, err
	}

	return &Tracer{
		cfg:    cfg,
		tracer: tr,
	}, nil
}

func (t *Tracer) Start(_ context.Context) error {
	otel.SetTracerProvider(t.tracer.GetTracerProvider())
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return nil
}

func (t *Tracer) Stop(ctx context.Context) error {
	if err := t.tracer.GetTracerProvider().Shutdown(ctx); err != nil {
		return err
	}

	if err := t.tracer.GetExporter().Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (t *Tracer) GetStartTimeout() time.Duration {
	return t.cfg.StartTimeout.Duration
}

func (t *Tracer) GetStopTimeout() time.Duration {
	return t.cfg.StopTimeout.Duration
}

func (t *Tracer) GetShutdownDelay() time.Duration {
	return t.cfg.ShutdownDelay.Duration
}

func (t *Tracer) GetName() string {
	return cmpName
}
