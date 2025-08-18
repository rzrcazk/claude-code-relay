import { request } from '@/utils/request';

const Api = {
  // 登录
  Login: '/api/v1/auth/login',
  Register: '/api/v1/auth/register',
  SendVerificationCode: '/api/v1/auth/send-verification-code',

  // 获取用户信息
  GetProfile: '/api/v1/user/profile',
  UpdateProfile: '/api/v1/user/profile',
  ChangeEmail: '/api/v1/user/change-email',
  ChangePassword: '/api/v1/user/change-password',

  // 获取用户列表
  GetUsers: '/api/v1/admin/users',
  AdminCreateUser: '/api/v1/admin/users',
  AdminUpdateUserStatus: '/api/v1/admin/users/{id}/status',
};

// 用户信息
export interface UserInfo {
  id: number;
  name: string;
  username: string;
  email: string;
  role: string;
  status: number;
}

// 用户信息
export interface UserProfile {
  id: number;
  name: string;
  username: string;
  email: string;
  role: string;
  status: number;
  created_at: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  status: number;
  created_at: string;
  updated_at: string;
}

// 用户列表
export interface UserListResult {
  users: User[];
  total: number;
  page: number;
  limit: number;
}

// 登录结果
export interface LoginResult {
  token: string;
  user: UserInfo;
}

// 登录请求
export interface LoginRequest {
  username?: string;
  email?: string;
  password?: string;
  verification_code?: string;
  login_type: 'password' | 'sms_code';
}

// 注册请求
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  verification_code: string;
}

// 发送验证码请求
export interface SendVerificationCodeRequest {
  email: string;
  type: 'register' | 'login' | 'reset_password' | 'change_email';
}

// 修改邮箱请求（暂时简化，实际后端需要密码和验证码）
export interface ChangeEmailRequest {
  email: string;
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

// 更新用户信息请求
export interface UpdateProfileRequest {
  name?: string;
  username?: string;
  email?: string;
}

export interface AdminCreateUserRequest {
  username: string;
  email: string;
  password: string;
  role: string;
}

export interface AdminUpdateUserStatusRequest {
  status: number;
}

// 登录
export function login(data: LoginRequest) {
  return request.post<LoginResult>({
    url: Api.Login,
    data,
  });
}

// 注册
export function register(data: RegisterRequest) {
  return request.post({
    url: Api.Register,
    data,
  });
}

// 发送验证码
export function sendVerificationCode(data: SendVerificationCodeRequest) {
  return request.post({
    url: Api.SendVerificationCode,
    data,
  });
}

// 获取用户信息
export function getUserProfile() {
  return request.get<UserProfile>({
    url: Api.GetProfile,
  });
}

// 更新用户信息
export function updateUserProfile(data: UpdateProfileRequest) {
  return request.put({
    url: Api.UpdateProfile,
    data,
  });
}

// 修改邮箱
export function changeUserEmail(data: ChangeEmailRequest) {
  return request.put({
    url: Api.ChangeEmail,
    data,
  });
}

// 修改密码
export function changeUserPassword(data: ChangePasswordRequest) {
  return request.put({
    url: Api.ChangePassword,
    data,
  });
}

// 获取用户列表
export function getUsers(params: { page?: number; limit?: number } = {}) {
  return request.get<UserListResult>({
    url: Api.GetUsers,
    params,
  });
}

// 创建用户
export function adminCreateUser(data: AdminCreateUserRequest) {
  return request.post({
    url: Api.AdminCreateUser,
    data,
  });
}

// 更新用户状态
export function adminUpdateUserStatus(id: number, data: AdminUpdateUserStatusRequest) {
  return request.put({
    url: Api.AdminUpdateUserStatus.replace('{id}', String(id)),
    data,
  });
}
