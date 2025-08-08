import { request } from '@/utils/request';

// API路径定义
const Api = {
  GetMyLogs: '/api/v1/logs/my',
  GetMyStats: '/api/v1/logs/stats/my',
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
