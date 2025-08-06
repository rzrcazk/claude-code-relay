package service

import (
	"claude-code-relay/model"
	"errors"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

// GetAccountList 获取账号列表
func (s *AccountService) GetAccountList(page, limit int, userID *uint) (*model.AccountListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	accounts, total, err := model.GetAccountList(page, limit, userID)
	if err != nil {
		return nil, errors.New("获取账号列表失败")
	}

	result := &model.AccountListResponse{
		Accounts: accounts,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}

	return result, nil
}

// CreateAccount 创建账号
func (s *AccountService) CreateAccount(req *model.CreateAccountRequest, userID uint) (*model.Account, error) {
	account := &model.Account{
		Name:         req.Name,
		PlatformType: req.PlatformType,
		RequestURL:   req.RequestURL,
		SecretKey:    req.SecretKey,
		GroupID:      req.GroupID,
		Priority:     req.Priority,
		Weight:       req.Weight,
		EnableProxy:  req.EnableProxy,
		ProxyURI:     req.ProxyURI,
		ActiveStatus: req.ActiveStatus,
		IsMax:        req.IsMax,
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
		ExpiresAt:    req.ExpiresAt,
		UserID:       userID,
	}

	if err := model.CreateAccount(account); err != nil {
		return nil, errors.New("创建账号失败")
	}

	return account, nil
}

// GetAccountByID 根据ID获取账号详情
func (s *AccountService) GetAccountByID(id uint, userID *uint) (*model.Account, error) {
	account, err := model.GetAccountByID(id)
	if err != nil {
		return nil, errors.New("账号不存在")
	}

	// 如果指定了用户ID，验证账号是否属于该用户
	if userID != nil && account.UserID != *userID {
		return nil, errors.New("无权访问此账号")
	}

	return account, nil
}

// UpdateAccount 更新账号
func (s *AccountService) UpdateAccount(id uint, req *model.UpdateAccountRequest, userID *uint) (*model.Account, error) {
	account, err := s.GetAccountByID(id, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	account.Name = req.Name
	account.PlatformType = req.PlatformType
	account.RequestURL = req.RequestURL
	account.GroupID = req.GroupID
	account.Priority = req.Priority
	account.Weight = req.Weight
	account.EnableProxy = req.EnableProxy
	account.ProxyURI = req.ProxyURI
	account.ActiveStatus = req.ActiveStatus
	account.IsMax = req.IsMax

	if req.SecretKey != "" {
		account.SecretKey = req.SecretKey
	}

	if err := model.UpdateAccount(account); err != nil {
		return nil, errors.New("更新账号失败")
	}

	return account, nil
}

// DeleteAccount 删除账号
func (s *AccountService) DeleteAccount(id uint, userID *uint) error {
	account, err := s.GetAccountByID(id, userID)
	if err != nil {
		return err
	}

	if err := model.DeleteAccount(account.ID); err != nil {
		return errors.New("删除账号失败")
	}

	return nil
}

// GetAccountsByUserID 获取指定用户的所有账号
func (s *AccountService) GetAccountsByUserID(userID uint) ([]model.Account, error) {
	accounts, err := model.GetAccountsByUserID(userID)
	if err != nil {
		return nil, errors.New("获取用户账号列表失败")
	}

	return accounts, nil
}
