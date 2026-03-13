package txmanager

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxManager struct {
	conn      *pgxpool.Conn
	manager   *manager.Manager
	ctxGetter *trmsqlx.CtxGetter
}

func NewPgxPool(
	conn *pgxpool.Conn,
) *pgxManager {
	return &pgxManager{
		conn:      conn,
		manager:   manager.Must(trmsqlx.NewDefaultFactory(conn)),
		ctxGetter: trmsqlx.DefaultCtxGetter,
	}
}

func (m *pgxManager) Do(ctx context.Context, fn Fn) (err error) {
	return m.manager.Do(ctx, fn)
}

func (m *pgxManager) GetTxOrDb(ctx context.Context) trmsqlx.Tr {
	return m.ctxGetter.DefaultTrOrDB(ctx, m.conn)
}
