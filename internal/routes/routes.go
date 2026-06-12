package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/dinesh/user-api/internal/handler"
	"github.com/dinesh/user-api/internal/middleware"
)

func Setup(app *fiber.App, h *handler.UserHandler) {
	// Global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// User routes
	users := app.Group("/users")
	users.Post("/", h.CreateUser)
	users.Get("/", h.ListUsers)
	users.Get("/:id", h.GetUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}
