<template>
  <div>
    <t-card class="list-card-container" :bordered="false">
      <t-row justify="space-between">
        <div class="left-operation-container">
          <t-button @click="handleCreate"> 创建密钥 </t-button>
        </div>
        <div class="search-input">
          <t-input v-model="searchValue" placeholder="搜索密钥名称" clearable @enter="handleSearch">
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
        <template #user_id="{ row }">
          <t-tag theme="default" variant="light"> {{ row.user_id }} </t-tag>
        </template>

        <template #status="{ row }">
          <t-tag v-if="row.status === 1" theme="success" variant="light"> 启用 </t-tag>
          <t-tag v-else theme="danger" variant="light"> 禁用 </t-tag>
        </template>

        <template #usage="{ row }">
          <div class="usage-info">
            <p>使用次数：{{ row.today_usage_count }}</p>
            <p>输入Token：{{ formatNumber(row.today_input_tokens) }}</p>
            <p>输出Token：{{ formatNumber(row.today_output_tokens) }}</p>
            <p>预计费用：${{ row.today_total_cost.toFixed(4) }}</p>
          </div>
        </template>

        <template #daily_limit="{ row }">
          <t-tag theme="default" variant="light"> ${{ row.daily_limit }} </t-tag>
        </template>

        <template #model_restriction="{ row }">
          <t-tag v-if="row.model_restriction" theme="default" variant="light"> 限制 </t-tag>
          <span v-else class="text-placeholder">不限制 </span>
        </template>

        <template #expires_at="{ row }">
          <span v-if="row.expires_at">
            {{ formatDateTime(row.expires_at) }}
          </span>
          <span v-else class="text-placeholder">永不过期</span>
        </template>

        <template #group="{ row }">
          <t-tag v-if="row.group" variant="outline">
            {{ row.group.name }}
          </t-tag>
          <span v-else class="text-placeholder">全局共享组</span>
        </template>

        <template #op="{ row }">
          <t-space size="2px">
            <t-button variant="text" size="small" @click="copyToClipboard(row.key)"> 复制 </t-button>
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

    <!-- 创建/编辑密钥弹窗 -->
    <t-dialog
      v-model:visible="formVisible"
      :header="editingItem ? '编辑密钥' : '创建密钥'"
      width="600px"
      placement="center"
      @confirm="handleFormConfirm"
      @cancel="handleFormCancel"
    >
      <t-form ref="formRef" :model="formData" :rules="formRules" label-align="top" label-width="120px">
        <t-form-item label="密钥名称" name="name">
          <t-input v-model="formData.name" placeholder="请输入密钥名称" />
        </t-form-item>

        <t-form-item label="自定义密钥" name="key">
          <t-input v-model="formData.key" placeholder="留空将自动生成" :disabled="!!editingItem" />
          <template #help>
            <span v-if="editingItem">密钥创建后不可修改</span>
            <span v-else>留空将自动生成 sk- 开头的密钥</span>
          </template>
        </t-form-item>

        <t-form-item label="过期时间" name="expires_at">
          <t-date-picker
            v-model="formData.expires_at"
            format="YYYY-MM-DD HH:mm:ss"
            :disable-date="disableDate"
            clearable
            placeholder="留空表示永不过期"
            style="width: 100%"
          />
        </t-form-item>

        <t-form-item label="状态" name="status">
          <t-radio-group v-model="formData.status">
            <t-radio :value="1">启用</t-radio>
            <t-radio :value="0">禁用</t-radio>
          </t-radio-group>
        </t-form-item>

        <t-form-item label="模型限制" name="model_restriction">
          <t-input v-model="formData.model_restriction" placeholder="多个模型用逗号分隔，留空表示不限制" />
          <template #help> 例如：claude-3-5-sonnet-20241022,claude-3-haiku-20240307 </template>
        </t-form-item>

        <t-form-item label="每日限额(美元)" name="daily_limit">
          <t-input-number
            v-model="formData.daily_limit"
            :min="0"
            :step="0.01"
            placeholder="0表示不限制"
            style="width: 100%"
          />
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

import type { ApiKey, CreateApiKeyRequest, UpdateApiKeyRequest } from '@/api/apikey';
import { createApiKey, deleteApiKey, getApiKeys, updateApiKey, updateApiKeyStatus } from '@/api/apikey';
import { prefix } from '@/config/global';
import { useSettingStore } from '@/store';

defineOptions({
  name: 'ApiKeysList',
});

const store = useSettingStore();

// 表格列定义
const COLUMNS: PrimaryTableCol<TableRowData>[] = [
  {
    title: 'ID',
    align: 'left',
    width: 100,
    colKey: 'id',
  },
  {
    title: '密钥名称',
    align: 'left',
    width: 180,
    colKey: 'name',
    ellipsis: true,
  },
  { title: '状态', colKey: 'status', width: 100 },
  { title: '用户ID', colKey: 'user_id', width: 100 },
  {
    title: '分组',
    colKey: 'group',
    width: 140,
  },
  {
    title: '今日使用情况',
    colKey: 'usage',
    width: 180,
  },
  {
    title: '过期时间',
    colKey: 'expires_at',
    width: 180,
  },
  {
    title: '每日限额',
    colKey: 'daily_limit',
    width: 120,
  },
  {
    title: '限制模型',
    colKey: 'model_restriction',
    width: 100,
  },
  {
    title: '最后使用时间',
    colKey: 'last_used_time',
    width: 180,
  },
  {
    title: '创建时间',
    colKey: 'created_at',
    width: 180,
    cell: (h, { row }) => formatDateTime(row.created_at),
  },
  {
    title: '操作',
    align: 'center',
    fixed: 'right',
    width: 180,
    colKey: 'op',
  },
];

