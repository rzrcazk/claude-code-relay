package model

import (
	"claude-code-relay/common"
	"os"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	Role      string         `json:"role" gorm:"default:user"`
	CreatedAt Time           `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt Time           `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}

func CreateUser(user *User) error {
	user.ID = 0
	return DB.Create(user).Error
}

func GetUserById(id uint) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

func DeleteUser(id uint) error {
	return DB.Delete(&User{}, id).Error
}

func GetUsers(page, limit int) ([]User, int64, error) {
	var users []User
	var total int64

	err := DB.Model(&User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = DB.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserProfile struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type UserListResult struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

func init() {
	// 创建默认管理员用户
	go func() {
		time.Sleep(1 * time.Second) // 等待数据库初始化完成
		if DB == nil {
			return
		}

		var count int64
		DB.Model(&User{}).Count(&count)
		if count == 0 {
			adminUsername := os.Getenv("DEFAULT_ADMIN_USERNAME")
			if adminUsername == "" {
				adminUsername = "admin"
			}

			adminPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
			if adminPassword == "" {
				adminPassword = "admin123"
			}

			adminEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
			if adminEmail == "" {
				adminEmail = "admin@example.com"
			}

			adminUser := &User{
				Username: adminUsername,
				Email:    adminEmail,
				Password: common.HashPassword(adminPassword),
				Role:     "admin",
				Status:   1,
			}
			if err := CreateUser(adminUser); err == nil {
				common.SysLog("Default admin user created: " + adminUsername + "/" + adminPassword)
			}
		}
	}()
}
