package collector

type Config struct {
	GroupID string   `validate:"required"`
	Topics  []string `validate:"required"`
}
