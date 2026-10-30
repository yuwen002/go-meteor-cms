import request from '@/utils/request'

export interface DeptTreeItem {
  id: number
  name: string
  sort: number
  is_active: boolean
  children: DeptTreeItem[]
}

export interface DeptDetail {
  id: number
  name: string
  parent_id: number
  sort: number
  is_active: boolean
  leader_id: number
}

export interface DeptForm {
  name: string
  parent_id: number
  sort: number
  is_active: boolean
  leader_id: number
}

export function getDeptTree() {
  return request({ url: '/admin/departments/tree', method: 'get' }) as Promise<DeptTreeItem[]>
}
export function getDeptDetail(id: number) {
  return request({ url: `/admin/departments/${id}`, method: 'get' }) as Promise<DeptDetail>
}
export function createDept(data: DeptForm) {
  return request({ url: '/admin/departments', method: 'post', data })
}
export function updateDept(id: number, data: DeptForm) {
  return request({ url: `/admin/departments/${id}`, method: 'put', data })
}
export function deleteDept(id: number) {
  return request({ url: `/admin/departments/${id}`, method: 'delete' })
}
