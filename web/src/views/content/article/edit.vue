<template>
  <div class="article-edit">
    <el-card shadow="never">
      <div class="header">
        <div class="header-left">
          <el-button link @click="goBack">← 返回列表</el-button>
          <span class="title">{{ isEdit ? '编辑文章' : '写文章' }}</span>
          <el-tag v-if="isEdit" size="small" :type="statusTagType">{{ statusText }}</el-tag>
        </div>
        <div class="actions">
          <el-button v-if="isEdit && canRevision" @click="openRevisions">修订历史</el-button>
          <el-button
            v-if="isEdit && status === 0 && canSubmit"
            :loading="submitting"
            @click="onSubmit"
          >
            提交审核
          </el-button>
          <el-button :loading="saving" @click="save">保存草稿</el-button>
          <el-button v-if="canPublish" type="primary" :loading="publishing" @click="publish">保存并发布</el-button>
        </div>
      </div>

      <el-form :model="form" label-width="90px" class="main-form">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" maxlength="255" show-word-limit placeholder="请输入文章标题" />
        </el-form-item>
        <el-form-item label="别名">
          <el-input v-model="form.slug" placeholder="留空自动生成；仅小写字母、数字和连字符" />
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="2" maxlength="512" show-word-limit />
        </el-form-item>
        <el-form-item label="正文" required>
          <MarkdownEditor v-model="form.content" :rows="20" />
        </el-form-item>
      </el-form>

      <el-divider content-position="left">发布设置</el-divider>
      <el-form :model="form" label-width="90px">
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="分类">
              <el-select v-model="form.category_id" clearable placeholder="未分类" style="width: 100%">
                <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="标签">
              <el-select v-model="form.tag_ids" multiple filterable placeholder="选择标签" style="width: 100%">
                <el-option v-for="t in tags" :key="t.id" :label="t.name" :value="t.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="置顶"><el-switch v-model="form.is_top" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="封面图">
          <div class="cover-row">
            <el-upload :show-file-list="false" accept="image/*" :http-request="uploadCover">
              <el-button :loading="coverUploading">{{ form.cover ? '更换封面' : '上传封面' }}</el-button>
            </el-upload>
            <el-button v-if="form.cover" link type="danger" @click="form.cover = ''">移除</el-button>
            <img v-if="form.cover" :src="form.cover" class="cover-preview" alt="封面预览" />
          </div>
        </el-form-item>
      </el-form>

      <el-collapse>
        <el-collapse-item title="SEO 设置（留空自动取标题与摘要）" name="seo">
          <el-form :model="form" label-width="90px">
            <el-form-item label="SEO 标题"><el-input v-model="form.seo_title" /></el-form-item>
            <el-form-item label="SEO 描述">
              <el-input v-model="form.seo_description" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="关键词">
              <el-input v-model="form.seo_keywords" placeholder="逗号分隔" />
            </el-form-item>
          </el-form>
        </el-collapse-item>
      </el-collapse>

      <p class="hint">
        定时发布请在文章列表页点「发布」并选择未来时间；到点后文章会自动对前台可见，无需后台任务。修改别名会保留旧地址的 301 跳转。
      </p>
    </el-card>

    <el-drawer v-model="revisionVisible" title="文章轨迹" size="540px">
      <el-tabs v-model="traceTab">
        <el-tab-pane label="修订历史" name="revisions">
          <div v-loading="revisionLoading" class="revision-list">
            <div v-for="r in revisions" :key="r.id" class="revision-item">
              <div class="revision-head">
                <span class="revision-time">{{ formatDate(r.created_at) }}</span>
                <span class="revision-title">{{ r.title }}</span>
              </div>
              <p class="revision-remark">{{ r.remark || '编辑前自动存档' }}</p>
              <div class="revision-actions">
                <el-button link type="primary" @click="toggleRevisionDetail(r.id)">
                  {{ currentRevision?.id === r.id ? '收起预览' : '查看内容' }}
                </el-button>
                <el-button v-if="canRevisionRestore" link type="warning" @click="onRestoreRevision(r.id)">
                  恢复此版本
                </el-button>
              </div>
              <!-- renderMarkdown 先转义 HTML 再套用有限语法，可安全用于 v-html -->
              <div
                v-if="currentRevision?.id === r.id"
                class="revision-preview"
                v-html="revisionRendered"
              />
            </div>
            <el-empty v-if="!revisionLoading && revisions.length === 0" description="暂无历史版本" :image-size="60" />
          </div>
        </el-tab-pane>

        <el-tab-pane v-if="canReview" label="审核记录" name="reviews">
          <div v-loading="reviewLoading" class="review-list">
            <div v-for="rw in reviews" :key="rw.id" class="review-item">
              <div class="review-head">
                <el-tag size="small" :type="reviewTagType(rw.action)">{{ rw.action_text }}</el-tag>
                <span class="review-time">{{ formatDate(rw.created_at) }}</span>
              </div>
              <p class="review-operator">
                {{ rw.operator_name || '系统' }}<span v-if="rw.operator_id">（#{{ rw.operator_id }}）</span>
              </p>
              <p v-if="rw.remark" class="review-remark">备注：{{ rw.remark }}</p>
            </div>
            <el-empty v-if="!reviewLoading && reviews.length === 0" description="暂无审核记录" :image-size="60" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from '@/utils/element'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useUserStore } from '@/stores/user'
