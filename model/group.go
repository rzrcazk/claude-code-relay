package model

import (
	"claude-code-relay/common"
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"type:varchar(100);not null;uniqueIndex:idx_groups_user_name"`
	Remark     string         `json:"remark" gorm:"type:text"`
	Status     int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	UserID     uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_groups_user_name"`
	InstanceID string         `json:"instance_id" gorm:"type:varchar(150)"`
	CreatedAt  Time           `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt  Time           `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"uniqueIndex:idx_groups_user_name"`

	// 统计字段，不存储在数据库中
	ApiKeyCount  int `json:"api_key_count" gorm:"-"`
	AccountCount int `json:"account_count" gorm:"-"`
}

type CreateGroupRequest struct {
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}

type UpdateGroupRequest struct {
	Name   string `json:"name"`
	Remark string `json:"remark"`
	Status *int   `json:"status"`
}

type GroupListResult struct {
	Groups []Group `json:"groups"`
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

func (g *Group) TableName() string {
	return "groups"
}

func CreateGroup(group *Group) error {
	group.ID = 0
	return DB.Create(group).Error
}

func GetGroupById(id int, userID uint) (*Group, error) {
	var group Group
	err := DB.Where("id = ? AND user_id = ?", id, userID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetGroupStatus 获取分组状态（带缓存），找不到返回0
func GetGroupStatus(id int) int {
	// 先尝试从缓存获取
	if common.RDB != nil {
		cacheKey := fmt.Sprintf("group_status:%d", id)
		cachedStatus, err := common.RDB.Get(context.Background(), cacheKey).Result()
		if err == nil {
			if status, parseErr := strconv.Atoi(cachedStatus); parseErr == nil {
				return status
			}
		}
	}

	// 缓存未命中，从数据库查询
	var group Group
	err := DB.Select("id,status").Where("id = ?", id).First(&group).Error
	if err != nil {
		// 查询失败时，缓存结果0也可以避免频繁查询
		if common.RDB != nil {
			cacheKey := fmt.Sprintf("group_status:%d", id)
			common.RDB.Set(context.Background(), cacheKey, "0", 5*time.Minute)
		}
		return 0
	}

	// 存储到缓存（5分钟）
	if common.RDB != nil {
		cacheKey := fmt.Sprintf("group_status:%d", id)
		common.RDB.Set(context.Background(), cacheKey, strconv.Itoa(group.Status), 5*time.Minute)
	}

	return group.Status
}

// clearGroupStatusCache 清理分组状态缓存
func clearGroupStatusCache(groupID int) {
	if common.RDB != nil {
		cacheKey := fmt.Sprintf("group_status:%d", groupID)
		common.RDB.Del(context.Background(), cacheKey)
	}
}

func GetGroupByName(name string, userID uint) (*Group, error) {
	var group Group
	err := DB.Where("name = ? AND user_id = ?", name, userID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func UpdateGroup(group *Group) error {
	err := DB.Save(group).Error
	if err != nil {
		return err
	}

	// 更新成功后清理相关缓存
	clearGroupStatusCache(int(group.ID))
	return nil
}

func DeleteGroup(id uint) error {
	err := DB.Delete(&Group{}, id).Error
	if err != nil {
		return err
	}

	// 删除成功后清理相关缓存
	clearGroupStatusCache(int(id))
	return nil
}

// GetAllGroups 获取所有分组（不分页）
func GetAllGroups(userID uint) ([]Group, error) {
	var groups []Group

	err := DB.Where("user_id = ? AND status = 1", userID).Find(&groups).Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func GetGroups(page, limit int, userID uint) ([]Group, int64, error) {
	var groups []Group
	var total int64

	err := DB.Model(&Group{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}

	// 为每个分组统计API密钥和账号数量
	for i := range groups {
		// 统计API密钥数量
		var apiKeyCount int64
		err = DB.Model(&ApiKey{}).Where("group_id = ? AND user_id = ?", groups[i].ID, userID).Count(&apiKeyCount).Error
		if err != nil {
			return nil, 0, err
		}
		groups[i].ApiKeyCount = int(apiKeyCount)

		// 统计账号数量
		var accountCount int64
		err = DB.Model(&Account{}).Where("group_id = ? AND user_id = ?", groups[i].ID, userID).Count(&accountCount).Error
		if err != nil {
			return nil, 0, err
		}
		groups[i].AccountCount = int(accountCount)
	}

	return groups, total, nil
}

// GetInstanceID 根据分组ID获取实例ID，如果不存在则生成并存储
func GetInstanceID(groupID uint) string {
	var group Group
	err := DB.Where("id = ?", groupID).First(&group).Error
	if err != nil {
		return common.GetInstanceID()
	}

	if group.InstanceID == "" {
		group.InstanceID = common.GenerateRandomInstanceID()

		// 保存到数据库
		DB.Model(&group).Update("instance_id", group.InstanceID)
	}
	return group.InstanceID
}
