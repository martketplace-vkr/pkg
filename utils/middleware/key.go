//nolint:godot,godox
package middleware

// TODO: переместить константы в правильный пакет/место
const (
	TraceIDKey    = "trace_id"
	TraceFlagsKey = "trace_flags"
)

// TODO: wait for 1.21
const (
	RPCServerRequestSizeKey  = "rpc.server.request.size"
	RPCServerResponseSizeKey = "rpc.server.response.size"
	RPCServerDurationKey     = "rpc.server.duration"
)

const (
	TypeKey = "logging/structure"
)
