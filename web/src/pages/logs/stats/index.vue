<template>
  <div class="stats-page">
    <!-- 筛选器 -->
    <t-card class="filter-card" :bordered="false">
      <t-form :data="queryParams" layout="inline">
        <t-form-item label="时间区间">
          <t-date-range-picker
            v-model="queryParams.dateRange"
            placeholder="选择时间区间"
            clearable
            format="YYYY-MM-DD"
            style="width: 280px"
            @change="handleDateRangeChange"
          />
        </t-form-item>

        <t-form-item label="账号">
          <t-input
            v-model="queryParams.account_filter"
            placeholder="输入账号ID或邮箱"
            clearable
            style="width: 180px"
            @blur="handleInputChange"
            @clear="handleInputChange"
          />
        </t-form-item>

        <t-form-item label="API Key">
          <t-input
            v-model="queryParams.api_key_filter"
            placeholder="输入Key ID或秘钥"
            clearable
            style="width: 180px"
            @blur="handleInputChange"
            @clear="handleInputChange"
          />
        </t-form-item>
      </t-form>
    </t-card>

    <!-- 统计卡片 -->
    <div v-if="statsData?.summary" class="stats-cards">
      <t-row :gutter="[16, 16]">
        <t-col :xs="12" :sm="6">
          <t-card hover>
            <div class="stat-item">
              <div class="stat-value">{{ formatNumber(statsData.summary.total_requests) }}</div>
              <div class="stat-label">总请求数</div>
            </div>
          </t-card>
        </t-col>

        <t-col :xs="12" :sm="6">
          <t-card hover>
            <div class="stat-item">
              <div class="stat-value">{{ formatNumber(statsData.summary.total_tokens) }}</div>
              <div class="stat-label">总Tokens</div>
            </div>
          </t-card>
        </t-col>

        <t-col :xs="12" :sm="6">
          <t-card hover>
            <div class="stat-item">
              <div class="stat-value">${{ formatCurrency(statsData.summary.total_cost) }}</div>
              <div class="stat-label">总费用</div>
            </div>
          </t-card>
        </t-col>

        <t-col :xs="12" :sm="6">
          <t-card hover>
            <div class="stat-item">
              <div class="stat-value">{{ Math.round(statsData.summary.avg_duration) }}ms</div>
              <div class="stat-label">平均响应时间</div>
            </div>
          </t-card>
        </t-col>
      </t-row>
    </div>

    <!-- Token使用和费用分布 -->
    <t-card v-if="statsData?.summary" title="Token使用和费用分布" :bordered="false" class="distribution-card">
      <t-row :gutter="[16, 16]">
        <t-col :span="12">
          <div class="chart-section">
            <h4 class="chart-title">Token 使用分布</h4>
            <div ref="tokenChartRef" class="chart-container"></div>
          </div>
        </t-col>

        <t-col :span="12">
          <div class="chart-section">
            <h4 class="chart-title">费用分布</h4>
            <div ref="costChartRef" class="chart-container"></div>
          </div>
        </t-col>
      </t-row>
    </t-card>

    <!-- 趋势图 -->
    <t-card v-if="statsData?.trend_data?.length" title="使用趋势" :bordered="false" class="trend-card">
      <div ref="trendChartRef" class="trend-chart"></div>
    </t-card>

    <!-- 加载状态 -->
    <t-card v-if="loading" class="loading-card">
      <t-loading size="large" text="正在加载统计数据..." />
    </t-card>

    <!-- 空状态 -->
    <t-card v-if="!loading && !statsData" class="empty-card">
      <t-empty icon="chart-bar" description="暂无统计数据，请调整筛选条件" />
    </t-card>
  </div>
</template>
<script setup lang="ts">
import * as echarts from 'echarts';
import { MessagePlugin } from 'tdesign-vue-next';
import { nextTick, onMounted, onUnmounted, reactive, ref, watch } from 'vue';

import type { StatsQueryParams, StatsResponse } from '@/api/logs';
import { getMyUsageStats } from '@/api/logs';

defineOptions({
  name: 'LogStats',
});

// 响应式数据
const loading = ref(false);
const statsData = ref<StatsResponse | null>(null);
const tokenChartRef = ref<HTMLDivElement>();
const costChartRef = ref<HTMLDivElement>();
const trendChartRef = ref<HTMLDivElement>();

// 初始化默认时间范围（最近7天）
const initDefaultDateRange = () => {
  const today = new Date();
  const sevenDaysAgo = new Date();
  sevenDaysAgo.setDate(today.getDate() - 6); // 包含今天，往前6天
  
  const formatDate = (date: Date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };
  
  return [formatDate(sevenDaysAgo), formatDate(today)];
};

// 查询参数
const queryParams = reactive({
  dateRange: initDefaultDateRange(),
  account_filter: '' as string,
  api_key_filter: '' as string,
});

// ECharts 实例
let tokenChart: echarts.ECharts | null = null;
let costChart: echarts.ECharts | null = null;
let trendChart: echarts.ECharts | null = null;

// 数字格式化
const formatNumber = (num: number) => {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(1)}M`;
  }
  if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}K`;
  }
  return num.toString();
};

// 货币格式化
const formatCurrency = (amount: number) => {
  return amount.toFixed(4);
};

