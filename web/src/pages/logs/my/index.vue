<template>
  <div>
    <t-card class="list-card-container" :bordered="false">
      <!-- 搜索和筛选区域 -->
      <t-row justify="space-between" class="search-row">
        <div class="left-operation-container">
          <t-space>
            <t-button variant="outline" @click="handleRefresh">
              <template #icon><refresh-icon /></template>
              刷新
            </t-button>
            <t-button variant="outline" @click="resetFilters"> 重置筛选 </t-button>
          </t-space>
        </div>
        <div class="right-operation-container">
          <t-space>
            <t-select
              v-model="searchFilters.model_name"
              placeholder="选择模型"
              clearable
              style="width: 200px"
              @change="handleSearch"
            >
              <t-option v-for="model in modelOptions" :key="model" :value="model" :label="model" />
            </t-select>

            <t-date-range-picker
              v-model="dateRange"
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择时间范围"
              clearable
              @change="handleDateRangeChange"
            />
          </t-space>
        </div>
      </t-row>

      <!-- 高级筛选 -->
      <t-row v-if="showAdvancedFilter" class="advanced-filter">
        <t-space>
          <t-input-number
            v-model="searchFilters.min_cost"
            placeholder="最小费用"
            :min="0"
            :step="0.001"
            :decimal-places="4"
            style="width: 200px"
            @blur="handleSearch"
          />
          <span>-</span>
          <t-input-number
            v-model="searchFilters.max_cost"
            placeholder="最大费用"
            :min="0"
            :step="0.001"
            :decimal-places="4"
            style="width: 200px"
            @blur="handleSearch"
          />
          <t-input
            v-model="searchFilters.account_id"
            type="number"
            placeholder="账号ID"
            style="width: 180px"
            clearable
            @blur="handleSearch"
          />
          <t-input
            v-model="searchFilters.api_key_id"
            type="number"
            placeholder="API Key ID"
            style="width: 180px"
            clearable
            @blur="handleSearch"
          />
        </t-space>
      </t-row>

      <t-row>
        <t-button variant="text" size="small" @click="showAdvancedFilter = !showAdvancedFilter">
          {{ showAdvancedFilter ? '收起' : '高级筛选' }}
          <template #suffix>
            <chevron-down-icon v-if="!showAdvancedFilter" />
            <chevron-up-icon v-else />
          </template>
        </t-button>
      </t-row>

      <!-- 数据表格 -->
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
        <template #account_id="{ row }">
          <t-tag theme="default" variant="light">{{ row.account_id }}</t-tag>
        </template>

        <template #model_name="{ row }">
          <t-tag theme="primary" variant="outline">{{ row.model_name }}</t-tag>
        </template>

        <template #tokens="{ row }">
          <div class="tokens-info">
            <p><strong>输入:</strong> {{ formatNumber(row.input_tokens) }}</p>
            <p><strong>输出:</strong> {{ formatNumber(row.output_tokens) }}</p>
            <p v-if="row.cache_read_input_tokens > 0">
              <strong>缓存读:</strong> {{ formatNumber(row.cache_read_input_tokens) }}
            </p>
            <p v-if="row.cache_creation_input_tokens > 0">
              <strong>缓存创建:</strong> {{ formatNumber(row.cache_creation_input_tokens) }}
            </p>
          </div>
        </template>

        <template #cost="{ row }">
          <div class="cost-info">
            <p class="total-cost">
              <strong>${{ row.total_cost.toFixed(4) }}</strong>
            </p>
            <div class="cost-breakdown">
              <small>输入: ${{ row.input_cost.toFixed(4) }}</small
              ><br />
              <small>输出: ${{ row.output_cost.toFixed(4) }}</small>
              <small v-if="row.cache_read_cost > 0"> <br />缓存读: ${{ row.cache_read_cost.toFixed(4) }} </small>
              <small v-if="row.cache_write_cost > 0"> <br />缓存写: ${{ row.cache_write_cost.toFixed(4) }} </small>
            </div>
          </div>
        </template>

        <template #duration="{ row }">
          <t-tag theme="default" variant="light">{{ formatDuration(row.duration) }}</t-tag>
        </template>

        <template #api_key="{ row }">
          <div v-if="row.api_key">
            <t-tag variant="outline">{{ row.api_key.name }}</t-tag>
            <br /><t-tag class="mt-2" size="small" theme="default" variant="light">ID: {{ row.api_key_id }}</t-tag>
          </div>
          <span v-else class="text-placeholder">-</span>
        </template>

        <template #op="{ row }">
          <t-space size="2px">
            <t-button variant="text" size="small" theme="primary" @click="handleViewDetail(row)"> 详情 </t-button>
          </t-space>
        </template>
      </t-table>
    </t-card>

    <!-- 详情弹窗 -->
    <t-dialog
      v-model:visible="detailVisible"
      header="日志详情"
      width="800px"
      placement="center"
      @cancel="handleDetailCancel"
    >
      <div v-if="detailData" class="log-detail">
        <t-row :gutter="16">
          <t-col :span="12">
            <div class="detail-item">
              <label>日志ID:</label>
              <span>{{ detailData.id }}</span>
            </div>
          </t-col>
          <t-col :span="12">
            <div class="detail-item">
              <label>模型名称:</label>
              <t-tag theme="primary" variant="outline">{{ detailData.model_name }}</t-tag>
            </div>
          </t-col>
        </t-row>

        <t-row :gutter="16">
          <t-col :span="12">
            <div class="detail-item">
              <label>请求类型:</label>
              <t-tag v-if="detailData.is_stream" theme="success" variant="light">流式</t-tag>
              <t-tag v-else theme="default" variant="light">非流式</t-tag>
            </div>
          </t-col>
          <t-col :span="12">
            <div class="detail-item">
              <label>耗时:</label>
              <span>{{ formatDuration(detailData.duration) }}</span>
            </div>
          </t-col>
        </t-row>

        <t-row :gutter="16">
          <t-col :span="12">
            <div class="detail-item">
              <label>Token使用:</label>
              <div class="token-detail">
                <p>输入Token: {{ formatNumber(detailData.input_tokens) }}</p>
                <p>输出Token: {{ formatNumber(detailData.output_tokens) }}</p>
                <p v-if="detailData.cache_read_input_tokens > 0">
                  缓存读取Token: {{ formatNumber(detailData.cache_read_input_tokens) }}
                </p>
                <p v-if="detailData.cache_creation_input_tokens > 0">
                  缓存创建Token: {{ formatNumber(detailData.cache_creation_input_tokens) }}
                </p>
              </div>
            </div>
          </t-col>
          <t-col :span="12">
            <div class="detail-item">
              <label>费用明细:</label>
              <div class="cost-detail">
                <p>
                  <strong>总费用: ${{ detailData.total_cost.toFixed(4) }}</strong>
                </p>
                <p>输入费用: ${{ detailData.input_cost.toFixed(4) }}</p>
                <p>输出费用: ${{ detailData.output_cost.toFixed(4) }}</p>
                <p v-if="detailData.cache_read_cost > 0">缓存读取费用: ${{ detailData.cache_read_cost.toFixed(4) }}</p>
                <p v-if="detailData.cache_write_cost > 0">
                  缓存写入费用: ${{ detailData.cache_write_cost.toFixed(4) }}
                </p>
              </div>
            </div>
          </t-col>
        </t-row>

        <div class="detail-item">
          <label>API Key:</label>
          <div v-if="detailData.api_key">
            <t-tag variant="outline">{{ detailData.api_key.name }}</t-tag>
            <span class="text-secondary"> (ID: {{ detailData.api_key_id }})</span>
          </div>
          <span v-else class="text-placeholder">-</span>
        </div>

        <div class="detail-item">
          <label>创建时间:</label>
          <span>{{ formatDateTime(detailData.created_at) }}</span>
        </div>
      </div>
    </t-dialog>
  </div>
