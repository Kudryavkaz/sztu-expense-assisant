package database

import (
	"fmt"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/config"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase 数据库初始化
//
//	@author warpzhang
//	@update 2024-11-12 13:17:20
func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.GetString("mysql.user"),
		config.Config.GetString("mysql.password"),
		config.Config.GetString("mysql.host"),
		config.Config.GetString("mysql.port"),
		config.Config.GetString("mysql.database"),
	)

	DB = lo.Must(gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{
		DryRun:         false,
		TranslateError: true,
	}))

	db := lo.Must(DB.DB())

	db.SetMaxIdleConns(10)

	db.SetMaxOpenConns(100)

	db.SetConnMaxLifetime(5 * time.Hour)

	log.Logger().Info("Database connected")
}
