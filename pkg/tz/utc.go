package tz

import (
	"fmt"
	"time"
)

// Config holds the configuration for time zone setup.
type Config struct {
	TZ string `env:"TZ" envDefault:"UTC"`
}

// Setup sets the local time zone based on the provided configuration.
func Setup(cfg Config) error {
	tz, err := time.LoadLocation(cfg.TZ)
	if err != nil {
		return fmt.Errorf("time.LoadLocation: %w", err)
	}

	time.Local = tz

	return nil
}
