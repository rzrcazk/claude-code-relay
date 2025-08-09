<template>
  <div>
    <!-- 顶部 card  -->
    <top-panel :dashboard-data="dashboardData" :loading="loading" class="row-container" />
    <!-- 今日数据概览 -->
    <output-overview :dashboard-data="dashboardData" :loading="loading" class="row-container" />
    <!-- 中部图表  -->
    <middle-chart :dashboard-data="dashboardData" :loading="loading" class="row-container" />
    <!-- 列表排名 -->
    <rank-list :dashboard-data="dashboardData" :loading="loading" class="row-container" />
  </div>
</template>
<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next';
import { onMounted, ref } from 'vue';

import type { DashboardStats } from '@/api/dashboard';
import { getDashboardStats } from '@/api/dashboard';

import MiddleChart from './components/MiddleChart.vue';
import OutputOverview from './components/OutputOverview.vue';
import RankList from './components/RankList.vue';
import TopPanel from './components/TopPanel.vue';

defineOptions({
  name: 'DashboardBase',
});

// 响应式数据
const loading = ref(false);
const dashboardData = ref<DashboardStats | null>(null);

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    loading.value = true;
    const response = await getDashboardStats();
    dashboardData.value = response;
  } catch (error) {
    console.error('获取仪表盘数据失败:', error);
    MessagePlugin.error('获取仪表盘数据失败');
  } finally {
    loading.value = false;
  }
};

// 组件挂载时获取数据
onMounted(() => {
  fetchDashboardData();
});
</script>
<style scoped>
.row-container:not(:last-child) {
  margin-bottom: 16px;
}
</style>
