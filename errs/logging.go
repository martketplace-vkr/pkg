package errs

import (
	"github.com/rs/zerolog"
)

type LoggingData struct {
	TraceID      string `json:"traceID"`
	RemoteIP     string `json:"remoteIP"`
	Host         string `json:"host"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Protocol     string `json:"protocol"`
	RequestBody  string `json:"requestBody"`
	ResponseBody string `json:"responseBody"`
	StatusCode   int    `json:"statusCode"`
	Latency      int64  `json:"latency"`
	User         any    `json:"user"`
	Error        string `json:"error"`
	Stack        []byte `json:"stack"`
}

func (d *LoggingData) MarshalZerologObject(e *zerolog.Event) {
	if e == nil {
		return
	}

	e.
		Str("traceID", d.TraceID).
		Str("remoteIP", d.RemoteIP).
		Str("host", d.Host).
		Str("method", d.Method).
		Str("path", d.Path).
		Str("protocol", d.Protocol).
		Str("requestBody", d.RequestBody).
		Str("responseBody", d.ResponseBody).
		Int("statusCode", d.StatusCode).
		Int64("latency", d.Latency).
		Any("user", &d.User).
		Str("error", d.Error)
}
