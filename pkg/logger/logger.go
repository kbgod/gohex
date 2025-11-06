package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Format   string        `env:"LOG_FORMAT" envDefault:"json"`
	LogLevel zerolog.Level `env:"LOG_LEVEL" envDefault:"info"`
}

func New(cfg Config) (*zerolog.Logger, error) {
	logger, err := newLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("newLogger: %w", err)
	}

	zerolog.DefaultContextLogger = logger

	return logger, nil
}

func newLogger(cfg Config) (*zerolog.Logger, error) {
	var writer io.Writer
	switch cfg.Format {
	case "console":
		writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	case "json":
		writer = os.Stderr
	default:
		return nil, fmt.Errorf("unknown log format: %s", cfg.Format)
	}

	logger := zerolog.New(writer).
		Level(cfg.LogLevel).
		With().
		Timestamp().
		Logger()

	return &logger, nil
}
