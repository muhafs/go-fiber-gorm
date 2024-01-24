package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/initializer"
	"github.com/muhafs/go-fiber-gorm/route"
)

var appURL string

func init() {
	// load env file data
	config, err := initializer.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	// connect to database
	initializer.ConnectDB(&config)

	appURL = config.AppURL
}

func main() {
	// create app
	app := fiber.New()

	// init routes
	route.Setup(app)

	// start server
	app.Listen(appURL)
}
