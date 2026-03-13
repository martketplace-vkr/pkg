package healthcheck

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/martketplace-vkr/pkg/logger/log"
)

const (
	cmpName = "healthcheck"
)

type HealthCheck struct {
	cfg   Config
	serve bool
	s     *fiber.App
}

func New(cfg Config, srv ...*fiber.App) *HealthCheck {
	serve := true
	s := fiber.New()

	if len(srv) > 0 {
		serve = false
		s = srv[0]
	}

	return &HealthCheck{
		cfg:   cfg,
		serve: serve,
		s:     s,
	}
}

func (hc *HealthCheck) Start(_ context.Context) error {
	hc.s.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	if hc.serve {
		go func() {
			if err := hc.s.Listen(hc.cfg.Address); err != nil {
				log.Fatal(err.Error())
			}
		}()
	}

	return nil
}

func (hc *HealthCheck) Stop(_ context.Context) error {
	return hc.s.ShutdownWithTimeout(hc.GetStopTimeout())
}

func (hc *HealthCheck) GetStartTimeout() time.Duration {
	return hc.cfg.StartTimeout.Duration
}

func (hc *HealthCheck) GetStopTimeout() time.Duration {
	return hc.cfg.StopTimeout.Duration
}

func (hc *HealthCheck) GetShutdownDelay() time.Duration {
	return hc.cfg.ShutdownDelay.Duration
}

func (p *HealthCheck) GetName() string {
	return cmpName
}
