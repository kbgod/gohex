package invoker

import (
	"app/config"
	"app/pkg/tz"
)

func SetupTimezone(cfg *config.Config) error {
	return tz.Setup(cfg.Time)
}
