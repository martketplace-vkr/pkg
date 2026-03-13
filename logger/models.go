package logger

import (
	"encoding/json"

	"github.com/rs/zerolog"
)

type HttpFields struct {
	TraceID      string `json:"traceID,omitempty"`
	RemoteIP     string `json:"remoteIP,omitempty"`
	Host         string `json:"host,omitempty"`
	Method       string `json:"method,omitempty"`
	Path         string `json:"path,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
	RequestBody  string `json:"requestBody,omitempty"`
	ResponseBody string `json:"responseBody,omitempty"`
	StatusCode   int    `json:"statusCode,omitempty"`
	UserAgent    string `json:"userAgent,omitempty"`
	Latency      int64  `json:"latency,omitempty"`
	Error        string `json:"error,omitempty"`
	Stack        []byte `json:"stack,omitempty"`
	APIKey       string `json:"APIKey,omitempty"`
	UserID       int    `json:"userID,omitempty"`
	ClientID     int    `json:"clientID,omitempty"`
	APIID        int    `json:"APIID,omitempty"`
}

func (lf *HttpFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("traceID", lf.TraceID).
		Str("remoteIP", lf.RemoteIP).
		Str("host", lf.Host).
		Str("method", lf.Method).
		Str("path", lf.Path).
		Str("protocol", lf.Protocol).
		Str("requestBody", lf.RequestBody).
		Str("responseBody", lf.ResponseBody).
		Int("statusCode", lf.StatusCode).
		Int64("latency", lf.Latency).
		Str("APIKey", lf.APIKey).
		Int("userID", lf.UserID).
		Int("clientID", lf.ClientID).
		Int("APIID", lf.APIID).
		Str("error", lf.Error)
}

type LogStructure struct {
	HTTP zerolog.LogObjectMarshaler
}

func (lf *LogStructure) MarshalZerologObject(e *zerolog.Event) {
	if UnsafeMarshalJSON(lf.HTTP) != nil {
		e.RawJSON("HTTP", UnsafeMarshalJSON(lf.HTTP))
	}
}

func UnsafeMarshalJSON(value interface{}) []byte {
	v, _ := json.Marshal(value)
	if string(v) == "null" || string(v) == "{}" {
		return nil
	}
	return v
}
