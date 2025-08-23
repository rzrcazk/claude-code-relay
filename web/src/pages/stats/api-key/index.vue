<template>
  <div class="api-key-stats-container">
    <!-- 顶部API Key输入区域 -->
    <div class="header-section">
      <div class="page-title">
        <h1>API Key 使用统计查询</h1>
        <p>输入API Key查看使用统计和调用日志</p>
      </div>

      <t-card class="api-key-input-card">
        <div class="api-key-form">
          <t-form :model="searchForm">
            <t-form-item label="API Key" required>
              <t-input
                v-model="searchForm.apiKey"
                placeholder="请输入API Key (sk-...)"
                clearable
                :loading="loading"
                @enter="handleSearch"
                @blur="handleInputBlur"
                @clear="handleClear"
              />
            </t-form-item>
          </t-form>
        </div>
      </t-card>
    </div>

    <!-- 统计数据展示区域 -->
    <div v-if="statsData" class="stats-section">
      <!-- API Key基本信息 -->
      <t-card class="api-key-info-card" title="API Key 信息">
        <div class="api-key-info">
          <div class="info-item">
            <span class="label">名称：</span>
            <span class="value">{{ statsData.api_key_info.name }}</span>
          </div>
          <div class="info-item">
            <span class="label">状态：</span>
            <t-tag v-if="statsData.api_key_info.status === 1" theme="success" variant="light"> 启用 </t-tag>
            <t-tag v-else theme="danger" variant="light"> 禁用 </t-tag>
          </div>
        </div>
      </t-card>

      <!-- 30天统计概览 -->
      <t-card class="overview-card" title="最近30天统计概览">
        <div class="overview-stats">
          <div class="stat-item">
            <div class="stat-value">{{ formatNumber(statsData.stats.summary.total_requests) }}</div>
            <div class="stat-label">总请求数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ formatNumber(statsData.stats.summary.total_tokens) }}</div>
            <div class="stat-label">总Token数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">${{ formatCost(statsData.stats.summary.total_cost) }}</div>
            <div class="stat-label">总费用</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ formatDuration(statsData.stats.summary.avg_duration) }}</div>
            <div class="stat-label">平均响应时间</div>
          </div>
        </div>
      </t-card>

      <!-- 趋势图表 -->
      <t-card v-if="chartData.length > 0" class="chart-card" title="使用趋势">
        <div ref="chartContainer" class="chart-container"></div>
      </t-card>
    </div>

    <!-- 日志列表 -->
    <div v-if="logsData" class="logs-section">
      <t-card class="logs-card" title="调用日志">
        <t-table
          :data="logsData.list"
          :columns="logColumns"
          :pagination="pagination"
          :loading="loading"
          @page-change="handlePageChange"
        >
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
            </div>
            <span v-else class="text-placeholder">-</span>
          </template>

          <template #created_at="{ row }">
            <span>{{ formatDateTime(row.created_at) }}</span>
          </template>
        </t-table>
      </t-card>
    </div>

    <!-- 空状态 -->
    <t-empty v-if="!loading && !statsData && searchAttempted" description="请输入有效的API Key进行查询" />
  </div>
</template>
<script setup lang="ts">
import { LineChart } from 'echarts/charts';
import { DataZoomComponent, GridComponent, LegendComponent, TooltipComponent } from 'echarts/components';
import * as echarts from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import type { PrimaryTableCol, TableRowData } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, nextTick, onMounted, onUnmounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { getApiKeyStats } from '@/api/stats';

// 注册ECharts组件
echarts.use([LineChart, GridComponent, TooltipComponent, LegendComponent, DataZoomComponent, CanvasRenderer]);

// 响应式数据
const loading = ref(false);
const searchAttempted = ref(false);
const chartContainer = ref<HTMLElement>();
let chartInstance: echarts.ECharts | null = null;

const searchForm = reactive({
  apiKey: '',
});

const statsData = ref<any>(null);
const logsData = ref<any>(null);

// 路由相关
const route = useRoute();
const router = useRouter();

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showJumper: true,
});

// 表格列配置
const logColumns: PrimaryTableCol<TableRowData>[] = [
  {
    title: 'ID',
    align: 'left',
    width: 140,
    colKey: 'id',
    ellipsis: true,
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
  },
];

// 图表数据
const chartData = computed(() => {
  if (!statsData.value?.stats?.trend_data) return [];
  return statsData.value.stats.trend_data;
});

// 初始化页面
onMounted(() => {
  // 检查URL参数
  const urlApiKey = route.query.api_key as string;
  if (urlApiKey) {
    searchForm.apiKey = urlApiKey;
    handleSearch();
  }
});

