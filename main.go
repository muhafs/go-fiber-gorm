package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/config"
	"github.com/muhafs/go-fiber-gorm/database"
	"github.com/muhafs/go-fiber-gorm/router"
)

var appURL string

func init() {
	// load env file data
	config, err := config.LoadEnv(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	// connect to database
	database.ConnectDB(&config)

	// set appUrl as global var
	appURL = config.AppURL
}

func main() {
	// create fiber app
	app := fiber.New()

	// setup static assets
	app.Static("/", config.Root+"/public")

	// setup routing
	router.Setup(app)

	// start server
	log.Fatal(app.Listen(appURL))
}
