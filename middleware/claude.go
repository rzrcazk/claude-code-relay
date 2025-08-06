package middleware

import "github.com/gin-gonic/gin"

// ClaudeCodeAuth 鉴权中间件
func ClaudeCodeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
