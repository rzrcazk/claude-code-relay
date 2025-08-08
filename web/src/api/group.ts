import { request } from '@/utils/request';

// API 分组接口
const Api = {
  GetList: '/v1/groups/list',
  GetAll: '/v1/groups/all',
  Create: '/v1/groups/create',
  Update: '/v1/groups/update',
  Delete: '/v1/groups/delete',
  GetDetail: '/v1/groups/detail',
};

// 分组列表
export interface Group {
  id: number;
  name: string;
  remark: string; // 对应后端的remark字段
  status: number; // 0: 禁用, 1: 启用
  user_id: number;
  created_at: string;
  updated_at: string;
  api_key_count: number; // API密钥数量
  account_count: number; // 账号数量
}

// 创建分组
export interface GroupCreateParams {
  name: string;
  remark?: string;
  status?: number;
}

export interface GroupUpdateParams extends GroupCreateParams {
  id: number;
}

// 分组列表参数
export interface GroupListParams {
  page?: number;
  size?: number;
  name?: string;
  status?: number;
}

// 分组列表响应
export interface GroupListResponse {
  groups: Group[];
  total: number;
  page: number;
  limit: number;
}

/**
 * 获取分组列表
 */
export function getGroupList(params?: GroupListParams) {
  return request.get<GroupListResponse>({
    url: Api.GetList,
    params,
  });
}

/**
 * 获取所有分组（用于下拉选择）
 */
export function getAllGroups() {
  return request.get<Group[]>({
    url: Api.GetAll,
  });
}

/**
 * 获取分组详情
 */
export function getGroupDetail(id: number) {
  return request.get<Group>({
    url: `${Api.GetDetail}/${id}`,
  });
}

/**
 * 创建分组
 */
export function createGroup(data: GroupCreateParams) {
  return request.post<Group>({
    url: Api.Create,
    data,
  });
}

/**
 * 更新分组
 */
export function updateGroup(data: GroupUpdateParams) {
  return request.put<Group>({
    url: `${Api.Update}/${data.id}`,
    data,
  });
}

/**
 * 删除分组
 */
export function deleteGroup(id: number) {
  return request.delete({
    url: `${Api.Delete}/${id}`,
  });
}

/**
 * 批量删除分组
 */
export function batchDeleteGroups(ids: number[]) {
  return request.delete({
    url: Api.Delete,
    data: { ids },
  });
}

/**
 * 更新分组状态 - 注意：后端分组路由中没有专门的状态更新接口，使用更新接口
 */
export function updateGroupStatus(id: number, status: number) {
  return request.put({
    url: `${Api.Update}/${id}`,
    data: { status },
  });
}

/**
 * 批量更新分组状态
 */
export function batchUpdateGroupStatus(ids: number[], status: number) {
  return request.put({
    url: Api.Update,
    data: { ids, status },
  });
}
