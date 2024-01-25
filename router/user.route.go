package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/controller"
)

func UserRoutes(v1 fiber.Router) {
	// set prefix for user
	user := v1.Group("user")

	// routes for user api
	user.Get("/", controller.GetListUser)
	user.Post("/", controller.CreateUser)
	user.Get("/:id", controller.GetUser)
	user.Put("/:id", controller.UpdateUser)
	user.Delete("/:id", controller.DeleteUser)
}
