package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"

	"github.com/dinesh/user-api/config"
	"github.com/dinesh/user-api/internal/handler"
	"github.com/dinesh/user-api/internal/logger"
	"github.com/dinesh/user-api/internal/repository"
	"github.com/dinesh/user-api/internal/routes"
	"github.com/dinesh/user-api/internal/service"
)

func main() {
	logger.Init()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	// Wire everything together
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if ok := false; !ok {
				_ = e
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	routes.Setup(app, userHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Log.Sugar().Infof("server starting on %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatal("server error:", err)
	}
}
