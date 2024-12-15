package model

import (
	"errors"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/resource/database"
	"gorm.io/gorm"
)

type Expense struct {
	BaseModel
	UserID         uint    `gorm:"not null"`
	User           User    `gorm:"foreignKey:UserID"`
	Sno            string  `gorm:"not null"`
	TranTypeDesc   string  `gorm:"not null"`
	TranMethodDesc string  `gorm:"not null"`
	OriBalance     float64 `gorm:"not null"`
	Amount         float64 `gorm:"not null"`
	Balance        float64 `gorm:"not null"`
	ItemName       string  `gorm:"not null"`
	FinishTime     int64   `gorm:"not null"`
	StrFinishTime  string  `gorm:"not null"`
	StatusStr      string  `gorm:"not null"`
}

type Expenses []*Expense

func (e *Expense) Create() (err error) {
	err = database.DB.Create(e).Error

	return
}

func (e *Expenses) Create() (err error) {
	err = database.DB.Create(e).Error

	return
}

func GetLatestExpenseByUserID(userID uint) (timeStamp int64, err error) {
	expense := Expense{}
	err = database.DB.Where("user_id = ?", userID).Order("finish_time desc").First(&expense).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		date, _ := time.Parse("2006-01-02", "2024-08-01")
		timeStamp = date.Unix()
		err = nil
	} else {
		timeStamp = expense.FinishTime / 1000
	}

	return
}
