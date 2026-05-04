package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txCtxKey struct{}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (m *TxManager) WithinTx(ctx context.Context, txFn func(ctx context.Context) error) (err error) {
	if _, ok := ctx.Value(txCtxKey{}).(pgx.Tx); ok {
		return txFn(ctx)
	}

	var opts pgx.TxOptions

	tx, err := m.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, txCtxKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			rollback(ctxWithTx, tx)
			panic(p)
		} else if err != nil {
			rollback(ctxWithTx, tx)
		}
	}()

	err = txFn(ctxWithTx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func rollback(ctx context.Context, tx pgx.Tx) {
	err := tx.Rollback(ctx)
	if err != nil {
		slog.Default().ErrorContext(ctx, "rolling back transaction", "error", err)
	}
}
