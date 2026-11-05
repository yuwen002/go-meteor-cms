<template>
  <div class="article-page">
    <el-card shadow="never">
      <div class="toolbar">
        <el-input
          v-model="query.keyword"
          placeholder="标题 / 摘要"
          clearable
          style="width: 200px"
          @keyup.enter="onSearch"
          @clear="onSearch"
        />
        <el-select v-model="query.status" placeholder="状态" clearable style="width: 120px" @change="onSearch">
          <el-option label="草稿" :value="0" />
          <el-option label="已发布" :value="1" />
          <el-option label="已下线" :value="2" />
        </el-select>
        <el-select v-model="query.category_id" placeholder="分类" clearable style="width: 150px" @change="onSearch">
          <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
        <el-button @click="onSearch">搜索</el-button>
        <el-button v-if="canAdd" type="primary" @click="openCreate">写文章</el-button>
      </div>

      <el-table v-loading="loading" :data="list" stripe>
        <el-table-column label="标题" min-width="240">
          <template #default="{ row }">
            <span class="article-title">{{ asRow(row).title }}</span>
            <el-tag v-if="asRow(row).is_top" type="warning" size="small" effect="plain">置顶</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="95">
          <template #default="{ row }">
            <el-tag :type="statusType(asRow(row).status)" size="small">
              {{ asRow(row).status_text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="category_name" label="分类" min-width="110" />
        <el-table-column label="标签" min-width="150">
          <template #default="{ row }">
            <el-tag
              v-for="t in asRow(row).tag_names"
              :key="t"
              size="small"
              effect="plain"
              class="tag-chip"
            >
              {{ t }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="80" />
        <el-table-column label="发布时间" width="170">
          <template #default="{ row }">{{ asRow(row).published_at || '—' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button v-if="canEdit" link type="primary" @click="openEdit(asRow(row))">编辑</el-button>
            <el-button v-if="canPublish" link type="success" @click="openPublish(asRow(row))">发布</el-button>
            <el-button
              v-if="canWithdraw && asRow(row).status === 1"
              link
              type="warning"
              @click="onWithdraw(asRow(row))"
            >
              撤回
            </el-button>
            <el-button v-if="canTop" link @click="onToggleTop(asRow(row))">
              {{ asRow(row).is_top ? '取消置顶' : '置顶' }}
            </el-button>
            <el-button link @click="onPreview(asRow(row))">预览</el-button>
            <el-button v-if="canDelete" link type="danger" @click="onDelete(asRow(row))">删除</el-button>
          </template>
        </el-table-column>
        <template #empty><el-empty description="暂无文章" :image-size="60" /></template>
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

    <el-dialog v-model="publishVisible" title="发布文章" width="440px">
      <p class="publish-tip">
        「{{ publishTarget?.title }}」
        <br />
        留空表示立即发布；选择未来时间则到点后自动对前台可见。
      </p>
      <el-date-picker
        v-model="publishAt"
        type="datetime"
        placeholder="发布时间（可留空）"
        value-format="YYYY-MM-DD HH:mm:ss"
        style="width: 100%"
      />
      <template #footer>
        <el-button @click="publishVisible = false">取消</el-button>
        <el-button type="primary" :loading="publishing" @click="doPublish">确定发布</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import {
  deleteArticle,
  getArticles,
  getCategories,
  publishArticle,
  toggleArticleTop,
  withdrawArticle,
  type ArticleItem,
  type CategoryItem,
} from '@/api/content'

const router = useRouter()
const userStore = useUserStore()
const canAdd = computed(() => userStore.hasPermission('cms:article:add'))
const canEdit = computed(() => userStore.hasPermission('cms:article:edit'))
const canDelete = computed(() => userStore.hasPermission('cms:article:delete'))
const canPublish = computed(() => userStore.hasPermission('cms:article:publish'))
const canWithdraw = computed(() => userStore.hasPermission('cms:article:withdraw'))
const canTop = computed(() => userStore.hasPermission('cms:article:toggle_top'))

const list = ref<ArticleItem[]>([])
const categories = ref<CategoryItem[]>([])
const total = ref(0)
const loading = ref(false)
const query = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  status: undefined as number | undefined,
  category_id: undefined as number | undefined,
})

const publishVisible = ref(false)
const publishing = ref(false)
const publishTarget = ref<ArticleItem | null>(null)
const publishAt = ref('')

const asRow = (row: unknown) => row as ArticleItem
const statusType = (status: number) => {
  if (status === 1) return 'success'
  if (status === 2) return 'info'
  return 'warning'
}

async function loadList() {
  loading.value = true
  try {
    const res = await getArticles({
      page: query.page,
      page_size: query.page_size,
      keyword: query.keyword || undefined,
      status: query.status,
      category_id: query.category_id,
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

async function loadCategories() {
  try {
    const res = await getCategories()
    categories.value = res.list ?? []
  } catch {
    categories.value = []
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

function openCreate() {
  void router.push('/content/articles/new')
}

function openEdit(row: ArticleItem) {
  void router.push(`/content/articles/${row.id}/edit`)
}

function onPreview(row: ArticleItem) {
  if (row.status !== 1) {
    ElMessage.info('文章尚未发布，请先发布后再预览前台页面')
    return
  }
  window.open(`/articles/${row.slug}`, '_blank')
}

function openPublish(row: ArticleItem) {
  publishTarget.value = row
  publishAt.value = ''
  publishVisible.value = true
}

async function doPublish() {
  if (!publishTarget.value) return
  publishing.value = true
  try {
    await publishArticle(publishTarget.value.id, publishAt.value || undefined)
    ElMessage.success(publishAt.value ? '已设置定时发布' : '已发布')
    publishVisible.value = false
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  } finally {
    publishing.value = false
  }
}

async function onWithdraw(row: ArticleItem) {
  try {
    await ElMessageBox.confirm(`确定撤回「${row.title}」吗？撤回后前台将不再展示。`, '提示', {
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    await withdrawArticle(row.id)
    ElMessage.success('已撤回')
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  }
}

async function onToggleTop(row: ArticleItem) {
  try {
    await toggleArticleTop(row.id)
    ElMessage.success(row.is_top ? '已取消置顶' : '已置顶')
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  }
}

async function onDelete(row: ArticleItem) {
  try {
    await ElMessageBox.confirm(`确定删除「${row.title}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteArticle(row.id)
    ElMessage.success('已删除')
    await loadList()
  } catch {
    // 统一拦截器已提示错误
  }
}

onMounted(() => {
  void loadList()
  void loadCategories()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 14px;
}
.article-title {
  margin-right: 6px;
}
.tag-chip {
  margin-right: 4px;
}
.pager {
  margin-top: 14px;
  justify-content: flex-end;
}
.publish-tip {
  margin: 0 0 14px;
  line-height: 1.7;
  color: var(--el-text-color-regular);
}
</style>
