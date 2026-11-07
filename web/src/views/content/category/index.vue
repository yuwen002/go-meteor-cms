<template>
  <div class="category-page">
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
        <el-button v-if="canAdd" type="primary" @click="openCreate">新增分类</el-button>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="name" label="分类名称" min-width="150" />
        <el-table-column prop="slug" label="别名" min-width="150" />
        <el-table-column prop="parent_name" label="上级分类" min-width="120" />
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
        <template #empty><el-empty description="暂无分类" :image-size="60" /></template>
      </el-table>
    </el-card>

    <el-dialog v-model="visible" :title="isEdit ? '编辑分类' : '新增分类'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="别名">
          <el-input v-model="form.slug" placeholder="留空自动生成，仅小写字母 / 数字 / 连字符" />
        </el-form-item>
        <el-form-item label="上级分类">
          <el-select v-model="form.parent_id" clearable placeholder="顶级分类" style="width: 100%">
            <el-option v-for="c in parentOptions" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
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
  createCategory,
  deleteCategory,
  getCategories,
  updateCategory,
  type CategoryItem,
  type CategoryPayload,
} from '@/api/content'

const userStore = useUserStore()
const canAdd = computed(() => userStore.hasPermission('cms:category:add'))
const canEdit = computed(() => userStore.hasPermission('cms:category:edit'))
const canDelete = computed(() => userStore.hasPermission('cms:category:delete'))

const list = ref<CategoryItem[]>([])
const loading = ref(false)
const keyword = ref('')
const visible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const editingId = ref(0)

interface CategoryForm {
  name: string
  slug: string
  parent_id: number | undefined
  description: string
  sort: number
  is_active: boolean
}
const emptyForm = (): CategoryForm => ({
  name: '',
  slug: '',
  parent_id: undefined,
  description: '',
  sort: 0,
  is_active: true,
})
const form = reactive<CategoryForm>(emptyForm())

const parentOptions = computed(() => list.value.filter((c) => c.id !== editingId.value))
const asRow = (row: unknown) => row as CategoryItem

async function loadList() {
  loading.value = true
  try {
    const res = await getCategories({ keyword: keyword.value || undefined })
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

function openEdit(row: CategoryItem) {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    slug: row.slug,
    parent_id: row.parent_id || undefined,
    description: row.description,
    sort: row.sort,
    is_active: row.is_active,
  } as CategoryForm)
  visible.value = true
}

async function submit() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分类名称')
    return
  }
  const payload: CategoryPayload = {
    name: form.name.trim(),
    slug: form.slug.trim() || undefined,
    parent_id: Number(form.parent_id ?? 0),
    description: form.description,
    sort: form.sort,
    is_active: form.is_active,
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateCategory(editingId.value, payload)
    } else {
      await createCategory(payload)
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

async function onDelete(row: CategoryItem) {
  try {
    await ElMessageBox.confirm(`确定删除分类「${row.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteCategory(row.id)
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
