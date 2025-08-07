package common

import (
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

// StreamCopyWriter 实现真正的流式转发，边转发边解析
type StreamCopyWriter struct {
	dst       io.Writer
	usage     *TokenUsage
	buffer    []byte
	remainder string
}

// Write 实现 io.Writer 接口，实现边转发边解析
func (w *StreamCopyWriter) Write(p []byte) (n int, err error) {
	// 先写入目标，实现真正的流式转发
	n, err = w.dst.Write(p)
	if err != nil {
		return n, err
	}

	// 智能刷新：检查数据中是否包含换行符，如果包含则刷新
	// 这样可以减少不必要的flush调用，同时保证及时性
	if strings.Contains(string(p[:n]), "\n") {
		if flusher, ok := w.dst.(interface{ Flush() }); ok {
			flusher.Flush()
		}
	}

	// 将数据添加到缓冲区进行解析
	w.buffer = append(w.buffer, p[:n]...)

	// 按行处理数据
	data := string(w.buffer)
	lines := strings.Split(w.remainder+data, "\n")

	// 保留最后一行（可能不完整）
	if len(lines) > 0 {
		w.remainder = lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	}

	// 处理完整的行
	for _, line := range lines {
		w.parseLine(line)
	}

	// 清空缓冲区
	w.buffer = w.buffer[:0]

	return n, nil
}

// parseLine 解析单行数据提取token使用量
func (w *StreamCopyWriter) parseLine(line string) {
	// 解析token使用量和模型信息
	if strings.HasPrefix(line, "data: ") {
		dataJSON := strings.TrimPrefix(line, "data: ")
		eventType := gjson.Get(dataJSON, "type").String()

		// 检查是否是message_start事件，解析model字段和使用量
		if eventType == "message_start" {
			model := gjson.Get(dataJSON, "message.model").String()
			if model != "" {
				w.usage.Model = model
			}

			// 从message_start中解析使用量
			usageJSON := gjson.Get(dataJSON, "message.usage")
			if usageJSON.Exists() {
				inputTokens := gjson.Get(dataJSON, "message.usage.input_tokens").Num
				outputTokens := gjson.Get(dataJSON, "message.usage.output_tokens").Num
				cacheReadInputTokens := gjson.Get(dataJSON, "message.usage.cache_read_input_tokens").Num
				cacheCreationInputTokens := gjson.Get(dataJSON, "message.usage.cache_creation_input_tokens").Num

				w.usage.InputTokens = int(inputTokens)
				w.usage.OutputTokens = int(outputTokens)
				w.usage.CacheReadInputTokens = int(cacheReadInputTokens)
				w.usage.CacheCreationInputTokens = int(cacheCreationInputTokens)
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
					w.usage.InputTokens = int(inputTokens)
				}
				w.usage.OutputTokens += int(outputTokens) // 累加output tokens
				if cacheReadInputTokens > 0 {
					w.usage.CacheReadInputTokens = int(cacheReadInputTokens)
				}
				if cacheCreationInputTokens > 0 {
					w.usage.CacheCreationInputTokens = int(cacheCreationInputTokens)
				}
			}
		}
	}
}

// ParseStreamResponse 解析流式响应并提取token使用量 - 实现真正的流式转发
func ParseStreamResponse(dst io.Writer, src io.Reader) (*TokenUsage, error) {
	usage := &TokenUsage{}

	// 创建流式拷贝写入器
	streamWriter := &StreamCopyWriter{
		dst:   dst,
		usage: usage,
	}

	// 使用优化的缓冲区大小，平衡性能和及时性
	buffer := make([]byte, 4096) // 使用4KB缓冲区
	for {
		n, err := src.Read(buffer)
		if n > 0 {
			_, writeErr := streamWriter.Write(buffer[:n])
			if writeErr != nil {
				return usage, writeErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return usage, err
		}
	}

	// 处理最后一行（如果有的话）
	if streamWriter.remainder != "" {
		streamWriter.parseLine(streamWriter.remainder)
	}

	return usage, nil
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
