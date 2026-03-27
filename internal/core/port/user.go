package port

import (
	"context"

	"app/internal/core/dto"
	"app/internal/core/entity"
	domainErrors "app/internal/core/error"
	"app/internal/types"
)

var (
	ErrUserNotFound      = domainErrors.New("user not found")
	ErrUserAlreadyExists = domainErrors.New("user already exists")
)

type UserService interface {
	Create(ctx context.Context, input dto.CreateUser) (*entity.User, error)
	GetByID(ctx context.Context, id types.ID) (*entity.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id types.ID) (*entity.User, error)
}
