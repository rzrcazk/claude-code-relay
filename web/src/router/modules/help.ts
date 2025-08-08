import { HelpCircleIcon } from 'tdesign-icons-vue-next';
import { shallowRef } from 'vue';

import Layout from '@/layouts/index.vue';

export default [
  {
    path: '/help',
    component: Layout,
    redirect: '/help/index',
    name: 'help',
    meta: {
      title: '帮助文档',
      icon: shallowRef(HelpCircleIcon),
      orderNo: 999,
      hidden: true,
    },
    children: [
      {
        path: 'index',
        name: 'HelpIndex',
        component: () => import('@/pages/help/index.vue'),
        meta: {
          title: '使用指南',
        },
      },
    ],
  },
];
