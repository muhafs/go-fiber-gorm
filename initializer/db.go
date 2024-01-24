package initializer

import (
	"fmt"
	"log"

	"github.com/muhafs/go-fiber-gorm/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	dbConnection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Couldn't connect to database:", err.Error())
	}

	DB = dbConnection

	log.Println("Running Migrations")
	DB.AutoMigrate(&entity.User{})

	log.Println("🚀 Connected Successfully to the Database")
}
