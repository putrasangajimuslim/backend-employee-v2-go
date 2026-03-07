// internal/delivery/http/handler/user_handler.go
package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/go-clean-arch/internal/delivery/http/dto"
	"github.com/yourname/go-clean-arch/internal/delivery/http/middleware"
	"github.com/yourname/go-clean-arch/internal/domain/entity"
	"github.com/yourname/go-clean-arch/internal/usecase"
	"github.com/yourname/go-clean-arch/pkg/response"
	"github.com/yourname/go-clean-arch/pkg/validator"
)

// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration"
// @Success 201 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid request body", nil))
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Validation failed", errors))
	}

	user, err := h.userUseCase.Register(c.Context(), usecase.RegisterRequest(req))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error(), nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success("User registered successfully", toUserResponse(user)))
}

// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid request body", nil))
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Validation failed", errors))
	}

	token, err := h.userUseCase.Login(c.Context(), usecase.LoginRequest(req))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("Login successful", dto.LoginResponse{Token: token}))
}

// @Summary Get all users
// @Description Retrieve paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	users, total, err := h.userUseCase.GetAll(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error(), nil))
	}

	if users == nil {
		users = make([]*entity.User, 0)
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = toUserResponse(user)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithMeta("Users retrieved successfully", userResponses, page, limit, total))
}

// @Summary Get user by ID
// @Description Retrieve a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid user ID", nil))
	}

	user, err := h.userUseCase.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("User retrieved successfully", toUserResponse(user)))
}

// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "User update"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid user ID", nil))
	}

	userID := middleware.GetUserID(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(response.Error("You can only update your own profile", nil))
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid request body", nil))
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Validation failed", errors))
	}

	user, err := h.userUseCase.Update(c.Context(), uint(id), usecase.UpdateUserRequest(req))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("User updated successfully", toUserResponse(user)))
}

// @Summary Delete user
// @Description Delete a user account
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid user ID", nil))
	}

	userID := middleware.GetUserID(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(response.Error("You can only delete your own profile", nil))
	}

	if err := h.userUseCase.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("User deleted successfully", nil))
}

func toUserResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}
