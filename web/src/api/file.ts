import request from '@/utils/request'

export interface FileItem {
  name: string
  path: string
  url: string
  size: number
  modified: number
  is_dir: boolean
}

export interface FileListResp {
  list: FileItem[]
  total: number
  current_dir: string
}

export function getFileList(params: {
  dir?: string
  keyword?: string
  page?: number
  page_size?: number
}) {
  return request({ url: '/admin/files', method: 'get', params }) as Promise<FileListResp>
}

export function deleteFile(path: string) {
  return request({ url: '/admin/files/delete', method: 'post', data: { path } })
}

export function renameFile(path: string, new_name: string) {
  return request({ url: '/admin/files/rename', method: 'put', data: { path, new_name } })
}
