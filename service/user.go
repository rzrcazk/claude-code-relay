package service

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"errors"
	"strconv"

	"github.com/gin-contrib/sessions"
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

func (s *UserService) Login(username, password string, c *gin.Context) (*LoginResult, error) {
	user, err := model.GetUserByUsername(username)
	if err != nil || user.Password != password {
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

	// 同时设置session以保持向后兼容
	session := sessions.Default(c)
	session.Set("user_id", strconv.Itoa(int(user.ID)))
	session.Save()

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
		Password: password, // 实际项目中应该加密
		Role:     constant.RoleUser,
		Status:   constant.UserStatusActive,
	}

	if err := model.CreateUser(user); err != nil {
		return errors.New("注册失败")
	}

	return nil
}

func (s *UserService) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()
}

func (s *UserService) GetProfile(user *model.User) *model.UserProfile {
	return &model.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
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
