package invoker

import (
	"app/config"
	"app/pkg/tz"

	"github.com/pkg/errors"
)

func SetupTimezone(cfg *config.Config) error {
	return errors.Wrap(tz.Setup(cfg.Time), "setup timezone")
}
