package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetClaudeCodeRouter(server *gin.Engine) {

	claude := server.Group("/claude-code/")
	{
		// 健康检查
		claude.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"time":   time.Now().Unix(),
			})
		})
	}

}
