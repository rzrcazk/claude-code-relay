<template>
  <div>
    <t-card class="list-card-container" :bordered="false">
      <t-table
        :data="data"
        :columns="COLUMNS"
        :row-key="rowKey"
        vertical-align="top"
        :hover="true"
        :pagination="pagination"
        :loading="dataLoading"
        :header-affixed-top="headerAffixedTop"
        @page-change="handlePageChange"
      >
        <template #method="{ row }">
          <t-tag :theme="getMethodColor(row.method)" variant="light">
            {{ row.method }}
          </t-tag>
        </template>

        <template #status_code="{ row }">
          <t-tag :theme="getStatusColor(row.status_code)" variant="light">
            {{ row.status_code }}
          </t-tag>
        </template>

        <template #path="{ row }">
          <t-tooltip :content="row.path" placement="top">
            <span class="path-text">{{ truncatePath(row.path) }}</span>
          </t-tooltip>
        </template>

        <template #user="{ row }">
          <span v-if="row.user">{{ row.user.username }}</span>
          <span v-else class="text-placeholder">-</span>
        </template>

        <template #duration="{ row }">
          <span>{{ row.duration }}ms</span>
        </template>

        <template #created_at="{ row }">
          <span>{{ formatDateTime(row.created_at) }}</span>
        </template>

        <template #ip="{ row }">
          <span>{{ row.ip || '-' }}</span>
        </template>
      </t-table>
    </t-card>
  </div>
</template>
<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next';
import { onMounted, ref } from 'vue';

import type { SystemLog } from '@/api/systemLogs';
import { getSystemLogs } from '@/api/systemLogs';
import { formatDateTime } from '@/utils/date';

const COLUMNS = [
  {
    title: '请求ID',
    colKey: 'request_id',
    width: 150,
    ellipsis: true,
  },
  {
    title: '方法',
    colKey: 'method',
    width: 80,
  },
  {
    title: '路径',
    colKey: 'path',
    width: 300,
  },
  {
    title: '状态码',
    colKey: 'status_code',
    width: 100,
  },
  {
    title: '用户',
    colKey: 'user',
    width: 120,
  },
  {
    title: 'IP地址',
    colKey: 'ip',
    width: 140,
  },
  {
    title: '耗时',
    colKey: 'duration',
    width: 100,
  },
  {
    title: '请求时间',
    colKey: 'created_at',
    width: 180,
  },
];

const data = ref<SystemLog[]>([]);
const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
});
const dataLoading = ref(false);
const headerAffixedTop = ref(false);
const rowKey = 'id';

const fetchData = async () => {
  dataLoading.value = true;
  try {
    const params = {
      page: pagination.value.current,
      limit: pagination.value.pageSize,
    };

    const result = await getSystemLogs(params);
    data.value = result.logs || [];
    pagination.value.total = result.total || 0;
  } catch (error) {
    console.error('获取系统日志失败:', error);
    MessagePlugin.error('获取系统日志失败');
  } finally {
    dataLoading.value = false;
  }
};

const handlePageChange = (pageInfo: any) => {
  pagination.value.current = pageInfo.current;
  pagination.value.pageSize = pageInfo.pageSize;
  fetchData();
};

const handleRefresh = () => {
  fetchData();
};

const getMethodColor = (method: string) => {
  const colors = {
    GET: 'success',
    POST: 'primary',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'warning',
  };
  // eslint-disable-next-line ts/ban-ts-comment
  // @ts-ignore
  return colors[method] || 'default';
};

const getStatusColor = (statusCode: number) => {
  if (statusCode >= 200 && statusCode < 300) return 'success';
  if (statusCode >= 300 && statusCode < 400) return 'warning';
  if (statusCode >= 400 && statusCode < 500) return 'danger';
  if (statusCode >= 500) return 'danger';
  return 'default';
};

const truncatePath = (path: string) => {
  return path.length > 40 ? `${path.substring(0, 40)}...` : path;
};

onMounted(() => {
  fetchData();
});
</script>
<style lang="less" scoped>
.list-card-container {
  padding: 20px;
}

.left-operation-container {
  .t-button + .t-button {
    margin-left: 8px;
  }
}

.path-text {
  cursor: pointer;
}

.text-placeholder {
  color: var(--td-text-color-placeholder);
}
</style>
