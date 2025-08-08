<template>
  <t-form
    ref="form"
    class="item-container register-email"
    :data="formData"
    :rules="FORM_RULES"
    label-width="0"
    @submit="onSubmit"
  >
    <t-form-item name="username">
      <t-input v-model="formData.username" size="large" placeholder="请输入用户名">
        <template #prefix-icon>
          <t-icon name="user" />
        </template>
      </t-input>
    </t-form-item>

    <t-form-item name="email">
      <t-input v-model="formData.email" size="large" placeholder="请输入邮箱">
        <template #prefix-icon>
          <t-icon name="mail" />
        </template>
      </t-input>
    </t-form-item>

    <t-form-item name="password">
      <t-input
        v-model="formData.password"
        size="large"
        :type="showPsw ? 'text' : 'password'"
        clearable
        placeholder="请输入登录密码"
      >
        <template #prefix-icon>
          <t-icon name="lock-on" />
        </template>
        <template #suffix-icon>
          <t-icon :name="showPsw ? 'browse' : 'browse-off'" @click="showPsw = !showPsw" />
        </template>
      </t-input>
    </t-form-item>

    <t-form-item class="verification-code" name="verifyCode">
      <t-input v-model="formData.verifyCode" size="large" placeholder="请输入验证码" />
      <t-button variant="outline" :disabled="countDown > 0 || !formData.email" @click="sendCode">
        {{ countDown === 0 ? '发送验证码' : `${countDown}秒后可重发` }}
      </t-button>
    </t-form-item>

    <t-form-item class="check-container" name="checked">
      <t-checkbox v-model="formData.checked">我已阅读并同意 </t-checkbox> <span>服务协议</span> 和
      <span>隐私声明</span>
    </t-form-item>

    <t-form-item>
      <t-button block size="large" type="submit" :loading="loading"> 注册 </t-button>
    </t-form-item>
  </t-form>
</template>
<script setup lang="ts">
import type { FormRule, SubmitContext } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { ref } from 'vue';

import type { RegisterRequest } from '@/api/user';
import { register, sendVerificationCode } from '@/api/user';
import { useCounter } from '@/hooks';

const emit = defineEmits(['register-success']);

const INITIAL_DATA = {
  username: '',
  email: '',
  password: '',
  verifyCode: '',
  checked: false,
};

const FORM_RULES: Record<string, FormRule[]> = {
  username: [
    { required: true, message: '用户名必填', type: 'error' },
    { min: 2, message: '用户名至少2个字符', type: 'warning' },
    { max: 20, message: '用户名不能超过20个字符', type: 'warning' },
  ],
  email: [
    { required: true, message: '邮箱必填', type: 'error' },
    { email: true, message: '请输入正确的邮箱格式', type: 'warning' },
  ],
  password: [
    { required: true, message: '密码必填', type: 'error' },
    { min: 6, message: '密码至少6个字符', type: 'warning' },
  ],
  verifyCode: [{ required: true, message: '验证码必填', type: 'error' }],
  checked: [{ required: true, message: '请同意服务协议', type: 'error' }],
};

const form = ref();
const formData = ref({ ...INITIAL_DATA });
const showPsw = ref(false);
const loading = ref(false);

const [countDown, handleCounter] = useCounter();

/**
 * 发送验证码
 */
const sendCode = async () => {
  const validateResult = await form.value.validate({ fields: ['email'] });
  if (validateResult === true) {
    try {
      await sendVerificationCode({
        email: formData.value.email,
        type: 'register',
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
    if (!formData.value.checked) {
      MessagePlugin.error('请同意服务协议和隐私声明');
      return;
    }

    loading.value = true;
    try {
      const registerData: RegisterRequest = {
        username: formData.value.username,
        email: formData.value.email,
        password: formData.value.password,
        verification_code: formData.value.verifyCode,
      };

      await register(registerData);
      MessagePlugin.success('注册成功');
      emit('register-success');

      // 重置表单
      form.value.reset();
      formData.value = { ...INITIAL_DATA };
    } catch (error: any) {
      console.error('注册失败:', error);
      MessagePlugin.error(error.message || '注册失败，请重试');
    } finally {
      loading.value = false;
    }
  }
};
</script>
<style lang="less" scoped>
@import '../index.less';
</style>
