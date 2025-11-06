package config

import (
	"app/pkg/logger"
	"app/pkg/postgres"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type MigrationConfig struct {
	Postgres postgres.Config
	Logger   logger.Config
}

func NewMigration() (MigrationConfig, error) {
	_ = godotenv.Load(".env")

	return env.ParseAs[MigrationConfig]()
}
