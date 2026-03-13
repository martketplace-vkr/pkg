package utils

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/jackc/pgtype"

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

	}

	switch v.(type) {
	case []byte:
		b = v.([]byte)
	case string:
		b = []byte(v.(string))
	case pgtype.JSONB:
		b = v.(pgtype.JSONB).Bytes
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
