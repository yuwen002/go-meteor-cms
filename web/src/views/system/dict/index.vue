<template>
  <div class="dict-page">
    <el-row :gutter="12">
      <el-col :span="10">
        <el-card shadow="never">
          <div class="toolbar">
            <el-input
              v-model="keyword"
              placeholder="编码/名称"
              clearable
              style="width: 180px"
              @keyup.enter="onSearch"
              @clear="onSearch"
            />
            <el-select v-model="status" placeholder="状态" clearable style="width: 110px" @change="onSearch">
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="0" />
            </el-select>
            <el-button type="primary" @click="onSearch">搜索</el-button>
            <el-button v-if="canAdd" type="primary" @click="openCreate">新增字典</el-button>
          </div>

          <el-table
            v-loading="loading"
            :data="list"
            highlight-current-row
            stripe
            @current-change="onSelectDict"
          >
            <el-table-column prop="code" label="编码" min-width="120" />
            <el-table-column prop="name" label="名称" min-width="120" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="asDict(row).status === 1 ? 'success' : 'info'" size="small">
                  {{ asDict(row).status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="70" />
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button v-if="canEdit" link type="primary" @click="openEdit(asDict(row))">编辑</el-button>
                <el-button v-if="canDelete" link type="danger" @click="onDeleteDict(asDict(row))">删除</el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty description="暂无字典" :image-size="60" />
            </template>
          </el-table>

          <el-pagination
            class="pager"
            layout="total, prev, pager, next"
            :total="total"
            :page-size="pageSize"
            :current-page="page"
            @current-change="onPageChange"
          />
        </el-card>
      </el-col>

      <el-col :span="14">
        <el-card shadow="never">
          <div class="toolbar">
            <span class="title">字典项{{ current ? `（${current.name}）` : '' }}</span>
            <div class="actions">
              <el-button :disabled="!current" @click="loadItems">刷新</el-button>
              <el-button v-if="canItemAdd" type="primary" :disabled="!current" @click="openItemCreate">
                新增字典项
              </el-button>
            </div>
          </div>

          <el-table v-loading="itemLoading" :data="items" stripe>
            <el-table-column prop="label" label="标签" min-width="120" />
            <el-table-column prop="value" label="值" min-width="120" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="asDictItem(row).status === 1 ? 'success' : 'info'" size="small">
                  {{ asDictItem(row).status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="70" />
            <el-table-column label="操作" width="140" fixed="right">
              <template #default="{ row }">
                <el-button v-if="canItemEdit" link type="primary" @click="openItemEdit(asDictItem(row))">
                  编辑
                </el-button>
                <el-button v-if="canItemDelete" link type="danger" @click="onDeleteItem(asDictItem(row))">
                  删除
                </el-button>
              </template>
            </el-table-column>
            <template #empty>
              <el-empty :description="current ? '暂无字典项' : '请先选择左侧字典'" :image-size="60" />
            </template>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="dictVisible" :title="dictTitle" width="460px">
      <el-form :model="dictForm" label-width="90px">
        <el-form-item label="编码" required>
          <el-input v-model="dictForm.code" :disabled="isEdit" placeholder="如 user_status" />
        </el-form-item>
        <el-form-item label="名称" required><el-input v-model="dictForm.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="dictForm.description" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="dictForm.sort" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="dictForm.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dictVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitDict">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="itemVisible" :title="itemTitle" width="460px">
      <el-form :model="itemForm" label-width="90px">
        <el-form-item label="标签" required><el-input v-model="itemForm.label" /></el-form-item>
        <el-form-item label="值" required><el-input v-model="itemForm.value" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="itemForm.sort" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="itemForm.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="itemVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitItem">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import {
  createDict,
  createDictItem,
  deleteDict,
  deleteDictItem,
  getDictItems,
  getDicts,
  updateDict,
  updateDictItem,
  type DictDataItem,
  type DictItem,
} from '@/api/dict'

// el-table 作用域插槽的 row 为 DefaultRow，这里显式收敛为具体类型
const asDict = (row: unknown) => row as DictItem
const asDictItem = (row: unknown) => row as DictDataItem

const userStore = useUserStore()
const canAdd = computed(() => userStore.hasPermission('system:dict:add'))
const canEdit = computed(() => userStore.hasPermission('system:dict:edit'))
const canDelete = computed(() => userStore.hasPermission('system:dict:delete'))
const canItemAdd = computed(() => userStore.hasPermission('system:dict_item:add'))
const canItemEdit = computed(() => userStore.hasPermission('system:dict_item:edit'))
const canItemDelete = computed(() => userStore.hasPermission('system:dict_item:delete'))

const loading = ref(false)
const itemLoading = ref(false)
const list = ref<DictItem[]>([])
const items = ref<DictDataItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const status = ref<number | undefined>(undefined)
const current = ref<DictItem | null>(null)

async function loadDicts() {
  loading.value = true
  try {
    const resp = await getDicts({
      keyword: keyword.value,
      status: status.value,
      page: page.value,
      page_size: pageSize.value,
    })
    list.value = resp.list || []
    total.value = resp.total || 0
    // 当前选中的字典若已不在列表中则清空右侧
    if (current.value && !list.value.some((d) => d.id === current.value?.id)) {
      current.value = null
      items.value = []
    }
  } finally {
    loading.value = false
  }
}

async function loadItems() {
  if (!current.value) {
    items.value = []
    return
  }
  itemLoading.value = true
  try {
    const resp = await getDictItems({ dict_code: current.value.code })
    items.value = resp.list || []
  } finally {
    itemLoading.value = false
  }
}

function onSelectDict(row: DictItem | null) {
  current.value = row ?? null
  loadItems()
}

function onSearch() {
  page.value = 1
  loadDicts()
}

function onPageChange(p: number) {
  page.value = p
  loadDicts()
}

const dictVisible = ref(false)
const dictTitle = ref('')
const isEdit = ref(false)
const editId = ref(0)
const submitting = ref(false)
const emptyDictForm = { code: '', name: '', description: '', sort: 0, enabled: true }
const dictForm = reactive({ ...emptyDictForm })

function openCreate() {
  isEdit.value = false
  dictTitle.value = '新增字典'
  Object.assign(dictForm, emptyDictForm)
  dictVisible.value = true
}

function openEdit(item: DictItem) {
  isEdit.value = true
  dictTitle.value = '编辑字典'
  Object.assign(dictForm, {
    code: item.code,
    name: item.name,
    description: item.description,
    sort: item.sort,
    enabled: item.status === 1,
  })
  editId.value = item.id
  dictVisible.value = true
}

async function submitDict() {
  if (!dictForm.code.trim() || !dictForm.name.trim()) {
    ElMessage.warning('编码与名称不能为空')
    return
  }
  submitting.value = true
  try {
    const statusVal = dictForm.enabled ? 1 : 0
    if (isEdit.value) {
      await updateDict(editId.value, {
        name: dictForm.name,
        description: dictForm.description,
        status: statusVal,
        sort: dictForm.sort,
      })
      ElMessage.success('更新成功')
    } else {
      await createDict({
        code: dictForm.code.trim(),
        name: dictForm.name.trim(),
        description: dictForm.description,
        status: statusVal,
        sort: dictForm.sort,
      })
      ElMessage.success('创建成功')
    }
    dictVisible.value = false
    await loadDicts()
  } finally {
    submitting.value = false
  }
}

async function onDeleteDict(item: DictItem) {
  try {
    await ElMessageBox.confirm(
      `确定删除字典「${item.name}」吗？其下所有字典项会一并删除，且不可恢复。`,
      '提示',
      { type: 'warning' },
    )
  } catch {
    return
  }
  await deleteDict(item.id)
  ElMessage.success('删除成功')
  if (current.value?.id === item.id) {
    current.value = null
    items.value = []
  }
  await loadDicts()
}

const itemVisible = ref(false)
const itemTitle = ref('')
const isItemEdit = ref(false)
const itemEditId = ref(0)
const emptyItemForm = { label: '', value: '', sort: 0, enabled: true }
const itemForm = reactive({ ...emptyItemForm })

function openItemCreate() {
  if (!current.value) return
  isItemEdit.value = false
  itemTitle.value = '新增字典项'
  Object.assign(itemForm, emptyItemForm)
  itemVisible.value = true
}

function openItemEdit(item: DictDataItem) {
  isItemEdit.value = true
  itemTitle.value = '编辑字典项'
  Object.assign(itemForm, {
    label: item.label,
    value: item.value,
    sort: item.sort,
    enabled: item.status === 1,
  })
  itemEditId.value = item.id
  itemVisible.value = true
}

async function submitItem() {
  if (!current.value) return
  if (!itemForm.label.trim() || !itemForm.value.trim()) {
    ElMessage.warning('标签与值不能为空')
    return
  }
  submitting.value = true
  try {
    const statusVal = itemForm.enabled ? 1 : 0
    if (isItemEdit.value) {
      await updateDictItem(itemEditId.value, {
        label: itemForm.label,
        value: itemForm.value,
        status: statusVal,
        sort: itemForm.sort,
      })
      ElMessage.success('更新成功')
    } else {
      await createDictItem({
        dict_code: current.value.code,
        label: itemForm.label.trim(),
        value: itemForm.value.trim(),
        status: statusVal,
        sort: itemForm.sort,
      })
      ElMessage.success('创建成功')
    }
    itemVisible.value = false
    await loadItems()
  } finally {
    submitting.value = false
  }
}

async function onDeleteItem(item: DictDataItem) {
  try {
    await ElMessageBox.confirm(`确定删除字典项「${item.label}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  await deleteDictItem(item.id)
  ElMessage.success('删除成功')
  await loadItems()
}

onMounted(loadDicts)
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}
.actions {
  display: flex;
  gap: 8px;
}
.title {
  font-weight: 600;
}
.pager {
  margin-top: 12px;
  justify-content: flex-end;
}
</style>
