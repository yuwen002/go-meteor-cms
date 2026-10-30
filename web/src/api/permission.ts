import request from '@/utils/request'

export const PERMISSION_TYPE = { MENU: 1, BUTTON: 2, API: 3 } as const

export interface PermissionDetailItem {
  id: number
  name: string
  parent_id: number
  type: number
  path: string
  component: string
  icon: string
  method: string
  api_path: string
  permission: string
  is_active: boolean
  sort: number
  created_at: string
  updated_at: string
}

export interface PermissionTreeItem {
  id: number
  name: string
  parent_id: number
  type: number
  path: string
  component: string
  icon: string
  method: string
  api_path: string
  permission: string
  is_active: boolean
  sort: number
  children: PermissionTreeItem[]
}

export interface PermissionListResp {
  total: number
  page: number
  page_size: number
  list: PermissionDetailItem[]
}

export interface PermissionForm {
  name: string
  parent_id: number
  type: number
  path: string
  component: string
  icon: string
  method: string
  api_path: string
  permission: string
  is_active: boolean
  sort: number
}

export function getPermissions(params: Record<string, unknown>) {
  return request({ url: '/admin/permissions', method: 'get', params }) as Promise<PermissionListResp>
}
export function getPermissionTree() {
  return request({ url: '/admin/permissions/tree', method: 'get' }) as Promise<PermissionTreeItem[]>
}
export function getPermissionDetail(id: number) {
  return request({ url: `/admin/permissions/${id}`, method: 'get' }) as Promise<PermissionDetailItem>
}
export function createPermission(data: PermissionForm) {
  return request({ url: '/admin/permissions', method: 'post', data })
}
export function updatePermission(id: number, data: Partial<PermissionForm>) {
  return request({ url: `/admin/permissions/${id}`, method: 'put', data })
}
export function deletePermission(id: number) {
  return request({ url: `/admin/permissions/${id}`, method: 'delete' })
}
