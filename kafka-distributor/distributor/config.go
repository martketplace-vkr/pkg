package distributor

type Config struct {
	WorkersCount int `validate:"required"`
}
