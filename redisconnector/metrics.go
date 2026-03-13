package redisconnector

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
)

type (
	Hook struct {
		singleCommands    *prometheus.HistogramVec
		pipelinedCommands *prometheus.CounterVec
		singleErrors      *prometheus.CounterVec
		pipelinedErrors   *prometheus.CounterVec
	}
)

var (
	registry   = prometheus.DefaultRegisterer
	labelNames = []string{"command"}
	buckets    = []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1}
)

func NewHook() *Hook {
	registerer := promauto.With(registry)

	singleCommands := registerer.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "redis",
			Name:      "single_commands",
			Help:      "Histogram of single Redis commands",
			Buckets:   buckets,
		},
		labelNames,
	)

	pipelinedCommands := registerer.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "redis",
			Name:      "pipelined_commands",
			Help:      "Number of pipelined Redis commands",
		},
		labelNames,
	)

	singleErrors := registerer.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "redis",
			Name:      "single_errors",
			Help:      "Number of single Redis commands that have failed",
		},
		labelNames,
	)

	pipelinedErrors := registerer.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "redis",
			Name:      "pipelined_errors",
			Help:      "Number of pipelined Redis commands that have failed",
		},
		labelNames,
	)

	return &Hook{
		singleCommands:    singleCommands,
		pipelinedCommands: pipelinedCommands,
		singleErrors:      singleErrors,
		pipelinedErrors:   pipelinedErrors,
	}
}

func (hook *Hook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (hook *Hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		duration := time.Since(start).Seconds()

		hook.singleCommands.WithLabelValues(cmd.Name()).Observe(duration)

		if isActualErr(err) {
			hook.singleErrors.WithLabelValues(cmd.Name()).Inc()
		}

		return err
	}
}

func (hook *Hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		duration := time.Since(start).Seconds()

		for _, cmd := range cmds {
			hook.pipelinedCommands.WithLabelValues(cmd.Name()).Inc()

			if isActualErr(cmd.Err()) {
				hook.pipelinedErrors.WithLabelValues(cmd.Name()).Inc()
			}
		}

		hook.singleCommands.WithLabelValues("pipeline").Observe(duration)

		return err
	}
}

func isActualErr(err error) bool {
	return err != nil && !errors.Is(err, redis.Nil)
}
