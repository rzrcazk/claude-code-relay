package controller

import (
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"claude-code-relay/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAccountList 获取账号列表
func GetAccountList(c *gin.Context) {
	var req model.AccountListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	// 获取当前用户信息
	user := c.MustGet("user").(*model.User)
	var userID *uint

	// 如果是普通用户，只能查看自己的账号
	if user.Role != "admin" {
		userID = &user.ID
	} else {
		// 管理员可以通过参数指定查看特定用户的账号
		userID = req.UserID
	}

	accountService := service.NewAccountService()
	result, err := accountService.GetAccountList(req.Page, req.Limit, userID)
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

// CreateAccount 创建账号
func CreateAccount(c *gin.Context) {
	var req model.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	user := c.MustGet("user").(*model.User)
	accountService := service.NewAccountService()

	// 验证 req.PlatformType 是否是有效值
	validPlatformTypes := map[string]bool{
		constant.PlatformClaude:        true,
		constant.PlatformClaudeConsole: true,
		constant.PlatformOpenAI:        true,
		constant.PlatformGemini:        true,
	}
	if !validPlatformTypes[req.PlatformType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的平台类型",
			"code":  constant.InvalidParams,
		})
		return
	}

	account, err := accountService.CreateAccount(&req, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"code":    constant.Success,
		"data":    account,
	})
}

// GetAccount 获取账号详情
func GetAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的账号ID",
			"code":  constant.InvalidParams,
		})
		return
	}

	user := c.MustGet("user").(*model.User)
	var userID *uint

	// 如果是普通用户，只能查看自己的账号
	if user.Role != "admin" {
		userID = &user.ID
	}

	accountService := service.NewAccountService()
	account, err := accountService.GetAccountByID(uint(id), userID)
	if err != nil {
		var statusCode int
		var code int
		if err.Error() == "账号不存在" {
			statusCode = http.StatusNotFound
			code = constant.NotFound
		} else if err.Error() == "无权访问此账号" {
			statusCode = http.StatusForbidden
			code = constant.Unauthorized
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
		"message": "获取成功",
		"code":    constant.Success,
		"data":    account,
	})
}

// UpdateAccount 更新账号
func UpdateAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的账号ID",
			"code":  constant.InvalidParams,
		})
		return
	}

	var req model.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	user := c.MustGet("user").(*model.User)
	var userID *uint

	// 如果是普通用户，只能更新自己的账号
	if user.Role != "admin" {
		userID = &user.ID
	}

	accountService := service.NewAccountService()
	account, err := accountService.UpdateAccount(uint(id), &req, userID)
	if err != nil {
		var statusCode int
		var code int
		if err.Error() == "账号不存在" {
			statusCode = http.StatusNotFound
			code = constant.NotFound
		} else if err.Error() == "无权访问此账号" {
			statusCode = http.StatusForbidden
			code = constant.Unauthorized
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
		"message": "更新成功",
		"code":    constant.Success,
		"data":    account,
	})
}

// DeleteAccount 删除账号
func DeleteAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的账号ID",
			"code":  constant.InvalidParams,
		})
		return
	}

	user := c.MustGet("user").(*model.User)
	var userID *uint

	// 如果是普通用户，只能删除自己的账号
	if user.Role != "admin" {
		userID = &user.ID
	}

	accountService := service.NewAccountService()
	err = accountService.DeleteAccount(uint(id), userID)
	if err != nil {
		var statusCode int
		var code int
		if err.Error() == "账号不存在" {
			statusCode = http.StatusNotFound
			code = constant.NotFound
		} else if err.Error() == "无权访问此账号" {
			statusCode = http.StatusForbidden
			code = constant.Unauthorized
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
		"message": "删除成功",
		"code":    constant.Success,
	})
}
