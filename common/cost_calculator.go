package common

import (
	"fmt"
)

// ModelPricing Claude模型价格配置 (USD per 1M tokens)
type ModelPricing struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheWrite float64 `json:"cache_write"`
	CacheRead  float64 `json:"cache_read"`
}

// CostDetails 费用详情
type CostDetails struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheWrite float64 `json:"cache_write"`
	CacheRead  float64 `json:"cache_read"`
	Total      float64 `json:"total"`
}

// FormattedCosts 格式化的费用字符串
type FormattedCosts struct {
	Input      string `json:"input"`
	Output     string `json:"output"`
	CacheWrite string `json:"cache_write"`
	CacheRead  string `json:"cache_read"`
	Total      string `json:"total"`
}

// UsageDetails 详细的token使用量数据（包含总计）
type UsageDetails struct {
	InputTokens       int `json:"input_tokens"`
	OutputTokens      int `json:"output_tokens"`
	CacheCreateTokens int `json:"cache_creation_input_tokens"`
	CacheReadTokens   int `json:"cache_read_input_tokens"`
	TotalTokens       int `json:"total_tokens"`
}

// CostCalculationResult 费用计算结果
type CostCalculationResult struct {
	Model     string         `json:"model"`
	Pricing   ModelPricing   `json:"pricing"`
	Usage     UsageDetails   `json:"usage"`
	Costs     CostDetails    `json:"costs"`
	Formatted FormattedCosts `json:"formatted"`
}

// SavingsResult 缓存节省信息
type SavingsResult struct {
	NormalCost        float64 `json:"normal_cost"`
	CacheCost         float64 `json:"cache_cost"`
	Savings           float64 `json:"savings"`
	SavingsPercentage float64 `json:"savings_percentage"`
	Formatted         struct {
		NormalCost        string `json:"normal_cost"`
		CacheCost         string `json:"cache_cost"`
		Savings           string `json:"savings"`
		SavingsPercentage string `json:"savings_percentage"`
	} `json:"formatted"`
}

// MODEL_PRICING Claude模型价格配置表
var MODEL_PRICING = map[string]ModelPricing{
	// Claude 3.5 Sonnet
	"claude-3-5-sonnet-20241022": {
		Input:      3.00,
		Output:     15.00,
		CacheWrite: 3.75,
		CacheRead:  0.30,
	},

	"claude-sonnet-4-20250514": {
		Input:      3.00,
		Output:     15.00,
		CacheWrite: 3.75,
		CacheRead:  0.30,
	},

	"claude-opus-4-20250514": {
		Input:      15.00,
		Output:     75.00,
		CacheWrite: 18.75,
		CacheRead:  1.50,
	},

	"claude-opus-4-1-20250805": {
		Input:      15.00,
		Output:     75.00,
		CacheWrite: 18.75,
		CacheRead:  1.50,
	},

	// Claude 3.5 Haiku
	"claude-3-5-haiku-20241022": {
		Input:      0.25,
		Output:     1.25,
		CacheWrite: 0.30,
		CacheRead:  0.03,
	},

	// Claude 3 Opus
	"claude-3-opus-20240229": {
		Input:      15.00,
		Output:     75.00,
		CacheWrite: 18.75,
		CacheRead:  1.50,
	},

	// Claude 3 Sonnet
	"claude-3-sonnet-20240229": {
		Input:      3.00,
		Output:     15.00,
		CacheWrite: 3.75,
		CacheRead:  0.30,
	},

	// Claude 3 Haiku
	"claude-3-haiku-20240307": {
		Input:      0.25,
		Output:     1.25,
		CacheWrite: 0.30,
		CacheRead:  0.03,
	},

	// 默认定价（用于未知模型）
	"unknown": {
		Input:      3.00,
		Output:     15.00,
		CacheWrite: 3.75,
		CacheRead:  0.30,
	},
}

// CostCalculator 费用计算器
type CostCalculator struct{}

// NewCostCalculator 创建费用计算器实例
func NewCostCalculator() *CostCalculator {
	return &CostCalculator{}
}

