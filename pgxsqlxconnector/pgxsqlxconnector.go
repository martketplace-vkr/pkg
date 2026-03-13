package pgxsqlxconnector

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

func NewDB(cfg Config) (result *sqlx.DB, err error) {
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	driverName := "pgx"

	// TODO: implement sqlhooks later.
	//if cfg.LogQuery || cfg.Extra.EnableMetrics {
	//	var hooks []sqlhooks.Hooks
	//
	//	if cfg.LogQuery {
	//		hooks = append(hooks, NewLogHook(&cfg))
	//	}
	//	if cfg.Extra.EnableMetrics {
	//		hooks = append(hooks, NewMetricHook(&cfg.Extra, cfg.DBName))
	//	}
	//
	//	composite := NewCompositeHook(hooks...)
	//	driverName = "sqlhooks_" + cfg.DBName
	//	sql.Register(driverName, sqlhooks.Wrap(&stdlib.Driver{}, composite))
	//}

	result, err = sqlx.Connect(
		driverName,
		connectionURL,
	)
	if err != nil {
		return nil, err
	}

	if cfg.Extra.EnableMetrics {
		prometheus.MustRegister(NewPgSqlxStatsCollector(result, cfg.DBName))
	}

	result.SetMaxOpenConns(cfg.Settings.MaxOpenConns)
	result.SetConnMaxLifetime(cfg.Settings.ConnMaxLifetime * time.Second)
	result.SetMaxIdleConns(cfg.Settings.MaxIdleConns)
	result.SetConnMaxIdleTime(cfg.Settings.ConnMaxIdleTime * time.Second)

	return result, nil
}
