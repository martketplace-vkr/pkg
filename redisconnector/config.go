package redisconnector

type Config struct {
	Host               string `validate:"required" default:"localhost"`
	Port               string `validate:"required" default:"6379"`
	MinIdleConns       int    `validate:"required" default:"10"`
	PoolSize           int    `validate:"required" default:"100"`
	PoolTimeout        int    `validate:"required" default:"30"`
	Password           string `default:"password"`
	UseCertificates    bool
	InsecureSkipVerify bool
	CertificatesPaths  struct {
		Cert string
		Key  string
		Ca   string `default:"./config/redis_ca.crt"`
	}
	DB              int  `default:"0"`
	WithMetricsHook bool `default:"true"`
	EnableTracing   bool `default:"true"`
}
