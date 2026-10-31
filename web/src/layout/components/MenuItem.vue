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

// 未在映射表中的 icon 键只在控制台提示一次，便于定位后端下发了哪些未配置的图标
const warnedKeys = new Set<string>()

const hasChildren = computed(() => !!props.item.children && props.item.children.length > 0)
const iconName = computed(() => {
  const key = props.item.icon || ''
  const name = iconMap[key]
  if (!name && !warnedKeys.has(key)) {
    warnedKeys.add(key)
    console.warn(`[menu] 未配置的菜单图标: "${key}"，已回落到 Menu（菜单：${props.item.name}）`)
  }
  return name || 'Menu'
})
</script>
