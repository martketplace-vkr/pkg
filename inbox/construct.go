package inbox

import (
	"time"

	ik "github.com/martketplace-vkr/pkg/inbox/internal/provider/kafka"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/inbox/internal/storage/postgres"

	pgxManager "github.com/avito-tech/go-transaction-manager/pgxv5"
	sqlxManager "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

func NewKafkaProvider(
	cfg ik.MessageReceiverConfig,
	cgFactory func(groupID string) (sarama.ConsumerGroup, error),
) EventsProvider {
	return ik.New(cfg, cgFactory)
}

// NewSqlxStorage creates a new storage with sqlx.
// a fabric method to hide the implementation details of the storage in internal.
func NewSqlxStorage(
	db *sqlx.DB,
	ctxGetter *sqlxManager.CtxGetter,
	eventLockTimeout time.Duration,
	processAttempts int,
	disableSequentialEntriesSelection bool,
) EventsStorage {
	return postgres.NewSqlxStorage(db, ctxGetter, eventLockTimeout, processAttempts, disableSequentialEntriesSelection)
}

// NewPgxStorage creates a new storage with pgx.
// a fabric method to hide the implementation details of the storage in internal.
func NewPgxStorage(
	db *pgxpool.Pool,
	ctxGetter *pgxManager.CtxGetter,
	eventLockTimeout time.Duration,
	processAttempts int,
	disableSequentialEntriesSelection bool,
) EventsStorage {
	return postgres.NewPgxPoolStorage(db, ctxGetter, eventLockTimeout, processAttempts, disableSequentialEntriesSelection)
}
