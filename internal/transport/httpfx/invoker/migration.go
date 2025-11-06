package invoker

import (
	"errors"
	"fmt"

	"app/config"
	"app/database/migrations"
	"app/pkg/logger/adapter/zerogoose"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

func RunMigrations(cfg *config.Config, logger *zerolog.Logger) error {
	log := logger.With().Str("component", "migrator").Logger()

	db, err := goose.OpenDBWithDriver("pgx", cfg.Postgres.DSN())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("db.Close")
		}
	}()

	goose.SetTableName("migrations")
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(zerogoose.NewLogger(&log))

	if err = goose.Up(db, "."); err != nil {
		if errors.Is(err, goose.ErrNoMigrationFiles) {
			log.Info().Msg("migrator - no migrations found")

			return nil
		} else if errors.Is(err, goose.ErrAlreadyApplied) {
			log.Info().Msg("migrator - no changes")

			return nil
		}

		return err
	}

	return nil
}
