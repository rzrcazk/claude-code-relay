<template>
  <t-card title="今日数据概览" :bordered="false" class="dashboard-overview-card">
    <div class="overview-stats">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-title">今日请求</div>
          <div class="stat-value">{{ loading ? '--' : formatNumber(todayStats.requests) }}</div>
          <div class="stat-change">
            较昨日
            <trend
              :type="getChangeType(todayStats.requests, yesterdayStats.requests)"
              :describe="getChangePercent(todayStats.requests, yesterdayStats.requests)"
            />
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-title">今日Tokens</div>
          <div class="stat-value">{{ loading ? '--' : formatNumber(todayStats.tokens) }}</div>
          <div class="stat-change">
            较昨日
            <trend
              :type="getChangeType(todayStats.tokens, yesterdayStats.tokens)"
              :describe="getChangePercent(todayStats.tokens, yesterdayStats.tokens)"
            />
          </div>
        </div>

        <div class="stat-item">
          <div class="stat-title">今日费用</div>
          <div class="stat-value">${{ loading ? '--' : todayStats.cost.toFixed(2) }}</div>
          <div class="stat-change">
            较昨日
            <trend
              :type="getChangeType(todayStats.cost, yesterdayStats.cost)"
              :describe="getChangePercent(todayStats.cost, yesterdayStats.cost)"
            />
          </div>
        </div>
      </div>
    </div>
  </t-card>
</template>
<script setup lang="ts">
import { computed } from 'vue';

import type { DashboardStats } from '@/api/dashboard';
import Trend from '@/components/trend/index.vue';

interface Props {
  dashboardData?: DashboardStats;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

// 今日数据
const todayStats = computed(() => {
  return props.dashboardData?.today_stats || { requests: 0, tokens: 0, cost: 0 };
});

// 昨日数据
const yesterdayStats = computed(() => {
  return props.dashboardData?.yesterday_stats || { requests: 0, tokens: 0, cost: 0 };
});

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

// 获取变化类型
const getChangeType = (current: number, previous: number) => {
  return current >= previous ? 'up' : 'down';
};

// 获取变化百分比
const getChangePercent = (current: number, previous: number) => {
  if (previous === 0) {
    return current > 0 ? '+100%' : '0%';
  }
  const percent = Math.abs(((current - previous) / previous) * 100);
  const sign = current >= previous ? '+' : '-';
  return `${sign}${percent.toFixed(1)}%`;
};
</script>
<style lang="less" scoped>
.dashboard-overview-card {
  :deep(.t-card__header) {
    padding-bottom: 0;
  }

  :deep(.t-card__title) {
    font: var(--td-font-title-large);
    font-weight: 400;
  }

  :deep(.t-card__body) {
    padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);
  }
}

.overview-stats {
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
  }

  .stat-item {
    text-align: center;

    .stat-title {
      font-size: 14px;
      color: var(--td-text-color-secondary);
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: var(--td-text-color-primary);
      margin-bottom: 4px;
    }

    .stat-change {
      font-size: 12px;
      color: var(--td-text-color-placeholder);
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 4px;
    }
  }
}
</style>
