import request from '@/utils/request'

// ---------- 文章 ----------

export interface ArticleItem {
  id: number
  title: string
  slug: string
  summary: string
  cover: string
  status: number
  status_text: string
  author_name: string
  category_id: number
  category_name: string
  tag_ids: number[]
  tag_names: string[]
  is_top: boolean
  view_count: number
  published_at: string
  created_at: string
  updated_at: string
}

export interface ArticleListResp {
  list: ArticleItem[]
  total: number
  page: number
  page_size: number
}

export interface ArticleDetail extends ArticleItem {
  content: string
  seo_title: string
  seo_description: string
  seo_keywords: string
}

export interface ArticlePayload {
  title: string
  slug?: string
  summary?: string
  cover?: string
  content?: string
  category_id?: number
  tag_ids?: number[]
  is_top?: boolean
  seo_title?: string
  seo_description?: string
  seo_keywords?: string
}

export interface ArticleListParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
  category_id?: number
  tag_id?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export function getArticles(params: ArticleListParams) {
  return request({ url: '/admin/articles', method: 'get', params }) as Promise<ArticleListResp>
}

export function getArticle(id: number) {
  return request({ url: `/admin/articles/${id}`, method: 'get' }) as Promise<ArticleDetail>
}

export function createArticle(data: ArticlePayload) {
  return request({ url: '/admin/articles', method: 'post', data })
}

export function updateArticle(id: number, data: ArticlePayload) {
  return request({ url: `/admin/articles/${id}`, method: 'put', data })
}

export function deleteArticle(id: number) {
  return request({ url: `/admin/articles/${id}`, method: 'delete' })
}

export function publishArticle(id: number, published_at?: string) {
  return request({ url: `/admin/articles/${id}/publish`, method: 'post', data: { published_at: published_at ?? '' } })
}

export function withdrawArticle(id: number) {
  return request({ url: `/admin/articles/${id}/withdraw`, method: 'post', data: {} })
}

export function toggleArticleTop(id: number) {
  return request({ url: `/admin/articles/${id}/toggle-top`, method: 'post', data: {} })
}

// ---------- 分类 ----------

export interface CategoryItem {
  id: number
  name: string
  slug: string
  parent_id: number
  parent_name: string
  description: string
  sort: number
  is_active: boolean
  article_count: number
  created_at: string
}

export interface CategoryPayload {
  name: string
  slug?: string
  parent_id?: number
  description?: string
  sort?: number
  is_active?: boolean
}

export function getCategories(params: { keyword?: string; is_active?: boolean } = {}) {
  return request({ url: '/admin/categories', method: 'get', params }) as Promise<{ list: CategoryItem[] }>
}

export function createCategory(data: CategoryPayload) {
  return request({ url: '/admin/categories', method: 'post', data })
}

export function updateCategory(id: number, data: CategoryPayload) {
  return request({ url: `/admin/categories/${id}`, method: 'put', data })
}

export function deleteCategory(id: number) {
  return request({ url: `/admin/categories/${id}`, method: 'delete' })
}

// ---------- 标签 ----------

export interface TagItem {
  id: number
  name: string
  slug: string
  description: string
  sort: number
  is_active: boolean
  article_count: number
  created_at: string
}

export interface TagPayload {
  name: string
  slug?: string
  description?: string
  sort?: number
  is_active?: boolean
}

export function getTags(params: { keyword?: string; is_active?: boolean } = {}) {
  return request({ url: '/admin/tags', method: 'get', params }) as Promise<{ list: TagItem[] }>
}

export function createTag(data: TagPayload) {
  return request({ url: '/admin/tags', method: 'post', data })
}

export function updateTag(id: number, data: TagPayload) {
  return request({ url: `/admin/tags/${id}`, method: 'put', data })
}

export function deleteTag(id: number) {
  return request({ url: `/admin/tags/${id}`, method: 'delete' })
}

// ---------- 批量操作 ----------

export interface ArticleBatchResp {
  count: number
}

export function batchOfflineArticles(ids: number[]) {
  return request({
    url: '/admin/articles/batch-offline',
    method: 'post',
    data: { ids },
  }) as Promise<ArticleBatchResp>
}

export function batchDeleteArticles(ids: number[]) {
  return request({
    url: '/admin/articles/batch-delete',
    method: 'post',
    data: { ids },
  }) as Promise<ArticleBatchResp>
}

// ---------- 审核记录 ----------

export interface ArticleReviewItem {
  id: number
  article_id: number
  action: number
  action_text: string
  remark: string
  operator_id: number
  operator_name: string
  created_at: string
}

export function getArticleReviews(id: number) {
  return request({ url: `/admin/articles/${id}/reviews`, method: 'get' }) as Promise<{
    list: ArticleReviewItem[]
  }>
}

// ---------- 审核流 ----------

export function submitArticle(id: number) {
  return request({ url: `/admin/articles/${id}/submit`, method: 'post', data: {} })
}

export function approveArticle(id: number, published_at?: string) {
  return request({
    url: `/admin/articles/${id}/approve`,
    method: 'post',
    data: { published_at: published_at ?? '' },
  })
}

export function rejectArticle(id: number, remark = '') {
  return request({ url: `/admin/articles/${id}/reject`, method: 'post', data: { remark } })
}

// ---------- 修订历史 ----------

export interface ArticleRevisionItem {
  id: number
  article_id: number
  title: string
  remark: string
  created_at: string
}

export interface ArticleRevisionDetail extends ArticleRevisionItem {
  summary: string
  cover: string
  content: string
}

export function getArticleRevisions(id: number) {
  return request({ url: `/admin/articles/${id}/revisions`, method: 'get' }) as Promise<{
    list: ArticleRevisionItem[]
  }>
}

export function getArticleRevision(id: number, rid: number) {
  return request({
    url: `/admin/articles/${id}/revisions/${rid}`,
    method: 'get',
  }) as Promise<ArticleRevisionDetail>
}

export function restoreArticleRevision(id: number, rid: number) {
  return request({ url: `/admin/articles/${id}/revisions/${rid}/restore`, method: 'post', data: {} })
}

// ---------- 回收站 ----------

export type TrashEntityType = 'article' | 'category' | 'tag'

export interface TrashItem {
  id: number
  entity_type: TrashEntityType
  title: string
  slug: string
  deleted_at: string
}

export interface TrashListResp {
  list: TrashItem[]
  total: number
  page: number
  page_size: number
}

export function getTrash(params: { type?: TrashEntityType; page?: number; page_size?: number } = {}) {
  return request({ url: '/admin/trash', method: 'get', params }) as Promise<TrashListResp>
}

export function restoreTrash(type: TrashEntityType, id: number) {
  return request({ url: `/admin/trash/${type}/${id}/restore`, method: 'post', data: {} }) as Promise<{
    id: number
    message: string
  }>
}

export function purgeTrash(type: TrashEntityType, id: number) {
  return request({ url: `/admin/trash/${type}/${id}/purge`, method: 'delete' })
}
