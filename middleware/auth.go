package middleware

import (
	"claude-code-relay/common"
	"claude-code-relay/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查Authorization header中的JWT token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少Authorization header",
				"code":  40001,
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := common.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token",
				"code":  40001,
			})
			c.Abort()
			return
		}

		// 获取用户信息并验证状态
		user, err := model.GetUserById(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户不存在",
				"code":  40001,
			})
			c.Abort()
			return
		}

		if user.Status != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户状态异常",
				"code":  40002,
			})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查Authorization header中的JWT token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少Authorization header",
				"code":  40001,
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := common.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token",
				"code":  40001,
			})
			c.Abort()
			return
		}

		// 重新查询用户信息以确保最新状态
		user, err := model.GetUserById(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户不存在",
				"code":  40001,
			})
			c.Abort()
			return
		}

		if user.Status != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户状态异常",
				"code":  40002,
			})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "权限不足",
				"code":  40003,
			})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}
