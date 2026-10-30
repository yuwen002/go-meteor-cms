<template>
  <div class="dept-page">
    <el-card shadow="never">
      <div class="toolbar">
        <el-button type="primary" @click="openCreate(null)" v-if="canAdd">新增根部门</el-button>
        <el-button :disabled="!selected" @click="openCreate(selected)" v-if="canAdd">新增子部门</el-button>
        <el-button :disabled="!selected" @click="openEdit" v-if="canEdit">编辑</el-button>
        <el-button type="danger" :disabled="!selected" @click="handleDelete" v-if="canDelete">删除</el-button>
        <span class="tip">点击树节点选中后可编辑或新增子部门</span>
      </div>
      <el-tree :data="tree" :props="{ label: 'name', children: 'children' }" node-key="id"
        default-expand-all highlight-current @current-change="onSelect" class="dept-tree" />
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="父部门">
          <el-tree-select v-model="form.parent_id" :data="parentOptions" :props="{ label: 'name', children: 'children' }"
            node-key="id" check-strictly default-expand-all style="width: 100%" />
        </el-form-item>
        <el-form-item label="负责人ID"><el-input-number v-model="form.leader_id" :min="0" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.is_active" /></el-form-item>
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
const canAdd = computed(() => userStore.hasPermission('system:dept:add'))
const canEdit = computed(() => userStore.hasPermission('system:dept:edit'))
const canDelete = computed(() => userStore.hasPermission('system:dept:delete'))
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDeptTree, getDeptDetail, createDept, updateDept, deleteDept } from '@/api/dept'
import type { DeptTreeItem, DeptForm } from '@/api/dept'

const tree = ref<DeptTreeItem[]>([])
const selected = ref<DeptTreeItem | null>(null)

// el-tree-select 需要一棵含「根部门」选项的树
const parentOptions = computed(() => [{ id: 0, name: '（根部门）', children: tree.value }])

async function loadTree() {
  tree.value = await getDeptTree()
}
function onSelect(node: DeptTreeItem | null) {
  selected.value = node
}

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const editId = ref<number>(0)
const submitting = ref(false)
const emptyForm: DeptForm = { name: '', parent_id: 0, sort: 0, is_active: true, leader_id: 0 }
const form = reactive<DeptForm>({ ...emptyForm })

function openCreate(parent: DeptTreeItem | null) {
  isEdit.value = false
  dialogTitle.value = parent ? `新增子部门（父：${parent.name}）` : '新增根部门'
  Object.assign(form, { ...emptyForm, parent_id: parent ? parent.id : 0 })
  editId.value = 0
  dialogVisible.value = true
}
async function openEdit() {
  if (!selected.value) return
  const detail = await getDeptDetail(selected.value.id)
  isEdit.value = true
  dialogTitle.value = '编辑部门'
  Object.assign(form, {
    name: detail.name, parent_id: detail.parent_id, sort: detail.sort,
    is_active: detail.is_active, leader_id: detail.leader_id,
  })
  editId.value = detail.id
  dialogVisible.value = true
}
async function submit() {
  if (!form.name) { ElMessage.warning('请填写名称'); return }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateDept(editId.value, { ...form })
      ElMessage.success('更新成功')
    } else {
      await createDept({ ...form })
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
    await ElMessageBox.confirm(`确认删除部门「${selected.value.name}」？`, '提示', { type: 'warning' })
  } catch { return }
  await deleteDept(selected.value.id)
  ElMessage.success('删除成功')
  selected.value = null
  await loadTree()
}

onMounted(loadTree)
</script>

<style scoped>
.toolbar { margin-bottom: 12px; }
.tip { color: #909399; font-size: 12px; margin-left: 8px; }
.dept-tree { margin-top: 8px; }
</style>
