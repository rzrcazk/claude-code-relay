package common

import (
	"bufio"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

// TokenUsage 表示token使用情况
type TokenUsage struct {
	InputTokens              int    `json:"input_tokens"`
	OutputTokens             int    `json:"output_tokens"`
	CacheReadInputTokens     int    `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int    `json:"cache_creation_input_tokens"`
	Model                    string `json:"model"`
}

// ParseStreamResponse 解析流式响应并提取token使用量
func ParseStreamResponse(dst io.Writer, src io.Reader) (*TokenUsage, error) {
	usage := &TokenUsage{}
	scanner := bufio.NewScanner(src)

	for scanner.Scan() {
		line := scanner.Text()

		// 写入原始数据到客户端
		if _, err := dst.Write([]byte(line + "\n")); err != nil {
			return usage, err
		}

		// 解析token使用量和模型信息
		if strings.HasPrefix(line, "data: ") {
			dataJSON := strings.TrimPrefix(line, "data: ")
			eventType := gjson.Get(dataJSON, "type").String()

			// 检查是否是message_start事件，解析model字段和使用量
			if eventType == "message_start" {
				model := gjson.Get(dataJSON, "message.model").String()
				if model != "" {
					usage.Model = model
				}

				// 从message_start中解析使用量
				usageJSON := gjson.Get(dataJSON, "message.usage")
				if usageJSON.Exists() {
					inputTokens := gjson.Get(dataJSON, "message.usage.input_tokens").Num
					outputTokens := gjson.Get(dataJSON, "message.usage.output_tokens").Num
					cacheReadInputTokens := gjson.Get(dataJSON, "message.usage.cache_read_input_tokens").Num
					cacheCreationInputTokens := gjson.Get(dataJSON, "message.usage.cache_creation_input_tokens").Num

					usage.InputTokens = int(inputTokens)
					usage.OutputTokens = int(outputTokens)
					usage.CacheReadInputTokens = int(cacheReadInputTokens)
					usage.CacheCreationInputTokens = int(cacheCreationInputTokens)
				}
			}

			// 检查是否是message_delta事件
			if eventType == "message_delta" {
				usageJSON := gjson.Get(dataJSON, "usage")
				if usageJSON.Exists() {
					// 解析各种token字段
					inputTokens := gjson.Get(dataJSON, "usage.input_tokens").Num
					outputTokens := gjson.Get(dataJSON, "usage.output_tokens").Num
					cacheReadInputTokens := gjson.Get(dataJSON, "usage.cache_read_input_tokens").Num
					cacheCreationInputTokens := gjson.Get(dataJSON, "usage.cache_creation_input_tokens").Num

					// 更新token使用量
					if inputTokens > 0 {
						usage.InputTokens = int(inputTokens)
					}
					usage.OutputTokens += int(outputTokens) // 累加output tokens
					if cacheReadInputTokens > 0 {
						usage.CacheReadInputTokens = int(cacheReadInputTokens)
					}
					if cacheCreationInputTokens > 0 {
						usage.CacheCreationInputTokens = int(cacheCreationInputTokens)
					}
				}
			}
		}
	}

	return usage, scanner.Err()
}

// ParseJSONResponse 解析非流式JSON响应中的token使用量
func ParseJSONResponse(responseBody []byte) (*TokenUsage, error) {
	usage := &TokenUsage{}

	// 从JSON响应中解析usage字段
	usageJSON := gjson.GetBytes(responseBody, "usage")
	if usageJSON.Exists() {
		usage.InputTokens = int(gjson.GetBytes(responseBody, "usage.input_tokens").Num)
		usage.OutputTokens = int(gjson.GetBytes(responseBody, "usage.output_tokens").Num)
		usage.CacheReadInputTokens = int(gjson.GetBytes(responseBody, "usage.cache_read_input_tokens").Num)
		usage.CacheCreationInputTokens = int(gjson.GetBytes(responseBody, "usage.cache_creation_input_tokens").Num)
	}

	// 解析model字段
	model := gjson.GetBytes(responseBody, "model").String()
	if model != "" {
		usage.Model = model
	}

	return usage, nil
}
