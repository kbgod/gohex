package entity

import (
	"time"

	"app/internal/types"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
}

func NewUser(username string) *User {
	return &User{
		ID:       types.NewID(),
		Username: username,
	}
}
