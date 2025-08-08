<template>
  <div>
    <t-card class="stats-container" :bordered="false">
      <t-empty icon="chart-bar" description="使用统计功能开发中...">
        <template #action>
          <t-button theme="primary" @click="goToMyLogs">查看我的日志</t-button>
        </template>
      </t-empty>
    </t-card>
  </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router';

defineOptions({
  name: 'LogStats',
});

const router = useRouter();

const goToMyLogs = () => {
  router.push('/logs/my');
};
</script>
<style lang="less" scoped>
.stats-container {
  padding: var(--td-comp-paddingTB-xxl) var(--td-comp-paddingLR-xxl);
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;

  :deep(.t-card__body) {
    padding: 0;
  }
}
</style>
