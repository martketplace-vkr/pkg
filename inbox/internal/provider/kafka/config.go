package kafka

type MessageReceiverConfig struct {
	ConsumersCount   int    `validate:"required" default:"1"`
	GroupID          string `validate:"required"`
	Topics           []string
	ConsumerJitterMs int64  `validate:"required" default:"100"`
	BufferSize       int64  `validate:"required" default:"1024"`
	MetricsNS        string `validate:"required" default:"inbox"`
	Debug            bool
	MaxLogValueBytes int `default:"1024"`
}
