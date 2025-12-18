package provider

import (
	"app/pkg/transactor"
	pgxTransactor "app/pkg/transactor/pgx"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type TransactorResult struct {
	fx.Out

	Transactor transactor.Transactor
	DBGetter   pgxTransactor.DBGetter
}

func NewPgxTransactor(pool *pgxpool.Pool) TransactorResult {
	tx, dbGetter := pgxTransactor.New(pool)

	return TransactorResult{
		Transactor: tx,
		DBGetter:   dbGetter,
	}
}
