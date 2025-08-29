package middleware

import (
	"claude-code-relay/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SystemMessage 系统消息结构体
type SystemMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// RequestBody Claude API请求体结构
type RequestBody struct {
	System interface{} `json:"system"`
}

// ClaudeCodeAuth API Key鉴权中间件
func ClaudeCodeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从多个可能的请求头中获取API Key
		apiKey := getApiKeyFromHeaders(c)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少API Key",
				"code":  40001,
			})
			c.Abort()
			return
		}

		// 从数据库查询API Key
		keyInfo, err := model.GetApiKeyByKey(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的API Key",
				"code":  40001,
			})
			c.Abort()
			return
		}

		// 判断是否达到每日限额
		if keyInfo.DailyLimit > 0 && keyInfo.TodayTotalCost >= keyInfo.DailyLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "API Key已达到每日使用限额",
				"code":  40004,
			})
			c.Abort()
			return
		}

		// 检查分组是否被禁用
		if keyInfo.GroupID > 0 {
			status := model.GetGroupStatus(keyInfo.GroupID)
			if status != 1 {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "API Key所属分组不可用",
					"code":  40005,
				})
				c.Abort()
				return
			}
		}

		// API Key已经在model层验证了状态和过期时间
		// 将API Key信息存储到上下文中供后续使用
		c.Set("api_key_id", keyInfo.ID)
		c.Set("api_key", keyInfo)
		c.Set("user_id", keyInfo.UserID)
		c.Set("group_id", keyInfo.GroupID)

		c.Next()
	}
}

// getApiKeyFromHeaders 从多个可能的请求头中提取API Key
func getApiKeyFromHeaders(c *gin.Context) string {
	// 1. 检查 X-API-Key
	if apiKey := c.GetHeader("x-api-key"); apiKey != "" {
		return apiKey
	}

	// 2. 检查 X-Goog-API-Key (Google Cloud API格式)
	if apiKey := c.GetHeader("X-Goog-API-Key"); apiKey != "" {
		return apiKey
	}

	// 3. 检查 Authorization Bearer Token
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return strings.TrimSpace(authHeader[7:])
		}
	}

	// 4. 检查 API-Key
	if apiKey := c.GetHeader("API-Key"); apiKey != "" {
		return apiKey
	}

	return ""
}