// CalculateCost 计算单次请求的费用
func (c *CostCalculator) CalculateCost(usage *TokenUsage) *CostCalculationResult {
	model := usage.Model
	if model == "" {
		model = "unknown"
	}

	// 获取定价信息
	pricing, exists := MODEL_PRICING[model]
	if !exists {
		pricing = MODEL_PRICING["unknown"]
	}

	// 计算各类型token的费用 (USD)
	inputCost := (float64(usage.InputTokens) / 1000000) * pricing.Input
	outputCost := (float64(usage.OutputTokens) / 1000000) * pricing.Output
	cacheWriteCost := (float64(usage.CacheCreationInputTokens) / 1000000) * pricing.CacheWrite
	cacheReadCost := (float64(usage.CacheReadInputTokens) / 1000000) * pricing.CacheRead

	totalCost := inputCost + outputCost + cacheWriteCost + cacheReadCost

	return &CostCalculationResult{
		Model:   model,
		Pricing: pricing,
		Usage: UsageDetails{
			InputTokens:       usage.InputTokens,
			OutputTokens:      usage.OutputTokens,
			CacheCreateTokens: usage.CacheCreationInputTokens,
			CacheReadTokens:   usage.CacheReadInputTokens,
			TotalTokens:       usage.InputTokens + usage.OutputTokens + usage.CacheCreationInputTokens + usage.CacheReadInputTokens,
		},
		Costs: CostDetails{
			Input:      inputCost,
			Output:     outputCost,
			CacheWrite: cacheWriteCost,
			CacheRead:  cacheReadCost,
			Total:      totalCost,
		},
		Formatted: FormattedCosts{
			Input:      c.FormatCost(inputCost),
			Output:     c.FormatCost(outputCost),
			CacheWrite: c.FormatCost(cacheWriteCost),
			CacheRead:  c.FormatCost(cacheReadCost),
			Total:      c.FormatCost(totalCost),
		},
	}
}

// CalculateAggregatedCost 计算聚合使用量的费用
func (c *CostCalculator) CalculateAggregatedCost(inputTokens, outputTokens, cacheCreateTokens, cacheReadTokens int, model string) *CostCalculationResult {
	usage := &TokenUsage{
		InputTokens:              inputTokens,
		OutputTokens:             outputTokens,
		CacheCreationInputTokens: cacheCreateTokens,
		CacheReadInputTokens:     cacheReadTokens,
		Model:                    model,
	}

	return c.CalculateCost(usage)
}

// GetModelPricing 获取模型定价信息
func (c *CostCalculator) GetModelPricing(model string) ModelPricing {
	if model == "" {
		model = "unknown"
	}

	pricing, exists := MODEL_PRICING[model]
	if !exists {
		pricing = MODEL_PRICING["unknown"]
	}
	return pricing
}

// GetAllModelPricing 获取所有支持的模型和定价
func (c *CostCalculator) GetAllModelPricing() map[string]ModelPricing {
	result := make(map[string]ModelPricing)
	for k, v := range MODEL_PRICING {
		result[k] = v
	}
	return result
}

// IsModelSupported 验证模型是否支持
func (c *CostCalculator) IsModelSupported(model string) bool {
	_, exists := MODEL_PRICING[model]
	return exists
}

// FormatCost 格式化费用显示
func (c *CostCalculator) FormatCost(cost float64) string {
	if cost >= 1 {
		return fmt.Sprintf("$%.2f", cost)
	} else if cost >= 0.001 {
		return fmt.Sprintf("$%.4f", cost)
	} else {
		return fmt.Sprintf("$%.6f", cost)
	}
}

// CalculateCacheSavings 计算费用节省（使用缓存的节省）
func (c *CostCalculator) CalculateCacheSavings(usage *TokenUsage) *SavingsResult {
	pricing := c.GetModelPricing(usage.Model)
	cacheReadTokens := usage.CacheReadInputTokens

	// 如果这些token不使用缓存，需要按正常input价格计费
	normalCost := (float64(cacheReadTokens) / 1000000) * pricing.Input
	cacheCost := (float64(cacheReadTokens) / 1000000) * pricing.CacheRead
	savings := normalCost - cacheCost
	savingsPercentage := 0.0
	if normalCost > 0 {
		savingsPercentage = (savings / normalCost) * 100
	}

	result := &SavingsResult{
		NormalCost:        normalCost,
		CacheCost:         cacheCost,
		Savings:           savings,
		SavingsPercentage: savingsPercentage,
	}

	result.Formatted.NormalCost = c.FormatCost(normalCost)
	result.Formatted.CacheCost = c.FormatCost(cacheCost)
	result.Formatted.Savings = c.FormatCost(savings)
	result.Formatted.SavingsPercentage = fmt.Sprintf("%.1f%%", savingsPercentage)

	return result
}

// 全局费用计算器实例
var GlobalCostCalculator = NewCostCalculator()

// 便利函数，使用全局实例
func CalculateCost(usage *TokenUsage) *CostCalculationResult {
	return GlobalCostCalculator.CalculateCost(usage)
}

func CalculateAggregatedCost(inputTokens, outputTokens, cacheCreateTokens, cacheReadTokens int, model string) *CostCalculationResult {
	return GlobalCostCalculator.CalculateAggregatedCost(inputTokens, outputTokens, cacheCreateTokens, cacheReadTokens, model)
}

func GetModelPricing(model string) ModelPricing {
	return GlobalCostCalculator.GetModelPricing(model)
}

func GetAllModelPricing() map[string]ModelPricing {
	return GlobalCostCalculator.GetAllModelPricing()
}

func IsModelSupported(model string) bool {
	return GlobalCostCalculator.IsModelSupported(model)
}

func FormatCost(cost float64) string {
	return GlobalCostCalculator.FormatCost(cost)
}

func CalculateCacheSavings(usage *TokenUsage) *SavingsResult {
	return GlobalCostCalculator.CalculateCacheSavings(usage)
}
