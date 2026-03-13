package pgxpoolcomponent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxpoolconnctor "github.com/martketplace-vkr/pkg/pgxpoolconnector"
	"github.com/martketplace-vkr/pkg/utils/migrate"
)

const (
	cmpName = "pgxpool"
)

type PgxPoolConnector struct {
	*pgxpool.Pool
	cfg Config
}

func New(cfg Config) *PgxPoolConnector {
	pool, err := pgxpoolconnctor.NewPgxPool(cfg.Config)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}

	return &PgxPoolConnector{
		Pool: pool,
		cfg:  cfg,
	}
}

func (p *PgxPoolConnector) Start(ctx context.Context) (err error) {
	if err = p.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping: %v", err)
	}

	if p.cfg.MigrateConfig != nil {
		migration, err := migrate.Migrate(
			p.cfg.MigrateConfig.FS,
			p.cfg.MigrateConfig.MigPath,
			createDSN(p.cfg.Config),
		)
		if err != nil {
			return fmt.Errorf("failed to migrate: %v", err)
		}

		err = migration.Up()
		if err != nil {
			return fmt.Errorf("failed to up migrations: %v", err)
		}
	}

	return err
}

func (p *PgxPoolConnector) Stop(_ context.Context) (err error) {
	p.Close()

	return err
}

func (p *PgxPoolConnector) GetStartTimeout() time.Duration {
	return p.cfg.StartTimeout.Duration
}

func (p *PgxPoolConnector) GetStopTimeout() time.Duration {
	return p.cfg.StopTimeout.Duration
}

func (p *PgxPoolConnector) GetShutdownDelay() time.Duration {
	return p.cfg.ShutdownDelay.Duration
}

func (p *PgxPoolConnector) GetName() string {
	return cmpName
}

func createDSN(cfg pgxpoolconnctor.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
		cfg.SSLMode,
	)
}
