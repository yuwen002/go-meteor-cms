import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// 开发期把后端接口与静态资源代理到 Go 服务（默认 8888 端口）。
// /admin 与 /api 为业务接口，/uploads 为上传文件静态目录。
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
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