import { uploadFile } from '@/api/upload'
import { renderMarkdown } from '@/utils/markdown'
import {
  createArticle,
  getArticle,
  getArticleRevision,
  getArticleReviews,
  getArticleRevisions,
  getCategories,
  getTags,
  publishArticle,
  restoreArticleRevision,
  submitArticle,
  updateArticle,
  type ArticlePayload,
  type ArticleReviewItem,
  type ArticleRevisionDetail,
  type ArticleRevisionItem,
  type CategoryItem,
  type TagItem,
} from '@/api/content'

type TagType = 'primary' | 'success' | 'info' | 'warning' | 'danger'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const canPublish = computed(() => userStore.hasPermission('cms:article:publish'))
const canSubmit = computed(() => userStore.hasPermission('cms:article:submit'))
const canRevision = computed(() => userStore.hasPermission('cms:article:revision'))
const canRevisionRestore = computed(() => userStore.hasPermission('cms:article:revision_restore'))

const articleId = computed(() => Number(route.params.id || 0))
const isEdit = computed(() => articleId.value > 0)
const status = ref(0)
const statusText = ref('草稿')
const statusTagType = computed<TagType>(() => {
  if (statusText.value === '定时中') return 'primary'
  if (status.value === 1) return 'success'
  if (status.value === 2) return 'info'
  if (status.value === 3) return 'danger'
  return 'warning'
})

interface ArticleForm {
  title: string
  slug: string
  summary: string
  cover: string
  content: string
  category_id: number | undefined
  tag_ids: number[]
  is_top: boolean
  seo_title: string
  seo_description: string
  seo_keywords: string
}
const emptyForm = (): ArticleForm => ({
  title: '',
  slug: '',
  summary: '',
  cover: '',
  content: '',
  category_id: undefined,
  tag_ids: [],
  is_top: false,
  seo_title: '',
  seo_description: '',
  seo_keywords: '',
})
const form = reactive<ArticleForm>(emptyForm())

const categories = ref<CategoryItem[]>([])
const tags = ref<TagItem[]>([])
const saving = ref(false)
const publishing = ref(false)
const submitting = ref(false)
const coverUploading = ref(false)

const revisionVisible = ref(false)
const revisionLoading = ref(false)
const revisions = ref<ArticleRevisionItem[]>([])
const currentRevision = ref<ArticleRevisionDetail | null>(null)
const revisionRendered = computed(() => renderMarkdown(currentRevision.value?.content ?? ''))

const canReview = computed(() => userStore.hasPermission('cms:article:review'))
const traceTab = ref('revisions')
const reviewLoading = ref(false)
const reviews = ref<ArticleReviewItem[]>([])
const reviewTagType = (action: number): TagType =>
  action === 2 ? 'success' : action === 3 ? 'danger' : 'primary'

function payload(): ArticlePayload {
  return {
    title: form.title.trim(),
    slug: form.slug.trim() || undefined,
    summary: form.summary,
    cover: form.cover,
    content: form.content,
    category_id: Number(form.category_id ?? 0),
    tag_ids: form.tag_ids,
    is_top: form.is_top,
    seo_title: form.seo_title,
    seo_description: form.seo_description,
    seo_keywords: form.seo_keywords,
  }
}

function validate(): boolean {
  if (!form.title.trim()) {
    ElMessage.warning('请填写文章标题')
    return false
  }
  return true
}

async function loadTaxonomy() {
  try {
    const [cs, ts] = await Promise.all([getCategories(), getTags()])
    categories.value = cs.list ?? []
    tags.value = ts.list ?? []
  } catch {
    categories.value = []
    tags.value = []
  }
}

async function loadDetail() {
  if (!isEdit.value) return
  try {
    const d = await getArticle(articleId.value)
    Object.assign(form, {
      title: d.title,
      slug: d.slug,
      summary: d.summary,
      cover: d.cover,
      content: d.content,
      category_id: d.category_id || undefined,
      tag_ids: d.tag_ids ?? [],
      is_top: d.is_top,
      seo_title: d.seo_title,
      seo_description: d.seo_description,
      seo_keywords: d.seo_keywords,
    } as ArticleForm)
    status.value = d.status
    statusText.value = d.status_text
  } catch {
    // 统一拦截器已提示错误
  }
}

async function save() {
  if (!validate()) return
  saving.value = true
  try {
    if (isEdit.value) {
      await updateArticle(articleId.value, payload())
      ElMessage.success('已保存')
    } else {
      await createArticle(payload())
      ElMessage.success('草稿已创建')
      void router.push('/content/articles')
    }
  } catch {
    // 统一拦截器已提示错误
  } finally {
    saving.value = false
  }
}

