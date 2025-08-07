package model

import (
	"claude-code-relay/common"
	"errors"
	"strconv"
	"time"
)

// Log 日志记录表 - 记录Claude Code调用的详细日志
type Log struct {
	ID                       string  `json:"id" gorm:"primaryKey;type:varchar(19)"`                     // 雪花算法ID，支持排序
	ModelName                string  `json:"model_name" gorm:"type:varchar(100);not null;index"`        // 模型名称，如claude-3-5-sonnet-20241022
	AccountID                uint    `json:"account_id" gorm:"index"`                                   // 账户ID
	UserID                   uint    `json:"user_id" gorm:"index"`                                      // 用户ID
	ApiKeyID                 uint    `json:"api_key_id" gorm:"index"`                                   // API Key ID
	InputTokens              int     `json:"input_tokens" gorm:"default:0"`                             // 输入tokens数量
	OutputTokens             int     `json:"output_tokens" gorm:"default:0"`                            // 输出tokens数量
	CacheReadInputTokens     int     `json:"cache_read_input_tokens" gorm:"default:0"`                  // 缓存读取输入tokens数量
	CacheCreationInputTokens int     `json:"cache_creation_input_tokens" gorm:"default:0"`              // 缓存创建输入tokens数量
	InputCost                float64 `json:"input_cost" gorm:"default:0"`                               // 输入费用(USD)
	OutputCost               float64 `json:"output_cost" gorm:"default:0"`                              // 输出费用(USD)
	CacheWriteCost           float64 `json:"cache_write_cost" gorm:"default:0"`                         // 缓存写入费用(USD)
	CacheReadCost            float64 `json:"cache_read_cost" gorm:"default:0"`                          // 缓存读取费用(USD)
	TotalCost                float64 `json:"total_cost" gorm:"default:0"`                               // 总费用(USD)
	IsStream                 bool    `json:"is_stream" gorm:"default:false"`                            // 是否为流式输出
	Duration                 int64   `json:"duration"`                                                  // 请求总耗时(毫秒)
	CreatedAt                Time    `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"` // 创建时间

	// 关联关系
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ApiKey ApiKey `json:"api_key,omitempty" gorm:"foreignKey:ApiKeyID"`
}

