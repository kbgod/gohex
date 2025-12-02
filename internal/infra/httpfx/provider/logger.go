package provider

import (
	"fmt"

	"app/config"
	"app/pkg/logger"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *config.Config) (*zerolog.Logger, error) {
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	return log, nil
}
