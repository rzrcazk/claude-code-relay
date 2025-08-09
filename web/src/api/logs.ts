import { request } from '@/utils/request';

// API路径定义
const Api = {
  GetMyLogs: '/api/v1/logs/my',
  GetMyStats: '/api/v1/logs/stats/my',
  GetMyUsageStats: '/api/v1/logs/usage-stats/my',
  GetLogDetail: '/api/v1/logs/detail',
};

// 日志记录类型
export interface Log {
  id: string;
  model_name: string;
  account_id: number;
  user_id: number;
  api_key_id: number;
  input_tokens: number;
  output_tokens: number;
  cache_read_input_tokens: number;
  cache_creation_input_tokens: number;
  input_cost: number;
  output_cost: number;
  cache_write_cost: number;
  cache_read_cost: number;
  total_cost: number;
  is_stream: boolean;
  duration: number;
  created_at: string;
  user?: {
    id: number;
    username: string;
  };
  api_key?: {
    id: number;
    name: string;
  };
}

// 日志查询参数
export interface LogQueryParams {
  page?: number;
  limit?: number;
  account_id?: number;
  api_key_id?: number;
  model_name?: string;
  is_stream?: boolean;
  start_time?: string; // 格式: 2024-01-01 15:04:05
  end_time?: string; // 格式: 2024-01-01 15:04:05
  min_cost?: number;
  max_cost?: number;
}

// 日志列表响应
export interface LogListResponse {
  logs: Log[];
  total: number;
  page: number;
  limit: number;
}

// 日志统计结果
export interface LogStatsResult {
  total_requests: number;
  total_tokens: number;
  total_cost: number;
  avg_duration: number;
  stream_requests: number;
  stream_percent: number;
}

// 统计查询参数
export interface StatsQueryParams {
  user_id?: number; // 用户ID筛选
  account_filter?: string; // 账号筛选（ID或邮箱/名称）
  api_key_filter?: string; // API Key筛选（ID或秘钥值）
  model_name?: string; // 模型名称筛选
  start_time?: string; // 开始时间
  end_time?: string; // 结束时间
}

// 详细统计结果
export interface DetailedStatsResult {
  total_requests: number; // 总请求数
  total_input_tokens: number; // 总输入tokens
  total_output_tokens: number; // 总输出tokens
  total_cache_read_tokens: number; // 总缓存读取tokens
  total_cache_creation_tokens: number; // 总缓存创建tokens
  total_tokens: number; // 总tokens数
  total_cost: number; // 总费用
  input_cost: number; // 输入费用
  output_cost: number; // 输出费用
  cache_write_cost: number; // 缓存写入费用
  cache_read_cost: number; // 缓存读取费用
  avg_duration: number; // 平均响应时间
  stream_requests: number; // 流式请求数
  stream_percent: number; // 流式请求比例
}

// 趋势数据项
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

// 统计响应结果
export interface StatsResponse {
  summary: DetailedStatsResult; // 汇总统计
  trend_data: TrendDataItem[]; // 趋势数据
}

/**
 * 获取我的日志列表
 */
export function getMyLogs(params?: LogQueryParams) {
  return request.get<LogListResponse>({
    url: Api.GetMyLogs,
    params,
  });
}

/**
 * 获取我的日志统计
 */
export function getMyLogStats(params?: Omit<LogQueryParams, 'page' | 'limit'>) {
  return request.get<LogStatsResult>({
    url: Api.GetMyStats,
    params,
  });
}

/**
 * 获取日志详情
 */
export function getLogDetail(id: string) {
  return request.get<Log>({
    url: `${Api.GetLogDetail}/${id}`,
  });
}

/**
 * 获取我的使用统计数据
 */
export function getMyUsageStats(params: StatsQueryParams) {
  return request.get<StatsResponse>({
    url: Api.GetMyUsageStats,
    params,
  });
}
