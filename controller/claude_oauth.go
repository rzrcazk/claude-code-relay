package controller

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetOAuthURL 获取OAuth授权URL
func GetOAuthURL(c *gin.Context) {
	helper := common.NewOAuthHelper(nil)
	// 生成OAuth参数
	params, err := helper.GenerateOAuthParams()
	if err != nil {
		// 处理错误
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"code":    constant.Success,
		"data":    params,
	})
}
