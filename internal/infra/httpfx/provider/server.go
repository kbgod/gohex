package provider

import (
	"fmt"

	"app/config"
	"app/internal/infra/httpfx/handler"
	"app/pkg/httpserver"

	fiberzerolog "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func NewServer(
	cfg *config.Config, logger *zerolog.Logger,
) (*fiber.App, error) {
	app, err := httpserver.New(cfg.HTTP, handler.ErrorHandler)
	if err != nil {
		return nil, fmt.Errorf("new http server: %w", err)
	}

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Fields: []string{
			fiberzerolog.FieldLatency,
			fiberzerolog.FieldStatus,
			fiberzerolog.FieldMethod,
			fiberzerolog.FieldURL,
			fiberzerolog.FieldIP,
			fiberzerolog.FieldRequestID,
			fiberzerolog.FieldError,
		},
		GetLogger: func(ctx fiber.Ctx) zerolog.Logger {
			return *logger
		},
	}))

	return app, nil
}
