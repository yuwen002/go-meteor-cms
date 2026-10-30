<template>
  <el-sub-menu v-if="hasChildren" :index="String(item.id)">
    <template #title>
      <el-icon><component :is="iconName" /></el-icon>
      <span>{{ item.name }}</span>
    </template>
    <menu-item v-for="child in item.children" :key="child.id" :item="child" />
  </el-sub-menu>

  <el-menu-item v-else :index="item.path">
    <el-icon><component :is="iconName" /></el-icon>
    <template #title>{{ item.name }}</template>
  </el-menu-item>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { MyMenuItem } from '@/api/user'

const props = defineProps<{ item: MyMenuItem }>()

// 后端菜单 icon 存的是语义键（与具体图标库解耦），这里映射为 Element Plus
// 全局注册的图标组件名；未配置的键回落到通用 Menu 图标。
const iconMap: Record<string, string> = {
  dashboard: 'Odometer',
  user: 'User',
  dept: 'OfficeBuilding',
  role: 'Avatar',
  permission: 'Lock',
  document: 'Document',
  clock: 'Clock',
  file: 'Folder',
  dict: 'Collection',
  setting: 'Setting',
}

const hasChildren = computed(() => !!props.item.children && props.item.children.length > 0)
const iconName = computed(() => iconMap[props.item.icon] || 'Menu')
</script>
