import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'
import Layout from '@/layout/index.vue'
import { useUserStore } from '@/stores/user'
import { getProfile, getMyPermissions } from '@/api/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', hidden: true },
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', permission: 'system:dashboard:view' },
      },
      {
        path: 'system/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/user/index.vue'),
        meta: { title: '管理员管理', permission: 'system:admin_user:list' },
      },
      {
        path: 'system/roles',
        name: 'AdminRoles',
        component: () => import('@/views/admin/role/index.vue'),
        meta: { title: '角色管理', permission: 'system:role:list' },
      },
      {
        path: 'system/depts',
        name: 'AdminDepts',
        component: () => import('@/views/admin/dept/index.vue'),
        meta: { title: '部门管理', permission: 'system:dept:list' },
      },
      {
        path: 'system/permissions',
        name: 'AdminPermissions',
        component: () => import('@/views/admin/permission/index.vue'),
        meta: { title: '权限管理', permission: 'system:permission:list' },
      },
      {
        path: 'system/files',
        name: 'FileManage',
        component: () => import('@/views/system/file/index.vue'),
        meta: { title: '文件管理', permission: 'system:file:list' },
      },
      {
        path: 'system/operation-log',
        name: 'OperationLog',
        component: () => import('@/views/system/log/operation.vue'),
        meta: { title: '操作日志', permission: 'system:log:operation_list' },
      },
      {
        path: 'system/login-log',
        name: 'LoginLog',
        component: () => import('@/views/system/log/login.vue'),
        meta: { title: '登录日志', permission: 'system:log:login_list' },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const whiteList = ['/login']

router.beforeEach(async (to, from, next) => {
  const token = getToken()
  if (!token) {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next('/login')
    }
    return
  }
  if (to.path === '/login') {
    next('/')
    return
  }

  const userStore = useUserStore()
  if (!userStore.accessLoaded) {
    try {
      const profile = await getProfile()
      // 注意：该接口返回的是对象 {is_super, menus, permissions}，此前被误当作数组使用
      const access = await getMyPermissions()
      userStore.setUserInfo(profile)
      userStore.setAccess(access.permissions ?? [], access.menus ?? [])
      next({ ...to, replace: true })
    } catch {
      userStore.logoutAction()
      next('/login')
    }
    return
  }

  // 路由级权限守卫：防止直接在地址栏输入未授权路径。
  // 无权限时回退到「第一个有权访问的菜单」而非固定首页，避免用户对首页也无权限时形成重定向死循环。
  const required = to.meta.permission as string | undefined
  if (required && !hasAccess(userStore, required)) {
    const fallback = firstAccessiblePath(userStore)
    if (fallback && fallback !== to.path) {
      next(fallback)
      return
    }
  }
  next()
})

// 是否有访问权限：超级管理员天然放行（同时校验后端下发的 is_super，避免权限数组异常时被误拦）
function hasAccess(userStore: ReturnType<typeof useUserStore>, code: string): boolean {
  return userStore.userInfo?.isSuper === true || userStore.hasPermission(code)
}

// 第一个可访问的菜单路径；后端下发的菜单已按权限过滤，取第一个即可。
// 一个菜单都没有时兜底到个人中心（无需权限），避免重定向死循环。
function firstAccessiblePath(userStore: ReturnType<typeof useUserStore>): string {
  return userStore.menus[0]?.path || '/profile'
}

export default router
