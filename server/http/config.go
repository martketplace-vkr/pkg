package http

import "github.com/martketplace-vkr/pkg/utils/duration"

type (
	Cors struct {
		AllowOrigins     string `validate:"required"`
		AllowHeaders     string `validate:"required"`
		ExposeHeaders    string `validate:"required"`
		AllowCredentials bool
	}
	Logging struct {
		SecureReqJsonPaths          []string
		SecureResJsonPaths          []string
		ShowUnknownErrorsInResponse bool
	}

	Serve struct {
		Host                  string `validate:"required" default:"127.0.0.1"`
		IpHeader              string `validate:"required" default:"X-Real-IP"`
		BodyLimit             int
		StopTimeout           duration.Seconds `default:"10"`
		DisableStartupMessage bool
	}

	Config struct {
		Cors    Cors    `validate:"required"`
		Serve   Serve   `validate:"required"`
		Logging Logging `validate:"required"`
	}
)
