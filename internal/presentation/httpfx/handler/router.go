package handler

import (
	swagger "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"

	_ "app/docs"
)

func ApplyRoutes(app *fiber.App, handler *Handler) {
	app.Get("/docs/*", swagger.HandlerDefault)

	app.Post("/users", handler.CreateUser)
	app.Get("/users/:id", handler.GetUserByID)
}
