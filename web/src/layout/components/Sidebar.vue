<template>
  <el-menu
    :default-active="activeMenu"
    background-color="#001529"
    text-color="#fff"
    active-text-color="#409eff"
    router
    class="sidebar-menu"
  >
    <!-- 动态菜单：由后端按当前用户权限下发，支持多级 -->
    <menu-item v-for="m in menus" :key="m.id" :item="m" />

    <!-- 个人中心不属于后台权限体系，固定展示 -->
    <el-menu-item index="/profile">
      <el-icon><Setting /></el-icon>
      <template #title>个人中心</template>
    </el-menu-item>
  </el-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import MenuItem from './MenuItem.vue'

const route = useRoute()
const userStore = useUserStore()

const menus = computed(() => userStore.menus)
const activeMenu = computed(() => route.path)
</script>

<style scoped>
.sidebar-menu {
  border-right: none;
}
</style>
