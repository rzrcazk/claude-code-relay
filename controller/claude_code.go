package controller

import (
	"claude-code-relay/common"
	"claude-code-relay/constant"
	"claude-code-relay/model"
	"claude-code-relay/relay"
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
	// 从上下文中获取API Key的详细信息
	apiKey, _ := c.Get("api_key")
	keyInfo := apiKey.(*model.ApiKey)

	// 根据API Key的分组ID查询可用账号列表
	accounts, err := model.GetAvailableAccountsByGroupID(keyInfo.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询账号列表失败",
			"code":    constant.InternalServerError,
		})
		return
	}

	if len(accounts) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "没有可用的账号",
			"code":    constant.NotFound,
		})
		return
	}

	// 选择第一个账号（已按优先级和使用次数排序）
	selectedAccount := accounts[0]

	// 根据平台类型路由到不同的处理器
	switch selectedAccount.PlatformType {
	case constant.PlatformClaude:
		relay.HandleClaudeRequest(c, &selectedAccount)
	case constant.PlatformClaudeConsole:
		relay.HandleClaudeConsoleRequest(c, &selectedAccount)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "不支持的平台类型: " + selectedAccount.PlatformType,
			"code":    constant.InvalidParams,
		})
	}
}
