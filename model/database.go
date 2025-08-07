package model

import (
	"claude-code-relay/common"
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	// MySQL 数据库连接配置
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = ""
	}

	database := os.Getenv("MYSQL_DATABASE")
	if database == "" {
		database = "claude_code_relay"
	}

	// 先连接到MySQL服务器（不指定数据库）
	adminDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port)

	adminDB, err := gorm.Open(mysql.Open(adminDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %v", err)
	}

	// 创建数据库（如果不存在）
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", database)
	if err := adminDB.Exec(createDBSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// 关闭管理连接
	adminSqlDB, _ := adminDB.DB()
	adminSqlDB.Close()

	// 构建应用数据库 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	// 连接到应用数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to application database: %v", err)
	}

	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置最大打开连接数（默认100）
	maxOpenConns := getIntEnv("MYSQL_MAX_OPEN_CONNS", 100)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// 设置最大空闲连接数（默认10）
	maxIdleConns := getIntEnv("MYSQL_MAX_IDLE_CONNS", 10)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// 设置连接最大生存时间（默认1小时）
	maxLifetimeMinutes := getIntEnv("MYSQL_MAX_LIFETIME_MINUTES", 60)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeMinutes) * time.Minute)

	// 设置连接最大空闲时间（默认30分钟）
	maxIdleTimeMinutes := getIntEnv("MYSQL_MAX_IDLE_TIME_MINUTES", 30)
	sqlDB.SetConnMaxIdleTime(time.Duration(maxIdleTimeMinutes) * time.Minute)

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&User{},
		&Task{},
		&ApiLog{},
		&Account{},
		&Group{},
		&ApiKey{},
		&Log{},
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

// getIntEnv 获取环境变量的整型值，如果不存在或无效则返回默认值
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}
