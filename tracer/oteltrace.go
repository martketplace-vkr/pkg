package tracer

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SpanSetErrWrapf(span trace.Span, err error, template string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	template = fmt.Sprintf("%s; traceID: %s", template, span.SpanContext().TraceID().String())
	err = errors.Wrapf(err, template, args...)
	span.SetStatus(codes.Error, err.Error())

	return err
}

//nolint:typecheck,nolintlint
func StartFiberTrace(c *fiber.Ctx, caller ...string) (context.Context, Span) {
	spanName := ""

	if len(caller) != 0 {
		spanName = caller[0]
	} else {
		pc, _, _, ok := runtime.Caller(1)

		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			detailsItems := strings.Split(details.Name(), slash)

			if len(detailsItems) != 0 {
				spanName = detailsItems[len(detailsItems)-1]
			} else {
				spanName = details.Name()
			}
		} else {
			spanName = Anonymous
		}
	}

	tracedCtx, ok := c.Locals(TracedCTX).(context.Context)
	if !ok {
		tracedCtx = c.Context()
	}

	ctx, span := otel.Tracer("").Start(tracedCtx, spanName)
	if !ok {
		c.Locals(TracedCTX, ctx)
	}

	traceID := span.SpanContext().TraceID().String()

	c.Response().Header.Set(TraceIDHeader, traceID)
	c.Locals(TraceIDHeader, traceID)

	return ctx, Span{
		Span: span,
		name: spanName,
	}
}

func NewSpan(c context.Context, caller ...string) (ctx context.Context, span Span, spanName string) {
	if len(caller) != 0 {
		spanName = caller[0]
	} else {
		pc, _, _, ok := runtime.Caller(1)

		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			detailsItems := strings.Split(details.Name(), slash)

			if len(detailsItems) != 0 {
				spanName = detailsItems[len(detailsItems)-1]
			} else {
				spanName = details.Name()
			}
		} else {
			spanName = Anonymous
		}
	}

	span.name = spanName

	ctx, span.Span = otel.Tracer("").Start(c, spanName)

	return ctx, span, spanName
}
