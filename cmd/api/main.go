package main

import (
	"app/config"
	"app/internal/infra/httpfx"

	"go.uber.org/fx"
)

// @title						gohex API
// @version					1.0
// @description				gohex – simple Go framework to build APIs
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Provide your Bearer token in the format: 'Bearer {token}'
func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	fx.New(httpfx.CreateApp(&cfg)).Run()
}
