package router

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// set prefix for api
	api := app.Group("api")
	// set prefix for api version 1
	v1 := api.Group("v1")

	// setup auth routes
	AuthRoutes(app)

	// setup user routes
	UserRoutes(v1)
}
