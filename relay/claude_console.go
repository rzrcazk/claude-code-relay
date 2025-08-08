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
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// HandleClaudeConsoleRequest 处理Claude Console平台的请求
func HandleClaudeConsoleRequest(c *gin.Context, account *model.Account) {
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
				"type":    "request_body_error",
				"message": "Failed to read request body: " + err.Error(),
			},
		})
		return
	}

	body, _ = sjson.SetBytes(body, "stream", true) // 强制流式输出

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, account.RequestURL+"/v1/messages?beta=true", bytes.NewBuffer(body))
	if nil != err {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": map[string]interface{}{
				"type":    "internal_server_error",
				"message": "Failed to create request: " + err.Error(),
			},
		})
		return
	}

	// 固定请求头配置，如果原请求中没有则使用固定值
	fixedHeaders := map[string]string{
		"x-api-key":                                 account.SecretKey,
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

	// 删除Authorization 请求头
	req.Header.Del("Authorization")
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

	// 如果提供了代理URI，配置代理
	if account.ProxyURI != "" {
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

		log.Println("request conversation failed:", err.Error())
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
			log.Printf("[Claude Console] 创建gzip解压缩器失败: %v", err)
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

	// 透传响应状态码
	c.Status(resp.StatusCode)

	// 透传响应头，但需要处理Content-Length以避免流式响应问题
	for name, values := range resp.Header {
		// 跳过Content-Length，让Gin自动处理流式响应
		if strings.ToLower(name) == "content-length" {
			continue
		}
		for _, value := range values {
			c.Header(name, value)
		}
	}

	var usageTokens *common.TokenUsage

	if resp.StatusCode >= 400 {
		// 错误响应，直接读取全部内容
		responseBody, readErr := io.ReadAll(responseReader)
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
		log.Printf("❌ Claude Console错误响应内容: %s", string(responseBody))
		c.Writer.Write(responseBody)
	} else {
		// 确保设置正确的流式响应头
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		if c.Writer.Header().Get("Content-Type") == "" {
			c.Header("Content-Type", "text/event-stream")
		}

		// 刷新响应头到客户端
		c.Writer.Flush()

		// 解析token使用量 - 现在使用真正的流式转发
		usageTokens, err = common.ParseStreamResponse(c.Writer, responseReader)
		if err != nil {
			log.Println("stream copy and parse failed:", err.Error())
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

// TestHandleClaudeConsoleRequest 测试处理Claude Console请求的函数
func TestHandleClaudeConsoleRequest(account *model.Account) (int, string) {
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
		"max_tokens": 64000,
		"stream": true
	}`

	body, _ := sjson.SetBytes([]byte(requestBody), "stream", true)

	req, err := http.NewRequest("POST", account.RequestURL+"/v1/messages?beta=true", bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, "Failed to create request: " + err.Error()
	}

	fixedHeaders := map[string]string{
		"x-api-key":                   account.SecretKey,
		"anthropic-version":           "2023-06-01",
		"Content-Type":                "application/json",
		"Accept":                      "text/event-stream",
		"Authorization":               "Bearer " + account.SecretKey, // 一些环境可能会检查这个头
		"X-Stainless-Retry-Count":     "0",
		"X-Stainless-Timeout":         "600",
		"X-Stainless-Lang":            "js",
		"X-Stainless-Package-Version": "0.55.1",
		"X-Stainless-OS":              "MacOS",
		"X-Stainless-Arch":            "arm64",
		"X-Stainless-Runtime":         "node",
		"x-stainless-helper-method":   "stream",
		"x-app":                       "cli",
		"User-Agent":                  "claude-cli/1.0.44 (external, cli)",
		"anthropic-beta":              "claude-code-20250219,oauth-2025-04-20,interleaved-thinking-2025-05-14,fine-grained-tool-streaming-2025-05-14",
		"X-Stainless-Runtime-Version": "v20.18.1",
		"anthropic-dangerous-direct-browser-access": "true",
	}

	for name, value := range fixedHeaders {
		req.Header.Set(name, value)
	}

	httpClientTimeout := 30 * time.Second
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if account.ProxyURI != "" {
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

	return resp.StatusCode, ""
}
