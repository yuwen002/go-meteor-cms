import request from '@/utils/request'

export function getOperationLogs(params: Record<string, unknown>) {
  return request({ url: '/admin/logs/operations', method: 'get', params })
}

export function getLoginLogs(params: Record<string, unknown>) {
  return request({ url: '/admin/logs/logins', method: 'get', params })
}
