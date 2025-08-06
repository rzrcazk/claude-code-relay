package relay

import (
	"bytes"
	"claude-code-relay/common"
	"claude-code-relay/model"
	"context"
	"crypto/tls"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// HandleClaudeConsoleRequest 处理Claude Console平台的请求
func HandleClaudeConsoleRequest(c *gin.Context, account *model.Account) {
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

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, account.RequestURL+"/v1/messages?beta=true", bytes.NewBuffer(body))
	if nil != err {
		c.AbortWithStatus(http.StatusInternalServerError)
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
		"anthropic-beta":                            "fine-grained-tool-streaming-2025-05-14",
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
	isStream := gjson.GetBytes(body, "stream").Bool()
	if isStream && c.Request.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/event-stream")
	}

	httpClientTimeout, _ := time.ParseDuration(os.Getenv("HTTP_CLIENT_TIMEOUT") + "s")

	// 创建基础Transport配置
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 如果提供了代理URI，配置代理
	if account.ProxyURI != "" {
		proxyURL, err := url.Parse(account.ProxyURI)
		if err != nil {
			log.Printf("invalid proxy URI: %s", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
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
			c.AbortWithStatus(http.StatusRequestTimeout)
			return
		}

		log.Println("request conversation failed:", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer closeIO(resp.Body)

	// 透传响应状态码
	c.Status(resp.StatusCode)

	// 透传响应头
	for name, values := range resp.Header {
		for _, value := range values {
			c.Header(name, value)
		}
	}

	// 检查是否是流式响应
	isStream = gjson.GetBytes(body, "stream").Bool()
	var usageTokens *common.TokenUsage

	if isStream && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// 流式响应，需要解析token使用量
		usageTokens, err = common.ParseStreamResponse(c.Writer, resp.Body)
		if err != nil {
			log.Println("stream copy and parse failed:", err.Error())
		}
	} else {
		// 非流式响应，直接转发
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			log.Println("response copy failed:", err.Error())
		}
	}

	// 处理响应状态码并更新账号状态
	go updateAccountStatus(account, resp.StatusCode, usageTokens)

	// 更新API Key统计信息
	if apiKey != nil {
		go updateApiKeyStatus(apiKey, resp.StatusCode, usageTokens)
	}
}

// updateAccountStatus 根据响应状态码更新账号状态
func updateAccountStatus(account *model.Account, statusCode int, usage *common.TokenUsage) {
	// 根据状态码设置CurrentStatus
	switch {
	case statusCode == 429:
		// 限流状态
		account.CurrentStatus = 3
	case statusCode > 400:
		// 接口异常
		account.CurrentStatus = 2
	case statusCode == 200 || statusCode == 201:
		// 正常状态
		account.CurrentStatus = 1

		// 请求成功时更新最后使用时间和今日使用次数
		now := time.Now()

		// 判断最后使用时间是否为当天
		if account.LastUsedTime != nil {
			lastUsedDate := time.Time(*account.LastUsedTime).Format("2006-01-02")
			todayDate := now.Format("2006-01-02")

			if lastUsedDate == todayDate {
				// 同一天，使用次数+1
				account.TodayUsageCount++
			} else {
				// 不同天，重置为1
				account.TodayUsageCount = 1
			}
		} else {
			// 首次使用，设置为1
			account.TodayUsageCount = 1
		}

		// 更新token使用量（如果有的话）
		if usage != nil {
			if account.LastUsedTime != nil {
				lastUsedDate := time.Time(*account.LastUsedTime).Format("2006-01-02")
				todayDate := now.Format("2006-01-02")

				if lastUsedDate == todayDate {
					// 同一天，累加各类tokens
					account.TodayInputTokens += usage.InputTokens
					account.TodayOutputTokens += usage.OutputTokens
					account.TodayCacheReadInputTokens += usage.CacheReadInputTokens
				} else {
					// 不同天，重置各类tokens
					account.TodayInputTokens = usage.InputTokens
					account.TodayOutputTokens = usage.OutputTokens
					account.TodayCacheReadInputTokens = usage.CacheReadInputTokens
				}
			} else {
				// 首次使用，设置各类tokens
				account.TodayInputTokens = usage.InputTokens
				account.TodayOutputTokens = usage.OutputTokens
				account.TodayCacheReadInputTokens = usage.CacheReadInputTokens
			}
		}

		// 更新最后使用时间
		nowTime := model.Time(now)
		account.LastUsedTime = &nowTime
	default:
		// 其他状态码保持原状态
		return
	}

	// 更新数据库
	if err := model.UpdateAccount(account); err != nil {
		log.Printf("failed to update account status: %v", err)
	}
}

// updateApiKeyStatus 根据响应状态码更新API Key统计信息
func updateApiKeyStatus(apiKey *model.ApiKey, statusCode int, usage *common.TokenUsage) {
	// 只在请求成功时更新API Key统计信息
	if statusCode != 200 && statusCode != 201 {
		return
	}

	now := time.Now()

	// 判断最后使用时间是否为当天
	if apiKey.LastUsedTime != nil {
		lastUsedDate := time.Time(*apiKey.LastUsedTime).Format("2006-01-02")
		todayDate := now.Format("2006-01-02")

		if lastUsedDate == todayDate {
			// 同一天，使用次数+1
			apiKey.TodayUsageCount++
		} else {
			// 不同天，重置为1
			apiKey.TodayUsageCount = 1
		}
	} else {
		// 首次使用，设置为1
		apiKey.TodayUsageCount = 1
	}

	// 更新token使用量（如果有的话）
	if usage != nil {
		if apiKey.LastUsedTime != nil {
			lastUsedDate := time.Time(*apiKey.LastUsedTime).Format("2006-01-02")
			todayDate := now.Format("2006-01-02")

			if lastUsedDate == todayDate {
				// 同一天，累加各类tokens
				apiKey.TodayInputTokens += usage.InputTokens
				apiKey.TodayOutputTokens += usage.OutputTokens
				apiKey.TodayCacheReadInputTokens += usage.CacheReadInputTokens
			} else {
				// 不同天，重置各类tokens
				apiKey.TodayInputTokens = usage.InputTokens
				apiKey.TodayOutputTokens = usage.OutputTokens
				apiKey.TodayCacheReadInputTokens = usage.CacheReadInputTokens
			}
		} else {
			// 首次使用，设置各类tokens
			apiKey.TodayInputTokens = usage.InputTokens
			apiKey.TodayOutputTokens = usage.OutputTokens
			apiKey.TodayCacheReadInputTokens = usage.CacheReadInputTokens
		}
	}

	// 更新最后使用时间
	nowTime := model.Time(now)
	apiKey.LastUsedTime = &nowTime

	// 更新数据库
	if err := model.UpdateApiKey(apiKey); err != nil {
		log.Printf("failed to update api key status: %v", err)
	}
}

func closeIO(c io.Closer) {
	err := c.Close()
	if nil != err {
		log.Println(err)
	}
}
