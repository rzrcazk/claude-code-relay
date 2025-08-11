<template>
  <t-form
    ref="form"
    class="item-container"
    :class="[`login-${type}`]"
    :data="formData"
    :rules="FORM_RULES"
    label-width="0"
    @submit="onSubmit"
  >
    <template v-if="type === 'password'">
      <t-form-item name="account">
        <t-input v-model="formData.account" size="large" :placeholder="`${t('pages.login.input.account')}：admin`">
          <template #prefix-icon>
            <t-icon name="user" />
          </template>
        </t-input>
      </t-form-item>

      <t-form-item name="password">
        <t-input
          v-model="formData.password"
          size="large"
          :type="showPsw ? 'text' : 'password'"
          clearable
          :placeholder="`${t('pages.login.input.password')}：admin`"
        >
          <template #prefix-icon>
            <t-icon name="lock-on" />
          </template>
          <template #suffix-icon>
            <t-icon :name="showPsw ? 'browse' : 'browse-off'" @click="showPsw = !showPsw" />
          </template>
        </t-input>
      </t-form-item>

      <div class="check-container remember-pwd">
        <t-checkbox>{{ t('pages.login.remember') }}</t-checkbox>
        <span class="tip">{{ t('pages.login.forget') }}</span>
      </div>
    </template>

    <!-- 邮箱登录 -->
    <template v-else>
      <t-form-item name="email">
        <t-input v-model="formData.email" size="large" placeholder="请输入邮箱">
          <template #prefix-icon>
            <t-icon name="mail" />
          </template>
        </t-input>
      </t-form-item>

      <t-form-item class="verification-code" name="verifyCode">
        <t-input v-model="formData.verifyCode" size="large" placeholder="请输入验证码" />
        <t-button size="large" variant="outline" :disabled="countDown > 0 || !formData.email" @click="sendCode">
          {{ countDown === 0 ? t('pages.login.sendVerification') : `${countDown}秒后可重发` }}
        </t-button>
      </t-form-item>
    </template>

    <t-form-item v-if="type !== 'qrcode'" class="btn-container">
      <t-button block size="large" type="submit" :loading="loading"> 登 录 </t-button>
    </t-form-item>

    <div class="switch-container">
      <div>
        <span v-if="type !== 'password'" class="tip" @click="switchType('password')">使用账号登录</span>
        <span v-if="type !== 'sms_code'" class="tip" @click="switchType('sms_code')">使用邮箱登录</span>
      </div>
      <div>
        <span class="tip" @click="toApiKey"><secured-icon /> <span class="tip-text">API KEY用量查询</span></span>
      </div>
    </div>
  </t-form>
</template>
<script setup lang="ts">
import { SecuredIcon } from 'tdesign-icons-vue-next';
import type { FormInstanceFunctions, FormRule, SubmitContext } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import type { LoginRequest } from '@/api/user';
import { login, sendVerificationCode } from '@/api/user';
import { useCounter } from '@/hooks';
import { t } from '@/locales';
import { useUserStore } from '@/store';

const userStore = useUserStore();

const INITIAL_DATA = {
  email: '',
  account: '',
  password: '',
  verifyCode: '',
  checked: false,
};

const FORM_RULES: Record<string, FormRule[]> = {
  email: [
    { required: true, message: '请输入邮箱', type: 'error' },
    { email: true, message: '请输入正确的邮箱格式', type: 'warning' },
  ],
  account: [{ required: true, message: t('pages.login.required.account'), type: 'error' }],
  password: [{ required: true, message: t('pages.login.required.password'), type: 'error' }],
  verifyCode: [{ required: true, message: '请输入验证码', type: 'error' }],
};

// 登录类型
const type = ref('password');

const form = ref<FormInstanceFunctions>();
const formData = ref({ ...INITIAL_DATA });
const showPsw = ref(false);
const loading = ref(false);

const [countDown, handleCounter] = useCounter();

const switchType = (val: string) => {
  type.value = val;
};

const router = useRouter();
const route = useRoute();

const toApiKey = () => {
  router.push('/stats/api-key');
};

/**
 * 发送验证码
 */
const sendCode = async () => {
  const validateResult = await form.value.validate({ fields: ['email'] });
  if (validateResult === true) {
    try {
      await sendVerificationCode({
        email: formData.value.email,
        type: 'login',
      });
      handleCounter();
      MessagePlugin.success('验证码已发送到您的邮箱');
    } catch (error: any) {
      MessagePlugin.error(error.message || '发送验证码失败');
    }
  }
};

const onSubmit = async (ctx: SubmitContext) => {
  if (ctx.validateResult === true) {
    loading.value = true;
    try {
      const loginData: LoginRequest = {
        login_type: type.value as 'password' | 'sms_code',
      };

      if (type.value === 'password') {
        // 账号密码登录
        loginData.username = formData.value.account;
        loginData.password = formData.value.password;
      } else {
        // 邮箱验证码登录
        loginData.email = formData.value.email;
        loginData.verification_code = formData.value.verifyCode;
      }

      const result = await login(loginData);

      // 存储登录信息到store
      await userStore.setUserInfo(result);

      MessagePlugin.success('登录成功');
      const redirect = route.query.redirect as string;
      const redirectUrl = redirect ? decodeURIComponent(redirect) : '/dashboard';
      router.push(redirectUrl);
    } catch (error: any) {
      console.error('登录失败:', error);
      MessagePlugin.error(error.message || '登录失败，请重试');
    } finally {
      loading.value = false;
    }
  }
};
</script>
<style lang="less" scoped>
@import '../index.less';
</style>
