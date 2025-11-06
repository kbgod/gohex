package httpserver

import (
	"time"

	"github.com/gofiber/fiber/v3"
	recoverMiddleware "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
)

type Config struct {
	Host            string        `env:"HTTP_HOST" envDefault:":8080"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"60s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"60s"`
	IdleTimeout     time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	ProxyHeaderName string        `env:"HTTP_PROXY_HEADER_NAME" envDefault:""`
	TrustProxy      bool          `env:"HTTP_TRUST_PROXY" envDefault:"false"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT_SECONDS" envDefault:"5s"`
	AppName         string        `env:"HTTP_APP_NAME" envDefault:"go-hex"`
	BodyLimit       int           `env:"HTTP_BODY_LIMIT" envDefault:"10485760"` // 10 MB
	TrustedProxies  []string      `env:"HTTP_TRUSTED_PROXIES" envSeparator:","`
}

func New(cfg Config, errorHandler fiber.ErrorHandler) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.ShutdownTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ProxyHeader:  cfg.ProxyHeaderName,
		TrustProxy:   cfg.TrustProxy,
		AppName:      cfg.AppName,
		BodyLimit:    cfg.BodyLimit,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies: cfg.TrustedProxies,
		},
		ErrorHandler: errorHandler,
	})

	app.Use(recoverMiddleware.New(recoverMiddleware.Config{
		EnableStackTrace: true,
	}))

	app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.Must(uuid.NewV7()).String()
		},
	}))

	return app, nil
}
