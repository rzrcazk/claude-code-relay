import { request } from '@/utils/request';

// API路径定义
const Api = {
  GetDashboardStats: '/api/v1/dashboard/stats',
};

// 模型使用统计项
export interface ModelUsageItem {
  model_name: string; // 模型名称
  requests: number; // 请求数
  tokens: number; // tokens数
  cost: number; // 费用
}

// 账号排名项
export interface AccountRankItem {
  account_id: number; // 账号ID
  account_name: string; // 账号名称
  platform_type: string; // 平台类型
  requests: number; // 请求数
  tokens: number; // tokens数
  cost: number; // 费用
  growth_rate: number; // 增长率(%)
}

// API Key排名项
export interface ApiKeyRankItem {
  api_key_id: number; // API Key ID
  api_key_name: string; // API Key名称
  requests: number; // 请求数
  tokens: number; // tokens数
  cost: number; // 费用
  growth_rate: number; // 增长率(%)
}

// 每日统计项
export interface DayStatsItem {
  date: string; // 日期
  requests: number; // 请求数
  tokens: number; // tokens数
  cost: number; // 费用
}

// 趋势数据项（复用logs.ts中的）
export interface TrendDataItem {
  date: string; // 日期
  requests: number; // 请求数
  tokens: number; // tokens数
  cost: number; // 费用
  avg_duration: number; // 平均响应时间
  cache_tokens: number; // 缓存tokens
  input_tokens: number; // 输入tokens
  output_tokens: number; // 输出tokens
}

// 仪表盘统计数据
export interface DashboardStats {
  // 顶部面板数据
  total_cost: number; // 总费用(USD)
  total_tokens: number; // 总Tokens
  user_count: number; // 用户数量
  api_key_count: number; // API Key数量
  
  // 趋势数据
  trend_data: TrendDataItem[]; // 使用趋势
  
  // 模型使用分布
  model_stats: ModelUsageItem[]; // 模型使用统计
  
  // 账号排名
  account_ranking: AccountRankItem[]; // 账号排名
  
  // API Key排名
  api_key_ranking: ApiKeyRankItem[]; // API Key排名
  
  // 今日vs昨日数据对比
  today_stats: DayStatsItem; // 今日统计
  yesterday_stats: DayStatsItem; // 昨日统计
}

/**
 * 获取仪表盘统计数据
 */
export function getDashboardStats() {
  return request.get<DashboardStats>({
    url: Api.GetDashboardStats,
  });
}