async function publish() {
  if (!isEdit.value) {
    ElMessage.info('请先保存草稿，再发布')
    return
  }
  if (!validate()) return
  publishing.value = true
  try {
    await updateArticle(articleId.value, payload())
    await publishArticle(articleId.value)
    status.value = 1
    statusText.value = '已发布'
    ElMessage.success('已发布')
  } catch {
    // 统一拦截器已提示错误
  } finally {
    publishing.value = false
  }
}

async function onSubmit() {
  if (!validate()) return
  submitting.value = true
  try {
    await updateArticle(articleId.value, payload())
    await submitArticle(articleId.value)
    status.value = 3
    statusText.value = '待审核'
    ElMessage.success('已提交审核')
  } catch {
    // 统一拦截器已提示错误
  } finally {
    submitting.value = false
  }
}

function formatDate(s: string): string {
  return s ? s.replace('T', ' ').slice(0, 19) : ''
}

function openRevisions() {
  revisionVisible.value = true
  void loadRevisions()
  if (canReview.value) {
    void loadReviews()
  }
}

async function loadReviews() {
  if (!isEdit.value) return
  reviewLoading.value = true
  try {
    const res = await getArticleReviews(articleId.value)
    reviews.value = res.list ?? []
  } catch {
    reviews.value = []
  } finally {
    reviewLoading.value = false
  }
}

async function loadRevisions() {
  if (!isEdit.value) return
  revisionLoading.value = true
  try {
    const res = await getArticleRevisions(articleId.value)
    revisions.value = res.list ?? []
  } catch {
    revisions.value = []
  } finally {
    revisionLoading.value = false
  }
}

async function toggleRevisionDetail(rid: number) {
  if (currentRevision.value?.id === rid) {
    currentRevision.value = null
    return
  }
  try {
    currentRevision.value = await getArticleRevision(articleId.value, rid)
  } catch {
    currentRevision.value = null
  }
}

async function onRestoreRevision(rid: number) {
  try {
    await ElMessageBox.confirm(
      '恢复该版本会覆盖当前标题、摘要、封面与正文；当前内容会先自动存为新版本，可再次恢复。确定继续吗？',
      '提示',
      { type: 'warning' },
    )
  } catch {
    return
  }
  try {
    await restoreArticleRevision(articleId.value, rid)
    ElMessage.success('已恢复该版本')
    currentRevision.value = null
    await Promise.all([loadDetail(), loadRevisions()])
  } catch {
    // 统一拦截器已提示错误
  }
}

async function uploadCover(options: unknown) {
  const opts = options as { file: File }
  coverUploading.value = true
  try {
    const res = await uploadFile(opts.file)
    form.cover = res.url
    ElMessage.success('封面已上传')
  } catch {
    ElMessage.error('封面上传失败')
  } finally {
    coverUploading.value = false
  }
}

function goBack() {
  void router.push('/content/articles')
}

onMounted(() => {
  void loadTaxonomy()
  void loadDetail()
  if (route.query.tab === 'revisions') {
    openRevisions()
  }
})
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}
.actions {
  display: flex;
  gap: 8px;
}
.title {
  font-size: 17px;
  font-weight: 600;
}
.main-form {
  margin-top: 8px;
}
.cover-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.cover-preview {
  width: 96px;
  height: 60px;
  object-fit: cover;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
}
.hint {
  margin: 16px 0 0;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.revision-list {
  min-height: 120px;
}
.revision-item {
  border-bottom: 1px solid var(--el-border-color-lighter);
  padding: 4px 0 14px;
  margin-bottom: 12px;
}
.revision-head {
  display: flex;
  align-items: center;
  gap: 10px;
}
.revision-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.revision-title {
  font-size: 14px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.revision-remark {
  margin: 6px 0 2px;
  font-size: 12px;
  color: var(--el-text-color-regular);
}
.revision-actions {
  display: flex;
  gap: 6px;
}
.revision-preview {
  margin-top: 10px;
  padding: 12px 14px;
  border-radius: 6px;
  background: var(--el-fill-color-light);
  font-size: 13px;
  line-height: 1.7;
  max-height: 320px;
  overflow: auto;
  word-break: break-word;
}
.revision-preview :deep(img) { max-width: 100%; }
.revision-preview :deep(pre) { overflow: auto; }
.review-list {
  min-height: 120px;
}
.review-item {
  border-bottom: 1px solid var(--el-border-color-lighter);
  padding: 4px 0 12px;
  margin-bottom: 10px;
}
.review-head {
  display: flex;
  align-items: center;
  gap: 10px;
}
.review-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.review-operator {
  margin: 8px 0 0;
  font-size: 13px;
  color: var(--el-text-color-regular);
}
.review-remark {
  margin: 6px 0 0;
  font-size: 13px;
  line-height: 1.6;
  color: var(--el-text-color-regular);
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
