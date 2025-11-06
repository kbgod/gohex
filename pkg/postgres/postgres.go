package postgres

import (
	"context"
	"fmt"

	"app/pkg/logger/adapter/pgxtracer"
	"app/pkg/postgres/tracer"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(cfg Config) (*pgxpool.Pool, error) {
	pgxCFG, err := pgxpool.ParseConfig(cfg.PGXDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx dsn: %w", err)
	}

	if cfg.QueryDebug {
		pgxCFG.ConnConfig.Tracer = tracer.NewLogTracer(pgxtracer.NewAdapter(cfg.SlowQueryThreshold))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCFG)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	return pool, nil
}
