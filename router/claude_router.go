package router

import (
	"claude-code-relay/controller"
	"claude-code-relay/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

func SetClaudeCodeRouter(server *gin.Engine) {
	claude := server.Group("/claude-code")
	claude.Use(middleware.RateLimit(20, time.Second))
	// api key 鉴权
	claude.Use(middleware.ClaudeCodeAuth())
	{
		// 对话接口
		claude.POST("/v1/messages", controller.GetMessages)
	}
}
