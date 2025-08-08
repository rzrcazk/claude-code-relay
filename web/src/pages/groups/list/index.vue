<template>
  <div>
    <t-card class="list-card-container" :bordered="false">
      <t-row justify="space-between">
        <div class="left-operation-container">
          <t-button @click="handleCreate"> 创建分组 </t-button>
          <t-button v-if="selectedRowKeys.length > 0" theme="danger" variant="outline" @click="handleBatchDelete">
            批量删除
          </t-button>
        </div>
        <div class="search-input">
          <t-input v-model="searchValue" placeholder="搜索分组名称" clearable @enter="handleSearch">
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

        <template #api_key_count="{ row }">
          <t-tag theme="default" variant="light"> {{ row.api_key_count || 0 }} </t-tag>
        </template>

        <template #account_count="{ row }">
          <t-tag theme="primary" variant="light"> {{ row.account_count || 0 }} </t-tag>
        </template>

        <template #op="{ row }">
          <t-space size="2px">
            <t-button variant="text" size="small" theme="primary" @click="handleEdit(row)"> 编辑 </t-button>
            <t-button
              variant="text"
              size="small"
              :theme="row.status === 1 ? 'warning' : 'success'"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </t-button>
            <t-button variant="text" size="small" theme="danger" @click="handleDelete([row])"> 删除 </t-button>
          </t-space>
        </template>
      </t-table>
    </t-card>

    <!-- 创建/编辑分组弹窗 -->
    <t-dialog
      v-model:visible="formVisible"
      :header="editingItem ? '编辑分组' : '创建分组'"
      width="500px"
      placement="center"
      @confirm="handleFormConfirm"
      @cancel="handleFormCancel"
    >
      <t-form ref="formRef" :model="formData" :rules="formRules" label-align="top" label-width="120px">
        <t-form-item label="分组名称" name="name">
          <t-input v-model="formData.name" placeholder="请输入分组名称" />
        </t-form-item>

        <t-form-item label="描述" name="remark">
          <t-textarea v-model="formData.remark" placeholder="请输入分组描述（可选）" :rows="3" />
        </t-form-item>

        <t-form-item label="状态" name="status">
          <t-radio-group v-model="formData.status">
            <t-radio :value="1">启用</t-radio>
            <t-radio :value="0">禁用</t-radio>
          </t-radio-group>
        </t-form-item>
      </t-form>
    </t-dialog>

    <!-- 删除确认弹窗 -->
    <t-dialog
      v-model:visible="deleteVisible"
      header="确认删除"
      @confirm="handleDeleteConfirm"
      @cancel="handleDeleteCancel"
    >
      <p>{{ deleteConfirmText }}</p>
    </t-dialog>
  </div>
</template>
<script setup lang="ts">
import { SearchIcon } from 'tdesign-icons-vue-next';
import type { FormInstanceFunctions, FormRules, PrimaryTableCol, TableRowData } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, onMounted, reactive, ref } from 'vue';

import type { Group, GroupCreateParams, GroupUpdateParams } from '@/api/group';
import { batchDeleteGroups, createGroup, deleteGroup, getGroupList, updateGroup, updateGroupStatus } from '@/api/group';
import { prefix } from '@/config/global';
import { useSettingStore } from '@/store';

defineOptions({
  name: 'GroupsList',
});

const store = useSettingStore();

// 表格列定义
const COLUMNS: PrimaryTableCol<TableRowData>[] = [
  {
    title: 'ID',
    align: 'left',
    width: 80,
    colKey: 'id',
  },
  {
    title: '分组名称',
    align: 'left',
    width: 200,
    colKey: 'name',
    ellipsis: true,
  },

  { title: '状态', colKey: 'status', width: 100 },
  {
    title: 'API密钥数量',
    colKey: 'api_key_count',
    width: 120,
  },
  {
    title: '账号数量',
    colKey: 'account_count',
    width: 120,
  },
  {
    title: '描述',
    align: 'left',
    width: 240,
    colKey: 'remark',
    ellipsis: true,
  },
  {
    title: '创建时间',
    colKey: 'created_at',
    width: 180,
    cell: (h, { row }) => formatDateTime(row.created_at),
  },
  {
    title: '更新时间',
    colKey: 'updated_at',
    width: 180,
    cell: (h, { row }) => formatDateTime(row.updated_at),
  },
  {
    title: '操作',
    align: 'center',
    fixed: 'right',
    width: 160,
    colKey: 'op',
  },
];

// 数据相关
const data = ref<Group[]>([]);
const dataLoading = ref(false);
const selectedRowKeys = ref<(string | number)[]>([]);
const searchValue = ref('');

// 分页
const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
  showJumper: true,
  showSizeChanger: true,
});

// 表单相关
const formVisible = ref(false);
const formRef = ref<FormInstanceFunctions>();
const editingItem = ref<Group | null>(null);
const formData = reactive<GroupCreateParams & GroupUpdateParams>({
  name: '',
  remark: '',
  status: 1,
  id: 0,
});

const formRules = reactive<FormRules<GroupCreateParams & GroupUpdateParams>>({
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur', type: 'error' }],
});

