<template>
  <div class="dashboard">
    <el-row :gutter="16">
      <el-col :span="6" v-for="card in cards" :key="card.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-label">{{ card.label }}</div>
          <div class="stat-value">{{ card.value }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="welcome-card" v-loading="loading">
      <template #header>欢迎</template>
      <p>
        你好，{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}（{{
          isSuper ? '超级管理员' : '普通管理员'
        }}）。
      </p>
      <p class="sub">账号状态：{{ userStore.userInfo?.isActive ? '已启用' : '已禁用' }}</p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { getDashboardStats } from '@/api/dashboard'
import type { DashboardStats } from '@/api/dashboard'

const userStore = useUserStore()
const isSuper = computed(() => userStore.userInfo?.isSuper || false)
const loading = ref(false)
const stats = ref<DashboardStats | null>(null)

const cards = computed(() => [
  { label: '管理员总数', value: stats.value?.admin_count ?? '-' },
  { label: '部门总数', value: stats.value?.dept_count ?? '-' },
  { label: '角色总数', value: stats.value?.role_count ?? '-' },
  { label: '权限总数', value: stats.value?.permission_count ?? '-' },
  { label: '今日登录', value: stats.value?.today_login_count ?? '-' },
  { label: '今日失败登录', value: stats.value?.today_fail_login ?? '-' },
  { label: '近7天操作数', value: stats.value?.week_op_count ?? '-' },
])

onMounted(async () => {
  loading.value = true
  try {
    stats.value = await getDashboardStats()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.stat-card { margin-bottom: 16px; }
.stat-label { color: #909399; font-size: 13px; }
.stat-value { font-size: 26px; font-weight: 600; margin-top: 8px; }
.welcome-card { margin-top: 8px; }
.sub { color: #909399; font-size: 13px; margin-top: 4px; }
</style>
