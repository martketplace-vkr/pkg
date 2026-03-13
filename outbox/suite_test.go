package outbox

import (
	"context"
	"testing"
	"time"

	"github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlxdb "github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/exporters/jaeger"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/martketplace-vkr/pkg/outbox/internal/broker/kafka"
	"github.com/martketplace-vkr/pkg/outbox/internal/storage/postgres"

	"github.com/martketplace-vkr/pkg/kafkaconnector"
	"github.com/martketplace-vkr/pkg/utils/duration"
)

const (
	StorageSqlx     StorageType = "sqlx"
	StoragePostgres StorageType = "pgx"
)

type (
	TestSuite struct {
		suite.Suite
		cfg            Config
		dbPgx          *pgxpool.Pool
		dbSqlx         *sqlxdb.DB
		kafka          kafkaconnector.Client
		jaegerExporter *jaeger.Exporter
		traceProvider  *tracesdk.TracerProvider
		txManager      *manager.Manager
		ctxGetter      *pgxv5.CtxGetter
		outbox         *outbox
	}
	StorageType string
)

var storageType StorageType

func (s *TestSuite) SetupTest() {
	var err error

	s.cfg = Config{
		Metrics: Metrics{
			Namespace: "outbox",
		},
		MessageCleaner: &MessageCleaner{
			// EventExpireTimeDays: 4,
		},
		MessagesScheduler: MessagesScheduler{
			BatchSize: 10,
			ScheduleTime: duration.Seconds{
				time.Second * 10,
			},
		},
	}

	s.dbPgx, err = InitPool()
	s.Require().Nil(err)
	s.Require().NotNil(s.dbPgx)

	s.dbSqlx, err = InitPostgres()
	s.Require().Nil(err)
	s.Require().NotNil(s.dbSqlx)

	s.kafka, err = InitKafka()
	s.Require().Nil(err)
	s.Require().NotNil(s.kafka)

	s.jaegerExporter, err = InitJaeger()
	s.Require().Nil(err)
	s.Require().NotNil(s.jaegerExporter)

	s.traceProvider = InitTraceProvider(s.jaegerExporter)

	var storage Storage
	switch storageType {
	case StorageSqlx:
		s.txManager = manager.Must(sqlx.NewDefaultFactory(s.dbSqlx))
		storage = postgres.NewSqlxStorage(
			s.dbSqlx,
			sqlx.DefaultCtxGetter,
			s.cfg.MessagesScheduler.EventLockTimeout.Duration,
			s.cfg.MessagesScheduler.DisableSequentialEntriesSelection,
		)
	case StoragePostgres:
		s.txManager = manager.Must(pgxv5.NewDefaultFactory(s.dbPgx))
		storage = postgres.NewPgxPoolStorage(
			s.dbPgx,
			pgxv5.DefaultCtxGetter,
			s.cfg.MessagesScheduler.EventLockTimeout.Duration,
			s.cfg.MessagesScheduler.DisableSequentialEntriesSelection,
		)
	}
	brok := kafka.NewSync(s.kafka.NewSyncProducer())

	s.outbox = New(s.cfg, storage, brok, s.txManager, prometheus.NewRegistry(), cron.New())
}

func (s *TestSuite) TearDownSuite() {
	s.dbPgx.Close()

	err := s.traceProvider.Shutdown(context.Background())
	s.Require().Nil(err)

	err = s.jaegerExporter.Shutdown(context.Background())
	s.Require().Nil(err)
}

func TestOutboxTestSuite(t *testing.T) {
	storageTypes := []StorageType{StoragePostgres, StorageSqlx}
	for _, sType := range storageTypes {
		storageType = sType
		t.Run(string(sType), func(t *testing.T) {
			suite.Run(t, new(TestSuite))
		})
	}
}
