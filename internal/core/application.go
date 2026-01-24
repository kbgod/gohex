package core

import (
	"app/internal/core/port"
)

type Application struct {
	UserService port.UserService
}

func NewApplication(userService port.UserService) *Application {
	return &Application{
		UserService: userService,
	}
}
