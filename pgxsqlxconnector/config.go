package pgxsqlxconnector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	Config struct {
		Host     string `validate:"required" default:"localhost"`
		Port     string `validate:"required" default:"5432"`
		User     string `validate:"required" default:"postgres"`
		Password string `validate:"required" default:"password"`
		DBName   string `validate:"required"`
		SSLMode  string `validate:"required" default:"disable"`
		Settings Settings
		Extra    Extra

		// for backward compatibility. Use Extra instead of these three fields.
		LogQuery bool
		AppName  string
		HookFunc func(query string, took int64, appName string)
	}

	Extra struct {
		LogQuery      bool
		HookFunc      func(query string, took int64, appName string)
		EnableMetrics bool `default:"true"`
		Registry      prometheus.Registerer
	}
	Settings struct {
		MaxOpenConns    int           `validate:"required,min=1"`
		ConnMaxLifetime time.Duration `validate:"required,min=1"`
		MaxIdleConns    int           `validate:"required,min=1"`
		ConnMaxIdleTime time.Duration `validate:"required,min=1"`
	}
)
