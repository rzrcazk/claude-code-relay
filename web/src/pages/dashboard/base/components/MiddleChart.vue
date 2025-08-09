<template>
  <t-row :gutter="16" class="row-container">
    <t-col :xs="12" :xl="9">
      <t-card title="使用趋势分析" subtitle="最近30天" class="dashboard-chart-card" :bordered="false">
        <div
          id="monitorContainer"
          class="dashboard-chart-container"
          :style="{ width: '100%', height: `${resizeTime * 326}px` }"
        />
      </t-card>
    </t-col>
    <t-col :xs="12" :xl="3">
      <t-card title="模型使用分布" subtitle="按费用排序" class="dashboard-chart-card" :bordered="false">
        <div
          id="countContainer"
          class="dashboard-chart-container"
          :style="{ width: `${resizeTime * 326}px`, height: `${resizeTime * 326}px`, margin: '0 auto' }"
        />
      </t-card>
    </t-col>
  </t-row>
</template>
<script setup lang="ts">
import { useWindowSize } from '@vueuse/core';
import { LineChart, PieChart } from 'echarts/charts';
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components';
import * as echarts from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { computed, nextTick, onDeactivated, onMounted, ref, watch } from 'vue';

import type { DashboardStats } from '@/api/dashboard';
import { t } from '@/locales';
import { useSettingStore } from '@/store';
import { changeChartsTheme } from '@/utils/color';
import { LAST_7_DAYS } from '@/utils/date';

import { getLineChartDataSet, getPieChartDataSet } from '../index';

interface Props {
  dashboardData?: DashboardStats;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

echarts.use([TooltipComponent, LegendComponent, PieChart, GridComponent, LineChart, CanvasRenderer]);

const getThisMonth = (checkedValues?: string[]) => {
  let date: Date;
  if (!checkedValues || checkedValues.length === 0) {
    date = new Date();
    return `${date.getFullYear()}-${date.getMonth() + 1}`;
  }
  date = new Date(checkedValues[0]);
  const date2 = new Date(checkedValues[1]);

  const startMonth = date.getMonth() + 1 > 9 ? date.getMonth() + 1 : `0${date.getMonth() + 1}`;
  const endMonth = date2.getMonth() + 1 > 9 ? date2.getMonth() + 1 : `0${date2.getMonth() + 1}`;
  return `${date.getFullYear()}-${startMonth}  至  ${date2.getFullYear()}-${endMonth}`;
};

const store = useSettingStore();
const resizeTime = ref(1);

const chartColors = computed(() => store.chartColors);

// monitorChart
let monitorContainer: HTMLElement;
let monitorChart: echarts.ECharts;
const renderMonitorChart = () => {
  if (!monitorContainer) {
    monitorContainer = document.getElementById('monitorContainer');
  }
  monitorChart = echarts.init(monitorContainer);
  monitorChart.setOption(
    getLineChartDataSet({
      trendData: props.dashboardData?.trend_data || [],
      ...chartColors.value,
    }),
  );
};

// monitorChart
let countContainer: HTMLElement;
let countChart: echarts.ECharts;
const renderCountChart = () => {
  if (!countContainer) {
    countContainer = document.getElementById('countContainer');
  }
  countChart = echarts.init(countContainer);
  countChart.setOption(
    getPieChartDataSet({
      modelStats: props.dashboardData?.model_stats || [],
      ...chartColors.value,
    }),
  );

  // 高亮第一个数据项（最高费用的模型）
  if (props.dashboardData?.model_stats && props.dashboardData.model_stats.length > 0) {
    countChart.dispatchAction({
      type: 'highlight',
      seriesIndex: 0,
      dataIndex: 0,
    });
  }
};

const renderCharts = () => {
  renderMonitorChart();
  renderCountChart();
};

// chartSize update
const updateContainer = () => {
  if (document.documentElement.clientWidth >= 1400 && document.documentElement.clientWidth < 1920) {
    resizeTime.value = Number((document.documentElement.clientWidth / 2080).toFixed(2));
  } else if (document.documentElement.clientWidth < 1080) {
    resizeTime.value = Number((document.documentElement.clientWidth / 1080).toFixed(2));
  } else {
    resizeTime.value = 1;
  }

  monitorChart.resize({
    width: monitorContainer.clientWidth,
    height: resizeTime.value * 326,
  });
  countChart.resize({
    width: resizeTime.value * 326,
    height: resizeTime.value * 326,
  });
};

onMounted(() => {
  renderCharts();
  nextTick(() => {
    updateContainer();
  });
});

const { width, height } = useWindowSize();
watch([width, height], () => {
  updateContainer();
});

// 监听数据变化，重新渲染图表
watch(
  () => props.dashboardData,
  (newData) => {
    if (newData && monitorChart && countChart) {
      // 更新趋势图
      monitorChart.setOption(
        getLineChartDataSet({
          trendData: newData.trend_data || [],
          ...chartColors.value,
        }),
      );

      // 更新饼图
      countChart.setOption(
        getPieChartDataSet({
          modelStats: newData.model_stats || [],
          ...chartColors.value,
        }),
      );

      // 高亮第一个数据项
      if (newData.model_stats && newData.model_stats.length > 0) {
        countChart.dispatchAction({
          type: 'highlight',
          seriesIndex: 0,
          dataIndex: 0,
        });
      }
    }
  },
  { deep: true },
);

onDeactivated(() => {
  storeModeWatch();
  storeBrandThemeWatch();
  storeSidebarCompactWatch();
});

const currentMonth = ref(getThisMonth());

const storeBrandThemeWatch = watch(
  () => store.brandTheme,
  () => {
    changeChartsTheme([monitorChart, countChart]);
  },
);

const storeSidebarCompactWatch = watch(
  () => store.isSidebarCompact,
  () => {
    if (store.isSidebarCompact) {
      nextTick(() => {
        updateContainer();
      });
    } else {
      setTimeout(() => {
        updateContainer();
      }, 180);
    }
  },
);

const storeModeWatch = watch(
  () => store.mode,
  () => {
    [monitorChart, countChart].forEach((item) => {
      item.dispose();
    });

    renderCharts();
  },
);
</script>
<style lang="less" scoped>
.dashboard-chart-card {
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);

  :deep(.t-card__header) {
    padding: 0;
  }

  :deep(.t-card__body) {
    padding: 0;
    margin-top: var(--td-comp-margin-xxl);
  }

  :deep(.t-card__title) {
    font: var(--td-font-title-large);
    font-weight: 400;
  }
}
</style>
