import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';


export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, '.', '');
    return {
      cacheDir: '/tmp/dmh-admin-vite-cache',
      server: {
        port: 3000,
        host: '0.0.0.0',
        proxy: {
          '/api': {
            target: 'http://localhost:8889',
            changeOrigin: true,
          },
        },
      },
      plugins: [vue()],
      define: {
        'process.env.API_KEY': JSON.stringify(env.GEMINI_API_KEY),
        'process.env.GEMINI_API_KEY': JSON.stringify(env.GEMINI_API_KEY),
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false
      },
      resolve: {
        alias: {
          '@': path.resolve(__dirname, '.'),
          vue: 'vue/dist/vue.esm-bundler.js',
        }
      }
    };
});
