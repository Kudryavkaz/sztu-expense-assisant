package model

import (
	"errors"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/resource/database"
	"gorm.io/gorm"
)

type Expense struct {
	BaseModel
	// UserID         uint    `gorm:"not null"`
	// User           User    `gorm:"foreignKey:UserID"`
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

func GetLatestExpenseBySno(sno string) (timeStamp int64, err error) {
	expense := Expense{}
	err = database.DB.Where("sno = ?", sno).Order("finish_time desc").First(&expense).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		date, _ := time.Parse("2006-01-02", "2024-08-01")
		timeStamp = date.Unix()
		err = nil
	} else {
		timeStamp = expense.FinishTime / 1000
	}

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

func (e *Expense) GetExpensesByTimeRange(startTime int64, endTime int64) (expenses Expenses, err error) {
	err = database.DB.Where("finish_time >= ? AND finish_time <= ?", startTime, endTime).Where(e).Find(&expenses).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (e *Expense) GetExpensesByPage(page int, perPage int) (expenses Expenses, err error) {
	offset := (page - 1) * perPage
	err = database.DB.Where(e).Limit(perPage).Offset(offset).Find(&expenses).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (e *Expense) GetTotalExpense() (total int64, err error) {
	err = database.DB.Model(&Expense{}).Where(e).Count(&total).Error

	return
}

type ExpenseTimeLine struct {
	EventDate string  `json:"event_date"`
	Amount    float64 `json:"amount"`
}

func (e *Expense) GetExpenseTimeLine(startTime int64, endTime int64) (expenseTimeLines []*ExpenseTimeLine, err error) {
	expenseTimeLines = make([]*ExpenseTimeLine, 0)
	// duration := (float64)(endTime-startTime) / (1000 * 24 * 60 * 60)
	// if duration > 30 {
	// 	// week
	// 	err = database.DB.Debug().
	// 		Model(&Expense{}).
	// 		Select(`STR_TO_DATE(CONCAT(YEARWEEK(FROM_UNIXTIME(finish_time/1000), 1), ' Monday'), '%X%V %W') as event_date, SUM(amount) AS amount`).
	// 		Where("amount < 0 AND finish_time >= ? AND finish_time <= ?", startTime, endTime).
	// 		Where(e).
	// 		Group("event_date").
	// 		Find(&expenseTimeLines).Error
	// } else if duration > 1 {
	// 	// day
	// 	err = database.DB.Debug().
	// 		Model(&Expense{}).
	// 		Select(`DATE(FROM_UNIXTIME(finish_time/1000)) as event_date, SUM(amount) AS amount`).
	// 		Where("amount < 0 AND finish_time >= ? AND finish_time <= ?", startTime, endTime).
	// 		Where(e).
	// 		Group("event_date").
	// 		Find(&expenseTimeLines).Error
	// } else {
	// 	// hour
	// 	err = database.DB.Debug().
	// 		Model(&Expense{}).
	// 		Select(`DATE_FORMAT(FROM_UNIXTIME(finish_time/1000), '%Y-%m-%d %H:00:00') as event_date, SUM(amount) AS amount`).
	// 		Where("amount < 0 AND finish_time >= ? AND finish_time <= ?", startTime, endTime).
	// 		Where(e).
	// 		Group("event_date").
	// 		Find(&expenseTimeLines).Error
	// }
	err = database.DB.Debug().
		Model(&Expense{}).
		Select(`DATE(FROM_UNIXTIME(finish_time/1000)) as event_date, -SUM(amount) AS amount`).
		Where("amount < 0 AND tran_method_desc != '后台人工'").
		Where(e).
		Group("event_date").
		Find(&expenseTimeLines).Error
	return
}
