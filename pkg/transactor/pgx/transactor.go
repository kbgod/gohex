package pgx

import (
	"context"
	"fmt"
)

func New(pool DB) (*Transactor, DBGetter) {
	pgxTxGetter := func(ctx context.Context) pgxTx {
		if tx := txFromContext(ctx); tx != nil {
			return tx
		}

		return pool
	}

	dbGetter := func(ctx context.Context) DB {
		if tx := txFromContext(ctx); tx != nil {
			return tx
		}

		return pool
	}

	return &Transactor{
		pgxTxGetter,
	}, dbGetter
}

type (
	pgxTxGetter func(context.Context) pgxTx
)

type Transactor struct {
	pgxTxGetter
}

func (t *Transactor) Do(ctx context.Context, txFunc func(context.Context) error) error {
	db := t.pgxTxGetter(ctx)

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx) // If rollback fails, there's nothing to do, the transaction will expire by itself
	}()

	txCtx := txToContext(ctx, tx)

	if err := txFunc(txCtx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (t *Transactor) Skip(ctx context.Context) context.Context {
	return context.WithValue(ctx, transactorKey{}, nil)
}

func IsWithinTransaction(ctx context.Context) bool {
	return ctx.Value(transactorKey{}) != nil
}
