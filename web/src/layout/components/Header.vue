<template>
  <div class="header-right">
    <el-dropdown @command="handleCommand">
      <span class="user-trigger">
        <el-avatar :size="32" :src="avatarUrl">{{ nicknameInitial }}</el-avatar>
        <span class="username">{{ nickname }}</span>
        <el-icon><ArrowDown /></el-icon>
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="profile">个人中心</el-dropdown-item>
          <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { logout as logoutApi } from '@/api/auth'
import { removeToken } from '@/utils/auth'

const router = useRouter()
const userStore = useUserStore()

const nickname = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '未登录')
const avatarUrl = computed(() => userStore.userInfo?.avatar || '')
const nicknameInitial = computed(() => (nickname.value ? nickname.value.charAt(0).toUpperCase() : 'U'))

function handleCommand(command: string) {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    logoutApi()
      .then(() => {
        removeToken()
        userStore.logoutAction()
        router.push('/login')
        ElMessage.success('已退出登录')
      })
      .catch(() => {
        removeToken()
        userStore.logoutAction()
        router.push('/login')
      })
  }
}
</script>

<style scoped>
.header-right {
  display: flex;
  align-items: center;
}
.user-trigger {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}
.username {
  font-size: 14px;
}
</style>
