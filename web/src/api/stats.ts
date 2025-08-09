import { request } from '@/utils/request';

// API Key统计查询接口
export interface ApiKeyStatsParams {
  api_key: string;
  page?: number;
  limit?: number;
}

export interface ApiKeyStatsResponse {
  api_key_info: {
    id: number;
    name: string;
    status: number;
  };
  stats: {
    summary: {
      total_requests: number;
      total_input_tokens: number;
      total_output_tokens: number;
      total_cache_read_tokens: number;
      total_cache_creation_tokens: number;
      total_tokens: number;
      total_cost: number;
      input_cost: number;
      output_cost: number;
      cache_write_cost: number;
      cache_read_cost: number;
      avg_duration: number;
      stream_requests: number;
      stream_percent: number;
    };
    trend_data: Array<{
      date: string;
      requests: number;
      tokens: number;
      cost: number;
      avg_duration: number;
      cache_tokens: number;
      input_tokens: number;
      output_tokens: number;
    }>;
  };
  logs: {
    list: Array<{
      id: string;
      model_name: string;
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
    }>;
    total: number;
    page: number;
    limit: number;
  };
}

const Api = {
  GetApiKeyStats: '/api/v1/auth/api-key',
};

/**
 * 获取API Key统计数据
 */
export function getApiKeyStats(params: ApiKeyStatsParams) {
  return request.get<ApiKeyStatsResponse>({
    url: Api.GetApiKeyStats,
    params,
  });
}