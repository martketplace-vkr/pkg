package http

import (
	"context"
	"time"

	"github.com/martketplace-vkr/pkg/server/http"
)

const (
	cmpName = "http-server"
)

type Server struct {
	cfg Config
	s   *http.FiberServer
}

func New(
	cfg Config,
	s *http.FiberServer,
) *Server {
	return &Server{
		cfg: cfg,
		s:   s,
	}
}

func (s *Server) Start(ctx context.Context) error {
	return s.s.Start(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.s.Stop(ctx)
}

func (s *Server) GetName() string {
	return cmpName
}

func (s *Server) GetStartTimeout() time.Duration {
	return s.cfg.StartTimeout.Duration
}

func (s *Server) GetStopTimeout() time.Duration {
	return s.cfg.StopTimeout.Duration
}

func (s *Server) GetShutdownDelay() time.Duration {
	return s.cfg.ShutdownDelay.Duration
}
