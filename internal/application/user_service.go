package application

import (
	"context"

	"app/internal/domain/entity"
	"app/internal/domain/repository"
	"app/internal/dto"
	"app/internal/types"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(ctx context.Context, input dto.CreateUser) (*entity.User, error) {
	user := entity.NewUser(input.Username)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id types.ID) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
