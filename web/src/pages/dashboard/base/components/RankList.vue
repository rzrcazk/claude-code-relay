<template>
  <t-row :gutter="16" class="row-container">
    <t-col :xs="12" :xl="6">
      <t-card title="账号费用排名" subtitle="最近7天" class="dashboard-rank-card" :bordered="false">
        <t-table
          :data="accountRankingList"
          :columns="ACCOUNT_COLUMNS"
          :loading="loading"
          row-key="account_id"
          max-height="600"
        >
          <template #index="{ rowIndex }">
            <span :class="getRankClass(rowIndex)">
              {{ rowIndex + 1 }}
            </span>
          </template>
          <template #cost="{ row }">
            <span>${{ row.cost.toFixed(2) }}</span>
          </template>
          <template #growth_rate="{ row }">
            <trend :type="row.growth_rate >= 0 ? 'up' : 'down'" :describe="formatGrowthRate(row.growth_rate)" />
          </template>
        </t-table>
      </t-card>
    </t-col>
    <t-col :xs="12" :xl="6">
      <t-card title="API Key使用排名" subtitle="最近7天" class="dashboard-rank-card" :bordered="false">
        <t-table
          :data="apiKeyRankingList"
          :columns="APIKEY_COLUMNS"
          :loading="loading"
          row-key="api_key_id"
          max-height="600"
        >
          <template #index="{ rowIndex }">
            <span :class="getRankClass(rowIndex)">
              {{ rowIndex + 1 }}
            </span>
          </template>
          <template #requests="{ row }">
            <span>{{ formatNumber(row.requests) }}</span>
          </template>
          <template #tokens="{ row }">
            <span>{{ formatNumber(row.tokens) }}</span>
          </template>
          <template #growth_rate="{ row }">
            <trend :type="row.growth_rate >= 0 ? 'up' : 'down'" :describe="formatGrowthRate(row.growth_rate)" />
          </template>
        </t-table>
      </t-card>
    </t-col>
  </t-row>
</template>
<script setup lang="ts">
import type { TdBaseTableProps } from 'tdesign-vue-next';
import { computed } from 'vue';

import type { DashboardStats } from '@/api/dashboard';
import Trend from '@/components/trend/index.vue';
import { t } from '@/locales';

interface Props {
  dashboardData?: DashboardStats;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

// 账号排名数据
const accountRankingList = computed(() => {
  return props.dashboardData?.account_ranking || [];
});

// API Key排名数据
const apiKeyRankingList = computed(() => {
  return props.dashboardData?.api_key_ranking || [];
});

const ACCOUNT_COLUMNS: TdBaseTableProps['columns'] = [
  {
    align: 'center',
    colKey: 'index',
    title: '排名',
    width: 70,
    fixed: 'left',
  },
  {
    align: 'left',
    ellipsis: true,
    colKey: 'account_name',
    title: '账号名称',
    width: 120,
  },
  {
    align: 'left',
    ellipsis: true,
    colKey: 'platform_type',
    title: '平台类型',
    width: 150,
  },
  {
    align: 'center',
    colKey: 'cost',
    title: '费用($)',
    width: 80,
  },
  {
    align: 'center',
    colKey: 'growth_rate',
    title: '增长率',
    width: 80,
  },
];

const APIKEY_COLUMNS: TdBaseTableProps['columns'] = [
  {
    align: 'center',
    colKey: 'index',
    title: '排名',
    width: 70,
    fixed: 'left',
  },
  {
    align: 'left',
    ellipsis: true,
    colKey: 'api_key_name',
    title: 'API Key名称',
    width: 150,
  },
  {
    align: 'center',
    colKey: 'requests',
    title: '请求数',
    width: 80,
  },
  {
    align: 'center',
    colKey: 'tokens',
    title: 'Tokens',
    width: 90,
  },
  {
    align: 'center',
    colKey: 'growth_rate',
    title: '增长率',
    width: 80,
  },
];

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

// 格式化增长率
const formatGrowthRate = (rate: number) => {
  const formattedRate = Math.abs(rate).toFixed(1);
  return rate >= 0 ? `+${formattedRate}%` : `-${formattedRate}%`;
};

const rehandleClickOp = (val: MouseEvent) => {
  console.log(val);
};
const getRankClass = (index: number) => {
  return ['dashboard-rank__cell', { 'dashboard-rank__cell--top': index < 3 }];
};
</script>
<style lang="less" scoped>
.dashboard-rank-card {
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);

  :deep(.t-card__header) {
    padding: 0;
  }

  :deep(.t-card__title) {
    font: var(--td-font-title-large);
    font-weight: 400;
  }

  :deep(.t-card__body) {
    padding: 0;
    margin-top: var(--td-comp-margin-xxl);
  }
}

.dashboard-rank__cell {
  display: inline-flex;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  color: white;
  font-size: 14px;
  background-color: var(--td-gray-color-5);
  align-items: center;
  justify-content: center;
  font-weight: 700;

  &--top {
    background: var(--td-brand-color);
  }
}
</style>
