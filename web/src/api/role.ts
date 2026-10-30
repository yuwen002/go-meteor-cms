import request from '@/utils/request'

export interface RoleItem {
  id: number
  name: string
  code: string
  desc: string
  data_scope: number
  is_system: boolean
  is_active: boolean
  sort: number
  created_at: string
  updated_at: string
}

export interface RoleListResp {
  total: number
  page: number
  page_size: number
  list: RoleItem[]
}

export interface RoleForm {
  name: string
  code: string
  desc: string
  data_scope: number
  is_active: boolean
  sort: number
}

export function getRoles(params: Record<string, unknown>) {
  return request({ url: '/admin/roles', method: 'get', params }) as Promise<RoleListResp>
}
export function getRoleDetail(id: number) {
  return request({ url: `/admin/roles/${id}`, method: 'get' })
}
export function createRole(data: RoleForm) {
  return request({ url: '/admin/roles', method: 'post', data })
}
export function updateRole(id: number, data: Partial<RoleForm>) {
  return request({ url: `/admin/roles/${id}`, method: 'put', data })
}
export function deleteRole(id: number) {
  return request({ url: `/admin/roles/${id}`, method: 'delete' })
}
export function getRolePermissions(id: number) {
  return request({ url: `/admin/roles/${id}/permissions`, method: 'get' }) as Promise<{ permission_ids: number[] }>
}
export function getRolePermissionsWithNames(id: number) {
  return request({
    url: `/admin/roles/${id}/permissions-with-names`,
    method: 'get',
  }) as Promise<{ permissions: { id: number; name: string }[] }>
}
export function assignRolePermissions(id: number, permission_ids: number[]) {
  return request({ url: `/admin/roles/${id}/permissions`, method: 'put', data: { permission_ids } })
}
