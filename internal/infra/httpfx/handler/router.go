package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger/v2"

	_ "app/docs"
)

func ApplyRoutes(app *fiber.App, handler *Handler) {
	app.Get("/docs/*", swagger.HandlerDefault)

	app.Post("/users", handler.CreateUser)
	app.Get("/users/:id", handler.GetUserByID)
}
