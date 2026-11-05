<template>
  <div class="portal">
    <header class="site-header">
      <div class="header-inner">
        <router-link to="/articles" class="brand">Meteor<span>CMS</span></router-link>
        <nav class="nav-links">
          <router-link to="/articles">全部文章</router-link>
          <router-link to="/login">管理后台</router-link>
        </nav>
      </div>
    </header>

    <main class="container">
      <el-breadcrumb class="crumb" separator="/">
        <el-breadcrumb-item><router-link to="/articles">文章</router-link></el-breadcrumb-item>
        <el-breadcrumb-item>{{ article?.title ?? '详情' }}</el-breadcrumb-item>
      </el-breadcrumb>

      <div v-loading="loading" class="detail-card">
        <template v-if="article">
          <h1 class="detail-title">{{ article.title }}</h1>
          <div class="detail-meta">
            <span v-if="article.category_name" class="chip">{{ article.category_name }}</span>
            <span>{{ formatDate(article.published_at) }}</span>
            <span>{{ article.view_count }} 次阅读</span>
          </div>
          <img v-if="article.cover" :src="article.cover" class="detail-cover" :alt="article.title" />
          <!-- renderMarkdown 先转义 HTML 再套用有限语法，链接协议已做白名单过滤 -->
          <div class="article-body" v-html="rendered" />
          <div v-if="article.tag_names.length" class="detail-tags">
            <span v-for="t in article.tag_names" :key="t" class="tag"># {{ t }}</span>
          </div>
        </template>
        <el-empty v-else-if="!loading" description="文章不存在或尚未发布" :image-size="90" />
      </div>

      <div class="back-row">
        <el-button @click="goBack">← 返回文章列表</el-button>
      </div>
    </main>

    <footer class="site-footer">Powered by Go Meteor CMS</footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPublicArticle, type PublicArticleDetail } from '@/api/public'
import { renderMarkdown } from '@/utils/markdown'

const route = useRoute()
const router = useRouter()
const slug = computed(() => String(route.params.slug ?? ''))
const article = ref<PublicArticleDetail | null>(null)
const loading = ref(false)
const rendered = computed(() => renderMarkdown(article.value?.content ?? ''))
const originalTitle = document.title

function formatDate(s: string): string {
  return s ? s.replace('T', ' ').slice(0, 16) : ''
}

async function load() {
  if (!slug.value) return
  loading.value = true
  try {
    article.value = await getPublicArticle(slug.value)
    document.title = article.value.seo_title || article.value.title
    window.scrollTo({ top: 0 })
  } catch {
    article.value = null
  } finally {
    loading.value = false
  }
}

function goBack() {
  void router.push('/articles')
}

onMounted(() => {
  void load()
})
watch(slug, () => {
  void load()
})
onUnmounted(() => {
  document.title = originalTitle
})
</script>

<style scoped>
.portal {
  min-height: 100vh;
  background: #f6f8fb;
  color: #1f2937;
}
.site-header {
  background: linear-gradient(120deg, #1d4ed8, #0ea5e9);
  color: #fff;
}
.header-inner {
  max-width: 860px;
  margin: 0 auto;
  padding: 18px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.brand {
  font-size: 21px;
  font-weight: 700;
  color: #fff;
  text-decoration: none;
}
.brand span {
  opacity: 0.75;
  margin-left: 4px;
  font-weight: 400;
}
.nav-links {
  display: flex;
  gap: 22px;
}
.nav-links a {
  color: rgba(255, 255, 255, 0.88);
  text-decoration: none;
  font-size: 14px;
}
.container {
  max-width: 860px;
  margin: 0 auto;
  padding: 0 24px 48px;
}
.crumb {
  margin: 26px 0 18px;
}
.crumb a {
  color: #2563eb;
  text-decoration: none;
}
.detail-card {
  background: #fff;
  border-radius: 10px;
  padding: 36px 40px 40px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.06);
  min-height: 260px;
}
.detail-title {
  margin: 0 0 14px;
  font-size: 30px;
  line-height: 1.4;
  color: #0f172a;
}
.detail-meta {
  display: flex;
  align-items: center;
  gap: 14px;
  font-size: 13px;
  color: #94a3b8;
  padding-bottom: 22px;
  border-bottom: 1px solid #f1f5f9;
}
.chip {
  color: #2563eb;
  background: #eff6ff;
  padding: 2px 10px;
  border-radius: 10px;
}
.detail-cover {
  width: 100%;
  max-height: 360px;
  object-fit: cover;
  border-radius: 8px;
  margin: 24px 0;
}
.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid #f1f5f9;
}
.tag {
  font-size: 13px;
  color: #2563eb;
  background: #eff6ff;
  padding: 3px 10px;
  border-radius: 12px;
}
.back-row {
  margin-top: 22px;
}
.site-footer {
  text-align: center;
  padding: 28px;
  color: #94a3b8;
  font-size: 12px;
}

/* v-html 注入的节点不带 scoped 属性，必须用 :deep 命中 */
.article-body {
  margin-top: 24px;
  line-height: 1.85;
  font-size: 15px;
  word-break: break-word;
}
.article-body :deep(h1),
.article-body :deep(h2),
.article-body :deep(h3) {
  margin: 28px 0 12px;
  font-weight: 600;
  line-height: 1.4;
  color: #0f172a;
}
.article-body :deep(h1) { font-size: 26px; }
.article-body :deep(h2) { font-size: 22px; }
.article-body :deep(h3) { font-size: 18px; }
.article-body :deep(p) { margin: 12px 0; }
.article-body :deep(ul),
.article-body :deep(ol) { padding-left: 26px; margin: 12px 0; }
.article-body :deep(blockquote) {
  margin: 16px 0;
  padding: 8px 16px;
  border-left: 4px solid #60a5fa;
  color: #64748b;
  background: #f8fafc;
}
.article-body :deep(pre) {
  padding: 14px 16px;
  overflow: auto;
  border-radius: 6px;
  background: #0f172a;
  color: #e2e8f0;
}
.article-body :deep(code) {
  padding: 2px 5px;
  border-radius: 3px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  background: #f1f5f9;
}
.article-body :deep(pre code) { padding: 0; background: transparent; }
.article-body :deep(img) { max-width: 100%; border-radius: 6px; }
.article-body :deep(a) { color: #2563eb; text-decoration: none; }
.article-body :deep(a:hover) { text-decoration: underline; }
.article-body :deep(hr) { margin: 24px 0; border: 0; border-top: 1px solid #e2e8f0; }

@media (max-width: 640px) {
  .detail-card { padding: 24px 20px; }
  .detail-title { font-size: 24px; }
}
</style>
