package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Database struct {
	Name  string `json:"name,omitempty"`
	Query string `json:"query,omitempty"`
	Took  string `json:"took,omitempty"`
}

func (lf *Database) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("name", lf.Name).
		Str("query", lf.Query).
		Str("took", lf.Took)
}

var Logger zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debug().Msgf(template, args...)
}

func Info(msg string) {
	Logger.Info().Msg(msg)
}

func Infof(template string, args ...interface{}) {
	Logger.Info().Msgf(template, args...)
}

func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

func Warnf(template string, args ...interface{}) {
	Logger.Warn().Msgf(template, args...)
}

func Error(err error) {
	Logger.Error().Err(err).Send()
}

func ErrorSentry(err error, traceID string) {
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureException(err)

	Logger.Error().Err(err).Str("trace_id", traceID).Send()
}

func ErrorSentryIgnoreCtx(err error, traceID string) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}

	Logger.Error().Err(err).Str("trace_id", traceID).Send()
}

func Errorf(template string, args ...interface{}) {
	Logger.Error().Msgf(template, args...)
}

func Panic(msg string) {
	Logger.Panic().Msg(msg)
}

func Panicf(template string, args ...interface{}) {
	Logger.Panic().Msgf(template, args...)
}

func Fatal(msg string) {
	Logger.Fatal().Msg(msg)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatal().Msgf(template, args...)

}

func DatabaseQuery(name string, query string, took time.Duration) {
	Logger.Info().EmbedObject(&Database{
		Name:  name,
		Query: query,
		Took:  fmt.Sprintf("%v", took),
	}).Msg("database query")
}
