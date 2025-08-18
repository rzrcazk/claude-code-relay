<template>
  <t-dialog v-model:visible="visible" header="个人中心" width="600px" :footer="false" @close="handleClose">
    <t-tabs v-model="activeTab" :lazy="false">
      <t-tab-panel value="profile" label="基本信息">
        <div class="profile-info">
          <div class="info-item">
            <div class="info-label">用户ID</div>
            <div class="info-value">{{ userInfo.id }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">用户名</div>
            <div class="info-value">{{ userInfo.name || userInfo.username }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">邮箱</div>
            <div class="info-value">{{ userInfo.email }}</div>
          </div>
          <div class="info-item">
            <div class="info-label">角色</div>
            <div class="info-value">
              <t-tag :theme="userInfo.role === 'admin' ? 'success' : 'default'">
                {{ userInfo.role === 'admin' ? '管理员' : '普通用户' }}
              </t-tag>
            </div>
          </div>
          <div class="info-item">
            <div class="info-label">状态</div>
            <div class="info-value">
              <t-tag :theme="userInfo.status === 1 ? 'success' : 'danger'">
                {{ userInfo.status === 1 ? '正常' : '禁用' }}
              </t-tag>
            </div>
          </div>
          <div class="info-item">
            <div class="info-label">注册时间</div>
            <div class="info-value">{{ formatDate(userInfo.created_at) }}</div>
          </div>
        </div>
      </t-tab-panel>

      <t-tab-panel value="email" label="修改邮箱">
        <t-form ref="emailFormRef" :model="emailForm" label-width="80px" @submit="handleChangeEmail">
          <t-form-item label="当前邮箱">
            <t-input :value="userInfo.email" disabled />
          </t-form-item>
          <t-form-item label="新邮箱" name="newEmail">
            <t-input v-model="emailForm.newEmail" placeholder="请输入新邮箱" />
          </t-form-item>
          <t-form-item label="确认邮箱" name="confirmEmail">
            <t-input v-model="emailForm.confirmEmail" placeholder="请确认新邮箱" />
          </t-form-item>
          <t-form-item>
            <t-space>
              <t-button theme="primary" type="submit" :loading="emailLoading"> 修改邮箱 </t-button>
              <t-button variant="outline" @click="resetEmailForm"> 重置 </t-button>
            </t-space>
          </t-form-item>
        </t-form>
      </t-tab-panel>

      <t-tab-panel value="password" label="修改密码">
        <t-form ref="passwordFormRef" :model="passwordForm" label-width="80px" @submit="handleChangePassword">
          <t-form-item label="当前密码" name="oldPassword">
            <t-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入当前密码" />
          </t-form-item>
          <t-form-item label="新密码" name="newPassword">
            <t-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" />
          </t-form-item>
          <t-form-item label="确认密码" name="confirmPassword">
            <t-input v-model="passwordForm.confirmPassword" type="password" placeholder="请确认新密码" />
          </t-form-item>
          <t-form-item>
            <t-space>
              <t-button theme="primary" type="submit" :loading="passwordLoading"> 修改密码 </t-button>
              <t-button variant="outline" @click="resetPasswordForm"> 重置 </t-button>
            </t-space>
          </t-form-item>
        </t-form>
      </t-tab-panel>
    </t-tabs>
  </t-dialog>
</template>
<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, reactive, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { changeUserEmail, changeUserPassword } from '@/api/user';
import { useUserStore } from '@/store';

interface Props {
  visible: boolean;
}

interface Emits {
  (e: 'update:visible', value: boolean): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const router = useRouter();
const userStore = useUserStore();
const userInfo = computed(() => userStore.userProfile);

// 弹窗状态
const visible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value),
});

// 当前激活的标签页
const activeTab = ref('profile');

// 加载状态
const emailLoading = ref(false);
const passwordLoading = ref(false);

// 表单引用
const emailFormRef = ref();
const passwordFormRef = ref();

// 修改邮箱表单
const emailForm = reactive({
  newEmail: '',
  confirmEmail: '',
});

// 修改密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});

// 日期格式化函数
const formatDate = (dateString: string) => {
  if (!dateString) return '暂无';
  try {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  } catch (error) {
    return '暂无';
  }
};

// 修改邮箱
const handleChangeEmail = async ({ validateResult }: any) => {
  if (validateResult === true) {
    // 手动验证邮箱确认
    if (emailForm.newEmail !== emailForm.confirmEmail) {
      MessagePlugin.error('两次输入的邮箱不一致');
      return;
    }

    try {
      emailLoading.value = true;
      await changeUserEmail({
        email: emailForm.newEmail,
      });

      // 更新本地用户信息
      await userStore.getUserInfo();

      MessagePlugin.success('邮箱修改成功');
      resetEmailForm();
      activeTab.value = 'profile';
    } catch (error) {
      MessagePlugin.error('邮箱修改失败');
    } finally {
      emailLoading.value = false;
    }
  }
};

// 修改密码
const handleChangePassword = async ({ validateResult }: any) => {
  if (validateResult === true) {
    // 手动验证密码确认
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      MessagePlugin.error('两次输入的密码不一致');
      return;
    }

    try {
      passwordLoading.value = true;
      await changeUserPassword({
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword,
      });

      MessagePlugin.success('密码修改成功，即将退出登录');

      // 关闭弹窗
      visible.value = false;

      // 延迟1秒后退出登录
      setTimeout(async () => {
        // 清除用户信息和token
        await userStore.logout();

        // 跳转到登录页
        router.push({
          path: '/login',
          query: { redirect: encodeURIComponent(router.currentRoute.value.fullPath) },
        });
      }, 1000);
    } catch (error) {
      MessagePlugin.error('密码修改失败');
    } finally {
      passwordLoading.value = false;
    }
  }
};

// 重置表单
const resetEmailForm = () => {
  emailForm.newEmail = '';
  emailForm.confirmEmail = '';
  emailFormRef.value?.clearValidate();
};

const resetPasswordForm = () => {
  passwordForm.oldPassword = '';
  passwordForm.newPassword = '';
  passwordForm.confirmPassword = '';
  passwordFormRef.value?.clearValidate();
};

// 关闭弹窗
const handleClose = () => {
  visible.value = false;
  activeTab.value = 'profile';
  resetEmailForm();
  resetPasswordForm();
};
</script>
<style lang="less" scoped>
:deep(.t-dialog__body) {
  padding: 20px;
}

:deep(.t-form-item__label) {
  font-weight: 500;
}

:deep(.t-tabs__header) {
  margin-bottom: 20px;
}

:deep(.t-tab-panel) {
  padding-top: 10px;
}

.profile-info {
  .info-item {
    display: flex;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--td-component-stroke);

    &:last-child {
      border-bottom: none;
    }

    .info-label {
      flex: 0 0 80px;
      font-weight: 500;
      color: var(--td-text-color-primary);
    }

    .info-value {
      flex: 1;
      color: var(--td-text-color-secondary);
      display: flex;
      align-items: center;
    }
  }
}
</style>
