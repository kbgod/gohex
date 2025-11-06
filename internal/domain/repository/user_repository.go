//go:generate mockgen -source=$GOFILE -destination=../../mocks/$GOFILE -package=mocks
package repository

import (
	"context"

	"app/internal/domain/entity"
	"app/internal/types"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id types.ID) (*entity.User, error)
}
