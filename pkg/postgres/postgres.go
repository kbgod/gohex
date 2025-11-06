package postgres

import (
	"context"

	"app/pkg/logger/adapter/pgxtracer"
	"app/pkg/postgres/tracer"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(cfg Config) (*pgxpool.Pool, error) {
	pgxCFG, err := pgxpool.ParseConfig(cfg.PGXDSN())
	if err != nil {
		return nil, err
	}
	if cfg.QueryDebug {
		pgxCFG.ConnConfig.Tracer = tracer.NewLogTracer(pgxtracer.NewAdapter(cfg.SlowQueryThreshold))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCFG)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
