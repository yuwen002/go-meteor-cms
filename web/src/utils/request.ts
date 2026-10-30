import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { getToken, getRefreshToken } from './auth'
import { useUserStore } from '@/stores/user'

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

// ---- 无感刷新：access token 失效（code=2001）时，用 refresh token 换新 token 并重试原请求 ----
let isRefreshing = false
let pendingQueue: Array<(token: string) => void> = []

function flushQueue(token: string): void {
  pendingQueue.forEach((cb) => cb(token))
  pendingQueue = []
}

function forceLogout(): void {
  useUserStore().logoutAction()
  router.replace('/login')
}

async function doRefresh(refreshToken: string): Promise<string> {
  const data = (await service({
    url: '/admin/refresh',
    method: 'post',
    data: { refresh_token: refreshToken },
  })) as { access_token: string; refresh_token: string }
  // 刷新令牌轮换：用服务端返回的新 refresh token 覆盖本地旧值
  useUserStore().setRefreshTokenAction(data.refresh_token)
  return data.access_token
}

function refreshAndRetry(config: InternalAxiosRequestConfig): Promise<unknown> {
  const refreshToken = getRefreshToken()
  // 刷新接口自身失败、或没有刷新令牌：直接登出
  if (!refreshToken || config.url?.includes('/admin/refresh')) {
    forceLogout()
    return Promise.reject(new Error('token refresh failed'))
  }
  if (!isRefreshing) {
    isRefreshing = true
    return doRefresh(refreshToken)
      .then((newAccess) => {
        useUserStore().setTokenAction(newAccess)
        flushQueue(newAccess)
        const retried = { ...config, __retry: true } as InternalAxiosRequestConfig & {
          __retry?: boolean
        }
        return service(retried)
      })
      .catch((e) => {
        flushQueue('')
        forceLogout()
        return Promise.reject(e)
      })
      .finally(() => {
        isRefreshing = false
      })
  }
  // 已有刷新在进行中，排队等待刷新完成后再用新 token 重试
  return new Promise((resolve, reject) => {
    pendingQueue.push((token: string) => {
      if (token) {
        const retried = { ...config, __retry: true } as InternalAxiosRequestConfig & {
          __retry?: boolean
        }
        resolve(service(retried))
      } else {
        reject(new Error('token refresh failed'))
      }
    })
  })
}

service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code === 0) {
      return res.data
    }
    // 登录失效 / access token 过期：尝试无感刷新（已重试过的请求不再重试，直接登出）
    if (res.code === 2001) {
      if ((response.config as InternalAxiosRequestConfig & { __retry?: boolean }).__retry) {
        forceLogout()
        return Promise.reject(new Error('token expired'))
      }
      return refreshAndRetry(response.config as InternalAxiosRequestConfig)
    }
    ElMessage.error(res.message || '请求失败')
    return Promise.reject(new Error(res.message || 'Error'))
  },
  (error) => {
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  },
)

// 响应拦截器已把响应解包为业务 data（即 {code,msg,data} 中的 data），运行时不再是 AxiosResponse。
// 这里显式对外承诺「返回业务数据类型」，避免调用方被推断成 AxiosResponse 而无法访问业务字段。
type RequestFunc = <T = any>(config: AxiosRequestConfig) => Promise<T>

export default service as unknown as RequestFunc
