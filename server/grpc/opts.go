package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	GatewayOptions struct {
		registersHandlers []RegisterHandler
		mux               *runtime.ServeMux
	}
	ServerInterceptorOptions struct {
		unaryInterceptors  []grpc.UnaryServerInterceptor
		streamInterceptors []grpc.StreamServerInterceptor
	}
	Options struct {
		gateway      GatewayOptions
		interceptors ServerInterceptorOptions
	}
	RegisterHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
)

func WithGatewayServer(opts ...func(*GatewayOptions)) func(o *Options) {
	gwopts := GatewayOptions{}
	for _, opt := range opts {
		opt(&gwopts)
	}

	return func(o *Options) {
		o.gateway = gwopts
	}
}

func WithRegisterHandlers(registers ...RegisterHandler) func(*GatewayOptions) {
	return func(o *GatewayOptions) {
		o.registersHandlers = registers
	}
}

func WithMux(mux *runtime.ServeMux) func(*GatewayOptions) {
	return func(o *GatewayOptions) {
		o.mux = mux
	}
}

func WithUnaryInterceptors(unaryInterceptor ...grpc.UnaryServerInterceptor) func(*Options) {
	return func(o *Options) {
		o.interceptors.unaryInterceptors = append(o.interceptors.unaryInterceptors, unaryInterceptor...)
	}
}

func WithStreamInterceptors(streamInterceptors ...grpc.StreamServerInterceptor) func(*Options) {
	return func(o *Options) {
		o.interceptors.streamInterceptors = append(o.interceptors.streamInterceptors, streamInterceptors...)
	}
}
