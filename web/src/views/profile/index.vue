<template>
  <el-card class="profile-card">
    <el-form :model="form" label-width="90px" style="max-width: 480px">
      <el-form-item label="头像">
        <el-avatar :size="64" :src="avatarPreview || ''">{{ form.nickname?.charAt(0) || 'U' }}</el-avatar>
        <el-upload
          class="avatar-uploader"
          :show-file-list="false"
          :before-upload="beforeAvatarUpload"
          :http-request="handleUpload"
          accept="image/*"
        >
          <el-button size="small" type="primary" style="margin-left: 12px">上传新头像</el-button>
        </el-upload>
      </el-form-item>
      <el-form-item label="用户名">
        <el-input v-model="form.username" disabled />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="form.phone" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { UploadRequestOptions } from 'element-plus'
import { getProfile, updateProfile } from '@/api/user'
import { uploadFile } from '@/api/upload'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const saving = ref(false)
const avatarPreview = ref('')
const form = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
})

async function loadProfile() {
  const res = await getProfile()
  Object.assign(form, {
    username: res.username,
    nickname: res.nickname,
    email: res.email,
    phone: res.phone,
    avatar: res.avatar,
  })
  avatarPreview.value = res.avatar
}

function beforeAvatarUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片')
    return false
  }
  if (file.size / 1024 / 1024 > 5) {
    ElMessage.error('图片不能超过 5MB')
    return false
  }
  return true
}

async function handleUpload(options: UploadRequestOptions) {
  try {
    const res = await uploadFile(options.file)
    form.avatar = res.url
    avatarPreview.value = res.url
    ElMessage.success('头像已上传，记得点保存')
  } catch {
    // 错误已在拦截器提示
  }
}

async function save() {
  saving.value = true
  try {
    await updateProfile({
      nickname: form.nickname,
      email: form.email,
      phone: form.phone,
      avatar: form.avatar,
    })
    userStore.setUserInfo({
      ...(userStore.userInfo as object),
      nickname: form.nickname,
      email: form.email,
      phone: form.phone,
      avatar: form.avatar,
    } as never)
    ElMessage.success('保存成功')
  } finally {
    saving.value = false
  }
}

onMounted(loadProfile)
</script>

<style scoped>
.profile-card {
  max-width: 600px;
}
.avatar-uploader {
  display: inline-block;
  margin-left: 12px;
}
</style>
