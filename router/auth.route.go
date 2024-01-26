package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/controller"
)

func AuthRoutes(app *fiber.App) {
	// routes for auth api
	app.Post("/signin", controller.SignIn)
}
