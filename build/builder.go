package build

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/martketplace-vkr/pkg/build/components/healthcheck"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/build/components/prometheus"
	"github.com/martketplace-vkr/pkg/build/components/rediscomponent"
	"github.com/martketplace-vkr/pkg/build/components/tracer"
	"github.com/martketplace-vkr/pkg/logger"
	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/samber/lo"
)

type ComponentBuilder struct {
	components map[CmpType]Component
}

func NewBuilder() *ComponentBuilder {
	return &ComponentBuilder{
		components: make(map[CmpType]Component),
	}
}

func (b *ComponentBuilder) WithDefaults(cfg DefaultsConfig) *ComponentBuilder {
	appLogger := logger.NewLogger(cfg.Logger)
	appLogger.InitLogger()

	b.components[Tracer] = must(tracer.New(cfg.Tracer))
	b.components[Prometheus] = prometheus.New(cfg.Prometheus, nil)

	return b
}

func (b *ComponentBuilder) WithPgxSqlx(cfg pgxsqlxcomponent.Config) *ComponentBuilder {
	b.components[PgxSqlxConnector] = pgxsqlxcomponent.New(cfg)

	return b
}

func (b *ComponentBuilder) WithLivenessProbe(cfg healthcheck.Config, srv ...*fiber.App) *ComponentBuilder {
	b.components[HealthCheck] = healthcheck.New(cfg, srv...)

	return b
}

func (b *ComponentBuilder) WithRedis(cfg rediscomponent.Config) *ComponentBuilder {
	b.components[RedisComponent] = rediscomponent.New(cfg)

	return b
}

func (b *ComponentBuilder) Build() Components {
	defer func() {
		fmt.Println()
		log.Info("all components are built 🫦, good luck, little creatures")
	}()
	cmps := make(Components, len(b.components))

	return lo.MapToSlice(b.components, func(cmpType CmpType, cmp Component) Component {
		cmps = append(cmps, cmp)
		log.Infof("component %s is built 🚀", cmp.GetName())

		return cmp
	})
}

func must(component Component, err error) Component {
	if err != nil {
		log.Fatalf("failed to create component: %v", err)
	}

	return component
}
