package common

import (
	"fmt"
	"strings"
)

// TestRequestBodyTemplate 测试用的标准请求体模板
const TestRequestBodyTemplate = `{
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
		"user_id": "user_6f30b450fd7e45f22fd40f61e2b2617d2d8aa13ececfb5c00c879caec1c0b00d_account__session_d09624bf-09fc-4cdd-afde-cee6abaafb23"
	},
	"max_tokens": %d,
	"stream": true
}`

// GetTestRequestBody 获取带指定max_tokens的测试请求体
func GetTestRequestBody(maxTokens int) string {
	return fmt.Sprintf(TestRequestBodyTemplate, maxTokens)
}

// TestRequestBody 默认测试请求体（64000 tokens）
var TestRequestBody = GetTestRequestBody(64000)

// getGlobalClaudeCodeHeaders 获取全局Claude Code请求头
func getGlobalClaudeCodeHeaders() map[string]string {
	return map[string]string{
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

// MergeHeaders 合并全局Claude Code请求头和用户提供的请求头
// 用户提供的头部优先级更高，可以覆盖全局头部
func MergeHeaders(customRequestHeaders map[string]string, anthropicBeta string) map[string]string {
	globalHeaders := getGlobalClaudeCodeHeaders()

	result := make(map[string]string, len(globalHeaders)+len(customRequestHeaders))

	for k, v := range globalHeaders {
		result[k] = v
	}

	// 用户提供的头部优先级更高，可以覆盖全局头部
	for k, v := range customRequestHeaders {
		result[k] = v
	}

	// 使用原始的 anthropic-beta 请求头（如果存在）
	if anthropicBeta != "" {
		if strings.Contains(anthropicBeta, "oauth-") {
			result["anthropic-beta"] = anthropicBeta
		} else {
			result["anthropic-beta"] = "oauth-2025-04-20," + anthropicBeta
		}
	}

	return result
}
