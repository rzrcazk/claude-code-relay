package relay

import (
	"claude-code-relay/model"
	"github.com/gin-gonic/gin"
)

// HandleClaudeRequest 处理Claude平台的请求
func HandleClaudeRequest(c *gin.Context, account *model.Account) {
	// TODO: 实现Claude平台的请求处理逻辑
	c.JSON(200, gin.H{
		"message":  "Claude请求处理中",
		"account":  account.Name,
		"platform": "claude",
	})
}
