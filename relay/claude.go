package relay

import (
	"bytes"
	"claude-code-relay/common"
	"claude-code-relay/model"
	"claude-code-relay/service"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	ClaudeAPIURL        = "https://api.anthropic.com/v1/messages"
	ClaudeOAuthTokenURL = "https://console.anthropic.com/v1/oauth/token"
	ClaudeOAuthClientID = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
)

// OAuthTokenResponse 表示OAuth token刷新响应
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// HandleClaudeRequest 处理Claude官方API平台的请求
func HandleClaudeRequest(c *gin.Context, account *model.Account) {
	// 记录请求开始时间用于计算耗时
	startTime := time.Now()

	// 从上下文中获取API Key信息
	var apiKey *model.ApiKey
	if keyInfo, exists := c.Get("api_key"); exists {
		apiKey = keyInfo.(*model.ApiKey)
	}
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if nil != err {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 获取有效的访问token
	accessToken, err := getValidAccessToken(account)
	if err != nil {
		log.Printf("获取有效访问token失败: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, ClaudeAPIURL, bytes.NewBuffer(body))
	if nil != err {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 固定请求头配置
	fixedHeaders := map[string]string{
		"Authorization":                             "Bearer " + accessToken,
		"anthropic-version":                         "2023-06-01",
		"X-Stainless-Retry-Count":                   "0",
		"X-Stainless-Timeout":                       "600",
		"X-Stainless-Lang":                          "js",
		"X-Stainless-Package-Version":               "0.55.1",
		"X-Stainless-OS":                            "MacOS",
		"X-Stainless-Arch":                          "arm64",
		"X-Stainless-Runtime":                       "node",
		"x-stainless-helper-method":                 "stream",
		"x-app":                                     "cli",
		"User-Agent":                                "claude-cli/1.0.44 (external, cli)",
		"anthropic-beta":                            "claude-code-20250219,oauth-2025-04-20,interleaved-thinking-2025-05-14,fine-grained-tool-streaming-2025-05-14",
		"X-Stainless-Runtime-Version":               "v20.18.1",
		"anthropic-dangerous-direct-browser-access": "true",
	}

	// 透传所有原始请求头
	for name, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// 设置或覆盖固定请求头
	for name, value := range fixedHeaders {
		req.Header.Set(name, value)
	}

	// 删除不需要的请求头
	req.Header.Del("X-Api-Key")
	req.Header.Del("Cookie")

	// 处理流式请求的Accept头
	isStream := true
	if c.Request.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/event-stream")
	}

	httpClientTimeout, _ := time.ParseDuration(os.Getenv("HTTP_CLIENT_TIMEOUT") + "s")
	if httpClientTimeout == 0 {
		httpClientTimeout = 120 * time.Second
	}

	// 创建基础Transport配置
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 如果启用了代理并配置了代理URI，配置代理
	if account.EnableProxy && account.ProxyURI != "" {
		proxyURL, err := url.Parse(account.ProxyURI)
		if err != nil {
			log.Printf("invalid proxy URI: %s", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		log.Printf("使用代理: %s", account.ProxyURI)
	}

	client := &http.Client{
		Timeout:   httpClientTimeout,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if nil != err {
		if errors.Is(err, context.Canceled) {
			c.AbortWithStatus(http.StatusRequestTimeout)
			return
		}

		log.Printf("❌ 请求失败: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer common.CloseIO(resp.Body)

	// 读取响应体
	var responseBody []byte
	var usageTokens *common.TokenUsage

	if resp.StatusCode >= 400 {
		// 错误响应，直接读取全部内容
		var readErr error
		responseBody, readErr = io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("❌ 读取错误响应失败: %v", readErr)
			c.AbortWithStatus(resp.StatusCode)
			return
		}

		// 调试日志：打印错误响应内容
		log.Printf("❌ 错误响应内容: %s", string(responseBody))
	}

	// 透传响应状态码
	c.Status(resp.StatusCode)

	// 透传响应头，但需要处理Content-Length以避免流式响应问题
	for name, values := range resp.Header {
		// 跳过Content-Length，让Gin自动处理
		if strings.ToLower(name) == "content-length" {
			continue
		}
		for _, value := range values {
			c.Header(name, value)
		}
	}

	if resp.StatusCode < 400 {
		// 成功响应，使用流式解析
		usageTokens, err = common.ParseStreamResponse(c.Writer, resp.Body)
		if err != nil {
			log.Println("stream copy and parse failed:", err.Error())
		}
	}

	// 如果是错误响应，写入响应体
	if resp.StatusCode >= 400 {
		c.Writer.Write(responseBody)

		// 如果是401或403错误，记录详细信息
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			log.Printf("认证错误 %d，账号: %s，错误详情: %s", resp.StatusCode, account.Name, string(responseBody))
		}
	}

	// 处理限流逻辑
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 请求成功，检查并清除可能的限流状态
		if account.CurrentStatus == 3 && account.RateLimitEndTime != nil {
			now := time.Now()
			if now.After(time.Time(*account.RateLimitEndTime)) {
				// 限流时间已过，重置状态
				account.CurrentStatus = 1
				account.RateLimitEndTime = nil
				if err := model.UpdateAccount(account); err != nil {
					log.Printf("重置账号限流状态失败: %v", err)
				} else {
					log.Printf("账号 %s 限流状态已自动重置", account.Name)
				}
			}
		}
	} else {
		// 处理限流检测
		isRateLimited := false
		var rateLimitResetTimestamp int64 = 0

		if resp.StatusCode == 429 {
			isRateLimited = true

			// 提取限流重置时间戳
			if resetHeader := resp.Header.Get("anthropic-ratelimit-unified-reset"); resetHeader != "" {
				if timestamp, err := strconv.ParseInt(resetHeader, 10, 64); err == nil {
					rateLimitResetTimestamp = timestamp
					resetTime := time.Unix(timestamp, 0)
					log.Printf("🕐 提取到限流重置时间戳: %d (%s)", timestamp, resetTime.Format(time.RFC3339))
				}
			}
		} else if len(responseBody) > 0 {
			// 检查响应体中的限流错误信息（对于非429错误）
			errorBodyStr := string(responseBody)

			// 尝试解析为JSON
			if errorData := gjson.Get(errorBodyStr, "error.message"); errorData.Exists() {
				if strings.Contains(strings.ToLower(errorData.String()), "exceed your account's rate limit") {
					isRateLimited = true
				}
			} else {
				// 直接检查字符串内容
				if strings.Contains(strings.ToLower(errorBodyStr), "exceed your account's rate limit") {
					isRateLimited = true
				}
			}
		}

		if isRateLimited {
			log.Printf("🚫 检测到账号 %s 被限流，状态码: %d", account.Name, resp.StatusCode)

			// 更新账号限流状态
			account.CurrentStatus = 3 // 限流状态

			if rateLimitResetTimestamp > 0 {
				// 使用API提供的准确重置时间
				resetTime := time.Unix(rateLimitResetTimestamp, 0)
				rateLimitEndTime := model.Time(resetTime)
				account.RateLimitEndTime = &rateLimitEndTime
				log.Printf("账号 %s 限流至 %s", account.Name, resetTime.Format(time.RFC3339))
			} else {
				// 使用默认5小时限流时间
				resetTime := time.Now().Add(5 * time.Hour)
				rateLimitEndTime := model.Time(resetTime)
				account.RateLimitEndTime = &rateLimitEndTime
				log.Printf("账号 %s 限流至 %s (默认5小时)", account.Name, resetTime.Format(time.RFC3339))
			}

			// 立即更新数据库
			if err := model.UpdateAccount(account); err != nil {
				log.Printf("更新账号限流状态失败: %v", err)
			}
		}
	}

	// 处理响应状态码并更新账号状态
	accountService := service.NewAccountService()
	go accountService.UpdateAccountStatus(account, resp.StatusCode, usageTokens)

	// 更新API Key统计信息
	if apiKey != nil {
		go service.UpdateApiKeyStatus(apiKey, resp.StatusCode, usageTokens)
	}

	// 保存日志记录（仅在请求成功时记录）
	if resp.StatusCode >= 200 && resp.StatusCode < 300 && usageTokens != nil && apiKey != nil {
		duration := time.Since(startTime).Milliseconds()
		logService := service.NewLogService()
		go func() {
			_, err := logService.CreateLogFromTokenUsage(usageTokens, apiKey.UserID, apiKey.ID, account.ID, duration, isStream)
			if err != nil {
				log.Printf("保存日志失败: %v", err)
			}
		}()
	}

}

// getValidAccessToken 获取有效的访问token，如果过期则自动刷新
func getValidAccessToken(account *model.Account) (string, error) {
	// 检查当前token是否存在
	if account.AccessToken == "" {
		return "", errors.New("账号缺少访问token")
	}

	// 检查token是否过期（提前5分钟刷新）
	now := time.Now().Unix()
	expiresAt := int64(account.ExpiresAt)

	// 如果过期时间存在且距离过期不到5分钟，或者已经过期，则需要刷新
	if expiresAt > 0 && now >= (expiresAt-300) {
		log.Printf("账号 %s 的token即将过期或已过期，尝试刷新", account.Name)

		if account.RefreshToken == "" {
			return "", errors.New("账号缺少刷新token，无法自动刷新")
		}

		// 刷新token
		newAccessToken, newRefreshToken, newExpiresAt, err := refreshToken(account)
		if err != nil {
			log.Printf("刷新token失败: %v", err)
			// 刷新失败时，如果当前token未完全过期，仍尝试使用
			if now < expiresAt {
				log.Printf("刷新失败但token未完全过期，尝试使用当前token")
				return account.AccessToken, nil
			}
			return "", fmt.Errorf("token已过期且刷新失败: %v", err)
		}

		// 更新账号信息
		account.AccessToken = newAccessToken
		account.RefreshToken = newRefreshToken
		account.ExpiresAt = int(newExpiresAt)

		// 保存到数据库
		if err := model.UpdateAccount(account); err != nil {
			log.Printf("更新账号token信息到数据库失败: %v", err)
			// 不返回错误，因为内存中的token已经更新
		}

		log.Printf("账号 %s token刷新成功", account.Name)
		return newAccessToken, nil
	}

	// token还有效，直接返回
	return account.AccessToken, nil
}

// refreshToken 使用refresh token获取新的access token
func refreshToken(account *model.Account) (accessToken, refreshToken string, expiresAt int64, err error) {
	payload := map[string]interface{}{
		"grant_type":    "refresh_token",
		"refresh_token": account.RefreshToken,
		"client_id":     ClaudeOAuthClientID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", 0, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	req, err := http.NewRequest("POST", ClaudeOAuthTokenURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", "", 0, fmt.Errorf("创建刷新请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "claude-cli/1.0.56 (external, cli)")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://claude.ai/")
	req.Header.Set("Origin", "https://claude.ai")

	// 创建HTTP客户端，配置代理（如果启用）
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if account.EnableProxy && account.ProxyURI != "" {
		proxyURL, err := url.Parse(account.ProxyURI)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", 0, fmt.Errorf("刷新token请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, fmt.Errorf("读取刷新响应失败: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", "", 0, fmt.Errorf("刷新token失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", "", 0, fmt.Errorf("解析token响应失败: %v", err)
	}

	if tokenResp.AccessToken == "" {
		return "", "", 0, errors.New("刷新响应中缺少access_token")
	}

	// 计算过期时间戳
	expiresAt = time.Now().Unix() + int64(tokenResp.ExpiresIn)

	log.Printf("Token刷新成功，新token: %s，将在 %d 秒后过期", maskToken(tokenResp.AccessToken), tokenResp.ExpiresIn)

	return tokenResp.AccessToken, tokenResp.RefreshToken, expiresAt, nil
}

// maskToken 遮蔽token用于安全日志输出
func maskToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}
