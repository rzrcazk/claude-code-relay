package service

import (
	"claude-code-relay/common"
	"claude-code-relay/model"
	"errors"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

// CreateLog 创建日志记录
func (s *LogService) CreateLog(req *model.LogCreateRequest) (*model.Log, error) {
	if req.ModelName == "" {
		return nil, errors.New("模型名称不能为空")
	}
	if req.UserID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	log, err := model.CreateLog(req)
	if err != nil {
		return nil, errors.New("创建日志失败: " + err.Error())
	}

	return log, nil
}

// CreateLogFromTokenUsage 根据TokenUsage创建日志记录（推荐使用）
func (s *LogService) CreateLogFromTokenUsage(usage *common.TokenUsage, userID, apiKeyID, accountID uint, duration int64, isStream bool) (*model.Log, error) {
	if usage == nil {
		return nil, errors.New("TokenUsage不能为空")
	}
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	log, err := model.CreateLogFromTokenUsage(usage, userID, apiKeyID, accountID, duration, isStream)
	if err != nil {
		return nil, errors.New("创建日志失败: " + err.Error())
	}

	return log, nil
}

// GetLogById 根据ID获取日志
func (s *LogService) GetLogById(id string) (*model.Log, error) {
	if id == "" {
		return nil, errors.New("日志ID不能为空")
	}

	log, err := model.GetLogById(id)
	if err != nil {
		return nil, errors.New("获取日志失败")
	}

	return log, nil
}

// GetLogs 获取日志列表(支持分页)
func (s *LogService) GetLogs(page, limit int) (*model.LogListResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logs, total, err := model.GetLogs(page, limit)
	if err != nil {
		return nil, errors.New("获取日志列表失败")
	}

	result := &model.LogListResult{
		Logs:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// GetLogsByUser 获取用户的日志记录
func (s *LogService) GetLogsByUser(userID uint, page, limit int) (*model.LogListResult, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logs, total, err := model.GetLogsByUser(userID, page, limit)
	if err != nil {
		return nil, errors.New("获取用户日志失败")
	}

	result := &model.LogListResult{
		Logs:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// GetLogsByModel 根据模型名称获取日志记录
func (s *LogService) GetLogsByModel(modelName string, page, limit int) (*model.LogListResult, error) {
	if modelName == "" {
		return nil, errors.New("模型名称不能为空")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logs, total, err := model.GetLogsByModel(modelName, page, limit)
	if err != nil {
		return nil, errors.New("获取模型日志失败")
	}

	result := &model.LogListResult{
		Logs:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// GetLogStats 获取日志统计信息
func (s *LogService) GetLogStats(userID *uint) (*model.LogStatsResult, error) {
	stats, err := model.GetLogStats(userID)
	if err != nil {
		return nil, errors.New("获取统计信息失败")
	}

	return stats, nil
}

// GetUserLogStats 获取用户日志统计信息
func (s *LogService) GetUserLogStats(userID uint) (*model.LogStatsResult, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	stats, err := model.GetLogStats(&userID)
	if err != nil {
		return nil, errors.New("获取用户统计信息失败")
	}

	return stats, nil
}

// DeleteLogById 删除指定ID的日志记录
func (s *LogService) DeleteLogById(id string) error {
	if id == "" {
		return errors.New("日志ID不能为空")
	}

	err := model.DeleteLogById(id)
	if err != nil {
		return errors.New("删除日志失败")
	}

	return nil
}

// DeleteLogsByUser 删除指定用户的所有日志记录
func (s *LogService) DeleteLogsByUser(userID uint) error {
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}

	err := model.DeleteLogsByUser(userID)
	if err != nil {
		return errors.New("删除用户日志失败")
	}

	return nil
}

// GetLogsWithFilters 根据过滤条件获取日志列表
func (s *LogService) GetLogsWithFilters(filters *model.LogFilters, page, limit int) (*model.LogListResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logs, total, err := model.GetLogsWithFilters(filters, page, limit)
	if err != nil {
		return nil, errors.New("获取日志列表失败: " + err.Error())
	}

	result := &model.LogListResult{
		Logs:  logs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// DeleteExpiredLogs 删除过期的日志记录
func (s *LogService) DeleteExpiredLogs(months int) (int64, error) {
	if months <= 0 {
		return 0, errors.New("保留月数必须大于0")
	}

	deletedCount, err := model.DeleteExpiredLogs(months)
	if err != nil {
		return 0, errors.New("删除过期日志失败: " + err.Error())
	}

	return deletedCount, nil
}

// GetDetailedStats 获取详细统计数据
func (s *LogService) GetDetailedStats(req *model.StatsQueryRequest) (*model.DetailedStatsResult, error) {
	stats, err := model.GetDetailedStats(req)
	if err != nil {
		return nil, errors.New("获取详细统计失败: " + err.Error())
	}

	return stats, nil
}

// GetTrendData 获取趋势数据
func (s *LogService) GetTrendData(req *model.StatsQueryRequest) ([]model.TrendDataItem, error) {
	trendData, err := model.GetTrendData(req)
	if err != nil {
		return nil, errors.New("获取趋势数据失败: " + err.Error())
	}

	return trendData, nil
}

// GetCompleteStats 获取完整统计数据（汇总+趋势）
func (s *LogService) GetCompleteStats(req *model.StatsQueryRequest) (*model.StatsResponse, error) {
	result, err := model.GetCompleteStats(req)
	if err != nil {
		return nil, errors.New("获取完整统计数据失败: " + err.Error())
	}

	return result, nil
}
