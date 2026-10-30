import request from '@/utils/request'

export interface LoginData {
  username: string
  password: string
  captcha_id: string
  captcha: string
}

export interface LoginResult {
  access_token: string
  refresh_token: string
}

export function login(data: LoginData): Promise<LoginResult> {
  return request({ url: '/admin/login', method: 'post', data }) as Promise<LoginResult>
}

export function getCaptcha() {
  return request({ url: '/api/captcha', method: 'get' })
}

export function logout() {
  return request({ url: '/admin/logout', method: 'post' })
}
