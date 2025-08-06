package controller

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExchangeRequest struct {
	AuthorizationCode string `json:"authorization_code" binding:"required"`
	CallbackUrl       string `json:"callback_url" binding:"required"`
	ProxyURI          string `json:"proxy_uri" binding:"omitempty,url"`
	CodeVerifier      string `json:"code_verifier" binding:"required"`
	State             string `json:"state" binding:"required"`
}

// GetOAuthURL 获取OAuth授权URL
func GetOAuthURL(c *gin.Context) {
	oauthHelper := common.NewOAuthHelper(nil)
	// 生成OAuth参数
	params, err := oauthHelper.GenerateOAuthParams()
	if err != nil {
		// 处理错误
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"code":    constant.Success,
		"data":    params,
	})
}

// ExchangeCode 验证授权码并返回token
func ExchangeCode(c *gin.Context) {
	var req ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"code":  constant.InvalidParams,
		})
		return
	}

	oauthHelper := common.NewOAuthHelper(nil)
	finalAuthCode, err := oauthHelper.ParseCallbackURL(req.AuthorizationCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的授权码",
			"code":  constant.InvalidParams,
		})
		return
	}

	// 生成访问令牌
	tokenResult, err := oauthHelper.ExchangeCodeForTokens(finalAuthCode, req.CodeVerifier, req.State, req.ProxyURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成访问令牌事变",
			"code":  constant.InternalServerError,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"code":    constant.Success,
		"data":    tokenResult,
	})
}

// GetMessages 获取对话消息
func GetMessages(c *gin.Context) {
	// 获取API Key对象
	apiKey, exists := c.Get("api_key")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到API Key信息",
			"code":  constant.Unauthorized,
		})
		return
	}
	keyInfo := apiKey.(*model.ApiKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "操作成功",
		"code":    constant.Success,
		"data":    keyInfo,
	})
}
