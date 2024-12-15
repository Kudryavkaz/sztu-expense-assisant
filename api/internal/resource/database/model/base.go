package model

import (
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/resource/database"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func InitModels() (err error) {
	err = database.DB.AutoMigrate(
		&User{},
		&Expense{},
	)

	return
}
