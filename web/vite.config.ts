import path from 'node:path';

import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import type { UserConfig } from 'vite';
import { viteMockServe } from 'vite-plugin-mock';
import svgLoader from 'vite-svg-loader';

import packageJson from './package.json';

// https://vitejs.dev/config/
export default (): UserConfig => {
  return {
    base: '/',
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },

    define: {
      __APP_VERSION__: JSON.stringify(packageJson.version),
    },

    css: {
      preprocessorOptions: {
        less: {
          modifyVars: {
            hack: `true; @import (reference) "${path.resolve('src/style/variables.less')}";`,
          },
          math: 'strict',
          javascriptEnabled: true,
        },
      },
    },

    plugins: [
      vue(),
      vueJsx(),
      viteMockServe({
        mockPath: 'mock',
        enable: true,
      }),
      svgLoader(),
    ],

    server: {
      port: 3001,
      host: '0.0.0.0',
      allowedHosts: true,
    },
  };
};
