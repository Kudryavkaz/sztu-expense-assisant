package model

import "github.com/Kudryavkaz/sztuea-api/internal/resource/database"

type User struct {
	BaseModel
	Account      string `gorm:"type:varchar(127);not null;unique" json:"account"`
	Password     string `gorm:"type:varchar(255);not null" json:"password"`
	SztuAccount  string `gorm:"type:varchar(63)" json:"sztu_account"`
	SztuPassword string `gorm:"type:varchar(63)" json:"sztu_password"`
	Cookie       string `gorm:"type:varchar(255)" json:"cookie"`
}

func (u *User) Create() (err error) {
	err = database.DB.Create(u).Error

	return
}

func GetUserByAccount(account string) (user User, err error) {
	err = database.DB.Where("account = ?", account).First(&user).Error

	return
}

func GetUserByID(id uint) (user User, err error) {
	err = database.DB.Where("id = ?", id).First(&user).Error

	return
}

func UpdateAccountByID(id uint, user User) (err error) {
	err = database.DB.Model(&User{}).Where("id = ?", id).Updates(user).Error

	return
}
