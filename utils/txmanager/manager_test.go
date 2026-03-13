package txmanager

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSqlxManager_GetTxOrDb(t *testing.T) {
	db := &sqlx.DB{}
	mgr := NewSqlx(db)

	ctx := context.Background()
	tr := mgr.GetTxOrDb(ctx)
	assert.Equal(t, db, tr)
}

func TestPgxManager_GetTxOrDb(t *testing.T) {
	conn := &pgxpool.Conn{}
	mgr := NewPgxPool(conn)

	ctx := context.Background()
	tr := mgr.GetTxOrDb(ctx)
	assert.Equal(t, conn, tr)
}
