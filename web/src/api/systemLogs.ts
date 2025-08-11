import { request } from '@/utils/request';

// API路径定义
const Api = {
  GetSystemLogs: '/api/v1/admin/logs',
};

// 系统日志记录类型
export interface SystemLog {
  id: number;
  method: string;
  path: string;
  status_code: number;
  user_id: number;
  ip: string;
  user_agent: string;
  request_id: string;
  duration: number; // 毫秒
  created_at: string;
  user?: {
    id: number;
    username: string;
  };
}

// 系统日志查询参数
export interface SystemLogQueryParams {
  page?: number;
  limit?: number;
}

// 系统日志列表响应
export interface SystemLogListResponse {
  logs: SystemLog[];
  total: number;
  page: number;
  limit: number;
}

/**
 * 获取系统日志列表
 */
export function getSystemLogs(params?: SystemLogQueryParams) {
  return request.get<SystemLogListResponse>({
    url: Api.GetSystemLogs,
    params,
  });
}
