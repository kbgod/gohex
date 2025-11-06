package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"app/config"
	"app/database/migrations"
	"app/pkg/logger"
	"app/pkg/logger/adapter/zerogoose"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Context context.Context
	Config  config.MigrationConfig
}

type Command struct {
	Name      string
	Arguments []string
	Directory string
}

func main() {
	app, log := initializeApp()

	db := setupDatabase(app.Config, log)

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal().Err(err).Msg("db.Close")
		}
	}()

	configureGoose(log)

	cmd := parseCommand(log)

	runMigrations(app.Context, cmd, db, log)
}

func initializeApp() (App, *zerolog.Logger) {
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

	app := App{
		Context: context.Background(),
		Config:  cfg,
	}

	return app, log
}

func setupDatabase(cfg config.MigrationConfig, log *zerolog.Logger) *sql.DB {
	db, err := goose.OpenDBWithDriver("pgx", cfg.Postgres.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("goose.OpenDBWithDriver")
	}

	return db
}

func configureGoose(log *zerolog.Logger) {
	goose.SetTableName("migrations")
	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(zerogoose.NewLogger(log))
}

func parseCommand(log *zerolog.Logger) Command {
	if len(os.Args) < 2 {
		log.Fatal().Msg("command is required")
	}

	cmd := Command{
		Name:      os.Args[1],
		Directory: ".",
	}

	if len(os.Args) > 2 {
		cmd.Arguments = os.Args[2:]
	}

	if cmd.Name == "create" {
		cmd.Directory = "./database/migrations"

		if len(cmd.Arguments) == 1 {
			cmd.Arguments = append(cmd.Arguments, "sql")
		}
	}

	return cmd
}

func runMigrations(ctx context.Context, cmd Command, db *sql.DB, log *zerolog.Logger) {
	err := goose.RunContext(ctx, cmd.Name, db, cmd.Directory, cmd.Arguments...)
	if err == nil {
		log.Info().Str("command", cmd.Name).Msg("migrator - command executed successfully")

		return
	}

	if errors.Is(err, goose.ErrNoMigrationFiles) {
		log.Info().Msg("migrator - no migrations found")

		return
	}

	if errors.Is(err, goose.ErrAlreadyApplied) {
		log.Info().Msg("migrator - no changes")

		return
	}

	log.Error().Err(err).Msg("goose.RunContext failed")
}
