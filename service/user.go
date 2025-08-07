package service

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

type LoginResult struct {
	Token string          `json:"token"`
	User  *model.UserInfo `json:"user"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Login 用户登录
func (s *UserService) Login(username, password string, c *gin.Context) (*LoginResult, error) {
	user, err := model.GetUserByUsername(username)
	if err != nil || user.Password != common.HashPassword(password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != constant.UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	// 生成JWT token
	token, err := common.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成token失败")
	}

	result := &LoginResult{
		Token: token,
		User: &model.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}

	return result, nil
}

// Register 注册新用户
func (s *UserService) Register(username, email, password string) error {
	// 检查用户名是否已存在
	if _, err := model.GetUserByUsername(username); err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := model.GetUserByEmail(email); err == nil {
		return errors.New("邮箱已存在")
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: common.HashPassword(password), // 实际项目中应该加密
		Role:     constant.RoleUser,
		Status:   constant.UserStatusActive,
	}

	if err := model.CreateUser(user); err != nil {
		return errors.New("注册失败")
	}

	return nil
}

func (s *UserService) GetProfile(user *model.User) *model.UserProfile {
	return &model.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.String(),
	}
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(currentUser *model.User, username, email, password string) error {
	// 如果要更新用户名，检查是否已存在
	if username != "" && username != currentUser.Username {
		if _, err := model.GetUserByUsername(username); err == nil {
			return errors.New("用户名已存在")
		}
		currentUser.Username = username
	}

	// 如果要更新邮箱，检查是否已存在
	if email != "" && email != currentUser.Email {
		if _, err := model.GetUserByEmail(email); err == nil {
			return errors.New("邮箱已存在")
		}
		currentUser.Email = email
	}

	// 如果要更新密码，进行加密
	if password != "" {
		currentUser.Password = common.HashPassword(password)
	}

	// 保存更新
	if err := model.UpdateUser(currentUser); err != nil {
		return errors.New("更新失败")
	}

	return nil
}

func (s *UserService) GetUsers(page, limit int) (*model.UserListResult, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := model.GetUsers(page, limit)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}

	result := &model.UserListResult{
		Users: users,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

// AdminCreateUser 管理员创建用户
func (s *UserService) AdminCreateUser(username, email, password, role string) error {
	// 检查用户名是否已存在
	if _, err := model.GetUserByUsername(username); err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := model.GetUserByEmail(email); err == nil {
		return errors.New("邮箱已存在")
	}

	// 验证角色是否有效
	if role != constant.RoleUser && role != constant.RoleAdmin {
		return errors.New("角色参数无效")
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: common.HashPassword(password),
		Role:     role,
		Status:   constant.UserStatusActive,
	}

	if err := model.CreateUser(user); err != nil {
		return errors.New("创建用户失败")
	}

	return nil
}

// AdminUpdateUserStatus 管理员更新用户状态
func (s *UserService) AdminUpdateUserStatus(userID uint, status int) error {
	// 获取用户
	user, err := model.GetUserById(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证状态参数
	if status != constant.UserStatusActive && status != constant.UserStatusInactive {
		return errors.New("状态参数无效")
	}

	// 更新状态
	user.Status = status
	if err := model.UpdateUser(user); err != nil {
		return errors.New("更新用户状态失败")
	}

	return nil
}
