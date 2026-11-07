<template>
  <div class="trash-page">
    <el-card shadow="never">
      <div class="toolbar">
        <el-select v-model="query.type" placeholder="全部类型" clearable style="width: 140px" @change="onSearch">
          <el-option label="文章" value="article" />
          <el-option label="分类" value="category" />
          <el-option label="标签" value="tag" />
        </el-select>
        <el-button @click="onSearch">搜索</el-button>
        <el-button @click="loadList">刷新</el-button>
        <span class="tip">删除的内容会先进入回收站，可恢复；彻底删除会同时清理其引用文件且不可撤销。</span>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="typeTag(asRow(row).entity_type)">
              {{ typeLabel(asRow(row).entity_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="标题 / 名称" min-width="240">
          <template #default="{ row }">{{ asRow(row).title }}</template>
        </el-table-column>
        <el-table-column label="别名" min-width="180">
          <template #default="{ row }">{{ asRow(row).slug || '—' }}</template>
        </el-table-column>
        <el-table-column label="删除时间" width="190">
          <template #default="{ row }">{{ formatDate(asRow(row).deleted_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button v-if="canRestore" link type="primary" @click="onRestore(asRow(row))">恢复</el-button>
            <el-button v-if="canPurge" link type="danger" @click="onPurge(asRow(row))">彻底删除</el-button>
          </template>
        </el-table-column>
        <template #empty><el-empty description="回收站是空的" :image-size="60" /></template>
      </el-table>

      <el-pagination
        class="pager"
        layout="total, prev, pager, next"
        :total="total"
        :page-size="query.page_size"
        :current-page="query.page"
        @current-change="onPageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from '@/utils/element'
import { useUserStore } from '@/stores/user'
import {
  getTrash,
  purgeTrash,
  restoreTrash,
  type TrashEntityType,
  type TrashItem,
} from '@/api/content'

type TagType = 'primary' | 'success' | 'info' | 'warning' | 'danger'

const userStore = useUserStore()
const canRestore = computed(() => userStore.hasPermission('cms:trash:restore'))
const canPurge = computed(() => userStore.hasPermission('cms:trash:purge'))

const list = ref<TrashItem[]>([])
const total = ref(0)
const loading = ref(false)
const query = reactive({
  page: 1,
  page_size: 20,
  type: undefined as TrashEntityType | undefined,
})

const asRow = (row: unknown) => row as TrashItem
const typeLabel = (t: TrashEntityType) =>
  t === 'article' ? '文章' : t === 'category' ? '分类' : t === 'tag' ? '标签' : t
const typeTag = (t: TrashEntityType): TagType =>
  t === 'article' ? 'primary' : t === 'category' ? 'success' : 'info'
const formatDate = (s: string) => (s ? s.replace('T', ' ').slice(0, 19) : '—')

async function loadList() {
  loading.value = true
  try {
    const res = await getTrash({
      page: query.page,
      page_size: query.page_size,
      type: query.type,
    })
    list.value = res.list ?? []
    total.value = res.total ?? 0
  } catch {
    list.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function onSearch() {
  query.page = 1
  void loadList()
}

function onPageChange(p: number) {
  query.page = p
  void loadList()
}

async function onRestore(row: TrashItem) {
  try {
    await ElMessageBox.confirm(`确定恢复「${row.title}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    const res = await restoreTrash(row.entity_type, row.id)
    ElMessage.success(res?.message || '已恢复')
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  }
}

async function onPurge(row: TrashItem) {
  try {
    await ElMessageBox.confirm(
      `彻底删除「${row.title}」后不可恢复，关联的标签、修订记录与不再被引用的图片也会被清理。确定继续吗？`,
      '危险操作',
      { type: 'warning', confirmButtonText: '彻底删除', confirmButtonClass: 'el-button--danger' },
    )
  } catch {
    return
  }
  try {
    await purgeTrash(row.entity_type, row.id)
    ElMessage.success('已彻底删除')
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  }
}

onMounted(() => {
  void loadList()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}
.tip {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.pager {
  margin-top: 14px;
  justify-content: flex-end;
}
</style>
