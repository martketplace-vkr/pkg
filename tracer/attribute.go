package tracer

import (
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type attr struct {
	val any
	key string
}

func pretty(value any) (output string) {
	defer func() {
		if r := recover(); r != nil {
			output = "failed to serialize attribute value"
		}
	}()

	out, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		return fmt.Sprintf("%+v", value)
	}

	return string(out)
}

func setAttr(span trace.Span, key string, val any) {
	switch value := val.(type) {
	case int:
		span.SetAttributes(attribute.Int(key, value))
	case int64:
		span.SetAttributes(attribute.Int64(key, value))
	case float64:
		span.SetAttributes(attribute.Float64(key, value))
	case string:
		span.SetAttributes(attribute.String(key, value))
	case bool:
		span.SetAttributes(attribute.Bool(key, value))
	case json.RawMessage:
		span.SetAttributes(attribute.String(key, string(value)))
	default:
		span.SetAttributes(attribute.String(key, pretty(val)))
	}
}
