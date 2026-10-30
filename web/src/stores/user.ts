import { defineStore } from 'pinia'
import { getToken, setToken, removeToken, getRefreshToken, setRefreshToken, removeRefreshToken } from '@/utils/auth'

export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  isSuper: boolean
  isActive: boolean
}

export interface PermissionItem {
  id: number
  name: string
  permission: string
  type: number
  parentName: string
  apiPath: string
  method: string
  sort: number
  children?: PermissionItem[]
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    refreshToken: getRefreshToken(),
    userInfo: null as UserInfo | null,
    permissions: [] as string[],
    permissionTree: [] as PermissionItem[],
  }),
  getters: {
    hasPermission: (state) => (code: string) => state.permissions.includes(code),
  },
  actions: {
    setTokenAction(token: string) {
      this.token = token
      setToken(token)
    },
    setRefreshTokenAction(token: string) {
      this.refreshToken = token
      setRefreshToken(token)
    },
    setUserInfo(info: UserInfo) {
      this.userInfo = info
    },
    setPermissions(codes: string[], tree: PermissionItem[]) {
      this.permissions = codes
      this.permissionTree = tree
    },
    logoutAction() {
      this.token = ''
      this.refreshToken = ''
      this.userInfo = null
      this.permissions = []
      this.permissionTree = []
      removeToken()
      removeRefreshToken()
    },
  },
})
