package config

import (
	"app/pkg/httpserver"
	"app/pkg/logger"
	"app/pkg/postgres"
	"app/pkg/tz"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP     httpserver.Config
	Postgres postgres.Config
	Logger   logger.Config
	Time     tz.Config
}

func New() (Config, error) {
	_ = godotenv.Load(".env")

	return env.ParseAs[Config]()
}
