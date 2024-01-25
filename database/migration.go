package database

import (
	"github.com/muhafs/go-fiber-gorm/model/entity"
	"gorm.io/gorm"
)

func MigrateUp(DB *gorm.DB) {
	// migrate Entities into database
	DB.AutoMigrate(
		&entity.User{},
	)
}
