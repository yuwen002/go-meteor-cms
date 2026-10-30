<template>
  <div class="perm-page">
    <el-card shadow="never">
      <div class="toolbar">
        <el-button type="primary" @click="openCreate(null)" v-if="canAdd">新增根权限</el-button>
        <el-button :disabled="!selected" @click="openCreate(selected)" v-if="canAdd">新增子权限</el-button>
        <el-button :disabled="!selected" @click="openEdit" v-if="canEdit">编辑</el-button>
        <el-button type="danger" :disabled="!selected" @click="handleDelete" v-if="canDelete">删除</el-button>
        <span class="tip">点击树节点选中后可编辑或新增子权限</span>
      </div>
      <el-tree :data="tree" :props="{ label: 'name', children: 'children' }" node-key="id"
        default-expand-all highlight-current @current-change="onSelect" class="perm-tree">
        <template #default="{ data }">
          <span class="tree-node">
            <span>{{ data.name }}</span>
            <el-tag size="small" :type="typeTag(data.type)">{{ typeLabel(data.type) }}</el-tag>
            <span class="perm-code">{{ data.permission }}</span>
          </span>
        </template>
      </el-tree>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form :model="form" label-width="92px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="类型" required>
          <el-radio-group v-model="form.type">
            <el-radio :value="1">菜单</el-radio>
            <el-radio :value="2">按钮</el-radio>
            <el-radio :value="3">API</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="父权限">
          <el-tree-select v-model="form.parent_id" :data="parentOptions" :props="{ label: 'name', children: 'children' }"
            node-key="id" check-strictly default-expand-all style="width: 100%" />
        </el-form-item>
        <el-form-item label="路由路径" v-if="form.type === 1"><el-input v-model="form.path" placeholder="前端路由 path" /></el-form-item>
        <el-form-item label="组件" v-if="form.type === 1"><el-input v-model="form.component" placeholder="如 views/xxx/index.vue" /></el-form-item>
        <el-form-item label="图标" v-if="form.type === 1"><el-input v-model="form.icon" /></el-form-item>
        <el-form-item label="请求方法" v-if="form.type === 3"><el-input v-model="form.method" placeholder="GET/POST..." /></el-form-item>
        <el-form-item label="API路径" v-if="form.type === 3"><el-input v-model="form.api_path" placeholder="如 /admin/roles" /></el-form-item>
        <el-form-item label="权限标识"><el-input v-model="form.permission" placeholder="如 system:role:list" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.is_active" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'

// 按钮级权限：种子中 Type=2 按钮与 Type=3 接口使用同一标识，授权一次即可同时生效
const userStore = useUserStore()
const canAdd = computed(() => userStore.hasPermission('system:permission:add'))
const canEdit = computed(() => userStore.hasPermission('system:permission:edit'))
const canDelete = computed(() => userStore.hasPermission('system:permission:delete'))
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPermissionTree, getPermissionDetail, createPermission, updatePermission, deletePermission } from '@/api/permission'
import type { PermissionTreeItem, PermissionForm } from '@/api/permission'

const tree = ref<PermissionTreeItem[]>([])
const selected = ref<PermissionTreeItem | null>(null)

const parentOptions = computed(() => [{ id: 0, name: '（根权限）', children: tree.value }])

function typeLabel(v: number) {
  return { 1: '菜单', 2: '按钮', 3: 'API' }[v as 1 | 2 | 3] || `未知(${v})`
}
function typeTag(v: number): 'success' | 'warning' | 'info' {
  return v === 1 ? 'success' : v === 2 ? 'warning' : 'info'
}

async function loadTree() {
  tree.value = await getPermissionTree()
}
function onSelect(node: PermissionTreeItem | null) {
  selected.value = node
}

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const editId = ref<number>(0)
const submitting = ref(false)
const emptyForm: PermissionForm = {
  name: '', parent_id: 0, type: 1, path: '', component: '', icon: '',
  method: '', api_path: '', permission: '', is_active: true, sort: 0,
}
const form = reactive<PermissionForm>({ ...emptyForm })

function openCreate(parent: PermissionTreeItem | null) {
  isEdit.value = false
  dialogTitle.value = parent ? `新增子权限（父：${parent.name}）` : '新增根权限'
  Object.assign(form, { ...emptyForm, parent_id: parent ? parent.id : 0 })
  editId.value = 0
  dialogVisible.value = true
}
async function openEdit() {
  if (!selected.value) return
  const detail = await getPermissionDetail(selected.value.id)
  isEdit.value = true
  dialogTitle.value = '编辑权限'
  Object.assign(form, {
    name: detail.name, parent_id: detail.parent_id, type: detail.type,
    path: detail.path, component: detail.component, icon: detail.icon,
    method: detail.method, api_path: detail.api_path, permission: detail.permission,
    is_active: detail.is_active, sort: detail.sort,
  })
  editId.value = detail.id
  dialogVisible.value = true
}
async function submit() {
  if (!form.name) { ElMessage.warning('请填写名称'); return }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updatePermission(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createPermission({ ...form })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await loadTree()
  } finally {
    submitting.value = false
  }
}
async function handleDelete() {
  if (!selected.value) return
  try {
    await ElMessageBox.confirm(`确认删除权限「${selected.value.name}」？`, '提示', { type: 'warning' })
  } catch { return }
  await deletePermission(selected.value.id)
  ElMessage.success('删除成功')
  selected.value = null
  await loadTree()
}

onMounted(loadTree)
</script>

<style scoped>
.toolbar { margin-bottom: 12px; }
.tip { color: #909399; font-size: 12px; margin-left: 8px; }
.perm-tree { margin-top: 8px; }
.tree-node { display: inline-flex; align-items: center; gap: 8px; }
.perm-code { color: #909399; font-size: 12px; }
</style>
