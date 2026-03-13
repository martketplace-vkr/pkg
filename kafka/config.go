package kafka

type Config struct {
	SASL         *SASLConfig `envconfig:"" validate:"required"`
	Brokers      []string    `envconfig:"" validate:"required"`
	Version      string      `envconfig:"default=2.8.0" validate:"required"`
	OffsetNewest bool        `envconfig:"default=false"`
}

type SASLConfig struct {
	User     string `envconfig:"optional"`
	Password string `envconfig:"optional"`
}
