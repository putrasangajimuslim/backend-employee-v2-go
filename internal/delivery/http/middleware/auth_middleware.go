// internal/delivery/http/middleware/auth_middleware.go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourname/backend-employee-v2-go/pkg/response"
	"github.com/yourname/backend-employee-v2-go/internal/usecase"
)

type AuthMiddleware struct {
	userUseCase *usecase.UserUseCase
}

func NewAuthMiddleware(userUseCase *usecase.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{userUseCase: userUseCase}
}

func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("Authorization header is required", nil))
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("Invalid authorization header format", nil))
	}

	token := parts[1]

	userID, err := m.userUseCase.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error("Invalid or expired token", nil))
	}

	c.Locals("userID", userID)
	return c.Next()
}

func GetUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}
