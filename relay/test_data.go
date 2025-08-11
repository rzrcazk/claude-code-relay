package relay

import "fmt"

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
		"user_id": "20b98a014e3182f9ce654e6c105432083cca392beb1416f6406508b56dc5f"
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
