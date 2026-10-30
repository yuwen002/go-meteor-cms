<template>
  <el-menu
    :default-active="activeMenu"
    background-color="#001529"
    text-color="#fff"
    active-text-color="#409eff"
    router
    class="sidebar-menu"
  >
    <template v-for="item in menus" :key="item.path">
      <el-menu-item :index="item.path">
        <el-icon v-if="item.icon"><component :is="item.icon" /></el-icon>
        <span>{{ item.title }}</span>
      </el-menu-item>
    </template>
  </el-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

interface MenuEntry {
  path: string
  title: string
  permission?: string
  icon?: string
}

const route = useRoute()
const userStore = useUserStore()

const allMenus: MenuEntry[] = [
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

const menus = computed(() =>
  allMenus.filter((m) => !m.permission || userStore.hasPermission(m.permission)),
)

const activeMenu = computed(() => route.path)
</script>

<style scoped>
.sidebar-menu {
  border-right: none;
}
</style>
