package pgxpoolconnctor

type (
	Config struct {
		Host     string `validate:"required" default:"localhost"`
		Port     int    `validate:"required" default:"5432"`
		User     string `validate:"required" default:"postgres"`
		Password string `validate:"required" default:"password"`
		DbName   string `validate:"required"`
		SSLMode  string `validate:"required" default:"disable"`
		Extra    Extra
	}
	Extra struct {
		MaxOpenConnections int32 `validate:"required" default:"3"`
		MinOpenConnections int32 `validate:"required" default:"1"`
		EnableMonitoring   bool
	}
)
