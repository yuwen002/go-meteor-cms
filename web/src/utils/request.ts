import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { getToken, removeToken } from './auth'

// 后端统一响应结构：{ code, message, data }，code=0 为成功。
const service: AxiosInstance = axios.create({
  baseURL: '',
  timeout: 10000,
})

service.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code === 0) {
      return res.data
    }
    ElMessage.error(res.message || '请求失败')
    // 未登录 / 登录失效，清理并跳登录
    if (res.code === 2001) {
      removeToken()
      router.replace('/login')
    }
    return Promise.reject(new Error(res.message || 'Error'))
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  },
)

export default service
