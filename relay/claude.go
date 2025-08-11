package relay

import (
	"bytes"
	"claude-code-relay/common"
	"claude-code-relay/model"
	"claude-code-relay/service"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/sjson"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"type":    "request_error",
				"message": "Incorrect request body",
			},
		})
		return
	}

	body, _ = sjson.SetBytes(body, "stream", true) // 强制流式输出

	modelName := gjson.GetBytes(body, "model").String()
	if modelName == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": map[string]interface{}{
				"type":    "request_error",
				"message": "The model field is missing in the request body",
			},
		})
		return
	}

	// 模型名称是否允许在apiKey的限制模型中
	if apiKey.ModelRestriction != "" {
		allowedModels := strings.Split(apiKey.ModelRestriction, ",")
		modelAllowed := false
		for _, allowedModel := range allowedModels {
			if strings.EqualFold(strings.TrimSpace(allowedModel), modelName) {
				modelAllowed = true
				break
			}
		}
		if !modelAllowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": map[string]interface{}{
					"type":    "request_error",
					"message": "This model is not allowed.",
				},
			})
			return
		}
	}

	// 获取有效的访问token
	accessToken, err := getValidAccessToken(account)
	if err != nil {
		log.Printf("获取有效访问token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"type":    "authentication_error",
				"message": "Failed to get valid access token: " + err.Error(),
			},
		})
		return
	}

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, ClaudeAPIURL, bytes.NewBuffer(body))
	if nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"type":    "internal_server_error",
				"message": "Failed to create request: " + err.Error(),
			},
		})
		return
	}

	// 使用公共的请求头构建方法
	fixedHeaders := buildClaudeAPIHeaders(accessToken)

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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": map[string]interface{}{
					"type":    "proxy_configuration_error",
					"message": "Invalid proxy URI: " + err.Error(),
				},
			})
			return
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   httpClientTimeout,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if nil != err {
		if errors.Is(err, context.Canceled) {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": map[string]interface{}{
					"type":    "timeout_error",
					"message": "Request was canceled or timed out",
				},
			})
			return
		}

		log.Printf("❌ 请求失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"type":    "network_error",
				"message": "Failed to execute request: " + err.Error(),
			},
		})
		return
	}
	defer common.CloseIO(resp.Body)

	// 检查响应是否需要解压缩
	var responseReader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")

	switch strings.ToLower(contentEncoding) {
	case "gzip":
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Printf("[Claude API] 创建gzip解压缩器失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": map[string]interface{}{
					"type":    "decompression_error",
					"message": "Failed to create gzip decompressor: " + err.Error(),
				},
			})
			return
		}
		defer gzipReader.Close()
		responseReader = gzipReader
	case "deflate":
		deflateReader := flate.NewReader(resp.Body)
		defer deflateReader.Close()
		responseReader = deflateReader
	}

	// 读取响应体
	var responseBody []byte
	var usageTokens *common.TokenUsage

	if resp.StatusCode >= 400 {
		// 错误响应，直接读取全部内容
		var readErr error
		responseBody, readErr = io.ReadAll(responseReader)
		if readErr != nil {
			log.Printf("❌ 读取错误响应失败: %v", readErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": map[string]interface{}{
					"type":    "response_read_error",
					"message": "Failed to read error response: " + readErr.Error(),
				},
			})
			return
		}

		// 调试日志：打印错误响应内容
		log.Printf("❌ 错误响应内容: %s", string(responseBody))
	}

	// 透传响应状态码
	c.Status(resp.StatusCode)

	// 透传所有响应头，但需要处理Content-Length以避免流式响应问题
	for name, values := range resp.Header {
		// 跳过Content-Length，让Gin自动处理流式响应
		if strings.ToLower(name) == "content-length" {
			continue
		}
		for _, value := range values {
			c.Header(name, value)
		}
	}

	if resp.StatusCode < 400 {
		// 成功响应：确保设置正确的流式响应头
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		if c.Writer.Header().Get("Content-Type") == "" {
			c.Header("Content-Type", "text/event-stream")
		}

		// 刷新响应头到客户端
		c.Writer.Flush()

		// 成功响应，使用流式解析 - 现在使用真正的流式转发
		usageTokens, err = common.ParseStreamResponse(c.Writer, responseReader)
		if err != nil {
			log.Println("stream copy and parse failed:", err.Error())
		}
	}

	// 如果是错误响应，写入固定503错误
	if resp.StatusCode >= 400 {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": map[string]interface{}{
				"type":    "response_error",
				"message": "Request failed with status " + strconv.Itoa(resp.StatusCode),
			},
		})
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
	accountService.UpdateAccountStatus(account, resp.StatusCode, usageTokens)

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

// TestsHandleClaudeRequest 用于测试的Claude请求处理函数，功能同HandleClaudeRequest但不更新日志和账号状态
// 主要用于单元测试和集成测试，避免对数据库和日志系统的
func TestsHandleClaudeRequest(account *model.Account) (int, string) {
	requestBody := `{
		"model": "claude-sonnet-4-20250514",
		"messages": [
			{
				"role": "user",
				"content": [
					{
						"type": "text",
						"text": "hi"
					}
				]
			}
		],
		"temperature": 1,
		"system": [
			{
				"type": "text",
				"text": "You are Claude Code, Anthropic's official CLI for Claude.",
				"cache_control": {
					"type": "ephemeral"
				}
			}
		],
		"metadata": {
			"user_id": "20b98a014e3182f9ce654e6c105432083cca392beb1416f6406508b56dc5f"
		},
		"max_tokens": 100,
		"stream": true
	}`

	body, _ := sjson.SetBytes([]byte(requestBody), "stream", true)

	// 获取有效的访问token
	accessToken, err := getValidAccessToken(account)
	if err != nil {
		return http.StatusInternalServerError, "Failed to get valid access token: " + err.Error()
	}

	req, err := http.NewRequest("POST", ClaudeAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, "Failed to create request: " + err.Error()
	}

	// 使用公共的请求头构建方法
	fixedHeaders := buildClaudeAPIHeaders(accessToken)

	for name, value := range fixedHeaders {
		req.Header.Set(name, value)
	}

	httpClientTimeout := 30 * time.Second
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if account.EnableProxy && account.ProxyURI != "" {
		proxyURL, err := url.Parse(account.ProxyURI)
		if err != nil {
			return http.StatusInternalServerError, "Invalid proxy URI: " + err.Error()
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   httpClientTimeout,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, "Request failed: " + err.Error()
	}
	defer common.CloseIO(resp.Body)

	// 打印响应内容
	if resp.StatusCode >= 400 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Println("Response Status:", resp.Status)
		log.Println("Response body:", string(responseBody))
	}
	return resp.StatusCode, ""
}

// buildClaudeAPIHeaders 构建Claude API请求头
func buildClaudeAPIHeaders(accessToken string) map[string]string {
	return map[string]string{
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

			// token已过期且刷新失败，禁用此账号
			log.Printf("token已过期且刷新失败，禁用账号: %s", account.Name)
			account.CurrentStatus = 2 // 设置为禁用状态
			if updateErr := model.UpdateAccount(account); updateErr != nil {
				log.Printf("禁用账号失败: %v", updateErr)
			} else {
				log.Printf("账号 %s 已被自动禁用", account.Name)
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
