import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'
import Layout from '@/layout/index.vue'
import { useUserStore } from '@/stores/user'
import { getProfile, getMyPermissions } from '@/api/user'
import type { PermissionItem } from '@/stores/user'

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
        meta: { title: '文件管理' },
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

// 从权限树递归提取当前用户拥有的权限码
function extractCodes(nodes: PermissionItem[]): string[] {
  const codes: string[] = []
  const walk = (list: PermissionItem[]) => {
    list.forEach((n) => {
      if (n.permission) codes.push(n.permission)
      if (n.children) walk(n.children)
    })
  }
  walk(nodes)
  return codes
}

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
  if (userStore.permissions.length === 0) {
    try {
      const profile = await getProfile()
      const tree = (await getMyPermissions()) as unknown as PermissionItem[]
      userStore.setUserInfo(profile)
      userStore.setPermissions(extractCodes(tree), tree)
      next({ ...to, replace: true })
    } catch {
      userStore.logoutAction()
      next('/login')
    }
  } else {
    next()
  }
})

export default router
