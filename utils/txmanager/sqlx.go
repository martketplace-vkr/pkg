package txmanager

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jmoiron/sqlx"
)

type sqlxManager struct {
	db        *sqlx.DB
	manager   *manager.Manager
	ctxGetter *trmsqlx.CtxGetter
}

func NewSqlx(
	db *sqlx.DB,
) *sqlxManager {
	return &sqlxManager{
		db:        db,
		manager:   manager.Must(trmsqlx.NewDefaultFactory(db)),
		ctxGetter: trmsqlx.DefaultCtxGetter,
	}
}

func (m *sqlxManager) Do(ctx context.Context, fn Fn) error {
	return m.manager.Do(ctx, fn)
}

func (m *sqlxManager) GetTxOrDb(ctx context.Context) trmsqlx.Tr {
	return m.ctxGetter.DefaultTrOrDB(ctx, m.db)
}
