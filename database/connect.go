package database

import (
	"fmt"
	"log"

	"github.com/muhafs/go-fiber-gorm/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(env *config.Env) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)
	dbConnection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't connect to database:", err.Error())
	}

	DB = dbConnection

	log.Println("Running Migrations")
	MigrateUp(DB)

	log.Println("🚀 Connected Successfully to the Database")
}
