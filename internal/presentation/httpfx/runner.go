package httpfx

import (
	"app/config"
	"app/internal/application"
	"app/internal/domain/port"
	"app/internal/domain/service/user"
	"app/internal/infra/repository/postgres"
	"app/internal/presentation/httpfx/handler"
	"app/internal/presentation/httpfx/invoker"
	"app/internal/presentation/httpfx/provider"

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
		fx.Provide(fx.Annotate(postgres.NewUserRepository, fx.As(new(port.UserRepository)))),

		// Provide services
		fx.Provide(fx.Annotate(user.NewService, fx.As(new(port.UserService)))),

		// Provide application
		fx.Provide(application.New),

		// Provide http handlers
		fx.Provide(handler.NewHandler),

		fx.Invoke(invoker.SetupTimezone),
		fx.Invoke(invoker.RunMigrations),

		fx.Invoke(handler.ApplyRoutes),
		fx.Invoke(invoker.StartHTTPServer),
	)
}
