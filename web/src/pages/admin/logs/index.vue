<template>
  <div>
    <t-card class="list-card-container" :bordered="false">
      <div class="filter-container">
        <t-form :data="filterForm" @submit="handleFilter" @reset="handleReset">
          <div class="filter-form">
            <div class="left-operation-container">
              <t-form-item label="用户ID" name="user_id">
                <t-input
                  v-model="filterForm.user_id"
                  placeholder="请输入用户ID"
                  :min="1"
                  clearable
                  style="width: 150px"
                />
              </t-form-item>
              <t-form-item label="状态码" name="status_code">
                <t-input
                  v-model="filterForm.status_code"
                  placeholder="请输入状态码"
                  :min="100"
                  :max="599"
                  clearable
                  style="width: 150px"
                />
              </t-form-item>
            </div>
            <div>
              <t-form-item>
                <t-button theme="primary" type="submit">搜索</t-button>
                <t-button type="reset" variant="base" theme="default">重置</t-button>
              </t-form-item>
            </div>
          </div>
        </t-form>
      </div>
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

import type { SystemLog, SystemLogQueryParams } from '@/api/systemLogs';
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

const filterForm = ref({
  user_id: undefined as number | undefined,
  status_code: undefined as number | undefined,
});

const fetchData = async () => {
  dataLoading.value = true;
  try {
    const params: SystemLogQueryParams = {
      page: pagination.value.current,
      limit: pagination.value.pageSize,
    };

    if (filterForm.value.user_id) {
      params.user_id = filterForm.value.user_id;
    }
    if (filterForm.value.status_code) {
      params.status_code = filterForm.value.status_code;
    }

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

const handleFilter = () => {
  pagination.value.current = 1;
  fetchData();
};

const handleReset = () => {
  filterForm.value.user_id = undefined;
  filterForm.value.status_code = undefined;
  pagination.value.current = 1;
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

.filter-container {
  width: 100%;
  background: var(--td-bg-color-container);
  border-radius: 6px;
  border: 1px solid #fff;
  margin-bottom: 10px;

  .filter-form {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .left-operation-container {
      display: flex;
      justify-content: left;
      align-items: center;

      ::v-deep(.t-form__item) {
        margin-bottom: 0 !important;
      }
    }
  }
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
