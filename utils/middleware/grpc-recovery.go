package middleware

import (
	"context"
	"fmt"
	"sync"

	"github.com/cockroachdb/errors"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	panicCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "grpc_server_panic_recovered_total",
			Help: "Total number of panics recovered in gRPC server",
		},
	)

	once = &sync.Once{}
)

const (
	ua = "user_agent.original"
)

func GRPCRecoverUnaryServerInterceptor(logger *zerolog.Logger, enableStack bool, registry ...prometheus.Registerer) grpc.UnaryServerInterceptor {
	reg := prometheus.DefaultRegisterer
	if len(registry) > 0 {
		reg = registry[0]
	}

	return grpcRecovery.UnaryServerInterceptor(
		grpcRecovery.WithRecoveryHandlerContext(newGRPCRecoveryHandler(logger, enableStack, reg)),
	)
}

func GRPCRecoverStreamServerInterceptor(logger *zerolog.Logger, enableStack bool, registry ...prometheus.Registerer) grpc.StreamServerInterceptor {
	reg := prometheus.DefaultRegisterer
	if len(registry) > 0 {
		reg = registry[0]
	}

	return grpcRecovery.StreamServerInterceptor(
		grpcRecovery.WithRecoveryHandlerContext(newGRPCRecoveryHandler(logger, enableStack, reg)),
	)
}

func newGRPCRecoveryHandler(logger *zerolog.Logger, _ bool, registry prometheus.Registerer) grpcRecovery.RecoveryHandlerFuncContext {
	once.Do(func() {
		registry.MustRegister(panicCounter)
	})

	return func(ctx context.Context, p interface{}) error {
		err := errors.Errorf("%v", p)
		span := trace.SpanFromContext(ctx)

		// no need, but just in case.
		if err == nil {
			return nil
		}
		span.SetAttributes(
			attribute.KeyValue{
				Key:   "error.message",
				Value: attribute.StringValue(err.Error()),
			},
			attribute.KeyValue{
				Key: "error.stack",
				Value: attribute.StringValue(
					fmt.Sprintf("%+v", errors.WithStack(err)),
				),
			},
		)

		logger.
			WithLevel(zerolog.PanicLevel).
			Stack().
			Err(err).
			Str(TypeKey, "panic").
			Str(TraceIDKey, span.SpanContext().TraceID().String()).
			Str(TraceFlagsKey, span.SpanContext().TraceFlags().String()).
			Str(ua, grpcUserAgent(ctx)).
			Str(ClientAddressKey, grpcRemoteIP(ctx)).
			Send()

		panicCounter.Inc()

		return status.Errorf(codes.Internal, "Internal Server Error")
	}
}
