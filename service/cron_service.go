package service

import (
	"claude-code-relay/common"
	"claude-code-relay/model"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

// CronService 定时任务服务
type CronService struct {
	cron *cron.Cron
}

// NewCronService 创建定时任务服务实例
func NewCronService() *CronService {
	// 使用带秒的cron解析器
	c := cron.New(cron.WithSeconds())
	return &CronService{cron: c}
}

// Start 启动定时任务
func (s *CronService) Start() {
	// 每天凌晨0点清理统计数据
	_, err := s.cron.AddFunc("0 0 0 * * *", s.resetDailyStats)
	if err != nil {
		log.Printf("Failed to add daily reset cron job: %v", err)
		return
	}

	// 每天凌晨1点清理过期日志
	_, err = s.cron.AddFunc("0 0 1 * * *", s.cleanExpiredLogs)
	if err != nil {
		log.Printf("Failed to add log cleanup cron job: %v", err)
		return
	}

	// 启动定时任务
	s.cron.Start()
	common.SysLog("Cron service started successfully")
}

// Stop 停止定时任务
func (s *CronService) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		common.SysLog("Cron service stopped")
	}
}

// resetDailyStats 重置每日统计数据
func (s *CronService) resetDailyStats() {
	startTime := time.Now()
	common.SysLog("Starting daily stats reset task")

	// 重置Account表的今日统计数据
	err := s.resetAccountStats()
	if err != nil {
		common.SysError("Failed to reset account daily stats: " + err.Error())
	} else {
		common.SysLog("Account daily stats reset successfully")
	}

	// 重置ApiKey表的今日统计数据
	err = s.resetApiKeyStats()
	if err != nil {
		common.SysError("Failed to reset api key daily stats: " + err.Error())
	} else {
		common.SysLog("API Key daily stats reset successfully")
	}

	duration := time.Since(startTime)
	common.SysLog("Daily stats reset task completed in " + duration.String())
}

// resetAccountStats 重置账户今日统计数据
func (s *CronService) resetAccountStats() error {
	result := model.DB.Model(&model.Account{}).Where("1 = 1").Updates(map[string]interface{}{
		"today_usage_count":                 0,
		"today_input_tokens":                0,
		"today_output_tokens":               0,
		"today_cache_read_input_tokens":     0,
		"today_cache_creation_input_tokens": 0,
		"today_total_cost":                  0,
	})

	if result.Error != nil {
		return result.Error
	}

	log.Printf("Reset daily stats for %d accounts", result.RowsAffected)
	return nil
}

// resetApiKeyStats 重置API Key今日统计数据
func (s *CronService) resetApiKeyStats() error {
	result := model.DB.Model(&model.ApiKey{}).Where("1 = 1").Updates(map[string]interface{}{
		"today_usage_count":                 0,
		"today_input_tokens":                0,
		"today_output_tokens":               0,
		"today_cache_read_input_tokens":     0,
		"today_cache_creation_input_tokens": 0,
		"today_total_cost":                  0,
	})

	if result.Error != nil {
		return result.Error
	}

	log.Printf("Reset daily stats for %d api keys", result.RowsAffected)
	return nil
}

// ManualResetStats 手动重置统计数据（用于测试或管理员操作）
func (s *CronService) ManualResetStats() error {
	common.SysLog("Manual daily stats reset triggered")

	err := s.resetAccountStats()
	if err != nil {
		return err
	}

	err = s.resetApiKeyStats()
	if err != nil {
		return err
	}

	common.SysLog("Manual daily stats reset completed")
	return nil
}

// GetNextResetTime 获取下次重置时间
func (s *CronService) GetNextResetTime() time.Time {
	now := time.Now()
	// 计算明天凌晨0点
	tomorrow := now.AddDate(0, 0, 1)
	nextReset := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())
	return nextReset
}

// 全局定时任务服务实例
var GlobalCronService *CronService

// InitCronService 初始化全局定时任务服务
func InitCronService() {
	GlobalCronService = NewCronService()
	GlobalCronService.Start()
}

// StopCronService 停止全局定时任务服务
func StopCronService() {
	if GlobalCronService != nil {
		GlobalCronService.Stop()
	}
}

// cleanExpiredLogs 清理过期日志
func (s *CronService) cleanExpiredLogs() {
	startTime := time.Now()
	common.SysLog("Starting expired logs cleanup task")

	// 从环境变量获取日志保留月数，默认为3个月
	retentionMonths := getLogRetentionMonths()

	logService := NewLogService()
	deletedCount, err := logService.DeleteExpiredLogs(retentionMonths)
	if err != nil {
		common.SysError("Failed to clean expired logs: " + err.Error())
	} else {
		common.SysLog("Cleaned expired logs successfully, deleted " + strconv.FormatInt(deletedCount, 10) + " records (older than " + strconv.Itoa(retentionMonths) + " months)")
	}

	duration := time.Since(startTime)
	common.SysLog("Expired logs cleanup task completed in " + duration.String())
}

// getLogRetentionMonths 从环境变量获取日志保留月数
func getLogRetentionMonths() int {
	monthsStr := os.Getenv("LOG_RETENTION_MONTHS")
	if monthsStr == "" {
		return 3 // 默认保留3个月
	}

	months, err := strconv.Atoi(monthsStr)
	if err != nil || months <= 0 {
		log.Printf("Invalid LOG_RETENTION_MONTHS value: %s, using default value 3", monthsStr)
		return 3
	}

	return months
}

// ManualCleanExpiredLogs 手动清理过期日志（用于测试或管理员操作）
func (s *CronService) ManualCleanExpiredLogs() (int64, error) {
	common.SysLog("Manual expired logs cleanup triggered")

	retentionMonths := getLogRetentionMonths()
	logService := NewLogService()
	deletedCount, err := logService.DeleteExpiredLogs(retentionMonths)
	if err != nil {
		return 0, err
	}

	common.SysLog("Manual expired logs cleanup completed, deleted " + strconv.FormatInt(deletedCount, 10) + " records")
	return deletedCount, nil
}
