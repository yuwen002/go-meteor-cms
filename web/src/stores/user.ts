import { defineStore } from 'pinia'
import { getToken, setToken, removeToken, getRefreshToken, setRefreshToken, removeRefreshToken } from '@/utils/auth'
import type { MyMenuItem } from '@/api/user'

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



export const useUserStore = defineStore('user', {
  state: () => ({
    token: getToken(),
    refreshToken: getRefreshToken(),
    userInfo: null as UserInfo | null,
    permissions: [] as string[],
    // 动态菜单：由 /admin/profile/permissions 下发，已按当前用户权限过滤
    menus: [] as MyMenuItem[],
    // 是否已拉取过菜单与权限；零权限用户据此避免每次导航重复拉取
    accessLoaded: false,
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
    setAccess(codes: string[], menus: MyMenuItem[]) {
      this.permissions = codes
      this.menus = menus
      this.accessLoaded = true
    },
    logoutAction() {
      this.token = ''
      this.refreshToken = ''
      this.userInfo = null
      this.permissions = []
      this.menus = []
      this.accessLoaded = false
      removeToken()
      removeRefreshToken()
    },
  },
})
