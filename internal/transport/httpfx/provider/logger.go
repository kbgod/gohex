package provider

import (
	"app/config"
	"app/pkg/logger"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.Config) (*zerolog.Logger, error) {
	return logger.New(cfg.Logger)
}
