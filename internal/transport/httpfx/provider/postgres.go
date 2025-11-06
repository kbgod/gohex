package provider

import (
	"context"
	"fmt"

	"app/config"
	"app/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func NewPgxPool(cfg *config.Config, logger *zerolog.Logger, lc fx.Lifecycle) (*pgxpool.Pool, error) {
	pool, err := postgres.NewPgxPool(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return pool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("postgres: closing connection pool")
			pool.Close()

			return nil
		},
	})

	return pool, nil
}
