package model

import (
	"claude-code-relay/common"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ApiKey struct {
	ID                            uint           `json:"id" gorm:"primaryKey"`
	Name                          string         `json:"name" gorm:"type:varchar(100);not null"`
	Key                           string         `json:"key" gorm:"type:varchar(100);uniqueIndex;not null"`
	ExpiresAt                     *Time          `json:"expires_at" gorm:"type:datetime"`
	Status                        int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	GroupID                       int            `json:"group_id" gorm:"default:0;index"`
	UserID                        uint           `json:"user_id" gorm:"not null;index"`
	TodayUsageCount               int            `json:"today_usage_count" gorm:"default:0;comment:今日使用次数"`
	TodayInputTokens              int            `json:"today_input_tokens" gorm:"default:0;comment:今日输入tokens"`
	TodayOutputTokens             int            `json:"today_output_tokens" gorm:"default:0;comment:今日输出tokens"`
	TodayCacheReadInputTokens     int            `json:"today_cache_read_input_tokens" gorm:"default:0;comment:今日缓存读取输入tokens"`
	TodayCacheCreationInputTokens int            `json:"today_cache_creation_input_tokens" gorm:"default:0;comment:今日缓存创建输入tokens"`
	TodayTotalCost                float64        `json:"today_total_cost" gorm:"default:0;comment:今日使用总费用(USD)"`
	ModelRestriction              string         `json:"model_restriction" gorm:"type:text;comment:模型限制,逗号分隔"`
	DailyLimit                    float64        `json:"daily_limit" gorm:"default:0;comment:日限额(美元),0表示不限制"`
	LastUsedTime                  *Time          `json:"last_used_time" gorm:"comment:最后使用时间;type:datetime"`
	CreatedAt                     Time           `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt                     Time           `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt                     gorm.DeletedAt `json:"-" gorm:"index"`
	// 关联查询
	Group *Group `json:"group" gorm:"-"`

	// 最近一周统计数据（不存储到数据库，运行时计算）
	WeeklyCost  float64 `json:"weekly_cost" gorm:"-"`  // 最近一周使用费用
	WeeklyCount int64   `json:"weekly_count" gorm:"-"` // 最近一周使用次数
}

type CreateApiKeyRequest struct {
	Name             string  `json:"name" binding:"required"`
	Key              string  `json:"key"`
	ExpiresAt        *Time   `json:"expires_at"`
	Status           int     `json:"status" binding:"oneof=1 2"`
	GroupID          int     `json:"group_id"`
	ModelRestriction string  `json:"model_restriction"`
	DailyLimit       float64 `json:"daily_limit"`
}

type UpdateApiKeyRequest struct {
	Name             string   `json:"name"`
	ExpiresAt        *Time    `json:"expires_at"`
	Status           *int     `json:"status"`
	GroupID          *int     `json:"group_id"`
	ModelRestriction *string  `json:"model_restriction"`
	DailyLimit       *float64 `json:"daily_limit"`
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

	// 如果有分组ID，查询分组信息
	if apiKey.GroupID > 0 {
		var group Group
		if err := DB.Where("id = ? AND user_id = ?", apiKey.GroupID, userID).First(&group).Error; err == nil {
			apiKey.Group = &group
		}
	}

	return &apiKey, nil
}

// GetApiKeyByKey 根据API Key获取（带缓存）
func GetApiKeyByKey(key string) (*ApiKey, error) {
	// 先尝试从缓存获取
	if common.RDB != nil {
		cacheKey := fmt.Sprintf("api_key:%s", key)
		cachedData, err := common.RDB.Get(context.Background(), cacheKey).Result()
		if err == nil {
			var apiKey ApiKey
			if json.Unmarshal([]byte(cachedData), &apiKey) == nil {
				// 检查缓存的数据是否过期
				if apiKey.ExpiresAt != nil && time.Time(*apiKey.ExpiresAt).Before(time.Now()) {
					// 缓存的数据已过期，删除缓存
					common.RDB.Del(context.Background(), cacheKey)
					return nil, gorm.ErrRecordNotFound
				}
				return &apiKey, nil
			}
		}
	}

	// 缓存未命中，从数据库查询
	var apiKey ApiKey
	err := DB.Where("`key` = ? AND status = 1", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if apiKey.ExpiresAt != nil && time.Time(*apiKey.ExpiresAt).Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}

	// 存储到缓存（5分钟）
	if common.RDB != nil {
		cacheKey := fmt.Sprintf("api_key:%s", key)
		cachedData, err := json.Marshal(apiKey)
		if err == nil {
			common.RDB.Set(context.Background(), cacheKey, cachedData, 5*time.Minute)
		}
	}

	return &apiKey, nil
}

// GetApiKeyByKeyForUpdate 根据API Key获取（无缓存，专用于统计更新）
func GetApiKeyByKeyForUpdate(key string) (*ApiKey, error) {
	var apiKey ApiKey
	err := DB.Where("`key` = ? AND status = 1", key).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if apiKey.ExpiresAt != nil && time.Time(*apiKey.ExpiresAt).Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}

	return &apiKey, nil
}

