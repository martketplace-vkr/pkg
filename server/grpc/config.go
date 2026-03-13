package grpc

import "github.com/martketplace-vkr/pkg/utils/duration"

type Config struct {
	Host                string           `validate:"required" default:"0.0.0.0:5000"`
	GatewayHost         string           `validate:"required" default:"0.0.0.0:8080"`
	StartTimeout        duration.Seconds `validate:"required" default:"5"`
	StopTimeout         duration.Seconds `validate:"required" default:"5"`
	MaxConnectionIdle   duration.Seconds `validate:"required" default:"5"`
	Timeout             duration.Seconds `validate:"required" default:"20"`
	MaxConnectionAge    duration.Seconds `validate:"required" default:"5"`
	Time                duration.Seconds `validate:"required" default:"1"`
	ServiceID           int              `validate:"required" default:"66"`
	ServiceName         string           `validate:"required" default:"default-service"`
	EnableMetrics       bool
	SecureConnection    bool
	MaxRecvMsgSize      int `validate:"required" default:"4194304"` // 4 MB
	MaxSendMsgSize      int `validate:"required" default:"4194304"` // 4 MB
	AccessLog           bool
	DisableRecovery     bool
	BucketsForHistogram []float64
}
