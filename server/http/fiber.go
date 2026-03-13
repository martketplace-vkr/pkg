package http

import "github.com/gofiber/fiber/v2"

type (
	Server struct {
		*fiber.App
	}
)

func New(cfg Serve) *Server {
	return &Server{
		App: fiber.New(
			fiber.Config{
				BodyLimit:             cfg.BodyLimit,
				ProxyHeader:           cfg.IpHeader,
				DisableStartupMessage: cfg.DisableStartupMessage,
			},
		),
	}
}
