package invoker

import (
	"context"

	"app/config"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func StartHTTPServer(
	cfg *config.Config, app *fiber.App, logger *zerolog.Logger, lc fx.Lifecycle, shutdowner fx.Shutdowner,
) error {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := app.Listen(cfg.HTTP.Host); err != nil {
					logger.Error().Err(err).Msg("httpserver.Listen")

					_ = shutdowner.Shutdown()
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info().Msg("httpserver.Stop")

			return app.ShutdownWithTimeout(cfg.HTTP.ShutdownTimeout)
		},
	})

	return nil
}
