package core

import (
	"app/internal/core/port"
)

type Application struct {
	UserService port.UserService
}

func New(userService port.UserService) *Application {
	return &Application{
		UserService: userService,
	}
}
