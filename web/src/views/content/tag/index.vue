<template>
  <div class="tag-page">
    <el-card shadow="never">
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="名称 / 别名"
          clearable
          style="width: 200px"
          @keyup.enter="loadList"
          @clear="loadList"
        />
        <el-button @click="loadList">搜索</el-button>
        <el-button v-if="canAdd" type="primary" @click="openCreate">新增标签</el-button>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="name" label="标签名称" min-width="160" />
        <el-table-column prop="slug" label="别名" min-width="170" />
        <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
        <el-table-column prop="article_count" label="文章数" width="90" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="asRow(row).is_active ? 'success' : 'info'" size="small">
              {{ asRow(row).is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button v-if="canEdit" link type="primary" @click="openEdit(asRow(row))">编辑</el-button>
            <el-button v-if="canDelete" link type="danger" @click="onDelete(asRow(row))">删除</el-button>
          </template>
        </el-table-column>
        <template #empty><el-empty description="暂无标签" :image-size="60" /></template>
      </el-table>
    </el-card>

    <el-dialog v-model="visible" :title="isEdit ? '编辑标签' : '新增标签'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="别名">
          <el-input v-model="form.slug" placeholder="留空自动生成，仅小写字母 / 数字 / 连字符" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.is_active" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from '@/utils/element'
import { useUserStore } from '@/stores/user'
import {
  createTag,
  deleteTag,
  getTags,
  updateTag,
  type TagItem,
  type TagPayload,
} from '@/api/content'

const userStore = useUserStore()
const canAdd = computed(() => userStore.hasPermission('cms:tag:add'))
const canEdit = computed(() => userStore.hasPermission('cms:tag:edit'))
const canDelete = computed(() => userStore.hasPermission('cms:tag:delete'))

const list = ref<TagItem[]>([])
const loading = ref(false)
const keyword = ref('')
const visible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const editingId = ref(0)

interface TagForm {
  name: string
  slug: string
  description: string
  sort: number
  is_active: boolean
}
const emptyForm = (): TagForm => ({
  name: '',
  slug: '',
  description: '',
  sort: 0,
  is_active: true,
})
const form = reactive<TagForm>(emptyForm())
const asRow = (row: unknown) => row as TagItem

async function loadList() {
  loading.value = true
  try {
    const res = await getTags({ keyword: keyword.value || undefined })
    list.value = res.list ?? []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  isEdit.value = false
  editingId.value = 0
  Object.assign(form, emptyForm())
  visible.value = true
}

function openEdit(row: TagItem) {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    slug: row.slug,
    description: row.description,
    sort: row.sort,
    is_active: row.is_active,
  } as TagForm)
  visible.value = true
}

async function submit() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写标签名称')
    return
  }
  const payload: TagPayload = {
    name: form.name.trim(),
    slug: form.slug.trim() || undefined,
    description: form.description,
    sort: form.sort,
    is_active: form.is_active,
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateTag(editingId.value, payload)
    } else {
      await createTag(payload)
    }
    ElMessage.success(isEdit.value ? '已更新' : '已创建')
    visible.value = false
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  } finally {
    submitting.value = false
  }
}

async function onDelete(row: TagItem) {
  try {
    await ElMessageBox.confirm(`确定删除标签「${row.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteTag(row.id)
    ElMessage.success('已删除')
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  }
}

onMounted(loadList)
</script>

<style scoped>
.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
}
</style>
