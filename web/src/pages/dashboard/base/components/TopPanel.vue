<template>
  <t-row :gutter="[16, 16]">
    <t-col v-for="(item, index) in PANE_LIST" :key="item.title" :xs="6" :xl="3">
      <t-card
        :title="item.title"
        :bordered="false"
        class="dashboard-item"
        :class="{ 'dashboard-item--main-color': index === 0 }"
      >
        <div class="dashboard-item-top">
          <span :style="{ fontSize: `${resizeTime * 28}px` }">{{ loading ? '--' : item.number }}</span>
        </div>
        <div class="dashboard-item-left">
          <div
            v-if="index === 0"
            id="moneyContainer"
            class="dashboard-chart-container"
            :style="{ width: `${resizeTime * 120}px`, height: '100px', marginTop: '-24px' }"
          ></div>
          <div
            v-else-if="index === 1"
            id="refundContainer"
            class="dashboard-chart-container"
            :style="{ width: `${resizeTime * 120}px`, height: '56px', marginTop: '-24px' }"
          ></div>
          <span v-else-if="index === 2" :style="{ marginTop: `-24px` }">
            <usergroup-icon />
          </span>
          <span v-else :style="{ marginTop: '-24px' }">
            <file-icon />
          </span>
        </div>
        <template #footer>
          <div class="dashboard-item-bottom">
            <div class="dashboard-item-block">
              较昨日
              <trend
                class="dashboard-item-trend"
                :type="item.upTrend ? 'up' : 'down'"
                :is-reverse-color="index === 0"
                :describe="loading ? '--' : item.trend"
              />
            </div>
            <t-icon name="chevron-right" />
          </div>
        </template>
      </t-card>
    </t-col>
  </t-row>
</template>
<script setup lang="ts">
import { useWindowSize } from '@vueuse/core';
import { BarChart, LineChart } from 'echarts/charts';
import * as echarts from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { FileIcon, UsergroupIcon } from 'tdesign-icons-vue-next';
import { computed, nextTick, onMounted, ref, watch } from 'vue';

import type { DashboardStats } from '@/api/dashboard';
// 导入样式
import Trend from '@/components/trend/index.vue';
import { t } from '@/locales';
import { useSettingStore } from '@/store';
import { changeChartsTheme } from '@/utils/color';

import { constructMiniChart } from '../index';

defineOptions({
  name: 'DashboardTopPanel',
});

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

interface Props {
  dashboardData?: DashboardStats;
  loading?: boolean;
}

echarts.use([LineChart, BarChart, CanvasRenderer]);

const store = useSettingStore();
const resizeTime = ref(1);

// 格式化数字
const formatNumber = (num: number) => {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(1)}M`;
  }
  if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}K`;
  }
  return num.toString();
};

// 格式化费用
const formatCurrency = (amount: number) => {
  return `$${amount.toFixed(2)}`;
};

// 计算增长率
const calculateGrowthRate = (today: number, yesterday: number) => {
  if (yesterday === 0) {
    return today > 0 ? '+100%' : '0%';
  }
  const rate = ((today - yesterday) / yesterday) * 100;
  return rate >= 0 ? `+${rate.toFixed(1)}%` : `${rate.toFixed(1)}%`;
};

// 判断是否为上升趋势
const isUpTrend = (today: number, yesterday: number) => {
  return today >= yesterday;
};

// 动态生成面板数据
const PANE_LIST = computed(() => {
  if (!props.dashboardData) {
    return [
      {
        title: '总费用',
        number: '$0.00',
        trend: '0%',
        upTrend: true,
        leftType: 'echarts-line',
      },
      {
        title: '总Tokens',
        number: '0',
        trend: '0%',
        upTrend: true,
        leftType: 'echarts-bar',
      },
      {
        title: '用户数量',
        number: '0',
        trend: '0%',
        upTrend: true,
        leftType: 'icon-usergroup',
      },
      {
        title: 'API Keys',
        number: '0',
        trend: '0%',
        upTrend: true,
        leftType: 'icon-file-paste',
      },
    ];
  }

  const { today_stats, yesterday_stats } = props.dashboardData;

  return [
    {
      title: '总费用',
      number: formatCurrency(props.dashboardData.total_cost),
      trend: calculateGrowthRate(today_stats.cost, yesterday_stats.cost),
      upTrend: isUpTrend(today_stats.cost, yesterday_stats.cost),
      leftType: 'echarts-line',
    },
    {
      title: '总Tokens',
      number: formatNumber(props.dashboardData.total_tokens),
      trend: calculateGrowthRate(today_stats.tokens, yesterday_stats.tokens),
      upTrend: isUpTrend(today_stats.tokens, yesterday_stats.tokens),
      leftType: 'echarts-bar',
    },
    {
      title: '用户数量',
      number: formatNumber(props.dashboardData.user_count),
      trend: calculateGrowthRate(today_stats.requests, yesterday_stats.requests),
      upTrend: isUpTrend(today_stats.requests, yesterday_stats.requests),
      leftType: 'icon-usergroup',
    },
    {
      title: 'API Keys',
      number: formatNumber(props.dashboardData.api_key_count),
      trend: '活跃中',
      upTrend: true,
      leftType: 'icon-file-paste',
    },
  ];
});

