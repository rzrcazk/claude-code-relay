import { request } from '@/utils/request';

const Api = {
  GetList: '/api/v1/api-keys/list',
  Create: '/api/v1/api-keys/create',
  GetDetail: '/api/v1/api-keys/detail',
  Update: '/api/v1/api-keys/update',
  UpdateStatus: '/api/v1/api-keys/update-status',
  Delete: '/api/v1/api-keys/delete',
};

// API Key
export interface ApiKey {
  id: number;
  name: string;
  key: string;
  expires_at?: string;
  status: number; // 1: 启用 0: 禁用
  group_id: number;
  user_id: number;
  today_usage_count: number;
  today_input_tokens: number;
  today_output_tokens: number;
  today_cache_read_input_tokens: number;
  today_cache_creation_input_tokens: number;
  today_total_cost: number;
  model_restriction: string;
  daily_limit: number;
  last_used_time?: string;
  created_at: string;
  updated_at: string;
  group?: {
    id: number;
    name: string;
  };
}

// API Key列表
export interface ApiKeyListResult {
  api_keys: ApiKey[];
  total: number;
  page: number;
  limit: number;
}

// 创建API Key
export interface CreateApiKeyRequest {
  name: string;
  key?: string;
  expires_at?: string;
  status?: number;
  group_id?: number;
  model_restriction?: string;
  daily_limit?: number;
}

// 更新API Key
export interface UpdateApiKeyRequest {
  name?: string;
  expires_at?: string;
  status?: number;
  group_id?: number;
  model_restriction?: string;
  daily_limit?: number;
}

// 更新API Key状态
export interface UpdateApiKeyStatusRequest {
  status: number;
}

// 获取API Key列表
export function getApiKeys(params: { page?: number; limit?: number; group_id?: number } = {}) {
  return request.get<ApiKeyListResult>({
    url: Api.GetList,
    params,
  });
}

// 创建API Key
export function createApiKey(data: CreateApiKeyRequest) {
  return request.post<{ key: string }>({
    url: Api.Create,
    data,
  });
}

// 获取API Key
export function getApiKey(id: number) {
  return request.get<ApiKey>({
    url: `${Api.GetDetail}/${id}`,
  });
}

// 更新API Key
export function updateApiKey(id: number, data: UpdateApiKeyRequest) {
  return request.put<ApiKey>({
    url: `${Api.Update}/${id}`,
    data,
  });
}

// 更新API Key状态
export function updateApiKeyStatus(id: number, data: UpdateApiKeyStatusRequest) {
  return request.put({
    url: `${Api.UpdateStatus}/${id}`,
    data,
  });
}

// 删除API Key
export function deleteApiKey(id: number) {
  return request.delete({
    url: `${Api.Delete}/${id}`,
  });
}