</template>
<script setup lang="ts">
import { ChevronDownIcon, ChevronUpIcon, RefreshIcon } from 'tdesign-icons-vue-next';
import type { PrimaryTableCol, TableRowData } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, onMounted, reactive, ref } from 'vue';

import type { Log, LogQueryParams } from '@/api/logs';
import { getLogDetail, getMyLogs } from '@/api/logs';
import { prefix } from '@/config/global';
import { useSettingStore } from '@/store';

defineOptions({
  name: 'MyLogs',
});

const store = useSettingStore();

// 表格列定义
const COLUMNS: PrimaryTableCol<TableRowData>[] = [
  {
    title: 'ID',
    align: 'left',
    width: 140,
    colKey: 'id',
    ellipsis: true,
  },
  {
    title: '账号ID',
    align: 'left',
    width: 100,
    colKey: 'account_id',
  },
  {
    title: '模型名称',
    colKey: 'model_name',
    width: 200,
  },
  {
    title: 'Token使用',
    colKey: 'tokens',
    width: 180,
  },
  {
    title: '费用',
    colKey: 'cost',
    width: 160,
  },
  {
    title: '耗时',
    colKey: 'duration',
    width: 100,
  },
  {
    title: 'API Key',
    colKey: 'api_key',
    width: 140,
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
    width: 80,
    colKey: 'op',
  },
];

