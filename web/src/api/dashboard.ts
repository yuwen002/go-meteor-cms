import request from '@/utils/request'

export interface DashboardStats {
  admin_count: number
  dept_count: number
  role_count: number
  permission_count: number
  today_login_count: number
  today_fail_login: number
  week_op_count: number
}

export function getDashboardStats() {
  return request({ url: '/admin/dashboard/stats', method: 'get' }) as Promise<DashboardStats>
}
