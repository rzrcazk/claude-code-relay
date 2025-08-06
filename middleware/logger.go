package middleware

import (
	"claude-code-relay/model"
	"time"

	"github.com/gin-gonic/gin"
)

func SetUpLogger(server *gin.Engine) {
	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	}))

	server.Use(ApiLogger())
}

func ApiLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		// 获取用户ID（如果已认证）
		var userID uint
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = id
			}
		}

		// 获取请求ID
		requestID, _ := c.Get("request_id")
		requestIDStr := ""
		if id, ok := requestID.(string); ok {
			requestIDStr = id
		}

		// 记录API日志到数据库
		apiLog := &model.ApiLog{
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			UserID:     userID,
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			RequestID:  requestIDStr,
			Duration:   duration.Milliseconds(),
		}

		// 异步记录日志，避免阻塞请求
		go func() {
			_ = model.CreateApiLog(apiLog)
		}()
	}
}