// 获取统计数据
const fetchStats = async () => {
  try {
    loading.value = true;
    const params: StatsQueryParams = {};

    // 处理日期区间
    if (queryParams.dateRange && queryParams.dateRange.length === 2) {
      // 有时间区间选择：直接传递开始和结束时间
      params.start_time = `${queryParams.dateRange[0]} 00:00:00`;
      params.end_time = `${queryParams.dateRange[1]} 23:59:59`;
    }
    // 没有时间区间选择：后端会默认显示当天数据

    if (queryParams.account_filter) {
      params.account_filter = queryParams.account_filter;
    }

    if (queryParams.api_key_filter) {
      params.api_key_filter = queryParams.api_key_filter;
    }

    console.log('请求参数:', params); // 调试日志

    const response = await getMyUsageStats(params);
    console.log('获取到的数据:', response); // 调试日志
    statsData.value = response;

    // 等待DOM更新后初始化图表
    await nextTick();
    console.log('准备更新图表，statsData:', statsData.value); // 调试日志
    initCharts();
  } catch (error) {
    console.error('获取统计数据失败:', error);
    MessagePlugin.error('获取统计数据失败');
  } finally {
    loading.value = false;
  }
};

// 初始化图表
const initCharts = () => {
  if (!statsData.value?.summary) return;

  initTokenChart();
  initCostChart();
  initTrendChart();
};

// Token使用分布饼图
const initTokenChart = () => {
  if (!tokenChartRef.value || !statsData.value?.summary) return;

  if (tokenChart) {
    tokenChart.dispose();
  }

  tokenChart = echarts.init(tokenChartRef.value);

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
    },
    series: [
      {
        name: 'Token使用',
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: statsData.value.summary.total_input_tokens, name: '输入Tokens' },
          { value: statsData.value.summary.total_output_tokens, name: '输出Tokens' },
          { value: statsData.value.summary.total_cache_read_tokens, name: '缓存读取Tokens' },
          { value: statsData.value.summary.total_cache_creation_tokens, name: '缓存创建Tokens' },
        ].filter((item) => item.value > 0),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  };

  tokenChart.setOption(option);
};

// 费用分布饼图
const initCostChart = () => {
  if (!costChartRef.value || !statsData.value?.summary) return;

  if (costChart) {
    costChart.dispose();
  }

  costChart = echarts.init(costChartRef.value);

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: ${c} ({d}%)',
    },
    series: [
      {
        name: '费用分布',
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: statsData.value.summary.input_cost, name: '输入费用' },
          { value: statsData.value.summary.output_cost, name: '输出费用' },
          { value: statsData.value.summary.cache_write_cost, name: '缓存写入费用' },
          { value: statsData.value.summary.cache_read_cost, name: '缓存读取费用' },
        ].filter((item) => item.value > 0),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
        },
      },
    ],
  };

  costChart.setOption(option);
};

// 趋势图
const initTrendChart = () => {
  if (!trendChartRef.value || !statsData.value?.trend_data?.length) return;

  if (trendChart) {
    trendChart.dispose();
  }

  trendChart = echarts.init(trendChartRef.value);

  const dates = statsData.value.trend_data.map((item) => item.date);

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
    },
    legend: {
      data: ['请求数', 'Tokens', '费用', '响应时间'],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
    },
    yAxis: [
      {
        type: 'value',
        name: '请求数/Tokens',
        position: 'left',
      },
      {
        type: 'value',
        name: '费用($)/响应时间(ms)',
        position: 'right',
      },
    ],
    series: [
      {
        name: '请求数',
        type: 'line',
        data: statsData.value.trend_data.map((item) => item.requests),
        smooth: true,
      },
      {
        name: 'Tokens',
        type: 'line',
        data: statsData.value.trend_data.map((item) => item.tokens),
        smooth: true,
      },
      {
        name: '费用',
        type: 'line',
        yAxisIndex: 1,
        data: statsData.value.trend_data.map((item) => item.cost),
        smooth: true,
      },
      {
        name: '响应时间',
        type: 'line',
        yAxisIndex: 1,
        data: statsData.value.trend_data.map((item) => item.avg_duration),
        smooth: true,
      },
    ],
  };

  trendChart.setOption(option);
};

// 日期范围变化处理
const handleDateRangeChange = () => {
  fetchStats();
};

// 输入框失去焦点或清除时处理
const handleInputChange = () => {
  fetchStats();
};


// 组件挂载
onMounted(() => {
  fetchStats();

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    tokenChart?.resize();
    costChart?.resize();
    trendChart?.resize();
  });
});

// 组件卸载时清理
onUnmounted(() => {
  tokenChart?.dispose();
  costChart?.dispose();
  trendChart?.dispose();
});
</script>
<style lang="less" scoped>
.stats-page {
  padding: 16px;

  .filter-card {
    margin-bottom: 16px;
  }

  .stats-cards {
    margin-bottom: 24px;

    .stat-item {
      text-align: center;

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        color: var(--td-text-color-primary);
        margin-bottom: 8px;
      }

      .stat-label {
        font-size: 14px;
        color: var(--td-text-color-secondary);
      }
    }
  }

  .distribution-card {
    margin-bottom: 24px;
  }

  .trend-card {
    margin-top: 24px;
  }

  .chart-section {
    .chart-title {
      font-size: 16px;
      font-weight: 500;
      color: var(--td-text-color-primary);
      margin-bottom: 16px;
      text-align: center;
    }
  }

  // 强制保持水平布局
  :deep(.t-row) {
    display: flex !important;
    flex-wrap: nowrap !important;

    .t-col {
      flex: 1 !important;
      min-width: 0 !important;
    }
  }

  .chart-container {
    height: 300px;
    width: 100%;
  }

  .trend-chart {
    height: 400px;
    width: 100%;
  }

  .loading-card,
  .empty-card {
    margin-top: 60px;
    text-align: center;
    min-height: 300px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .stats-page {
    padding: 8px;

    .stats-cards {
      .stat-value {
        font-size: 20px;
      }
    }

    .chart-container {
      height: 250px;
    }

    .trend-chart {
      height: 300px;
    }
  }
}
</style>
