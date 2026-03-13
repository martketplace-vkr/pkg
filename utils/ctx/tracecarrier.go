package ctx

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TraceCarrier map[string]string

func (t TraceCarrier) Get(key string) string {
	return t[key]
}

func (t TraceCarrier) Set(key, value string) {
	t[key] = value
}

func (t TraceCarrier) Keys() []string {
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}
	return keys
}

func NewTraceCarrierFromTraceCarrier(traceCarrier propagation.TextMapCarrier) TraceCarrier {
	carrier := make(TraceCarrier)

	for _, key := range traceCarrier.Keys() {
		value := traceCarrier.Get(key)

		carrier.Set(key, value)
	}

	return carrier
}

func (t *TraceCarrier) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("failed to assert type")
	}

	return json.Unmarshal(b, &t)
}

func (t *TraceCarrier) Context() context.Context {
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), t)

	return ctx
}

func (t *TraceCarrier) String() string {
	result, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		return ""
	}

	return string(result)
}