// moneyCharts
let moneyContainer: HTMLElement;
let moneyChart: echarts.ECharts;
const renderMoneyChart = () => {
  if (!moneyContainer) {
    moneyContainer = document.getElementById('moneyContainer');
  }
  moneyChart = echarts.init(moneyContainer);
  moneyChart.setOption(constructMiniChart('line', props.dashboardData?.trend_data || []));
};

// refundCharts
let refundContainer: HTMLElement;
let refundChart: echarts.ECharts;
const renderRefundChart = () => {
  if (!refundContainer) {
    refundContainer = document.getElementById('refundContainer');
  }
  refundChart = echarts.init(refundContainer);
  refundChart.setOption(constructMiniChart('bar', props.dashboardData?.trend_data || []));
};

const renderCharts = () => {
  renderMoneyChart();
  renderRefundChart();
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
  moneyChart.resize({
    width: resizeTime.value * 120,
    // height: resizeTime.value * 100,
  });
  refundChart.resize({
    width: resizeTime.value * 120,
    // height: resizeTime.value * 56,
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

watch(
  () => store.brandTheme,
  () => {
    changeChartsTheme([refundChart]);
  },
);

watch(
  () => store.mode,
  () => {
    [moneyChart, refundChart].forEach((item) => {
      item.dispose();
    });

    renderCharts();
  },
);

// 监听数据变化，重新渲染图表
watch(
  () => props.dashboardData,
  (newData) => {
    if (newData && moneyChart && refundChart) {
      moneyChart.setOption(constructMiniChart('line', newData.trend_data || []));
      refundChart.setOption(constructMiniChart('bar', newData.trend_data || []));
    }
  },
  { deep: true },
);
</script>
<style lang="less" scoped>
.dashboard-item {
  padding: var(--td-comp-paddingTB-xl) var(--td-comp-paddingLR-xxl);

  :deep(.t-card__header) {
    padding: 0;
  }

  :deep(.t-card__footer) {
    padding: 0;
  }

  :deep(.t-card__title) {
    font: var(--td-font-body-medium);
    color: var(--td-text-color-secondary);
  }

  :deep(.t-card__body) {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    flex: 1;
    position: relative;
    padding: 0;
    margin-top: var(--td-comp-margin-s);
    margin-bottom: var(--td-comp-margin-xxl);
  }

  &:hover {
    cursor: pointer;
  }

  &-top {
    display: flex;
    flex-direction: row;
    align-items: flex-start;

    > span {
      display: inline-block;
      color: var(--td-text-color-primary);
      font-size: var(--td-font-size-headline-medium);
      line-height: var(--td-line-height-headline-medium);
    }
  }

  &-bottom {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;

    > .t-icon {
      cursor: pointer;
      font-size: var(--td-comp-size-xxxs);
    }
  }

  &-block {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--td-text-color-placeholder);
  }

  &-trend {
    margin-left: var(--td-comp-margin-s);
  }

  &-left {
    position: absolute;
    top: 0;
    right: 0;

    > span {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: var(--td-comp-size-xxxl);
      height: var(--td-comp-size-xxxl);
      background: var(--td-brand-color-light);
      border-radius: 50%;

      .t-icon {
        font-size: 24px;
        color: var(--td-brand-color);
      }
    }
  }

  /* 针对第一个卡片需要反色处理 */
  &--main-color {
    background: var(--td-brand-color);
    color: var(--td-text-color-primary);

    :deep(.t-card__title),
    .dashboard-item-top span,
    .dashboard-item-bottom {
      color: var(--td-text-color-anti);
    }

    .dashboard-item-block {
      color: var(--td-text-color-anti);
      opacity: 0.6;
    }

    .dashboard-item-bottom {
      color: var(--td-text-color-anti);
    }
  }
}
</style>
