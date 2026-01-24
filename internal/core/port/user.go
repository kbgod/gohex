package port

import (
	"context"

	"app/internal/core/dto"
	"app/internal/core/entity"
	"app/internal/types"
)

type UserService interface {
	Create(ctx context.Context, input dto.CreateUser) (*entity.User, error)
	GetByID(ctx context.Context, id types.ID) (*entity.User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id types.ID) (*entity.User, error)
}