// ClearApiKeyCache 清理API Key缓存
func ClearApiKeyCache(key string) {
	if common.RDB != nil {
		cacheKey := fmt.Sprintf("api_key:%s", key)
		common.RDB.Del(context.Background(), cacheKey)
	}
}

func UpdateApiKey(apiKey *ApiKey) error {
	err := DB.Save(apiKey).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteApiKey(id uint) error {
	// 先获取API Key信息用于清理缓存
	var apiKey ApiKey
	if err := DB.First(&apiKey, id).Error; err == nil {
		defer ClearApiKeyCache(apiKey.Key)
	}

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

	// 批量查询分组信息
	groupIDs := make(map[int]bool)
	for _, apiKey := range apiKeys {
		if apiKey.GroupID > 0 {
			groupIDs[apiKey.GroupID] = true
		}
	}

	if len(groupIDs) > 0 {
		var ids []int
		for id := range groupIDs {
			ids = append(ids, id)
		}

		var groups []Group
		DB.Where("id IN ? AND user_id = ?", ids, userID).Find(&groups)

		// 创建分组映射
		groupMap := make(map[int]*Group)
		for i := range groups {
			groupMap[int(groups[i].ID)] = &groups[i]
		}

		// 为每个API Key设置对应的分组信息
		for i := range apiKeys {
			if apiKeys[i].GroupID > 0 {
				if group, exists := groupMap[apiKeys[i].GroupID]; exists {
					apiKeys[i].Group = group
				}
			}
		}
	}

	// 批量查询最近一周的统计数据
	if err := setWeeklyStatsForApiKeys(apiKeys); err != nil {
		// 如果查询失败，只记录错误但不影响主要功能
		// 统计字段将保持默认值0
	}

	return apiKeys, total, nil
}

// setWeeklyStatsForApiKeys 为API Key列表设置最近一周的统计数据
func setWeeklyStatsForApiKeys(apiKeys []ApiKey) error {
	if len(apiKeys) == 0 {
		return nil
	}

	// 计算本周开始时间（周一00:00:00）到现在
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // 如果是周日，调整为7
		weekday = 7
	}

	// 计算本周一的日期
	weekStart := now.AddDate(0, 0, -(weekday - 1))
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	// 提取所有API Key ID
	apiKeyIDs := make([]uint, len(apiKeys))
	for i, apiKey := range apiKeys {
		apiKeyIDs[i] = apiKey.ID
	}

	// 查询最近一周的统计数据
	type WeeklyStats struct {
		ApiKeyID   uint    `json:"api_key_id"`
		TotalCost  float64 `json:"total_cost"`
		TotalCount int64   `json:"total_count"`
	}

	var weeklyStats []WeeklyStats
	err := DB.Table("logs").
		Select("api_key_id, SUM(total_cost) as total_cost, COUNT(*) as total_count").
		Where("api_key_id IN ? AND created_at >= ? AND created_at <= ?", apiKeyIDs, weekStart, now).
		Group("api_key_id").
		Scan(&weeklyStats).Error

	if err != nil {
		return err
	}

	// 创建统计数据映射
	statsMap := make(map[uint]WeeklyStats)
	for _, stat := range weeklyStats {
		statsMap[stat.ApiKeyID] = stat
	}

	// 为每个API Key设置统计数据
	for i := range apiKeys {
		if stat, exists := statsMap[apiKeys[i].ID]; exists {
			apiKeys[i].WeeklyCost = stat.TotalCost
			apiKeys[i].WeeklyCount = stat.TotalCount
		} else {
			apiKeys[i].WeeklyCost = 0
			apiKeys[i].WeeklyCount = 0
		}
	}

	return nil
}
