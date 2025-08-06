package relay

import (
	"bytes"
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

	// 原样转发响应体（支持流式输出）
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Println("response copy failed:", err.Error())
	}
}

func closeIO(c io.Closer) {
	err := c.Close()
	if nil != err {
		log.Println(err)
	}
}
