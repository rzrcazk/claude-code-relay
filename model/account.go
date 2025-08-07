package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID                            uint           `json:"id" gorm:"primaryKey"`
	Name                          string         `json:"name" gorm:"not null;comment:账号名称"`
	PlatformType                  string         `json:"platform_type" gorm:"not null;comment:平台类型(claude/claude_console)"`
	RequestURL                    string         `json:"request_url" gorm:"comment:请求地址"`
	SecretKey                     string         `json:"-" gorm:"comment:请求秘钥"`
	AccessToken                   string         `json:"-" gorm:"comment:claude的官方token"`
	RefreshToken                  string         `json:"-" gorm:"comment:claude的官方刷新token"`
	ExpiresAt                     int            `json:"expires_at" gorm:"default:0;comment:token过期时间戳"`
	IsMax                         bool           `json:"is_max" gorm:"default:false;comment:是否是max账号"`
	GroupID                       int            `json:"group_id" gorm:"default:0;comment:分组ID"`
	Priority                      int            `json:"priority" gorm:"default:100;comment:优先级(数字越小越高)"`
	Weight                        int            `json:"weight" gorm:"default:100;comment:权重(数字越大越高)"`
	TodayUsageCount               int            `json:"today_usage_count" gorm:"default:0;comment:今日使用次数"`
	TodayInputTokens              int            `json:"today_input_tokens" gorm:"default:0;comment:今日输入tokens"`
	TodayOutputTokens             int            `json:"today_output_tokens" gorm:"default:0;comment:今日输出tokens"`
	TodayCacheReadInputTokens     int            `json:"today_cache_read_input_tokens" gorm:"default:0;comment:今日缓存读取输入tokens"`
	TodayCacheCreationInputTokens int            `json:"today_cache_creation_input_tokens" gorm:"default:0;comment:今日缓存创建输入tokens"`
	TodayTotalCost                float64        `json:"today_total_cost" gorm:"default:0;comment:今日使用总费用(USD)"`
	EnableProxy                   bool           `json:"enable_proxy" gorm:"default:false;comment:是否启用代理"`
	ProxyURI                      string         `json:"proxy_uri" gorm:"comment:代理URI字符串"`
	LastUsedTime                  *Time          `json:"last_used_time" gorm:"comment:最后使用时间;type:timestamp"`
	RateLimitEndTime              *Time          `json:"rate_limit_end_time" gorm:"comment:限流结束时间;type:timestamp"`
	CurrentStatus                 int            `json:"current_status" gorm:"default:1;comment:当前状态(1:正常,2:接口异常,3:账号异常/限流)"`
	ActiveStatus                  int            `json:"active_status" gorm:"default:1;comment:激活状态(1:激活,2:禁用)"`
	UserID                        uint           `json:"user_id" gorm:"not null;comment:所属用户ID"`
	CreatedAt                     Time           `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt                     Time           `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt                     gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联用户
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// 账号列表请求参数
type AccountListRequest struct {
	Page   int   `json:"page" form:"page" binding:"min=1"`
	Limit  int   `json:"limit" form:"limit" binding:"min=1,max=100"`
	UserID *uint `json:"user_id" form:"user_id"`
}

// 账号列表响应结构
type AccountListResponse struct {
	Accounts []Account `json:"accounts"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

// 账号创建请求参数
type CreateAccountRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	PlatformType string `json:"platform_type" binding:"required,oneof=claude claude_console gemini openai"`
	RequestURL   string `json:"request_url"`
	SecretKey    string `json:"secret_key"`
	GroupID      int    `json:"group_id"`
	Priority     int    `json:"priority"`
	Weight       int    `json:"weight" binding:"min=1"`
	EnableProxy  bool   `json:"enable_proxy"`
	ProxyURI     string `json:"proxy_uri"`
	ActiveStatus int    `json:"active_status" binding:"oneof=1 2"`
	IsMax        bool   `json:"is_max"` // 是否是max账号
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int    `json:"expires_at" binding:"min=0"`
}

// 账号更新请求参数
type UpdateAccountRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	PlatformType string `json:"platform_type" binding:"required,oneof=claude claude_console"`
	RequestURL   string `json:"request_url" binding:"required,url"`
	SecretKey    string `json:"secret_key" binding:"required"`
	GroupID      int    `json:"group_id"`
	Priority     int    `json:"priority"`
	Weight       int    `json:"weight" binding:"min=1"`
	EnableProxy  bool   `json:"enable_proxy"`
	ProxyURI     string `json:"proxy_uri"`
	ActiveStatus int    `json:"active_status" binding:"oneof=1 2"`
	IsMax        bool   `json:"is_max"` // 是否是max账号
}

func (a *Account) TableName() string {
	return "accounts"
}

// 创建账号
func CreateAccount(account *Account) error {
	account.ID = 0
	return DB.Create(account).Error
}

// 根据ID获取账号
func GetAccountByID(id uint) (*Account, error) {
	var account Account
	err := DB.Preload("User").First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// 更新账号
func UpdateAccount(account *Account) error {
	return DB.Save(account).Error
}

// 删除账号（软删除）
func DeleteAccount(id uint) error {
	return DB.Delete(&Account{}, id).Error
}

// 分页获取账号列表
func GetAccountList(page, limit int, userID *uint) ([]Account, int64, error) {
	var accounts []Account
	var total int64

	query := DB.Model(&Account{})

	// 如果指定了用户ID，则只查询该用户的账号
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// 统计总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&accounts).Error

	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// 根据用户ID获取账号列表
func GetAccountsByUserID(userID uint) ([]Account, error) {
	var accounts []Account
	err := DB.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// 根据分组ID获取可用账号列表（按优先级和使用次数排序）
func GetAvailableAccountsByGroupID(groupID int) ([]Account, error) {
	var accounts []Account
	err := DB.Where("group_id = ? AND active_status = 1 AND (current_status = 1 OR (current_status = 3 AND (rate_limit_end_time IS NULL OR rate_limit_end_time < ?)))", groupID, time.Now()).
		Order("priority ASC, today_usage_count ASC").
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
