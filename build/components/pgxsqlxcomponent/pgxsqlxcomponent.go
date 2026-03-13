package pgxsqlxcomponent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/martketplace-vkr/pkg/pgxsqlxconnector"
	"github.com/martketplace-vkr/pkg/utils/migrate"
)

const (
	cmpName = "postgres"
)

type PgxSqlxConnector struct {
	*sqlx.DB
	cfg Config
}

func New(cfg Config) *PgxSqlxConnector {
	pgSQLDB, err := pgxsqlxconnector.NewDB(cfg.Config)
	if err != nil {
		log.Fatalf("failed to create pgx sqlx: %v", err)
	}

	return &PgxSqlxConnector{
		DB:  pgSQLDB,
		cfg: cfg,
	}
}

func (p *PgxSqlxConnector) Start(_ context.Context) (err error) {
	if err = p.Ping(); err != nil {
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

func (p *PgxSqlxConnector) Stop(_ context.Context) (err error) {
	p.Close()

	return err
}

func (p *PgxSqlxConnector) GetStartTimeout() time.Duration {
	return p.cfg.StartTimeout.Duration
}

func (p *PgxSqlxConnector) GetStopTimeout() time.Duration {
	return p.cfg.StopTimeout.Duration
}

func (p *PgxSqlxConnector) GetShutdownDelay() time.Duration {
	return p.cfg.ShutdownDelay.Duration
}

func (p *PgxSqlxConnector) GetName() string {
	return cmpName
}

func createDSN(cfg pgxsqlxconnector.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
}
