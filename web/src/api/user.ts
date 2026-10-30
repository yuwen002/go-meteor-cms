import request from '@/utils/request'

export function getProfile() {
  return request({ url: '/admin/profile', method: 'get' })
}

export function updateProfile(data: {
  nickname?: string
  email?: string
  phone?: string
  avatar?: string
}) {
  return request({ url: '/admin/profile', method: 'put', data })
}

export function getAdminList(params: {
  page?: number
  page_size?: number
  keyword?: string
  active?: number
}) {
  return request({ url: '/admin/admin-users', method: 'get', params })
}

export function getAdminDetail(id: number) {
  return request({ url: `/admin/admin-users/${id}`, method: 'get' })
}

export function createAdmin(data: Record<string, unknown>) {
  return request({ url: '/admin/admin-users/create', method: 'post', data })
}

export function updateAdmin(id: number, data: Record<string, unknown>) {
  return request({ url: `/admin/admin-users/${id}`, method: 'put', data })
}

export function deleteAdmin(id: number) {
  return request({ url: `/admin/admin-users/${id}`, method: 'delete' })
}

export function disableAdmin(id: number) {
  return request({ url: `/admin/admin-users/${id}/disable`, method: 'put' })
}

export function enableAdmin(id: number) {
  return request({ url: `/admin/admin-users/${id}/enable`, method: 'put' })
}

export function resetAdminPassword(id: number, newPassword: string) {
  return request({
    url: '/admin/admin-users/reset-password',
    method: 'put',
    data: { id, new_password: newPassword },
  })
}

export function getAdminRoles(id: number) {
  return request({ url: `/admin/admin-users/${id}/roles`, method: 'get' })
}

export function assignAdminRoles(id: number, roleIds: number[]) {
  return request({
    url: `/admin/admin-users/${id}/roles`,
    method: 'put',
    data: { role_ids: roleIds },
  })
}

// 登录后拉取「当前用户拥有的权限树」，用于菜单/按钮级权限控制
export function getMyPermissions() {
  return request({ url: '/admin/profile/permissions', method: 'get' })
}
