<template>
  <div>
    <t-card class="list-card-container" :bordered="false">
      <t-row justify="space-between">
        <div class="left-operation-container">
          <t-button @click="handleCreate"> 创建用户 </t-button>
        </div>
        <div class="search-input">
          <t-input v-model="searchValue" placeholder="搜索用户名或邮箱" clearable @enter="handleSearch">
            <template #suffix-icon>
              <search-icon size="16px" />
            </template>
          </t-input>
        </div>
      </t-row>

      <t-table
        :data="data"
        :columns="COLUMNS"
        :row-key="rowKey"
        vertical-align="top"
        :hover="true"
        :pagination="pagination"
        :selected-row-keys="selectedRowKeys"
        :loading="dataLoading"
        :header-affixed-top="headerAffixedTop"
        @page-change="handlePageChange"
        @select-change="handleSelectChange"
      >
        <template #status="{ row }">
          <t-tag v-if="row.status === 1" theme="success" variant="light"> 启用 </t-tag>
          <t-tag v-else theme="danger" variant="light"> 禁用 </t-tag>
        </template>

        <template #role="{ row }">
          <t-tag v-if="row.role === 'admin'" theme="warning" variant="light"> 管理员 </t-tag>
          <t-tag v-else theme="primary" variant="light"> 普通用户 </t-tag>
        </template>

        <template #created_at="{ row }">
          <span>{{ formatDateTime(row.created_at) }}</span>
        </template>

        <template #op="{ row }">
          <t-space size="2px">
            <t-link v-if="!isAdminUser(row)" theme="primary" @click="handleStatusToggle(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </t-link>
            <span v-else class="text-placeholder">管理员账户</span>
          </t-space>
        </template>
      </t-table>
    </t-card>

    <!-- 创建用户对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建用户"
      width="500px"
      :confirm-btn="{ content: '创建', theme: 'primary' }"
      @confirm="handleCreateConfirm"
    >
      <t-form ref="createForm" :data="createFormData" label-width="80px">
        <t-form-item label="用户名" name="username">
          <t-input v-model="createFormData.username" placeholder="请输入用户名" />
        </t-form-item>
        <t-form-item label="邮箱" name="email">
          <t-input v-model="createFormData.email" placeholder="请输入邮箱" />
        </t-form-item>
        <t-form-item label="密码" name="password">
          <t-input v-model="createFormData.password" type="password" placeholder="请输入密码" />
        </t-form-item>
        <t-form-item label="角色" name="role">
          <t-select v-model="createFormData.role" placeholder="请选择角色">
            <t-option value="user" label="普通用户" />
            <t-option value="admin" label="管理员" />
          </t-select>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>
<script setup lang="ts">
import { SearchIcon } from 'tdesign-icons-vue-next';
import type { PrimaryTableCol, TableRowData } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { onMounted, ref } from 'vue';

import type { CreateUserRequest, UserInfo } from '@/api/users';
import { createUser, getUserList, updateUserStatus, UserRole, UserStatus, UserStatusLabels } from '@/api/users';

const COLUMNS: PrimaryTableCol<TableRowData>[] = [
  {
    title: 'ID',
    colKey: 'id',
    width: 80,
    fixed: 'left',
  },
  {
    title: '用户名',
    colKey: 'username',
    width: 150,
  },
  {
    title: '邮箱',
    colKey: 'email',
    width: 200,
  },
  {
    title: '状态',
    colKey: 'status',
    width: 100,
  },
  {
    title: '角色',
    colKey: 'role',
    width: 100,
  },
  {
    title: '创建时间',
    colKey: 'created_at',
    width: 180,
  },
  {
    title: '操作',
    colKey: 'op',
    width: 120,
    fixed: 'right',
  },
];

const data = ref<UserInfo[]>([]);
const selectedRowKeys = ref<Array<string | number>>([]);
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
});
const dataLoading = ref(false);
const searchValue = ref('');
const headerAffixedTop = ref(false);
const rowKey = 'id';

// 创建用户相关
const createDialogVisible = ref(false);
const createForm = ref();
const createFormData = ref<CreateUserRequest>({
  username: '',
  email: '',
  password: '',
  role: UserRole.USER,
});

// 格式化时间
const formatDateTime = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN');
};

// 判断是否为管理员用户
const isAdminUser = (user: UserInfo) => {
  return user.role === UserRole.ADMIN;
};

// 获取用户列表
const fetchData = async () => {
  try {
    dataLoading.value = true;
    const { users, total } = await getUserList({
      page: pagination.value.current,
      limit: pagination.value.pageSize,
    });
    data.value = users;
    pagination.value.total = total;
  } catch (error) {
    console.error('获取用户列表失败:', error);
    MessagePlugin.error('获取用户列表失败');
  } finally {
    dataLoading.value = false;
  }
};

// 分页变化
const handlePageChange = (pageInfo: any) => {
  pagination.value.current = pageInfo.current;
  pagination.value.pageSize = pageInfo.pageSize;
  fetchData();
};

// 选择变化
const handleSelectChange = (value: Array<string | number>) => {
  selectedRowKeys.value = value;
};

// 搜索
const handleSearch = () => {
  fetchData();
};

// 创建用户
const handleCreate = () => {
  createDialogVisible.value = true;
  // 重置表单
  createFormData.value = {
    username: '',
    email: '',
    password: '',
    role: UserRole.USER,
  };
};

// 确认创建用户
const handleCreateConfirm = async () => {
  const result = await createForm.value?.validate();
  if (result !== true) return false;

  try {
    await createUser(createFormData.value);
    MessagePlugin.success('用户创建成功');
    createDialogVisible.value = false;
    fetchData();
    return true;
  } catch (error) {
    console.error('创建用户失败:', error);
    MessagePlugin.error('创建用户失败');
    return false;
  }
};

// 切换用户状态
const handleStatusToggle = async (row: UserInfo) => {
  // 检查是否为管理员用户
  if (isAdminUser(row)) {
    MessagePlugin.warning('不能操作管理员账户');
    return;
  }

  try {
    const newStatus = row.status === UserStatus.ENABLED ? UserStatus.DISABLED : UserStatus.ENABLED;
    const statusLabel = UserStatusLabels[newStatus];

    await updateUserStatus(row.id, { status: newStatus });

    MessagePlugin.success(`用户${statusLabel}成功`);
    fetchData();
  } catch (error) {
    console.error('更新用户状态失败:', error);
    MessagePlugin.error('更新用户状态失败');
  }
};

onMounted(() => {
  fetchData();
});
</script>
<style lang="less" scoped>
.list-card-container {
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);

  .left-operation-container {
    .t-button + .t-button {
      margin-left: var(--td-comp-margin-s);
    }
  }

  .search-input {
    width: 360px;
  }
}
</style>
