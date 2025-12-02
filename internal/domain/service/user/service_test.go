package user

import (
	"context"
	"errors"
	"testing"

	"app/internal/domain/dto"
	"app/internal/domain/entity"
	"app/internal/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	inputDTO := dto.CreateUser{Username: "testuser"}
	errRepoFailed := errors.New("repository failed")

	testCases := []struct {
		name        string
		setupMock   func(m *mocks.MockUserRepository)
		expectedErr error
	}{
		{
			name: "Success",
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, user *entity.User) error {
						assert.Equal(t, inputDTO.Username, user.Username)
						assert.NotEqual(t, uuid.Nil, user.ID)
						return nil
					},
				).Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "Repo Error",
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().Create(ctx, gomock.Any()).Return(errRepoFailed).Times(1)
			},
			expectedErr: errRepoFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			tc.setupMock(mockUserRepo)

			service := NewService(mockUserRepo)
			user, err := service.Create(ctx, inputDTO)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, inputDTO.Username, user.Username)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mockID := uuid.New()
	mockUser := &entity.User{
		ID:       mockID,
		Username: "founduser",
	}
	errUserNotFound := errors.New("user not found")

	testCases := []struct {
		name         string
		inputID      uuid.UUID
		setupMock    func(m *mocks.MockUserRepository)
		expectedUser *entity.User
		expectedErr  error
	}{
		{
			name:    "Success",
			inputID: mockID,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetByID(ctx, mockID).Return(mockUser, nil).Times(1)
			},
			expectedUser: mockUser,
			expectedErr:  nil,
		},
		{
			name:    "Not Found",
			inputID: mockID,
			setupMock: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetByID(ctx, mockID).Return(nil, errUserNotFound).Times(1)
			},
			expectedUser: nil,
			expectedErr:  errUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			tc.setupMock(mockUserRepo)

			service := NewService(mockUserRepo)
			user, err := service.GetByID(ctx, tc.inputID)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedUser, user)
		})
	}
}
