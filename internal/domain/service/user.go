//go:generate mockgen -source=$GOFILE -destination=../../mocks/$GOFILE -package=mocks
package service

import (
	"context"

	"app/internal/domain/entity"
	"app/internal/dto"
	"app/internal/types"
)

type UserService interface {
	Create(ctx context.Context, input dto.CreateUser) (*entity.User, error)
	GetByID(ctx context.Context, id types.ID) (*entity.User, error)
}
