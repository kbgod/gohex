package handler

import (
	"fmt"
	"time"

	"app/internal/core/dto"
	"app/internal/core/entity"
	"app/internal/types"

	"github.com/gofiber/fiber/v3"
)

type createUserRequest struct {
	Username string `json:"username"`
}

type userResponse struct {
	ID        types.ID  `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user *entity.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

// CreateUser
//
//	@Summary		CreateUser a new user
//	@Description	CreateUser a user with a username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserRequest	true	"CreateUser user payload"
//	@Success		200		{object}	userResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/users [post]
func (h *Handler) CreateUser(ctx fiber.Ctx) error {
	req := new(createUserRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return newBindError(err)
	}

	user, err := h.app.UserService.Create(ctx.Context(), dto.CreateUser{
		Username: req.Username,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return ctx.JSON(newUserResponse(user))
}

type getUserByIDRequest struct {
	ID types.ID `uri:"id"`
}

// GetUserByID
//
//	@Summary		Get a user by ID
//	@Description	Retrieve a user by its ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	userResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/users/{id} [get]
func (h *Handler) GetUserByID(ctx fiber.Ctx) error {
	req := new(getUserByIDRequest)
	if err := ctx.Bind().All(req); err != nil {
		return newBindError(err)
	}

	user, err := h.app.UserService.GetByID(ctx.Context(), req.ID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	return ctx.JSON(newUserResponse(user))
}