// LogCreateRequest 创建日志请求结构
type LogCreateRequest struct {
	ModelName                string  `json:"model_name" binding:"required"`
	AccountID                uint    `json:"account_id"`
	UserID                   uint    `json:"user_id" binding:"required"`
	ApiKeyID                 uint    `json:"api_key_id"`
	InputTokens              int     `json:"input_tokens"`
	OutputTokens             int     `json:"output_tokens"`
	CacheReadInputTokens     int     `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int     `json:"cache_creation_input_tokens"`
	InputCost                float64 `json:"input_cost"`
	OutputCost               float64 `json:"output_cost"`
	CacheWriteCost           float64 `json:"cache_write_cost"`
	CacheReadCost            float64 `json:"cache_read_cost"`
	TotalCost                float64 `json:"total_cost"`
	IsStream                 bool    `json:"is_stream"`
	Duration                 int64   `json:"duration"`
}

// LogListResult 日志列表响应结构
type LogListResult struct {
	Logs  []Log `json:"logs"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// LogStatsResult 日志统计结果
type LogStatsResult struct {
	TotalRequests  int64   `json:"total_requests"`
	TotalTokens    int64   `json:"total_tokens"`
	TotalCost      float64 `json:"total_cost"`
	AvgDuration    float64 `json:"avg_duration"`
	StreamRequests int64   `json:"stream_requests"`
	StreamPercent  float64 `json:"stream_percent"`
}

func (l *Log) TableName() string {
	return "logs"
}

// generateSnowflakeID 生成类雪花算法ID (简化版，基于时间戳+递增序列)
// 格式: 时间戳(13位) + 机器ID(2位) + 序列号(4位) = 19位数字字符串
var sequenceNum int64 = 0

func generateSnowflakeID() string {
	timestamp := time.Now().UnixMilli()
	machineID := int64(1) // 简化处理，可配置

	sequenceNum++
	if sequenceNum > 9999 {
		sequenceNum = 1
	}

	// 组合生成19位ID: 时间戳(13位) + 机器ID(2位) + 序列号(4位)
	id := timestamp*1000000 + machineID*10000 + sequenceNum
	return strconv.FormatInt(id, 10)
}

// CreateLog 创建日志记录
func CreateLog(logReq *LogCreateRequest) (*Log, error) {
	log := &Log{
		ID:                       generateSnowflakeID(),
		ModelName:                logReq.ModelName,
		AccountID:                logReq.AccountID,
		UserID:                   logReq.UserID,
		ApiKeyID:                 logReq.ApiKeyID,
		InputTokens:              logReq.InputTokens,
		OutputTokens:             logReq.OutputTokens,
		CacheReadInputTokens:     logReq.CacheReadInputTokens,
		CacheCreationInputTokens: logReq.CacheCreationInputTokens,
		InputCost:                logReq.InputCost,
		OutputCost:               logReq.OutputCost,
		CacheWriteCost:           logReq.CacheWriteCost,
		CacheReadCost:            logReq.CacheReadCost,
		TotalCost:                logReq.TotalCost,
		IsStream:                 logReq.IsStream,
		Duration:                 logReq.Duration,
	}

	err := DB.Create(log).Error
	if err != nil {
		common.SysError("创建日志记录失败: " + err.Error())
		return nil, err
	}

	return log, nil
}

// CreateLogFromTokenUsage 根据TokenUsage创建日志记录
func CreateLogFromTokenUsage(usage *common.TokenUsage, userID, apiKeyID, accountID uint, duration int64, isStream bool) (*Log, error) {
	// 使用费用计算器计算详细费用
	costResult := common.CalculateCost(usage)

	logReq := &LogCreateRequest{
		ModelName:                usage.Model,
		AccountID:                accountID,
		UserID:                   userID,
		ApiKeyID:                 apiKeyID,
		InputTokens:              usage.InputTokens,
		OutputTokens:             usage.OutputTokens,
		CacheReadInputTokens:     usage.CacheReadInputTokens,
		CacheCreationInputTokens: usage.CacheCreationInputTokens,
		InputCost:                costResult.Costs.Input,
		OutputCost:               costResult.Costs.Output,
		CacheWriteCost:           costResult.Costs.CacheWrite,
		CacheReadCost:            costResult.Costs.CacheRead,
		TotalCost:                costResult.Costs.Total,
		IsStream:                 isStream,
		Duration:                 duration,
	}

	return CreateLog(logReq)
}

// GetLogById 根据ID获取日志
func GetLogById(id string) (*Log, error) {
	var log Log
	err := DB.Preload("User").Preload("ApiKey").First(&log, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetLogs 获取日志列表(支持分页)
func GetLogs(page, limit int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	err := DB.Model(&Log{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = DB.Preload("User").Preload("ApiKey").
		Offset(offset).Limit(limit).
		Order("id DESC"). // 按ID倒序，新记录在前
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogsByUser 获取用户的日志记录
func GetLogsByUser(userID uint, page, limit int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	query := DB.Model(&Log{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Preload("User").Preload("ApiKey").
		Offset(offset).Limit(limit).
		Order("id DESC").
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogsByModel 根据模型名称获取日志记录
func GetLogsByModel(modelName string, page, limit int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	query := DB.Model(&Log{}).Where("model_name = ?", modelName)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Preload("User").Preload("ApiKey").
		Offset(offset).Limit(limit).
		Order("id DESC").
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLogStats 获取日志统计信息
func GetLogStats(userID *uint) (*LogStatsResult, error) {
	var stats LogStatsResult

	query := DB.Model(&Log{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// 统计总请求数
	err := query.Count(&stats.TotalRequests).Error
	if err != nil {
		return nil, err
	}

	// 统计其他数据
	var result struct {
		TotalTokens    int64
		TotalCost      float64
		AvgDuration    float64
		StreamRequests int64
	}

	err = query.Select(
		"SUM(input_tokens + output_tokens + cache_read_input_tokens + cache_creation_input_tokens) as total_tokens",
		"SUM(total_cost) as total_cost",
		"AVG(duration) as avg_duration",
		"SUM(CASE WHEN is_stream = true THEN 1 ELSE 0 END) as stream_requests",
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	stats.TotalTokens = result.TotalTokens
	stats.TotalCost = result.TotalCost
	stats.AvgDuration = result.AvgDuration
	stats.StreamRequests = result.StreamRequests

	// 计算流式请求百分比
	if stats.TotalRequests > 0 {
		stats.StreamPercent = float64(stats.StreamRequests) / float64(stats.TotalRequests) * 100
	}

	return &stats, nil
}

// DeleteLogById 删除指定ID的日志记录
func DeleteLogById(id string) error {
	return DB.Delete(&Log{}, "id = ?", id).Error
}

// DeleteLogsByUser 删除指定用户的所有日志记录
func DeleteLogsByUser(userID uint) error {
	return DB.Where("user_id = ?", userID).Delete(&Log{}).Error
}

// DeleteExpiredLogs 删除过期的日志记录
func DeleteExpiredLogs(months int) (int64, error) {
	if months <= 0 {
		return 0, errors.New("月数必须大于0")
	}

	// 计算cutoff时间：当前时间减去指定月数
	cutoffTime := time.Now().AddDate(0, -months, 0)

	// 删除创建时间早于cutoff时间的日志
	result := DB.Where("created_at < ?", cutoffTime).Delete(&Log{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