// 删除相关
const deleteVisible = ref(false);
const deleteItems = ref<Group[]>([]);

// 计算属性
const headerAffixedTop = computed(
  () =>
    ({
      offsetTop: store.isUseTabsRouter ? 48 : 0,
      container: `.${prefix}-layout`,
    }) as any,
);

const rowKey = 'id';

const deleteConfirmText = computed(() => {
  if (deleteItems.value.length === 1) {
    return `确认删除分组 "${deleteItems.value[0].name}" 吗？删除后该分组下的所有API密钥将移动到默认分组。`;
  }
  return `确认删除选中的 ${deleteItems.value.length} 个分组吗？删除后这些分组下的所有API密钥将移动到默认分组。`;
});

// 方法
const fetchData = async () => {
  dataLoading.value = true;
  try {
    const params = {
      page: pagination.value.current,
      limit: pagination.value.pageSize,
      name: searchValue.value || undefined,
    };
    const result = await getGroupList(params);
    data.value = result.groups || [];
    pagination.value = {
      ...pagination.value,
      total: result.total || 0,
    };
  } catch (error) {
    console.error('获取分组列表失败:', error);
    MessagePlugin.error('获取分组列表失败');
  } finally {
    dataLoading.value = false;
  }
};

const handleSearch = () => {
  pagination.value.current = 1;
  fetchData();
};

const handlePageChange = (pageInfo: any) => {
  pagination.value.current = pageInfo.current;
  pagination.value.pageSize = pageInfo.pageSize;
  fetchData();
};

const handleSelectChange = (value: (string | number)[]) => {
  selectedRowKeys.value = value;
};

// 工具函数
const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleString('zh-CN');
};

// 操作相关方法
const handleCreate = () => {
  editingItem.value = null;
  Object.assign(formData, {
    name: '',
    remark: '',
    status: 1,
    id: 0,
  });
  formVisible.value = true;
};

const handleEdit = (item: Group) => {
  editingItem.value = item;
  Object.assign(formData, {
    name: item.name,
    remark: item.remark || '',
    status: item.status,
    id: item.id,
  });
  formVisible.value = true;
};

const handleFormConfirm = async () => {
  const valid = await formRef.value?.validate();
  if (!valid) return;

  try {
    if (editingItem.value) {
      // 编辑
      const updateData: GroupUpdateParams = {
        id: formData.id,
        name: formData.name,
        remark: formData.remark,
        status: formData.status,
      };
      await updateGroup(updateData);
      MessagePlugin.success('更新成功');
    } else {
      // 创建
      const createData: GroupCreateParams = {
        name: formData.name,
        remark: formData.remark,
        status: formData.status,
      };
      await createGroup(createData);
      MessagePlugin.success('创建成功');
    }

    formVisible.value = false;
    await fetchData();
  } catch (error) {
    console.error('操作失败:', error);
    MessagePlugin.error('操作失败');
  }
};

const handleFormCancel = () => {
  formVisible.value = false;
};

const handleToggleStatus = async (item: Group) => {
  try {
    const newStatus = item.status === 1 ? 0 : 1;
    await updateGroupStatus(item.id, newStatus);
    MessagePlugin.success(newStatus === 1 ? '已启用' : '已禁用');
    await fetchData();
  } catch (error) {
    console.error('状态更新失败:', error);
    MessagePlugin.error('状态更新失败');
  }
};

const handleDelete = (items: Group[]) => {
  deleteItems.value = items;
  deleteVisible.value = true;
};

const handleBatchDelete = () => {
  const selectedItems = data.value.filter((item) => selectedRowKeys.value.includes(item.id));
  if (selectedItems.length === 0) {
    MessagePlugin.warning('请先选择要删除的分组');
    return;
  }
  handleDelete(selectedItems);
};

const handleDeleteConfirm = async () => {
  try {
    if (deleteItems.value.length === 1) {
      await deleteGroup(deleteItems.value[0].id);
    } else {
      const ids = deleteItems.value.map((item) => item.id);
      await batchDeleteGroups(ids);
    }
    MessagePlugin.success('删除成功');
    selectedRowKeys.value = [];
    await fetchData();
  } catch (error) {
    console.error('删除失败:', error);
    MessagePlugin.error('删除失败');
  } finally {
    deleteVisible.value = false;
  }
};

const handleDeleteCancel = () => {
  deleteVisible.value = false;
};

// 生命周期
onMounted(() => {
  fetchData();
});
</script>
<style lang="less" scoped>
.list-card-container {
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);

  :deep(.t-card__body) {
    padding: 0;
  }
}

.left-operation-container {
  display: flex;
  align-items: center;
  margin-bottom: var(--td-comp-margin-xxl);
  gap: var(--td-comp-margin-m);

  .selected-count {
    display: inline-block;
    margin-left: var(--td-comp-margin-l);
    color: var(--td-text-color-secondary);
  }
}

.search-input {
  width: 360px;
}

.text-placeholder {
  color: var(--td-text-color-placeholder);
}
</style>
