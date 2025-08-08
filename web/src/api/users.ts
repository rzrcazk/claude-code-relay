import { request } from '@/utils/request';

const Api = {
  GetList: '/api/v1/admin/users',
  Create: '/api/v1/admin/users',
  UpdateStatus: '/api/v1/admin/users/:id/status',
};

// 用户状态枚举
export const UserStatus = {
  DISABLED: 0, // 禁用
  ENABLED: 1, // 启用
} as const;

// 用户状态标签
export const UserStatusLabels = {
  [UserStatus.DISABLED]: '禁用',
  [UserStatus.ENABLED]: '启用',
} as const;

// 用户角色枚举
export const UserRole = {
  ADMIN: 'admin',
  USER: 'user',
} as const;

// 用户基本信息接口
export interface UserInfo {
  id: number;
  username: string;
  email: string;
  status: number;
  role: string;
  created_at: string;
  updated_at: string;
}

// 用户列表响应接口
export interface UserListResult {
  users: UserInfo[];
  total: number;
  page: number;
  limit: number;
}

// 用户列表查询参数
export interface UserListParams {
  page?: number;
  limit?: number;
}

// 创建用户请求参数
export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  role: string;
}

// 更新用户状态请求参数
export interface UpdateUserStatusRequest {
  status: number;
}

/**
 * 获取用户列表
 */
export function getUserList(params: UserListParams = {}) {
  return request.get<UserListResult>({
    url: Api.GetList,
    params,
  });
}

/**
 * 创建用户
 */
export function createUser(data: CreateUserRequest) {
  return request.post({
    url: Api.Create,
    data,
  });
}

/**
 * 更新用户状态
 */
export function updateUserStatus(id: number, data: UpdateUserStatusRequest) {
  return request.put({
    url: Api.UpdateStatus.replace(':id', String(id)),
    data,
  });
}
