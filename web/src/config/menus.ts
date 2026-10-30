export interface MenuEntry {
  path: string
  title: string
  permission?: string
  icon?: string
}

// 侧边栏菜单与路由权限守卫共用同一份定义，避免两处不一致导致
// 「菜单看得见但路由被拦」或「能进页面但菜单不显示」。
//
// 约定：
//   - permission 为空表示无需权限，所有登录用户可见（如个人中心）；
//   - permission 取值必须与权限菜单种子的 Type=1 项一致。
export const allMenus: MenuEntry[] = [
  { path: '/dashboard', title: '仪表盘', permission: 'system:dashboard:view', icon: 'Odometer' },
  { path: '/system/users', title: '管理员管理', permission: 'system:admin_user:list', icon: 'User' },
  { path: '/system/roles', title: '角色管理', permission: 'system:role:list', icon: 'Avatar' },
  { path: '/system/depts', title: '部门管理', permission: 'system:dept:list', icon: 'OfficeBuilding' },
  { path: '/system/permissions', title: '权限管理', permission: 'system:permission:list', icon: 'Lock' },
  { path: '/system/files', title: '文件管理', permission: 'system:file:list', icon: 'Folder' },
  { path: '/system/operation-log', title: '操作日志', permission: 'system:log:operation_list', icon: 'Document' },
  { path: '/system/login-log', title: '登录日志', permission: 'system:log:login_list', icon: 'Clock' },
  { path: '/profile', title: '个人中心', icon: 'Setting' },
]
