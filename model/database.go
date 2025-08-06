package model

import (
	"claude-scheduler/common"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/scheduler.db"
	}

	// 确保数据目录存在
	if err := os.MkdirAll("./data", 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&User{},
		&Task{},
		&ApiLog{},
	)
	if err != nil {
		return err
	}

	common.SysLog("Database initialized successfully")
	return nil
}

func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}