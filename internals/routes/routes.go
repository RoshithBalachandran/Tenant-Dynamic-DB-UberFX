package routes

import (
	"tenant-Dynamin-DB/internals/config"
	"tenant-Dynamin-DB/internals/handlers"
	"tenant-Dynamin-DB/internals/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUPRouter(app *fiber.App, h *handlers.Handlers, cfg *config.Config) {
	app.Post("/register", h.Registration)
	app.Post("/login", h.LoginRequest)

	protect := app.Group("/api", middleware.AuthMiddleware(cfg.JWT_SECRET))
	protect.Get("/all", h.ListAll)

	protect.Get("/profile", h.Profile)
}
