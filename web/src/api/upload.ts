import request from '@/utils/request'

export interface UploadResult {
  url: string
  filename: string
  size: number
}

export function uploadFile(file: File): Promise<UploadResult> {
  const formData = new FormData()
  formData.append('file', file)
  return request({ url: '/admin/upload', method: 'post', data: formData })
}
