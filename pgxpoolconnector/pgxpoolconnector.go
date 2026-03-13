package pgxpoolconnctor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

func NewPgxPool(cfg Config) (pool *pgxpool.Pool, err error) {
	connectionURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.SSLMode,
	)

	return NewWithDSN(connectionURL, cfg.Extra)
}

func NewWithDSN(dsn string, extra Extra) (pool *pgxpool.Pool, err error) {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cfg: %v", err)
	}
	dbName := getDbNameFromDSN(dsn)

	poolConfig.MaxConns = extra.MaxOpenConnections
	poolConfig.MinConns = extra.MinOpenConnections

	
	if extra.EnableMonitoring {
		poolConfig.ConnConfig.Tracer = NewPgxTracer(dbName)
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %v", err)
	}

	if extra.EnableMonitoring {
		prometheus.MustRegister(NewPgxPoolStatsCollector(pool, dbName))
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping: %v", err)
	}

	return pool, err
}


