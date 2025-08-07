package controller

import (
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"claude-code-relay/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateApiKey 创建API Key
func CreateApiKey(c *gin.Context) {
	var req model.CreateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 从认证中获取用户ID
	user := c.MustGet("user").(*model.User)
	userID := user.ID

	apiKey, err := service.CreateApiKey(userID, &req)
	if err != nil {
		var statusCode int
		var code int
		switch err.Error() {
		case "API Key名称不能为空", "指定的分组不存在", "过期时间不能早于当前时间":
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
		"code": constant.Success,
		"data": gin.H{
			"key": apiKey.Key,
		},
	})
}

// GetApiKey 获取API Key详情
func GetApiKey(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的API Key ID",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 从认证中获取用户ID
	user := c.MustGet("user").(*model.User)
	userID := user.ID

	apiKey, err := service.GetApiKeyById(uint(idInt), userID)
	if err != nil {
		var statusCode int
		var code int
		if err.Error() == "API Key不存在" {
			statusCode = http.StatusNotFound
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
		"message": "获取API Key成功",
		"code":    constant.Success,
		"data":    apiKey,
	})
}

// UpdateApiKey 更新API Key
func UpdateApiKey(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的API Key ID",
			"code":  constant.InvalidParams,
		})
		return
	}

	var req model.UpdateApiKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 从认证中获取用户ID
	user := c.MustGet("user").(*model.User)
	userID := user.ID

	apiKey, err := service.UpdateApiKey(uint(idInt), userID, &req)
	if err != nil {
		var statusCode int
		var code int
		switch err.Error() {
		case "API Key不存在", "指定的分组不存在", "过期时间不能早于当前时间":
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
		"message": "更新API Key成功",
		"code":    constant.Success,
		"data":    apiKey,
	})
}

// DeleteApiKey 删除API Key
func DeleteApiKey(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的API Key ID",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 从认证中获取用户ID
	user := c.MustGet("user").(*model.User)
	userID := user.ID

	err = service.DeleteApiKey(uint(idInt), userID)
	if err != nil {
		var statusCode int
		var code int
		if err.Error() == "API Key不存在" {
			statusCode = http.StatusNotFound
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
		"message": "删除API Key成功",
		"code":    constant.Success,
	})
}

// GetApiKeys 获取API Key列表
func GetApiKeys(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 可选的分组ID过滤
	var groupID *uint
	if groupIDStr := c.Query("group_id"); groupIDStr != "" {
		if groupIDInt, err := strconv.ParseUint(groupIDStr, 10, 32); err == nil {
			groupIDValue := uint(groupIDInt)
			groupID = &groupIDValue
		}
	}

	// 从认证中获取用户ID
	user := c.MustGet("user").(*model.User)
	userID := user.ID

	result, err := service.GetApiKeys(page, limit, userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  constant.InternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取API Key列表成功",
		"code":    constant.Success,
		"data":    result,
	})
}
