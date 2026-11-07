import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { fileURLToPath, URL } from 'node:url'

// 开发期把后端接口与静态资源代理到 Go 服务（默认 8888 端口）。
// /admin 与 /api 为业务接口，/uploads 为上传文件静态目录。
export default defineConfig({
  plugins: [
    vue(),
    // 组件与样式按需解析（含 v-loading 指令），不再整包引入 dist/index.css。
    // 服务式 API（ElMessage / ElMessageBox）走 src/utils/element.ts 的深路径入口，
    // 菜单图标在 MenuItem.vue 里按映射表精确引入。
    Components({
      resolvers: [ElementPlusResolver({ importStyle: 'css', directives: true })],
      dts: 'src/components.d.ts',
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    // 把稳定的第三方依赖拆成独立 chunk：业务代码改动时不会让整包缓存失效，
    // 同时避免单个入口 chunk 过大拖慢首屏。图标单独拆出以便浏览器并行下载。
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) return
          if (id.includes('@element-plus/icons-vue')) return 'element-icons'
          if (id.includes('element-plus')) return 'element-plus'
          if (
            id.includes('/vue/') ||
            id.includes('/@vue/') ||
            id.includes('vue-router') ||
            id.includes('pinia')
          ) {
            return 'vue-vendor'
          }
          return 'vendor'
        },
      },
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/admin': { target: 'http://localhost:8888', changeOrigin: true },
      '/api': { target: 'http://localhost:8888', changeOrigin: true },
      '/uploads': { target: 'http://localhost:8888', changeOrigin: true },
    },
  },
})
