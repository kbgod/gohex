package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger/v2"

	_ "app/docs"
)

func ApplyRoutes(app *fiber.App, userHandler *UserHandler) {
	app.Get("/docs/*", swagger.HandlerDefault)

	app.Post("/users", userHandler.Create)
	app.Get("/users/:id", userHandler.GetByID)
}
