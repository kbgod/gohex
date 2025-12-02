package application

import (
	"app/internal/domain/port"
)

type Application struct {
	UserService port.UserService
}

func New(userService port.UserService) *Application {
	return &Application{
		UserService: userService,
	}
}
