package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"path"
	"strings"
	"time"
)

const (
	TypeKey = "logging/structure"

	RPCServerRequestSizeKey  = "rpc.server.request.size"
	RPCServerResponseSizeKey = "rpc.server.response.size"
	RPCServerDurationKey     = "rpc.server.duration"

	TraceIDKey    = "trace_id"
	TraceFlagsKey = "trace_flags"

	UserAgentOriginalKey = "user_agent.original"
	RPCMethodKey         = "rpc.method"
	RPCGRPCStatusCodeKey = "rpc.transport.status_code"
	ClientAddressKey     = "client.address"

	ClientRequest  = "request"
	ServerResponse = "response"
)

func (l *LogManager) GRPCMiddleware(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		logErr := ctx.Err()
		if err != nil {
			logErr = err
		}

		grpcAccessLogEventWithBody(
			ctx,
			logger,
			info.FullMethod,
			start,
			logErr,
			req,
			resp,
			l.cfg.MaxGRPCBodySize,
			l.cfg.SecureReqGRPCPaths,
			l.cfg.SecureResJsonPaths).Send()

		return resp, err
	}
}

func grpcAccessLogEventWithBody(
	ctx context.Context,
	logger *zerolog.Logger,
	method string,
	start time.Time,
	err error,
	req, resp interface{},
	maxBodySize int,
	secureReqPath []string,
	secureResPath []string,
) *zerolog.Event {
	if strings.Contains(method, "grpc.health.v1.Health") {
		return logger.WithLevel(zerolog.NoLevel)
	}

	ev := logger.WithLevel(zerolog.InfoLevel).
		Err(err).
		Str(TypeKey, "access").
		Int64(RPCServerDurationKey, time.Since(start).Milliseconds()).
		Str(UserAgentOriginalKey, grpcUserAgent(ctx)).
		Str(RPCMethodKey, grpcMethodName(method)).
		Int(RPCGRPCStatusCodeKey, grpcStatusCode(err)).
		Str(ClientAddressKey, grpcRemoteIP(ctx)).
		Int(RPCServerRequestSizeKey, grpcMessageBytesCount(req)).
		Int(RPCServerResponseSizeKey, grpcMessageBytesCount(resp)).
		Str(ClientRequest, marshalAndTruncate(req, maxBodySize, secureReqPath)).
		Str(ServerResponse, marshalAndTruncate(resp, maxBodySize, secureResPath))

	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		ev = ev.
			Str(TraceIDKey, span.TraceID().String()).
			Str(TraceFlagsKey, span.TraceFlags().String())
	}

	return ev
}

func marshalAndTruncate(v interface{}, max int, securePath []string) string {
	if v == nil {
		return ""
	}

	msg, ok := v.(proto.Message)
	if !ok {
		return "unsupported type"
	}

	if proto.Size(msg) > max {
		notice, _ := sjson.SetBytes([]byte(`{}`), "hidden",
			fmt.Sprintf("too big (body length: %d, max size: %d)", proto.Size(msg), max))
		return string(notice)
	}

	bytes, err := protojson.Marshal(msg)
	if err != nil {
		return fmt.Sprintf("marshal error: %v", err)
	}

	for _, p := range securePath {
		if gjson.GetBytes(bytes, p).Exists() {
			bytes, _ = sjson.SetBytes(bytes, p, "secured in middleware")
		}
	}

	return string(bytes)
}

func grpcRemoteIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}

	return ""
}

func grpcMethodName(method string) string {
	return path.Base(method)
}

func grpcUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("user-agent"), "")
	}

	return ""
}

func grpcStatusCode(err error) int {
	return int(status.Convert(err).Code())
}

func grpcMessageBytesCount(message interface{}) int {
	if pb, ok := message.(proto.Message); ok {
		if b, err := protojson.Marshal(pb); err == nil {
			return len(b)
		}
	}

	return 0
}
