package handler

import "app/internal/core"

type Handler struct {
	app *core.Application
}

func NewHandler(app *core.Application) *Handler {
	return &Handler{app: app}
}
