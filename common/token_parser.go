package common

import (
	"bufio"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

// TokenUsage 表示token使用情况
type TokenUsage struct {
	InputTokens          int `json:"input_tokens"`
	OutputTokens         int `json:"output_tokens"`
	CacheReadInputTokens int `json:"cache_read_input_tokens"`
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

		// 解析token使用量
		if strings.HasPrefix(line, "data: ") {
			dataJSON := strings.TrimPrefix(line, "data: ")

			// 检查是否是message_delta事件
			if gjson.Get(dataJSON, "type").String() == "message_delta" {
				usageJSON := gjson.Get(dataJSON, "usage")
				if usageJSON.Exists() {
					// 解析各种token字段
					inputTokens := gjson.Get(dataJSON, "usage.input_tokens").Num
					outputTokens := gjson.Get(dataJSON, "usage.output_tokens").Num
					cacheReadInputTokens := gjson.Get(dataJSON, "usage.cache_read_input_tokens").Num

					// 设置token使用量（input_tokens和cache_read_input_tokens通常只在最后的delta中出现）
					usage.InputTokens = int(inputTokens)
					usage.OutputTokens += int(outputTokens) // 累加output tokens
					usage.CacheReadInputTokens = int(cacheReadInputTokens)
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
	}

	return usage, nil
}
