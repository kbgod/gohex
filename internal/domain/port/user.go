package port

import (
	"context"

	"app/internal/domain/dto"
	"app/internal/domain/entity"
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
