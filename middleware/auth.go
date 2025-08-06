package middleware

import (
	"claude-scheduler/model"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权访问",
				"code":  40001,
			})
			c.Abort()
			return
		}

		// 检查用户是否存在且状态正常
		uid, _ := strconv.ParseUint(userID.(string), 10, 32)
		user, err := model.GetUserById(uint(uid))
		if err != nil || user.Status != 1 {
			session.Delete("user_id")
			session.Save()
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
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权访问",
				"code":  40001,
			})
			c.Abort()
			return
		}

		u := user.(*model.User)
		if u.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "权限不足",
				"code":  40003,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}