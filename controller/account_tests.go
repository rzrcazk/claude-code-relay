package controller

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"claude-code-relay/tests"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TestAccountRequest 测试账号请求参数
type TestAccountRequest struct {
	AccountID uint `json:"account_id" binding:"required" form:"account_id"`
}

// TestAccountResponse 测试账号响应结构
type TestAccountResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	StatusCode   int    `json:"status_code,omitempty"`
	PlatformType string `json:"platform_type"`
}

// TestAccount 测试指定账号的连通性
func TestAccount(c *gin.Context) {
	var req TestAccountRequest

	// 从URL参数中获取账号ID
	accountIDStr := c.Param("id")
	if accountIDStr == "" {
		// 尝试从查询参数或表单中获取
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "参数错误: " + err.Error(),
				"code":    constant.InvalidParams,
			})
			return
		}
	} else {
		// 从URL参数解析账号ID
		accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "无效的账号ID",
				"code":    constant.InvalidParams,
			})
			return
		}
		req.AccountID = uint(accountID)
	}

	// 获取账号信息
	account, err := model.GetAccountByID(req.AccountID)
	if err != nil {
		common.SysError(fmt.Sprintf("测试账号时获取账号信息失败: %v", err))
		c.JSON(http.StatusNotFound, gin.H{
			"message": "账号不存在",
			"code":    constant.NotFound,
		})
		return
	}

	var testResult TestAccountResponse
	testResult.PlatformType = account.PlatformType

	// 根据平台类型调用不同的测试函数
	switch account.PlatformType {
	case constant.PlatformClaude:
		statusCode, err := tests.TestClaudeAccount(account)
		testResult.StatusCode = statusCode
		if err != nil {
			common.SysError(fmt.Sprintf("测试Claude账号失败: %v", err))
			testResult.Success = false
			testResult.Message = fmt.Sprintf("测试失败: %v", err)
		} else if statusCode == http.StatusOK {
			testResult.Success = true
			testResult.Message = "测试成功，账号连接正常"
		} else {
			testResult.Success = false
			testResult.Message = fmt.Sprintf("测试失败，HTTP状态码: %d", statusCode)
		}

	case constant.PlatformClaudeConsole:
		statusCode, err := tests.TestClaudeConsoleAccount(account)
		testResult.StatusCode = statusCode
		if err != nil {
			common.SysError(fmt.Sprintf("测试Claude Console账号失败: %v", err))
			testResult.Success = false
			testResult.Message = fmt.Sprintf("测试失败: %v", err)
		} else if statusCode == http.StatusOK {
			testResult.Success = true
			testResult.Message = "测试成功，账号连接正常"
		} else {
			testResult.Success = false
			testResult.Message = fmt.Sprintf("测试失败，HTTP状态码: %d", statusCode)
		}

	default:
		testResult.Success = false
		testResult.Message = "不支持的平台类型: " + account.PlatformType
		c.JSON(http.StatusBadRequest, gin.H{
			"message": testResult.Message,
			"code":    constant.InvalidParams,
		})
		return
	}

	// 记录测试日志
	common.SysLog(fmt.Sprintf("账号测试完成 - 账号ID: %d, 平台: %s, 结果: %v, 状态码: %d",
		req.AccountID, account.PlatformType, testResult.Success, testResult.StatusCode))

	// 返回测试结果
	c.JSON(http.StatusOK, gin.H{
		"message": "账号测试完成",
		"code":    constant.Success,
		"data":    testResult,
	})
}
