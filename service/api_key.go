package service

import (
	"claude-code-relay/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

func CreateApiKey(userID uint, req *model.CreateApiKeyRequest) (*model.ApiKey, error) {
	if req.Name == "" {
		return nil, errors.New("API Key名称不能为空")
	}

	// 如果指定了分组ID，验证分组是否存在且属于用户
	if req.GroupID > 0 {
		_, err := model.GetGroupById(req.GroupID, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("指定的分组不存在")
			}
			return nil, err
		}
	}

	// 如果指定了过期时间，验证时间有效性
	if req.ExpiresAt != nil && time.Time(*req.ExpiresAt).Before(time.Now()) {
		return nil, errors.New("过期时间不能早于当前时间")
	}

	apiKey := &model.ApiKey{
		Name:      req.Name,
		Key:       req.Key,
		ExpiresAt: req.ExpiresAt,
		Status:    req.Status,
		GroupID:   req.GroupID,
		UserID:    userID,
	}

	if apiKey.Status == 0 {
		apiKey.Status = 1 // 默认启用
	}

	err := model.CreateApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func GetApiKeyById(id, userID uint) (*model.ApiKey, error) {
	return model.GetApiKeyById(id, userID)
}

func UpdateApiKey(id, userID uint, req *model.UpdateApiKeyRequest) (*model.ApiKey, error) {
	apiKey, err := model.GetApiKeyById(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("API Key不存在")
		}
		return nil, err
	}

	// 如果指定了分组ID，验证分组是否存在且属于用户
	if req.GroupID != 0 {
		_, err := model.GetGroupById(req.GroupID, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("指定的分组不存在")
			}
			return nil, err
		}
	}

	// 如果指定了过期时间，验证时间有效性
	if req.ExpiresAt != nil && time.Time(*req.ExpiresAt).Before(time.Now()) {
		return nil, errors.New("过期时间不能早于当前时间")
	}

	// 更新字段
	if req.Name != "" {
		apiKey.Name = req.Name
	}
	if req.ExpiresAt != nil {
		apiKey.ExpiresAt = req.ExpiresAt
	}
	if req.Status != nil {
		apiKey.Status = *req.Status
	}
	if req.GroupID != 0 {
		apiKey.GroupID = req.GroupID
	}

	err = model.UpdateApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func DeleteApiKey(id, userID uint) error {
	apiKey, err := model.GetApiKeyById(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("API Key不存在")
		}
		return err
	}

	return model.DeleteApiKey(apiKey.ID)
}

func GetApiKeys(page, limit int, userID uint, groupID *uint) (*model.ApiKeyListResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	apiKeys, total, err := model.GetApiKeys(page, limit, userID, groupID)
	if err != nil {
		return nil, err
	}

	return &model.ApiKeyListResult{
		ApiKeys: apiKeys,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// ValidateApiKey 验证API Key是否有效
func ValidateApiKey(key string) (*model.ApiKey, error) {
	if key == "" {
		return nil, errors.New("API Key不能为空")
	}

	apiKey, err := model.GetApiKeyByKey(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("API Key无效或已过期")
		}
		return nil, err
	}

	return apiKey, nil
}
