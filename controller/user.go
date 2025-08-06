package controller

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	user, err := model.GetUserByUsername(req.Username)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
			"code":  constant.Unauthorized,
		})
		return
	}

	if user.Status != constant.UserStatusActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "账户已被禁用",
			"code":  constant.UserStatusAbnormal,
		})
		return
	}

	// 生成JWT token
	token, err := common.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成token失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	// 同时设置session以保持向后兼容
	session := sessions.Default(c)
	session.Set("user_id", strconv.Itoa(int(user.ID)))
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"code":    constant.Success,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 检查用户名是否已存在
	if _, err := model.GetUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名已存在",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 检查邮箱是否已存在
	if _, err := model.GetUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱已存在",
			"code":  constant.InvalidParams,
		})
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 实际项目中应该加密
		Role:     constant.RoleUser,
		Status:   constant.UserStatusActive,
	}

	if err := model.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "注册失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"code":    constant.Success,
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "退出成功",
		"code":    constant.Success,
	})
}

func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data": gin.H{
			"user": gin.H{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"status":     user.Status,
				"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		},
	})
}

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := model.GetUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户列表失败",
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data": gin.H{
			"users": users,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
