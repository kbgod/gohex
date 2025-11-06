package httpfx

import (
	"app/config"
	"app/internal/application"
	"app/internal/domain/repository"
	"app/internal/domain/service"
	"app/internal/infrastructure/database/postgres"
	"app/internal/transport/httpfx/handler"
	"app/internal/transport/httpfx/invoker"
	"app/internal/transport/httpfx/provider"

	"go.uber.org/fx"
)

func CreateApp(cfg *config.Config) fx.Option {
	return fx.Options(
		// fx.NopLogger,
		fx.Supply(cfg),

		// Provide infrastructure
		fx.Provide(provider.NewLogger),
		fx.Provide(provider.NewPgxPool),
		fx.Provide(provider.NewPgxTransactor),
		fx.Provide(provider.NewServer),

		// Provide ports
		fx.Provide(fx.Annotate(postgres.NewUserRepository, fx.As(new(repository.UserRepository)))),

		// Provide services
		fx.Provide(fx.Annotate(application.NewUserService, fx.As(new(service.UserService)))),

		// Provide http handlers
		fx.Provide(handler.NewUserHandler),

		fx.Invoke(invoker.SetupTimezone),
		fx.Invoke(invoker.RunMigrations),

		fx.Invoke(handler.ApplyRoutes),
		fx.Invoke(invoker.StartHTTPServer),
	)
}
