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
import { allMenus } from '@/config/menus'

const route = useRoute()
const userStore = useUserStore()

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
