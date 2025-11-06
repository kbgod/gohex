package handler

import (
	"time"

	"app/internal/domain/entity"
	"app/internal/domain/service"
	"app/internal/dto"
	"app/internal/types"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

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

// Create
//
//	@Summary		Create a new user
//	@Description	Create a user with a username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserRequest	true	"Create user payload"
//	@Success		200		{object}	userResponse
//	@Failure		400		{object}	map[string]string
//	@Router			/users [post]
func (h *UserHandler) Create(ctx fiber.Ctx) error {
	req := new(createUserRequest)
	if err := ctx.Bind().JSON(req); err != nil {
		return newBindError(err)
	}

	user, err := h.userService.Create(ctx.Context(), dto.CreateUser{
		Username: req.Username,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(newUserResponse(user))
}

type getUserByIDRequest struct {
	ID types.ID `uri:"id"`
}

// GetByID
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
func (h *UserHandler) GetByID(ctx fiber.Ctx) error {
	req := new(getUserByIDRequest)
	if err := ctx.Bind().All(req); err != nil {
		return newBindError(err)
	}

	user, err := h.userService.GetByID(ctx.Context(), req.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(newUserResponse(user))
}
