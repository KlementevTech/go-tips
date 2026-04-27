package sqlite

import (
	"context"
	"database/sql"
	"log/slog"
)

type txCtxKey struct{}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithinTransaction(ctx context.Context, txFn func(ctx context.Context) error) (err error) {
	if _, ok := ctx.Value(txCtxKey{}).(*sql.Tx); ok {
		return txFn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, nil)
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

	if _, err = tx.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		return err
	}

	err = txFn(ctxWithTx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func rollback(ctx context.Context, tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		slog.Default().ErrorContext(ctx, "rolling back transaction", slog.Any("error", err))
	}
}
