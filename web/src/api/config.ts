import request from '@/utils/request'

export interface ConfigItem {
  id: number
  config_key: string
  config_value: string
  name: string
  remark: string
  is_system: boolean
  created_at: string
}

export function getConfigs(params: { keyword?: string }) {
  return request({ url: '/admin/configs', method: 'get', params }) as Promise<{
    list: ConfigItem[]
  }>
}

export function createConfig(data: {
  config_key: string
  config_value?: string
  name: string
  remark?: string
}) {
  return request({ url: '/admin/configs', method: 'post', data })
}

export function updateConfig(
  id: number,
  data: { config_value?: string; name?: string; remark?: string },
) {
  return request({ url: `/admin/configs/${id}`, method: 'put', data })
}

export function deleteConfig(id: number) {
  return request({ url: `/admin/configs/${id}`, method: 'delete' })
}