// 处理搜索
async function handleSearch() {
  if (!searchForm.apiKey.trim()) {
    MessagePlugin.warning('请输入API Key');
    return;
  }

  loading.value = true;
  searchAttempted.value = true;

  try {
    const params = {
      api_key: searchForm.apiKey,
      page: pagination.current,
      limit: pagination.pageSize,
    };

    const response = await getApiKeyStats(params);

    statsData.value = response;
    logsData.value = response.logs;

    // 更新分页信息
    pagination.total = response.logs.total;
    pagination.current = response.logs.page;

    // 更新URL参数
    router.replace({
      query: { ...route.query, api_key: searchForm.apiKey },
    });

    // 渲染图表
    nextTick(() => {
      renderChart();
    });
  } catch (error: any) {
    console.error('查询失败:', error);
    MessagePlugin.error(error.message || '查询失败');
    statsData.value = null;
    logsData.value = null;
  } finally {
    loading.value = false;
  }
}

// 处理分页变化
async function handlePageChange(pageInfo: any) {
  pagination.current = pageInfo.current;
  pagination.pageSize = pageInfo.pageSize;
  await handleSearch();
}

// 处理输入框失去焦点
function handleInputBlur() {
  if (searchForm.apiKey.trim() && searchForm.apiKey.trim() !== '') {
    handleSearch();
  }
}

// 处理清除输入框
function handleClear() {
  statsData.value = null;
  logsData.value = null;
  searchAttempted.value = false;
  // 清除URL参数
  router.replace({
    query: { ...route.query, api_key: undefined },
  });
}

// 渲染图表
function renderChart() {
  if (!chartContainer.value || chartData.value.length === 0) return;

  if (chartInstance) {
    chartInstance.dispose();
  }

  chartInstance = echarts.init(chartContainer.value);

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
    },
    legend: {
      data: ['请求数', 'Token数', '费用'],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: chartData.value.map((item: any) => item.date),
    },
    yAxis: [
      {
        type: 'value',
        name: '请求数/Token数',
        position: 'left',
      },
      {
        type: 'value',
        name: '费用($)',
        position: 'right',
      },
    ],
    series: [
      {
        name: '请求数',
        type: 'line',
        data: chartData.value.map((item: any) => item.requests),
        smooth: true,
      },
      {
        name: 'Token数',
        type: 'line',
        data: chartData.value.map((item: any) => item.tokens),
        smooth: true,
      },
      {
        name: '费用',
        type: 'line',
        yAxisIndex: 1,
        data: chartData.value.map((item: any) => item.cost),
        smooth: true,
      },
    ],
  };

  chartInstance.setOption(option);
}

// 格式化数字
function formatNumber(value: number): string {
  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(1)}M`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}K`;
  }
  return value.toString();
}

// 格式化费用
function formatCost(value: number): string {
  return value.toFixed(4);
}

// 格式化时长
function formatDuration(duration: number): string {
  if (duration < 1000) return `${duration}ms`;
  return `${(duration / 1000).toFixed(2)}s`;
}

// 格式化日期时间
function formatDateTime(dateStr: string): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '-';
  return date.toLocaleString('zh-CN');
}

// 组件销毁时清理图表
onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose();
  }
});
</script>
<style lang="less" scoped>
.api-key-stats-container {
  padding: 24px;

  .header-section {
    margin-bottom: 24px;

    .page-title {
      margin-bottom: 16px;
      text-align: center;

      h1 {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        color: var(--td-text-color-primary);
      }

      p {
        margin: 0;
        color: var(--td-text-color-secondary);
        font-size: 14px;
      }
    }

    .api-key-input-card {
      max-width: 600px;
      margin: 0 auto;

      .api-key-form {
        :deep(.t-form-item) {
          margin-bottom: 0;
        }
      }
    }
  }

  .stats-section {
    margin-bottom: 24px;

    .api-key-info-card {
      margin-bottom: 16px;

      .api-key-info {
        display: flex;
        gap: 24px;

        .info-item {
          .label {
            color: var(--td-text-color-secondary);
          }

          .value {
            font-weight: 500;
            color: var(--td-text-color-primary);
          }
        }
      }
    }

    .overview-card {
      margin-bottom: 16px;

      .overview-stats {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 24px;

        .stat-item {
          text-align: center;
          padding: 16px;
          background: var(--td-bg-color-container-hover);
          border-radius: 8px;

          .stat-value {
            font-size: 24px;
            font-weight: 600;
            color: var(--td-color-primary);
            margin-bottom: 4px;
          }

          .stat-label {
            color: var(--td-text-color-secondary);
            font-size: 14px;
          }
        }
      }
    }

    .chart-card {
      .chart-container {
        height: 400px;
        width: 100%;
      }
    }
  }

  .logs-section {
    .logs-card {
      :deep(.t-table) {
        .t-table__cell {
          padding: 12px 16px;
        }
      }
    }
  }

  .logs-section {
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

    .mt-2 {
      margin-top: 4px;
    }
  }
}

@media (max-width: 768px) {
  .api-key-stats-container {
    padding: 16px;

    .stats-section {
      .api-key-info {
        flex-direction: column;
        gap: 12px;
      }

      .overview-stats {
        grid-template-columns: repeat(2, 1fr);
        gap: 12px;

        .stat-item {
          padding: 12px;
        }
      }
    }
  }
}
</style>
