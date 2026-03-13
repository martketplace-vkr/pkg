package tracer

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const (
	argAttributeName  = "span.argument"
	initialSkipFrames = 2
	attrSlicePrealloc = 2
	traceIDKey        = "traceID"
)

type SpanBuilder struct {
	name       string
	attributes []attr
	skipFrames int
}

// Name sets span name.
// If Name is not invoked before Start span name is set
// automatically using runtime.Caller().
func Name(name string) SpanBuilder {
	attrs := make([]attr, 0, attrSlicePrealloc)

	return SpanBuilder{
		name:       name,
		attributes: attrs,
		skipFrames: initialSkipFrames,
	}
}

// With sets span attribute.
func With(key string, value any) SpanBuilder {
	builder := Name("")

	builder.attributes = append(builder.attributes, attr{
		val: value,
		key: key,
	})

	return builder
}

// With sets span attribute.
func (sb SpanBuilder) With(key string, value any) SpanBuilder {
	sb.attributes = append(sb.attributes, attr{
		val: value,
		key: key,
	})

	return sb
}

// WithArg sets span attribute under key "richspan.argument".
// Primariy meant to be for functions that take only one argument.
// Consecutive calls overwrite previous value.
func WithArg(value any) SpanBuilder {
	builder := Name("")

	builder.attributes = append(builder.attributes, attr{
		val: value,
		key: argAttributeName,
	})

	return builder
}

// WithArg sets span attribute under key "richspan.argument".
// Primariy meant to be for functions that take only one argument.
// Consecutive calls overwrite previous value.
func (sb SpanBuilder) WithArg(value any) SpanBuilder {
	sb.attributes = append(sb.attributes, attr{
		val: value,
		key: argAttributeName,
	})

	return sb
}

// Start starts the span.
func (sb SpanBuilder) Start(ctx context.Context) (context.Context, Span) {
	if sb.name == "" {
		pc, _, _, ok := runtime.Caller(sb.skipFrames)

		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			detailsItems := strings.Split(details.Name(), slash)

			if len(detailsItems) != 0 {
				sb.name = detailsItems[len(detailsItems)-1]
			} else {
				sb.name = details.Name()
			}
		} else {
			sb.name = Anonymous
		}
	}

	ctx, span := otel.Tracer("").Start(ctx, sb.name)

	deadline, ok := ctx.Deadline()
	span.SetAttributes(attribute.Bool("ctx.deadline", ok))

	if ok {
		span.SetAttributes(attribute.String("ctx.deadline_in", fmt.Sprint(time.Until(deadline))))
	}

	for _, attr := range sb.attributes {
		setAttr(span, attr.key, attr.val)
	}

	return ctx, Span{
		Span: span,
		name: sb.name,
	}
}

func GetTraceID(span trace.Span) string {
	return span.SpanContext().TraceID().String()
}
