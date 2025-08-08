import { defineStore } from 'pinia';

import type { LoginRequest, LoginResult } from '@/api/user';
import { getUserProfile, login } from '@/api/user';
import { usePermissionStore } from '@/store';
import type { UserInfo } from '@/types/interface';

const InitUserInfo: UserInfo = {
  name: '', // 用户名，用于展示在页面右上角头像处
  roles: [], // 前端权限模型使用 如果使用请配置modules/permission-fe.ts使用
};

export const useUserStore = defineStore('user', {
  state: () => ({
    token: 'main_token', // 默认token不走权限
    userInfo: { ...InitUserInfo },
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
        name: result.user.username,
        roles: result.user.role === 'admin' ? ['all'] : ['user'],
      };
      return result;
    },

    async setUserInfo(loginResult: LoginResult) {
      this.token = loginResult.token;
      this.userInfo = {
        name: loginResult.user.username,
        roles: loginResult.user.role === 'admin' ? ['all'] : ['user'],
      };
    },

    async getUserInfo() {
      if (!this.token) {
        throw new Error('请先登录');
      }

      try {
        const result = await getUserProfile();
        this.userInfo = {
          name: result.username,
          roles: result.role === 'admin' ? ['all'] : ['user'],
        };
        return result;
      } catch (error) {
        // 如果获取用户信息失败，清除token
        this.token = '';
        this.userInfo = { ...InitUserInfo };
        throw error;
      }
    },

    async logout() {
      this.token = '';
      this.userInfo = { ...InitUserInfo };
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
