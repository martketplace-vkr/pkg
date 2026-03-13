package txmanager

import (
	"context"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
)

type (
	TxManager interface {
		Do(ctx context.Context, fn Fn) (err error)
		GetTxOrDb(ctx context.Context) trmsqlx.Tr
	}
	Fn func(ctx context.Context) error
)