// 数据相关
const data = ref<Log[]>([]);
const dataLoading = ref(false);
const modelOptions = ref<string[]>([]);

// 分页
const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
  showJumper: true,
  showSizeChanger: true,
});

// 搜索筛选
const searchFilters = reactive<LogQueryParams>({});
const showAdvancedFilter = ref(false);
const dateRange = ref<string[]>([]);

// 详情弹窗
const detailVisible = ref(false);
const detailData = ref<Log | null>(null);

// 计算属性
const headerAffixedTop = computed(
  () =>
    ({
      offsetTop: store.isUseTabsRouter ? 48 : 0,
      container: `.${prefix}-layout`,
    }) as any,
);

const rowKey = 'id';

// 方法
const fetchData = async () => {
  dataLoading.value = true;
  try {
    const params: LogQueryParams = {
      page: pagination.value.current,
      limit: pagination.value.pageSize,
      ...searchFilters,
    };

    const result = await getMyLogs(params);
    data.value = result.logs || [];
    pagination.value = {
      ...pagination.value,
      total: result.total || 0,
    };

    // 提取模型选项
    const models = Array.from(new Set(result.logs?.map((log) => log.model_name) || []));
    modelOptions.value = models.sort();
  } catch (error) {
    console.error('获取日志列表失败:', error);
    MessagePlugin.error('获取日志列表失败');
  } finally {
    dataLoading.value = false;
  }
};

const handleSearch = () => {
  pagination.value.current = 1;
  fetchData();
};

const handleRefresh = () => {
  fetchData();
};

const resetFilters = () => {
  Object.keys(searchFilters).forEach((key) => {
    delete searchFilters[key as keyof LogQueryParams];
  });
  dateRange.value = [];
  handleSearch();
};

const handleDateRangeChange = (value: any) => {
  if (value && Array.isArray(value) && value.length === 2) {
    // 确保将日期值转换为字符串格式
    searchFilters.start_time = typeof value[0] === 'string' ? value[0] : String(value[0]);
    searchFilters.end_time = typeof value[1] === 'string' ? value[1] : String(value[1]);
  } else {
    delete searchFilters.start_time;
    delete searchFilters.end_time;
  }
  handleSearch();
};

const handlePageChange = (pageInfo: any) => {
  pagination.value.current = pageInfo.current;
  pagination.value.pageSize = pageInfo.pageSize;
  fetchData();
};

const handleViewDetail = async (row: Log) => {
  try {
    detailData.value = await getLogDetail(row.id);
    detailVisible.value = true;
  } catch (error) {
    console.error('获取日志详情失败:', error);
    MessagePlugin.error('获取日志详情失败');
  }
};

const handleDetailCancel = () => {
  detailVisible.value = false;
  detailData.value = null;
};

// 工具函数
const formatNumber = (num: number): string => {
  if (num < 1000) return num.toString();
  if (num < 1000000) return `${(num / 1000).toFixed(1)}K`;
  return `${(num / 1000000).toFixed(1)}M`;
};

const formatDuration = (duration: number): string => {
  if (duration < 1000) return `${duration}ms`;
  return `${(duration / 1000).toFixed(2)}s`;
};

const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleString('zh-CN');
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

.search-row {
  margin-bottom: var(--td-comp-margin-l);
}

.left-operation-container {
  display: flex;
  align-items: center;
}

.right-operation-container {
  display: flex;
  align-items: center;
}

.advanced-filter {
  margin-bottom: var(--td-comp-margin-l);
  padding: var(--td-comp-paddingTB-m) var(--td-comp-paddingLR-m);
  background: var(--td-bg-color-container);
  border-radius: var(--td-radius-default);
}

.tokens-info,
.cost-info {
  font-size: 12px;
  line-height: 1.4;

  p {
    margin: 2px 0;
  }
}

.total-cost {
  color: var(--td-text-color-primary);
  font-weight: 600;
}

.cost-breakdown {
  color: var(--td-text-color-secondary);
  margin-top: 4px;
}

.text-placeholder {
  color: var(--td-text-color-placeholder);
}

.text-secondary {
  color: var(--td-text-color-secondary);
  font-size: 12px;
}

.log-detail {
  .detail-item {
    margin-bottom: 16px;
    display: flex;
    align-items: flex-start;

    label {
      font-weight: 500;
      min-width: 80px;
      margin-right: 12px;
      color: var(--td-text-color-secondary);
    }

    .token-detail,
    .cost-detail {
      p {
        margin: 2px 0;
      }
    }
  }
}
</style>
