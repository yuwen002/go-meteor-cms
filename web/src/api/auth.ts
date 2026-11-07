import request from '@/utils/request'
import { getRefreshToken } from '@/utils/auth'

export interface LoginData {
  username: string
  password: string
  // 后端 LOGIN_CAPTCHA_ENABLED=false 时这两个字段可省略
  captcha_id?: string
  captcha?: string
}

export interface LoginResult {
  access_token: string
  refresh_token: string
}

export interface CaptchaResult {
  enabled: boolean
  captcha_id: string
  captcha_base64: string
}

export function login(data: LoginData): Promise<LoginResult> {
  return request({ url: '/admin/login', method: 'post', data }) as Promise<LoginResult>
}

export function getCaptcha(): Promise<CaptchaResult> {
  return request({ url: '/api/captcha', method: 'get' }) as Promise<CaptchaResult>
}

export function logout() {
  return request({
    url: '/admin/logout',
    method: 'post',
    headers: { 'X-Refresh-Token': getRefreshToken() },
  })
}
