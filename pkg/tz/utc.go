package tz

import (
	"fmt"
	"time"
)

type Config struct {
	TZ string `env:"TZ" envDefault:"UTC"`
}

func Setup(cfg Config) error {
	tz, err := time.LoadLocation(cfg.TZ)
	if err != nil {
		return fmt.Errorf("time.LoadLocation: %w", err)
	}

	time.Local = tz

	return nil
}
