package controller

import (
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"claude-code-relay/service"
	"net/http"
	"strconv"

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

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

type AdminCreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type AdminUpdateUserStatusRequest struct {
	Status int `json:"status" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	userService := service.NewUserService()
	result, err := userService.Login(req.Username, req.Password, c)
	if err != nil {
		var code int
		switch err.Error() {
		case "用户名或密码错误":
			code = constant.Unauthorized
		case "账户已被禁用":
			code = constant.UserStatusAbnormal
		default:
			code = constant.InternalServerError
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
			"code":  code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"code":    constant.Success,
		"data":    result,
	})
}

// Register 新用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	userService := service.NewUserService()
	err := userService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		var statusCode int
		var code int
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" {
			statusCode = http.StatusBadRequest
			code = constant.InvalidParams
		} else {
			statusCode = http.StatusInternalServerError
			code = constant.InternalServerError
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"code":    constant.Success,
	})
}

// GetProfile 获取用户信息
func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	userService := service.NewUserService()
	profile := userService.GetProfile(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data":    profile,
	})
}

// UpdateProfile 更新用户信息
func UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	user := c.MustGet("user").(*model.User)
	userService := service.NewUserService()

	err := userService.UpdateProfile(user, req.Username, req.Email, req.Password)
	if err != nil {
		var statusCode int
		var code int
		switch err.Error() {
		case "用户名已存在", "邮箱已存在":
			statusCode = http.StatusBadRequest
			code = constant.InvalidParams
		default:
			statusCode = http.StatusInternalServerError
			code = constant.InternalServerError
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"code":    constant.Success,
	})
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	userService := service.NewUserService()
	result, err := userService.GetUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data":    result,
	})
}

// AdminCreateUser 管理员创建用户
func AdminCreateUser(c *gin.Context) {
	var req AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	userService := service.NewUserService()
	err := userService.AdminCreateUser(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		var statusCode int
		var code int
		switch err.Error() {
		case "用户名已存在", "邮箱已存在", "角色参数无效":
			statusCode = http.StatusBadRequest
			code = constant.InvalidParams
		default:
			statusCode = http.StatusInternalServerError
			code = constant.InternalServerError
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户创建成功",
		"code":    constant.Success,
	})
}

// AdminUpdateUserStatus 管理员更新用户状态
func AdminUpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户ID参数无效",
			"code":  constant.InvalidParams,
		})
		return
	}

	var req AdminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	userService := service.NewUserService()
	err = userService.AdminUpdateUserStatus(uint(userID), req.Status)
	if err != nil {
		var statusCode int
		var code int
		switch err.Error() {
		case "用户不存在":
			statusCode = http.StatusNotFound
			code = constant.ResourceNotFound
		case "状态参数无效":
			statusCode = http.StatusBadRequest
			code = constant.InvalidParams
		default:
			statusCode = http.StatusInternalServerError
			code = constant.InternalServerError
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"code":  code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户状态更新成功",
		"code":    constant.Success,
	})
}
