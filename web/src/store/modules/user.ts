import { defineStore } from 'pinia';

import type { LoginRequest, LoginResult, UserProfile } from '@/api/user';
import { getUserProfile, login } from '@/api/user';
import { usePermissionStore } from '@/store';
import type { UserInfo } from '@/types/interface';

const InitUserInfo: UserInfo = {
  name: '', // 用户名，用于展示在页面右上角头像处
  roles: [], // 前端权限模型使用 如果使用请配置modules/permission-fe.ts使用
};

const InitUserProfile: UserProfile = {
  id: 0,
  name: '',
  username: '',
  email: '',
  role: '',
  status: 0,
  created_at: '',
};

export const useUserStore = defineStore('user', {
  state: () => ({
    token: 'main_token', // 默认token不走权限
    userInfo: { ...InitUserInfo },
    userProfile: { ...InitUserProfile }, // 完整的用户信息
  }),
  getters: {
    roles: (state) => {
      return state.userInfo?.roles;
    },
  },
  actions: {
    async login(loginData: LoginRequest) {
      const result = await login(loginData);
      this.token = result.token;
      this.userInfo = {
        name: result.user.name || result.user.username,
        roles: result.user.role === 'admin' ? ['all'] : ['user'],
      };
      // 保存完整的用户信息
      this.userProfile = {
        id: result.user.id,
        name: result.user.name || result.user.username,
        username: result.user.username,
        email: result.user.email,
        role: result.user.role,
        status: result.user.status,
        created_at: '',
      };
      return result;
    },

    async setUserInfo(loginResult: LoginResult) {
      this.token = loginResult.token;
      this.userInfo = {
        name: loginResult.user.name || loginResult.user.username,
        roles: loginResult.user.role === 'admin' ? ['all'] : ['user'],
      };
      // 保存完整的用户信息
      this.userProfile = {
        id: loginResult.user.id,
        name: loginResult.user.name || loginResult.user.username,
        username: loginResult.user.username,
        email: loginResult.user.email,
        role: loginResult.user.role,
        status: loginResult.user.status,
        created_at: '',
      };
    },

    async getUserInfo() {
      if (!this.token) {
        throw new Error('请先登录');
      }

      try {
        const result = await getUserProfile();
        this.userInfo = {
          name: result.name || result.username,
          roles: result.role === 'admin' ? ['all'] : ['user'],
        };
        // 保存完整的用户信息
        this.userProfile = {
          id: result.id,
          name: result.name || result.username,
          username: result.username,
          email: result.email,
          role: result.role,
          status: result.status,
          created_at: result.created_at,
        };
        return result;
      } catch (error) {
        // 如果获取用户信息失败，清除token
        this.token = '';
        this.userInfo = { ...InitUserInfo };
        this.userProfile = { ...InitUserProfile };
        throw error;
      }
    },

    async logout() {
      this.token = '';
      this.userInfo = { ...InitUserInfo };
      this.userProfile = { ...InitUserProfile };
    },
  },
  persist: {
    afterRestore: () => {
      const permissionStore = usePermissionStore();
      permissionStore.initRoutes();
    },
    key: 'user',
    paths: ['token'],
  },
});
