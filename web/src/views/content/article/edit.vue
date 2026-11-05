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
        定时发布请在文章列表页点「发布」并选择未来时间；到点后文章会自动对前台可见，无需后台任务。
      </p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useUserStore } from '@/stores/user'
import { uploadFile } from '@/api/upload'
import {
  createArticle,
  getArticle,
  getCategories,
  getTags,
  publishArticle,
  updateArticle,
  type ArticlePayload,
  type CategoryItem,
  type TagItem,
} from '@/api/content'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const canPublish = computed(() => userStore.hasPermission('cms:article:publish'))

const articleId = computed(() => Number(route.params.id || 0))
const isEdit = computed(() => articleId.value > 0)
const status = ref(0)
const statusText = ref('草稿')
const statusTagType = computed(() => (status.value === 1 ? 'success' : status.value === 2 ? 'info' : 'warning'))

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
const coverUploading = ref(false)

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
</style>
