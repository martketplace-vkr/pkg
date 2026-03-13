package logger

type Config struct {
	Level                       string `default:"info" validate:"required"`
	SkipFrameCount              int    `default:"3"`
	PrettyLogging               bool
	ShowUnknownErrorsInResponse bool
	SecureReqJsonPaths          []string
	SecureResJsonPaths          []string
	SecureResGRPCPaths          []string
	SecureReqGRPCPaths          []string
	MaxHTTPBodySize             int `default:"500" validate:"required"`
	MaxGRPCBodySize             int `default:"500" validate:"required"`
}
