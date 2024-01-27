package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// model for storing into database
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null;unique"`
	Password  string         `json:"-" gorm:"not null"`
	Avatar    sql.NullString `json:"avatar"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return
}
