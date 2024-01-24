package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/controller"
	"github.com/muhafs/go-fiber-gorm/initializer"
)

func Setup(app *fiber.App) {
	// setup static
	app.Static("/", initializer.Root+"/public/assets")

	// set prefix for api v1
	v1 := app.Group("/api/v1")

	// routes for user api
	v1.Get("/user", controller.GetListUser)
	v1.Post("/user", controller.CreateUser)
	v1.Get("/user/:id", controller.GetUser)
	v1.Put("/user/:id", controller.UpdateUser)
	v1.Delete("/user/:id", controller.DeleteUser)
}
