package model

import (
	"claude-code-relay/common"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
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

// DetailedStatsResult 详细统计结果
type DetailedStatsResult struct {
	TotalRequests            int64   `json:"total_requests"`              // 总请求数
	TotalInputTokens         int64   `json:"total_input_tokens"`          // 总输入tokens
	TotalOutputTokens        int64   `json:"total_output_tokens"`         // 总输出tokens
	TotalCacheReadTokens     int64   `json:"total_cache_read_tokens"`     // 总缓存读取tokens
	TotalCacheCreationTokens int64   `json:"total_cache_creation_tokens"` // 总缓存创建tokens
	TotalTokens              int64   `json:"total_tokens"`                // 总tokens数
	TotalCost                float64 `json:"total_cost"`                  // 总费用
	InputCost                float64 `json:"input_cost"`                  // 输入费用
	OutputCost               float64 `json:"output_cost"`                 // 输出费用
	CacheWriteCost           float64 `json:"cache_write_cost"`            // 缓存写入费用
	CacheReadCost            float64 `json:"cache_read_cost"`             // 缓存读取费用
	AvgDuration              float64 `json:"avg_duration"`                // 平均响应时间
	StreamRequests           int64   `json:"stream_requests"`             // 流式请求数
	StreamPercent            float64 `json:"stream_percent"`              // 流式请求比例
}

// StatsQueryRequest 统计查询请求
type StatsQueryRequest struct {
	UserID        *uint      `form:"user_id"`        // 用户ID筛选
	AccountID     *uint      `form:"-"`              // 账号ID筛选(内部转换后使用)
	ApiKeyID      *uint      `form:"-"`              // API Key ID筛选(内部转换后使用)
	AccountFilter string     `form:"account_filter"` // 账号筛选（ID或邮箱/名称）
	ApiKeyFilter  string     `form:"api_key_filter"` // API Key筛选（ID或秘钥值）
	ModelName     string     `form:"model_name"`     // 模型名称筛选
	StartTime     *time.Time `form:"-"`              // 开始时间(不从form绑定)
	EndTime       *time.Time `form:"-"`              // 结束时间(不从form绑定)
}

// TrendDataItem 趋势数据项
type TrendDataItem struct {
	Date         string  `json:"date"`          // 日期
	Requests     int64   `json:"requests"`      // 请求数
	Tokens       int64   `json:"tokens"`        // tokens数
	Cost         float64 `json:"cost"`          // 费用
	AvgDuration  float64 `json:"avg_duration"`  // 平均响应时间
	CacheTokens  int64   `json:"cache_tokens"`  // 缓存tokens
	InputTokens  int64   `json:"input_tokens"`  // 输入tokens
	OutputTokens int64   `json:"output_tokens"` // 输出tokens
}

// StatsResponse 统计响应结果
type StatsResponse struct {
	Summary   *DetailedStatsResult `json:"summary"`    // 汇总统计
	TrendData []TrendDataItem      `json:"trend_data"` // 趋势数据
}

// LogFilters 日志查询过滤条件
type LogFilters struct {
	UserID    *uint      `json:"user_id"`    // 用户ID筛选
	AccountID *uint      `json:"account_id"` // 账号ID筛选
	ApiKeyID  *uint      `json:"api_key_id"` // API Key ID筛选
	ModelName *string    `json:"model_name"` // 模型名称筛选
	IsStream  *bool      `json:"is_stream"`  // 是否流式请求筛选
	StartTime *time.Time `json:"start_time"` // 开始时间
	EndTime   *time.Time `json:"end_time"`   // 结束时间
	MinCost   *float64   `json:"min_cost"`   // 最小费用
	MaxCost   *float64   `json:"max_cost"`   // 最大费用
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

// GetLogsWithFilters 根据过滤条件获取日志列表
func GetLogsWithFilters(filters *LogFilters, page, limit int) ([]Log, int64, error) {
	var logs []Log
	var total int64

	// 构建查询条件
	query := DB.Model(&Log{})
	countQuery := DB.Model(&Log{})

	// 应用过滤条件
	if filters != nil {
		// 用户ID筛选
		if filters.UserID != nil {
			query = query.Where("user_id = ?", *filters.UserID)
			countQuery = countQuery.Where("user_id = ?", *filters.UserID)
		}

		// 账号ID筛选
		if filters.AccountID != nil {
			query = query.Where("account_id = ?", *filters.AccountID)
			countQuery = countQuery.Where("account_id = ?", *filters.AccountID)
		}

		// API Key ID筛选
		if filters.ApiKeyID != nil {
			query = query.Where("api_key_id = ?", *filters.ApiKeyID)
			countQuery = countQuery.Where("api_key_id = ?", *filters.ApiKeyID)
		}

		// 模型名称筛选
		if filters.ModelName != nil {
			query = query.Where("model_name = ?", *filters.ModelName)
			countQuery = countQuery.Where("model_name = ?", *filters.ModelName)
		}

		// 是否流式请求筛选
		if filters.IsStream != nil {
			query = query.Where("is_stream = ?", *filters.IsStream)
			countQuery = countQuery.Where("is_stream = ?", *filters.IsStream)
		}

		// 时间范围筛选
		if filters.StartTime != nil {
			query = query.Where("created_at >= ?", *filters.StartTime)
			countQuery = countQuery.Where("created_at >= ?", *filters.StartTime)
		}

		if filters.EndTime != nil {
			query = query.Where("created_at <= ?", *filters.EndTime)
			countQuery = countQuery.Where("created_at <= ?", *filters.EndTime)
		}

		// 费用范围筛选
		if filters.MinCost != nil {
			query = query.Where("total_cost >= ?", *filters.MinCost)
			countQuery = countQuery.Where("total_cost >= ?", *filters.MinCost)
		}

		if filters.MaxCost != nil {
			query = query.Where("total_cost <= ?", *filters.MaxCost)
			countQuery = countQuery.Where("total_cost <= ?", *filters.MaxCost)
		}
	}

	// 先获取总数
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Preload("User").Preload("ApiKey").
		Offset(offset).Limit(limit).
		Order("id DESC"). // 按ID倒序，新记录在前
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
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

// GetDetailedStats 获取详细统计数据
func GetDetailedStats(req *StatsQueryRequest) (*DetailedStatsResult, error) {
	var stats DetailedStatsResult

	// 构建基础查询
	query := DB.Model(&Log{})

	// 应用过滤条件
	query = applyStatsFilters(query, req)

	// 计算时间范围
	startTime, endTime := calculateTimeRange(req)
	query = query.Where("created_at >= ? AND created_at <= ?", startTime, endTime)

	// 查询统计数据
	var result struct {
		TotalRequests            int64
		TotalInputTokens         int64
		TotalOutputTokens        int64
		TotalCacheReadTokens     int64
		TotalCacheCreationTokens int64
		TotalCost                float64
		InputCost                float64
		OutputCost               float64
		CacheWriteCost           float64
		CacheReadCost            float64
		AvgDuration              float64
		StreamRequests           int64
	}

	err := query.Select(
		"COUNT(*) as total_requests",
		"SUM(input_tokens) as total_input_tokens",
		"SUM(output_tokens) as total_output_tokens",
		"SUM(cache_read_input_tokens) as total_cache_read_tokens",
		"SUM(cache_creation_input_tokens) as total_cache_creation_tokens",
		"SUM(total_cost) as total_cost",
		"SUM(input_cost) as input_cost",
		"SUM(output_cost) as output_cost",
		"SUM(cache_write_cost) as cache_write_cost",
		"SUM(cache_read_cost) as cache_read_cost",
		"AVG(duration) as avg_duration",
		"SUM(CASE WHEN is_stream = true THEN 1 ELSE 0 END) as stream_requests",
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// 填充结果
	stats.TotalRequests = result.TotalRequests
	stats.TotalInputTokens = result.TotalInputTokens
	stats.TotalOutputTokens = result.TotalOutputTokens
	stats.TotalCacheReadTokens = result.TotalCacheReadTokens
	stats.TotalCacheCreationTokens = result.TotalCacheCreationTokens
	stats.TotalTokens = result.TotalInputTokens + result.TotalOutputTokens + result.TotalCacheReadTokens + result.TotalCacheCreationTokens
	stats.TotalCost = result.TotalCost
	stats.InputCost = result.InputCost
	stats.OutputCost = result.OutputCost
	stats.CacheWriteCost = result.CacheWriteCost
	stats.CacheReadCost = result.CacheReadCost
	stats.AvgDuration = result.AvgDuration
	stats.StreamRequests = result.StreamRequests

	// 计算流式请求比例
	if stats.TotalRequests > 0 {
		stats.StreamPercent = float64(stats.StreamRequests) / float64(stats.TotalRequests) * 100
	}

	return &stats, nil
}

// GetTrendData 获取趋势数据
func GetTrendData(req *StatsQueryRequest) ([]TrendDataItem, error) {
	// 构建基础查询
	query := DB.Model(&Log{})

	// 应用过滤条件
	query = applyStatsFilters(query, req)

	// 计算时间范围
	startTime, endTime := calculateTimeRange(req)
	query = query.Where("created_at >= ? AND created_at <= ?", startTime, endTime)

	// 根据时间范围自动选择分组方式
	var groupBy string

	// 计算时间跨度（天数）
	daysDiff := int(endTime.Sub(startTime).Hours() / 24)

	if daysDiff <= 7 {
		// 7天以内：按天分组
		groupBy = "DATE(created_at)"
	} else if daysDiff <= 60 {
		// 60天以内：按天分组
		groupBy = "DATE(created_at)"
	} else {
		// 60天以上：按月分组
		groupBy = "DATE_FORMAT(created_at, '%Y-%m')"
	}

	var trendData []TrendDataItem

	rows, err := query.Select(
		groupBy+" as date_group",
		"COUNT(*) as requests",
		"SUM(input_tokens + output_tokens + cache_read_input_tokens + cache_creation_input_tokens) as tokens",
		"SUM(total_cost) as cost",
		"AVG(duration) as avg_duration",
		"SUM(cache_read_input_tokens + cache_creation_input_tokens) as cache_tokens",
		"SUM(input_tokens) as input_tokens",
		"SUM(output_tokens) as output_tokens",
	).Group(groupBy).Order(groupBy).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item TrendDataItem
		var dateGroup string
		err := rows.Scan(
			&dateGroup,
			&item.Requests,
			&item.Tokens,
			&item.Cost,
			&item.AvgDuration,
			&item.CacheTokens,
			&item.InputTokens,
			&item.OutputTokens,
		)
		if err != nil {
			return nil, err
		}

		// 格式化日期显示
		item.Date = dateGroup

		trendData = append(trendData, item)
	}

	return trendData, nil
}

// applyStatsFilters 应用统计查询过滤条件
func applyStatsFilters(query *gorm.DB, req *StatsQueryRequest) *gorm.DB {
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.AccountID != nil {
		query = query.Where("account_id = ?", *req.AccountID)
	}
	if req.ApiKeyID != nil {
		query = query.Where("api_key_id = ?", *req.ApiKeyID)
	}
	if req.ModelName != "" {
		query = query.Where("model_name = ?", req.ModelName)
	}
	return query
}

// calculateTimeRange 计算时间范围
func calculateTimeRange(req *StatsQueryRequest) (time.Time, time.Time) {
	// 如果提供了具体的开始和结束时间，直接使用（时间区间选择器）
	if req.StartTime != nil && req.EndTime != nil {
		return *req.StartTime, *req.EndTime
	}

	// 否则显示当天数据（无时间区间选择时的默认行为）
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	return startTime, endTime
}

// resolveFilters 解析筛选条件，将字符串筛选转换为ID筛选
func resolveFilters(req *StatsQueryRequest) error {
	// 解析账号筛选
	if req.AccountFilter != "" {
		// 先尝试作为ID解析
		if accountID, err := strconv.ParseUint(req.AccountFilter, 10, 32); err == nil {
			id := uint(accountID)
			req.AccountID = &id
		} else {
			// 作为邮箱或名称查询
			var account Account
			err := DB.Where("email = ? OR name = ?", req.AccountFilter, req.AccountFilter).First(&account).Error
			if err == nil {
				req.AccountID = &account.ID
			}
			// 如果找不到账号，不报错，只是不会匹配到任何记录
		}
	}

	// 解析API Key筛选
	if req.ApiKeyFilter != "" {
		// 先尝试作为ID解析
		if apiKeyID, err := strconv.ParseUint(req.ApiKeyFilter, 10, 32); err == nil {
			id := uint(apiKeyID)
			req.ApiKeyID = &id
		} else {
			// 作为秘钥值查询（通过key字段）
			var apiKey ApiKey
			err := DB.Where("`key` = ?", req.ApiKeyFilter).First(&apiKey).Error
			if err == nil {
				req.ApiKeyID = &apiKey.ID
			}
			// 如果找不到API Key，不报错，只是不会匹配到任何记录
		}
	}

	return nil
}

// GetCompleteStats 获取完整统计数据（汇总+趋势）
func GetCompleteStats(req *StatsQueryRequest) (*StatsResponse, error) {
	// 先解析筛选条件
	if err := resolveFilters(req); err != nil {
		return nil, err
	}

	// 获取汇总统计
	summary, err := GetDetailedStats(req)
	if err != nil {
		return nil, err
	}

	// 获取趋势数据
	trendData, err := GetTrendData(req)
	if err != nil {
		return nil, err
	}

	return &StatsResponse{
		Summary:   summary,
		TrendData: trendData,
	}, nil
}
