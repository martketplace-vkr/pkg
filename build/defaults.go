package build

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/martketplace-vkr/pkg/build/components/healthcheck"
	"github.com/martketplace-vkr/pkg/build/components/pgxpoolcomponent"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/build/components/prometheus"
	"github.com/martketplace-vkr/pkg/build/components/rediscomponent"
	"github.com/martketplace-vkr/pkg/build/components/tracer"
	"github.com/martketplace-vkr/pkg/logger"
	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	DefaultsConfig struct {
		Logger      logger.Config
		Tracer      tracer.Config
		HealthCheck healthcheck.Config
		Prometheus  prometheus.Config
	}
	ComponentConfig struct {
		CmpType CmpType
		Config  CmpConfig
	}
)

func SetupDefaultComponents(
	cfg DefaultsConfig,
	addConfigs ...ComponentConfig,
) ([]Component, error) {
	appLogger := logger.NewLogger(cfg.Logger)
	appLogger.InitLogger()

	hc := healthcheck.New(cfg.HealthCheck)

	tr, err := tracer.New(cfg.Tracer)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %v", err)
	}

	prom := prometheus.New(cfg.Prometheus, nil)

	result := []Component{
		hc,
		tr,
		prom,
	}

	addComponents := make([]Component, len(addConfigs))
	for i, addConfig := range addConfigs {
		addComponents[i], err = componentFabric(addConfig.CmpType, addConfig.Config)
		if err != nil {
			return nil, err
		}
	}

	return append(
		result,
		addComponents...,
	), nil
}

type (
	CmpType   int
	CmpConfig any
)

const (
	RedisComponent CmpType = iota
	PgxpoolConnector
	PgxSqlxConnector
	Tracer
	Sentry
	Prometheus
	HealthCheck
	ClickHouseComponent
	MongoComponent
)

func componentFabric(cmpName CmpType, config CmpConfig) (Component, error) {
	switch cmpName {
	case RedisComponent:
		if cfg, ok := config.(rediscomponent.Config); !ok {
			return nil, errors.New("invalid redis component config type")
		} else {
			return rediscomponent.New(cfg), nil
		}
	case PgxSqlxConnector:
		if cfg, ok := config.(pgxsqlxcomponent.Config); !ok {
			return nil, errors.New("invalid pgxsqlx component config type")
		} else {
			return pgxsqlxcomponent.New(cfg), nil
		}
	case PgxpoolConnector:
		if cfg, ok := config.(pgxpoolcomponent.Config); !ok {
			return nil, errors.New("invalid pgxpool component config type")
		} else {
			return pgxpoolcomponent.New(cfg), nil
		}
	case Tracer:
		if cfg, ok := config.(tracer.Config); !ok {
			return nil, errors.New("invalid tracer component config type")
		} else {
			return tracer.New(cfg)
		}
	case Prometheus:
		if cfg, ok := config.(prometheus.Config); !ok {
			return nil, errors.New("invalid prometheus component config type")
		} else {
			return prometheus.New(cfg, nil), nil
		}
	case HealthCheck:
		if cfg, ok := config.(healthcheck.Config); !ok {
			return nil, errors.New("invalid healthcheck component config type")
		} else {
			return healthcheck.New(cfg), nil
		}
	default:
		log.Fatalf("unknown component type: %d", cmpName)
	}

	return nil, nil
}
