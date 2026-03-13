package build

import (
	"github.com/martketplace-vkr/pkg/build/components/healthcheck"
	"github.com/martketplace-vkr/pkg/build/components/pgxpoolcomponent"
	"github.com/martketplace-vkr/pkg/build/components/pgxsqlxcomponent"
	"github.com/martketplace-vkr/pkg/build/components/prometheus"
	"github.com/martketplace-vkr/pkg/build/components/rediscomponent"
	"github.com/martketplace-vkr/pkg/build/components/tracer"
	"github.com/martketplace-vkr/pkg/logger/log"
)

func (b *ComponentBuilder) GetHealthCheck() *healthcheck.HealthCheck {
	if cmp, ok := b.components[HealthCheck].(*healthcheck.HealthCheck); ok {
		return cmp
	}

	log.Fatal("HealthCheck component not found")
	return nil
}

// GetTracer returns the Tracer component or panics if not found
func (b *ComponentBuilder) GetTracer() *tracer.Tracer {
	if cmp, ok := b.components[Tracer].(*tracer.Tracer); ok {
		return cmp
	}

	log.Fatal("Tracer component not found")
	return nil
}

// GetPrometheus returns the Prometheus component or panics if not found
func (b *ComponentBuilder) GetPrometheus() *prometheus.Prometheus {
	if cmp, ok := b.components[Prometheus].(*prometheus.Prometheus); ok {
		return cmp
	}

	log.Fatal("Prometheus component not found")
	return nil
}

// GetPgxSqlx returns the PgxSqlx component or panics if not found
func (b *ComponentBuilder) GetPgxSqlx() *pgxsqlxcomponent.PgxSqlxConnector {
	if cmp, ok := b.components[PgxSqlxConnector].(*pgxsqlxcomponent.PgxSqlxConnector); ok {
		return cmp
	}

	log.Fatal("PgxSqlx component not found")
	return nil
}

// GetPgxPool returns the Pgxpool component or panics if not found
func (b *ComponentBuilder) GetPgxPool() *pgxpoolcomponent.PgxPoolConnector {
	if cmp, ok := b.components[PgxpoolConnector].(*pgxpoolcomponent.PgxPoolConnector); ok {
		return cmp
	}

	log.Fatal("Pgxpool component not found")
	return nil
}

// GetRedis returns the Redis component or panics if not found
func (b *ComponentBuilder) GetRedis() *rediscomponent.RedisConnector {
	if cmp, ok := b.components[RedisComponent].(*rediscomponent.RedisConnector); ok {
		return cmp
	}

	log.Fatal("Redis component not found")
	return nil
}
