<template>
  <div class="config-page">
    <el-card shadow="never">
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="配置键/名称"
          clearable
          style="width: 220px"
          @keyup.enter="load"
          @clear="load"
        />
        <el-button type="primary" @click="load">搜索</el-button>
        <el-button v-if="canAdd" type="primary" @click="openCreate">新增配置</el-button>
        <span class="tip">系统内置配置（标记为「内置」）不可删除</span>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column prop="config_key" label="配置键" min-width="180" />
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column prop="config_value" label="配置值" min-width="200" show-overflow-tooltip />
        <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
        <el-table-column label="类型" width="90">
          <template #default="{ row }">
            <el-tag v-if="asRow(row).is_system" type="warning" size="small">内置</el-tag>
            <el-tag v-else type="info" size="small">自定义</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button v-if="canEdit" link type="primary" @click="openEdit(asRow(row))">编辑</el-button>
            <el-button
              v-if="canDelete"
              link
              type="danger"
              :disabled="asRow(row).is_system"
              @click="onDelete(asRow(row))"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无配置" />
        </template>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="配置键" required>
          <el-input v-model="form.config_key" :disabled="isEdit" placeholder="如 site.name" />
        </el-form-item>
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="配置值">
          <el-input v-model="form.config_value" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { createConfig, deleteConfig, getConfigs, updateConfig, type ConfigItem } from '@/api/config'

// el-table 作用域插槽的 row 为 DefaultRow，这里显式收敛为 ConfigItem
const asRow = (row: unknown) => row as ConfigItem

const userStore = useUserStore()
const canAdd = computed(() => userStore.hasPermission('system:config:add'))
const canEdit = computed(() => userStore.hasPermission('system:config:edit'))
const canDelete = computed(() => userStore.hasPermission('system:config:delete'))

const loading = ref(false)
const list = ref<ConfigItem[]>([])
const keyword = ref('')

async function load() {
  loading.value = true
  try {
    const resp = await getConfigs({ keyword: keyword.value })
    list.value = resp.list || []
  } finally {
    loading.value = false
  }
}

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const editId = ref(0)
const submitting = ref(false)
const emptyForm = { config_key: '', name: '', config_value: '', remark: '' }
const form = reactive({ ...emptyForm })

function openCreate() {
  isEdit.value = false
  dialogTitle.value = '新增配置'
  Object.assign(form, emptyForm)
  dialogVisible.value = true
}

function openEdit(item: ConfigItem) {
  isEdit.value = true
  dialogTitle.value = '编辑配置'
  Object.assign(form, {
    config_key: item.config_key,
    name: item.name,
    config_value: item.config_value,
    remark: item.remark,
  })
  editId.value = item.id
  dialogVisible.value = true
}

async function submit() {
  if (!form.config_key.trim() || !form.name.trim()) {
    ElMessage.warning('配置键与名称不能为空')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateConfig(editId.value, {
        name: form.name,
        config_value: form.config_value,
        remark: form.remark,
      })
      ElMessage.success('更新成功')
    } else {
      await createConfig({
        config_key: form.config_key.trim(),
        config_value: form.config_value,
        name: form.name.trim(),
        remark: form.remark,
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    await load()
  } finally {
    submitting.value = false
  }
}

async function onDelete(item: ConfigItem) {
  if (item.is_system) {
    ElMessage.warning('系统内置配置不可删除')
    return
  }
  try {
    await ElMessageBox.confirm(`确定删除配置「${item.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  await deleteConfig(item.id)
  ElMessage.success('删除成功')
  await load()
}

onMounted(load)
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}
.tip {
  color: #909399;
  font-size: 12px;
}
</style>
