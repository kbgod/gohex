package user

import (
	"context"

	"app/internal/domain/dto"
	"app/internal/domain/entity"
	"app/internal/domain/port"
	"app/internal/types"
)

type Service struct {
	userRepo port.UserRepository
}

func NewService(userRepo port.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) Create(ctx context.Context, input dto.CreateUser) (*entity.User, error) {
	user := entity.NewUser(input.Username)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id types.ID) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
