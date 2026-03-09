// cmd/api/main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/yourname/backend-employee-v2-go/internal/delivery/http/router"
	"github.com/yourname/backend-employee-v2-go/internal/infrastructure/config"
	"github.com/yourname/backend-employee-v2-go/internal/infrastructure/database"
	repo "github.com/yourname/backend-employee-v2-go/internal/infrastructure/repository"
	"github.com/yourname/backend-employee-v2-go/internal/usecase"
	"github.com/yourname/backend-employee-v2-go/pkg/response"
	"github.com/yourname/backend-employee-v2-go/pkg/validator"
)

func main() {
	cfg := config.Load()

	validator.Init()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repo.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JWTSecret)

	app := fiber.New(fiber.Config{
		AppName:      "backend-employee-v2-go",
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(requestid.New())

	router.SetupRoutes(app, userUseCase)

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(response.Error(err.Error(), nil))
}
