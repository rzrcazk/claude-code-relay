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

func Logout(c *gin.Context) {
	userService := service.NewUserService()
	userService.Logout(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "退出成功",
		"code":    constant.Success,
	})
}

func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	userService := service.NewUserService()
	profile := userService.GetProfile(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"code":    constant.Success,
		"data": gin.H{
			"user": profile,
		},
	})
}

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
