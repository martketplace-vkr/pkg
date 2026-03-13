package http

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/utils/middleware/metrics"
)

var (
	ErrStartTimeout = errors.New("start http server timeout")
	defaultBuckets  = []float64{
		0.000000001,
		0.000000002,
		0.000000005,
		0.00000001,
		0.00000002,
		0.00000005,
		0.0000001,
		0.0000002,
		0.0000005,
		0.000001,
		0.000002,
		0.000005,
		0.00001,
		0.00002,
		0.00005,
		0.0001,
		0.0002,
		0.0005,
		0.001,
		0.002,
		0.005,
		0.01,
		0.02,
		0.05,
		0.1,
		0.2,
		0.5,
		1.0,
		2.0,
		5.0,
		10.0,
		15.0,
		20.0,
		30.0,
	}
)

type (
	FiberServer struct {
		cfg      Config
		server   Server
		binder   Binder
		registry prometheus.Registerer
	}

	Binder interface {
		BindRoutes(ctx context.Context)
	}
	UserNotFoundErr struct {
		Msg string `json:"msg"`
	}
)

func NewWithBinder(
	cfg Config,
	server *Server,
	binder Binder,
	registry prometheus.Registerer,
) *FiberServer {
	return &FiberServer{
		cfg:      cfg,
		server:   *server,
		binder:   binder,
		registry: registry,
	}
}

func (f *FiberServer) Start(_ context.Context) (err error) {
	if f.registry == nil {
		f.registry = prometheus.DefaultRegisterer
	}

	f.server.Use(cors.New(cors.Config{
		AllowOrigins:     f.cfg.Cors.AllowOrigins,
		AllowCredentials: f.cfg.Cors.AllowCredentials,
		AllowHeaders:     f.cfg.Cors.AllowHeaders,
		ExposeHeaders:    f.cfg.Cors.ExposeHeaders,
	}))

	f.server.Use(TraceMiddleware())

	f.server.Use(
		ErrorsMiddleware(
			log.Logger,
			f.cfg.Serve.IpHeader,
			getUserFromCtxFunc,
			f.cfg.Logging.SecureReqJsonPaths,
			f.cfg.Logging.SecureResJsonPaths,
			f.cfg.Logging.ShowUnknownErrorsInResponse,
		),
	)

	metric := metrics.NewHTTPMetrics(f.registry, defaultBuckets)
	f.server.Use(metric.MetricsMiddleware)

	f.withDefaultBinder(f.binder.BindRoutes)

	go func() {
		if err := f.server.Listen(f.cfg.Serve.Host); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return nil
}

func (f *FiberServer) Stop(_ context.Context) (err error) {
	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		if err := f.server.Shutdown(); err != nil {
			errCh <- err
		}
		okCh <- struct{}{}
	}()

	select {
	case <-okCh:
		return nil
	case err := <-errCh:
		return err
	case <-time.After(f.cfg.Serve.StopTimeout.Duration):
		return ErrStartTimeout
	}
}

func (f *FiberServer) withDefaultBinder(bindRoutes func(context.Context)) {
	bindRoutes(context.Background())

	f.server.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	f.server.Get("/lk/monitoring/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}
