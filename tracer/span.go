package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	trace.Span
	name string
}

func (r Span) TraceID() string {
	return r.Span.SpanContext().TraceID().String()
}

// Name returns name of span.
func (r Span) Name() string {
	return r.name
}

// Set sets attribute to span.
func (r Span) Set(key string, val any) Span {
	setAttr(r.Span, key, val)

	return r
}

// Event is a wrapper around otel/trace.span.AddEvent.
func (r Span) Event(name string) {
	r.Span.AddEvent(name)
}

func (r Span) Err(err error) error {
	if err == nil {
		return nil
	}

	r.Span.SetStatus(codes.Error, fmt.Sprintf("%+v", err))
	r.Span.RecordError(err)

	return err
}

func (r Span) ErrTrace(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w (trace: %s)", r.Err(err), r.TraceID())
}
