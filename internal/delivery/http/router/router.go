// internal/delivery/http/router/router.go
package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yourname/go-clean-arch/internal/delivery/http/handler"
	"github.com/yourname/go-clean-arch/internal/delivery/http/middleware"
	"github.com/yourname/go-clean-arch/internal/usecase"
	"github.com/yourname/go-clean-arch/pkg/response"
)

type RouteConfig struct {
	userUseCase *usecase.UserUseCase
}

func SetupRoutes(app *fiber.App, userUseCase *usecase.UserUseCase) {
	userHandler := handler.NewUserHandler(userUseCase)
	authMiddleware := middleware.NewAuthMiddleware(userUseCase)

	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(response.Success("ok", nil))
	})

	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	users := api.Group("/users")
	users.Use(authMiddleware.Authenticate)
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
}
