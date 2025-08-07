package model

import (
	"gorm.io/gorm"
)

type Group struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null;uniqueIndex:idx_groups_user_name"`
	Remark    string         `json:"remark" gorm:"type:text"`
	Status    int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex:idx_groups_user_name"`
	CreatedAt Time           `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt Time           `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"uniqueIndex:idx_groups_user_name"`
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

func GetGroupByName(name string, userID uint) (*Group, error) {
	var group Group
	err := DB.Where("name = ? AND user_id = ?", name, userID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func UpdateGroup(group *Group) error {
	return DB.Save(group).Error
}

func DeleteGroup(id uint) error {
	return DB.Delete(&Group{}, id).Error
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

	return groups, total, nil
}
