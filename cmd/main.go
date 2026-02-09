package main

import (
	"log"
	"tenant-Dynamin-DB/internals/config"
	"tenant-Dynamin-DB/internals/handlers"
	"tenant-Dynamin-DB/internals/repository"
	"tenant-Dynamin-DB/internals/routes"
	"tenant-Dynamin-DB/internals/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			repository.NewUserRepository,
			service.NewUserService,
			handlers.NewHandlers,
			fiber.New,
		),
		fx.Invoke(
			func(app *fiber.App, h *handlers.Handlers, cfg *config.Config) {
				routes.SetUPRouter(app, h, cfg)
				log.Println("Server Running on port :", cfg.APP_PORT)
				app.Listen(":" + cfg.APP_PORT)
			}),
	).Run()
}
