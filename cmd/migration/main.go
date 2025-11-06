package main

import (
	"context"
	"errors"
	"os"
	"time"

	"app/config"
	"app/database/migrations"
	"app/pkg/logger"
	"app/pkg/logger/adapter/zerogoose"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewMigration()
	if err != nil {
		panic(err)
	}
	log, err := logger.New(cfg.Logger)
	if err != nil {
		panic(err)
	}

	utc, err := time.LoadLocation(time.UTC.String())
	if err != nil {
		log.Fatal().Err(err).Msg("time.LoadLocation")
	}
	time.Local = utc

	if len(os.Args) < 2 {
		log.Fatal().Msg("command is required")
	}

	command := os.Args[1]
	var arguments []string
	if len(os.Args) > 2 {
		arguments = os.Args[2:]
	}

	db, err := goose.OpenDBWithDriver("pgx", cfg.Postgres.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("goose.OpenDBWithDriver")
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal().Err(err).Msg("db.Close")
		}
	}()

	goose.SetTableName("migrations")
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(zerogoose.NewLogger(log))

	dir := "."
	if command == "create" {
		dir = "./database/migrations"
		if len(arguments) == 1 {
			arguments = append(arguments, "sql")
		}
	}

	if err := goose.RunContext(ctx, command, db, dir, arguments...); err != nil {
		if errors.Is(err, goose.ErrNoMigrationFiles) {
			log.Info().Msg("migrator - no migrations found")

			return
		}
		if errors.Is(err, goose.ErrAlreadyApplied) {
			log.Info().Msg("migrator - no changes")

			return
		}
		log.Fatal().Err(err).Msg("goose.RunContext")
	}
}