// 数据相关
const data = ref<ApiKey[]>([]);
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
const editingItem = ref<ApiKey | null>(null);
const formData = reactive<CreateApiKeyRequest & UpdateApiKeyRequest>({
  name: '',
  key: '',
  expires_at: '',
  status: 1,
  group_id: 0,
  model_restriction: '',
  daily_limit: 0,
});

const formRules = reactive<FormRules<CreateApiKeyRequest & UpdateApiKeyRequest>>({
  name: [{ required: true, message: '请输入密钥名称', trigger: 'blur', type: 'error' }],
});

// 删除相关
const deleteVisible = ref(false);
const deleteItems = ref<ApiKey[]>([]);

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
    return `确认删除密钥 "${deleteItems.value[0].name}" 吗？删除后无法恢复。`;
  }
  return `确认删除选中的 ${deleteItems.value.length} 个密钥吗？删除后无法恢复。`;
});

// 方法
const fetchData = async () => {
  dataLoading.value = true;
  try {
    const params = {
      page: pagination.value.current,
      limit: pagination.value.pageSize,
    };
    const result = await getApiKeys(params);
    data.value = result.api_keys || [];
    pagination.value = {
      ...pagination.value,
      total: result.total || 0,
    };
  } catch (error) {
    console.error('获取密钥列表失败:', error);
    MessagePlugin.error('获取密钥列表失败');
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
const maskApiKey = (key: string): string => {
  if (!key) return '';
  if (key.length <= 8) return key;
  return `${key.substring(0, 8)}${'*'.repeat(Math.min(key.length - 8, 20))}`;
};

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    MessagePlugin.success('已复制到剪贴板');
  } catch {
    MessagePlugin.error('复制失败');
  }
};

const formatNumber = (num: number): string => {
  if (num < 1000) return num.toString();
  if (num < 1000000) return `${(num / 1000).toFixed(1)}K`;
  return `${(num / 1000000).toFixed(1)}M`;
};

const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleString('zh-CN');
};

const disableDate = (date: Date): boolean => {
  return date < new Date();
};

// 操作相关方法

const handleCreate = () => {
  editingItem.value = null;
  Object.assign(formData, {
    name: '',
    key: '',
    expires_at: '',
    status: 1,
    group_id: 0,
    model_restriction: '',
    daily_limit: 0,
  });
  formVisible.value = true;
};

const handleEdit = (item: ApiKey) => {
  editingItem.value = item;
  Object.assign(formData, {
    name: item.name,
    key: item.key,
    expires_at: item.expires_at || '',
    status: item.status,
    group_id: item.group_id,
    model_restriction: item.model_restriction || '',
    daily_limit: item.daily_limit,
  });
  formVisible.value = true;
};

const handleFormConfirm = async () => {
  const valid = await formRef.value?.validate();
  if (!valid) return;

  try {
    if (editingItem.value) {
      // 编辑
      const updateData: UpdateApiKeyRequest = {
        name: formData.name,
        expires_at: formData.expires_at || undefined,
        status: formData.status,
        group_id: formData.group_id,
        model_restriction: formData.model_restriction,
        daily_limit: formData.daily_limit,
      };
      await updateApiKey(editingItem.value.id, updateData);
      MessagePlugin.success('更新成功');
    } else {
      // 创建
      const createData: CreateApiKeyRequest = {
        name: formData.name,
        key: formData.key || undefined,
        expires_at: formData.expires_at || undefined,
        status: formData.status,
        group_id: formData.group_id,
        model_restriction: formData.model_restriction,
        daily_limit: formData.daily_limit,
      };
      await createApiKey(createData);
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

const handleToggleStatus = async (item: ApiKey) => {
  try {
    const newStatus = item.status === 1 ? 0 : 1;
    await updateApiKeyStatus(item.id, { status: newStatus });
    MessagePlugin.success(newStatus === 1 ? '已启用' : '已禁用');
    await fetchData();
  } catch (error) {
    console.error('状态更新失败:', error);
    MessagePlugin.error('状态更新失败');
  }
};

const handleDelete = (items: ApiKey[]) => {
  deleteItems.value = items;
  deleteVisible.value = true;
};

const handleBatchDelete = () => {
  const selectedItems = data.value.filter((item) => selectedRowKeys.value.includes(item.id));
  if (selectedItems.length === 0) {
    MessagePlugin.warning('请先选择要删除的密钥');
    return;
  }
  handleDelete(selectedItems);
};

const handleDeleteConfirm = async () => {
  try {
    await Promise.all(deleteItems.value.map((item) => deleteApiKey(item.id)));
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

  .selected-count {
    display: inline-block;
    margin-left: var(--td-comp-margin-l);
    color: var(--td-text-color-secondary);
  }
}

.search-input {
  width: 360px;
}

.key-display {
  display: flex;
  align-items: center;
  gap: 8px;

  .key-text {
    font-family: 'Monaco', 'Consolas', monospace;
    font-size: 12px;
    color: var(--td-text-color-primary);
  }
}

.usage-info {
  font-size: 12px;
  line-height: 1.4;

  p {
    margin: 2px 0;
  }
}

.text-placeholder {
  color: var(--td-text-color-placeholder);
}
</style>
