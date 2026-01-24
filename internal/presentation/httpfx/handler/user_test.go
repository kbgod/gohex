package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"app/internal/core"
	"app/internal/core/dto"
	"app/internal/core/entity"
	"app/internal/mocks"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Create(t *testing.T) {
	mockUser := &entity.User{
		ID:       uuid.Must(uuid.NewV7()),
		Username: "testuser",
	}

	dtoInput := dto.CreateUser{
		Username: "testuser",
	}

	serviceErr := errors.New("something went wrong")

	successBody, _ := json.Marshal(newUserResponse(mockUser))

	testCases := []struct {
		name            string
		body            []byte
		setupMock       func(m *mocks.MockUserService)
		expectedStatus  int
		expectedBody    string
		bodyMustBeEmpty bool
	}{
		{
			name: "Success",
			body: []byte(`{"username": "testuser"}`),
			setupMock: func(m *mocks.MockUserService) {
				m.EXPECT().Create(gomock.Any(), dtoInput).Return(mockUser, nil).Times(1)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   string(successBody),
		},
		{
			name: "Binding Error - Invalid JSON",
			body: []byte(`{"username": 123}`),
			setupMock: func(m *mocks.MockUserService) {
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   `cannot unmarshal number into Go struct field`,
		},
		{
			name: "Service Error",
			body: []byte(`{"username": "testuser"}`),
			setupMock: func(m *mocks.MockUserService) {
				m.EXPECT().Create(gomock.Any(), dtoInput).Return(nil, serviceErr).Times(1)
			},
			expectedStatus:  fiber.StatusInternalServerError,
			bodyMustBeEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocks.NewMockUserService(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockUserService)
			}

			app := core.NewApplication(mockUserService)

			handler := NewHandler(app)

			router := fiber.New(fiber.Config{
				ErrorHandler: ErrorHandler,
			})
			router.Post("/users", handler.CreateUser)

			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "core/json")

			resp, err := router.Test(req)
			require.NoError(t, err)

			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if tc.bodyMustBeEmpty {
				assert.Empty(t, bodyBytes)

				return
			} else {
				assert.Contains(t, string(bodyBytes), tc.expectedBody)
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	mockUserID := uuid.Must(uuid.NewV7())
	mockUser := &entity.User{
		ID:       mockUserID,
		Username: "get-user",
	}
	errSomething := errors.New("something went wrong")

	successBody, _ := json.Marshal(newUserResponse(mockUser))

	testCases := []struct {
		name            string
		userIDString    string
		setupMock       func(m *mocks.MockUserService)
		expectedStatus  int
		expectedBody    string
		bodyMustBeEmpty bool
	}{
		{
			name:         "Success",
			userIDString: mockUserID.String(),
			setupMock: func(m *mocks.MockUserService) {
				m.EXPECT().GetByID(gomock.Any(), mockUserID).Return(mockUser, nil).Times(1)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   string(successBody),
		},
		{
			name:         "Binding Error - Invalid UUID",
			userIDString: "not-a-uuid",
			setupMock: func(m *mocks.MockUserService) {
			},
			expectedStatus: fiber.StatusUnprocessableEntity,
			expectedBody:   "invalid UUID",
		},
		{
			name:         "Service Error - Not Found",
			userIDString: mockUserID.String(),
			setupMock: func(m *mocks.MockUserService) {
				m.EXPECT().GetByID(gomock.Any(), mockUserID).Return(nil, errSomething).Times(1)
			},
			expectedStatus:  fiber.StatusInternalServerError,
			bodyMustBeEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocks.NewMockUserService(ctrl)
			if tc.setupMock != nil {
				tc.setupMock(mockUserService)
			}

			app := core.NewApplication(mockUserService)

			handler := NewHandler(app)

			router := fiber.New(fiber.Config{
				ErrorHandler: ErrorHandler,
			})
			router.Get("/users/:id", handler.GetUserByID)

			req := httptest.NewRequest("GET", "/users/"+tc.userIDString, nil)

			resp, err := router.Test(req)
			require.NoError(t, err)

			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if tc.bodyMustBeEmpty {
				assert.Empty(t, bodyBytes)

				return
			} else {
				assert.Contains(t, string(bodyBytes), tc.expectedBody)
			}
		})
	}
}
