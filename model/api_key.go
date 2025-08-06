package model

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ApiKey struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Key       string         `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt *Time          `json:"expires_at" gorm:"type:timestamp"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	GroupID   int            `json:"group_id" gorm:"default:0;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	CreatedAt Time           `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt Time           `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateApiKeyRequest struct {
	Name      string `json:"name" binding:"required"`
	Key       string `json:"key"`
	ExpiresAt *Time  `json:"expires_at"`
	Status    int    `json:"status" binding:"oneof=1 2"`
	GroupID   int    `json:"group_id"`
}

type UpdateApiKeyRequest struct {
	Name      string `json:"name"`
	ExpiresAt *Time  `json:"expires_at"`
	Status    *int   `json:"status"`
	GroupID   int    `json:"group_id"`
}

type ApiKeyListResult struct {
	ApiKeys []ApiKey `json:"api_keys"`
	Total   int64    `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

func (a *ApiKey) TableName() string {
	return "api_keys"
}

func (a *ApiKey) BeforeCreate(tx *gorm.DB) error {
	if a.Key == "" {
		key, err := generateApiKey()
		if err != nil {
			return err
		}
		a.Key = key
	}
	return nil
}

func generateApiKey() (string, error) {
	bytes := make([]byte, 30)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sk-%x", bytes)[:30], nil
}

func CreateApiKey(apiKey *ApiKey) error {
	apiKey.ID = 0
	return DB.Create(apiKey).Error
}

func GetApiKeyById(id uint, userID uint) (*ApiKey, error) {
	var apiKey ApiKey
	err := DB.Where("id = ? AND user_id = ?", id, userID).First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetApiKeyByKey 根据API Key获取
func GetApiKeyByKey(key string) (*ApiKey, error) {
	var apiKey ApiKey
	err := DB.Where("key = ? AND status = 1", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if apiKey.ExpiresAt != nil && time.Time(*apiKey.ExpiresAt).Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}

	return &apiKey, nil
}

func UpdateApiKey(apiKey *ApiKey) error {
	return DB.Save(apiKey).Error
}

func DeleteApiKey(id uint) error {
	return DB.Delete(&ApiKey{}, id).Error
}

// GetApiKeys 分页获取API Keys
func GetApiKeys(page, limit int, userID uint, groupID *uint) ([]ApiKey, int64, error) {
	var apiKeys []ApiKey
	var total int64

	query := DB.Model(&ApiKey{}).Where("user_id = ?", userID)
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&apiKeys).Error
	if err != nil {
		return nil, 0, err
	}

	return apiKeys, total, nil
}
