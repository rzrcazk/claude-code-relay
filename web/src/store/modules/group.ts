import { defineStore } from 'pinia';
import type { Group, GroupListResponse } from '@/api/group';
import { getGroupList } from '@/api/group';

interface GroupState {
  groupList: Group[];
  total: number;
  loading: boolean;
  searchParams: {
    page: number;
    size: number;
    name?: string;
    status?: number;
  };
}

export const useGroupStore = defineStore('group', {
  state: (): GroupState => ({
    groupList: [],
    total: 0,
    loading: false,
    searchParams: {
      page: 1,
      size: 20,
    },
  }),

  getters: {
    // 获取启用状态的分组列表（用于下拉选择）
    enabledGroups: (state): Group[] => state.groupList.filter((group) => group.status === 1),
    
    // 根据ID获取分组信息
    getGroupById: (state) => (id: number): Group | undefined => 
      state.groupList.find((group) => group.id === id),
    
    // 获取分组选项（用于表单选择）
    groupOptions: (state) => 
      state.groupList
        .filter((group) => group.status === 1)
        .map((group) => ({
          label: group.name,
          value: group.id,
        })),
  },

  actions: {
    // 获取分组列表
    async fetchGroupList(params?: Partial<GroupState['searchParams']>) {
      this.loading = true;
      try {
        if (params) {
          this.searchParams = { ...this.searchParams, ...params };
        }
        
        const response: GroupListResponse = await getGroupList(this.searchParams);
        this.groupList = response.list || [];
        this.total = response.total || 0;
        
        return response;
      } catch (error) {
        console.error('获取分组列表失败:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // 重置搜索参数
    resetSearchParams() {
      this.searchParams = {
        page: 1,
        size: 20,
      };
    },

    // 添加分组到列表（创建后）
    addGroup(group: Group) {
      this.groupList.unshift(group);
      this.total += 1;
    },

    // 更新分组信息
    updateGroup(updatedGroup: Group) {
      const index = this.groupList.findIndex((group) => group.id === updatedGroup.id);
      if (index !== -1) {
        this.groupList[index] = updatedGroup;
      }
    },

    // 从列表中移除分组
    removeGroup(groupId: number) {
      const index = this.groupList.findIndex((group) => group.id === groupId);
      if (index !== -1) {
        this.groupList.splice(index, 1);
        this.total = Math.max(0, this.total - 1);
      }
    },

    // 批量移除分组
    removeGroups(groupIds: number[]) {
      this.groupList = this.groupList.filter((group) => !groupIds.includes(group.id));
      this.total = Math.max(0, this.total - groupIds.length);
    },

    // 更新分组状态
    updateGroupStatus(groupId: number, status: number) {
      const group = this.groupList.find((g) => g.id === groupId);
      if (group) {
        group.status = status;
      }
    },
  },

  // 持久化存储配置
  persist: {
    key: 'group-store',
    storage: localStorage,
    // 只持久化必要的数据，避免存储过期数据
    paths: ['searchParams'],
  },
});