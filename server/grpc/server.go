package grpc

import (
	"context"
	"net/http"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/utils/middleware"
)

type Server struct {
	Grpc    *grpc.Server
	Gateway *http.Server
}

func New(
	ctx context.Context,
	cfg Config,
	registry prometheus.Registerer,
	opts ...func(*Options),
) (server *Server, err error) {
	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	if registry == nil {
		registry = prometheus.DefaultRegisterer
	}

	if cfg.BucketsForHistogram == nil {
		cfg.BucketsForHistogram = []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets(cfg.BucketsForHistogram),
		),
	)

	err = registry.Register(srvMetrics)
	if err != nil {
		log.Errorf("failed register grpc metrics: %v", err)
	}

	unaryOpts := []grpc.UnaryServerInterceptor{
		middleware.GRPCContextLoggerUnaryInterceptor(ctx),
		// errs.GetGRPCInterceptor(cfg.ServiceID),
		middleware.Auth(),
		srvMetrics.UnaryServerInterceptor(),
	}
	streamOpts := []grpc.StreamServerInterceptor{
		middleware.GRPCContextLoggerStreamInterceptor(ctx),
		srvMetrics.StreamServerInterceptor(),
	}

	if cfg.AccessLog {
		unaryOpts = append(
			unaryOpts, middleware.GRPCAccessLogUnaryServerInterceptor(&log.Logger),
		)
		streamOpts = append(
			streamOpts, middleware.GRPCAccessLogStreamServerInterceptor(&log.Logger),
		)
	}
	if !cfg.DisableRecovery {
		unaryOpts = append(
			unaryOpts, middleware.GRPCRecoverUnaryServerInterceptor(&log.Logger, true),
		)
		streamOpts = append(
			streamOpts, middleware.GRPCRecoverStreamServerInterceptor(&log.Logger, true),
		)
	}

	streamOpts = append(streamOpts, options.interceptors.streamInterceptors...)
	unaryOpts = append(unaryOpts, options.interceptors.unaryInterceptors...)

	srv := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle.Duration,
			MaxConnectionAge:  cfg.MaxConnectionAge.Duration,
			Timeout:           cfg.Timeout.Duration,
			Time:              cfg.Time.Duration,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			unaryOpts...,
		),
		grpc.ChainStreamInterceptor(
			streamOpts...,
		),
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
	)

	var (
		gateway *http.Server
		mux     *runtime.ServeMux
	)

	if options.gateway.mux != nil {
		mux = options.gateway.mux
	} else {
		mux = runtime.NewServeMux()
	}

	if len(options.gateway.registersHandlers) > 0 {
		gateway = &http.Server{
			Addr: cfg.GatewayHost,
		}

		for _, register := range options.gateway.registersHandlers {
			err = register(
				ctx,
				mux,
				cfg.Host,
				[]grpc.DialOption{
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				})
			if err != nil {
				return server, err
			}
		}

		gateway.Handler = mux
	}

	return &Server{
		Grpc:    srv,
		Gateway: gateway,
	}, nil
}
