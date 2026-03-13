package tracer

type (
	Config struct {
		URL          string `validate:"required" default:"http://localhost:14268/api/traces"`
		ServiceName  string `validate:"required" default:"default-service"`
		Auth         *Auth
		Sampler      *Sampler
		IgnoreErrors bool
	}
	Auth struct {
		Password string `validate:"required"`
		Username string `validate:"required"`
	}
	Sampler struct {
		Ratio float64
	}
)
