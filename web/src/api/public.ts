import request from '@/utils/request'

export interface PublicArticle {
  id: number
  title: string
  slug: string
  summary: string
  cover: string
  author_name: string
  category_id: number
  category_name: string
  category_slug: string
  tag_names: string[]
  is_top: boolean
  view_count: number
  published_at: string
}

export interface PublicArticleListResp {
  list: PublicArticle[]
  total: number
  page: number
  page_size: number
}

export interface PublicArticleDetail extends PublicArticle {
  content: string
  seo_title: string
  seo_description: string
  seo_keywords: string
  canonical_url: string
  related: PublicArticle[]
  prev: PublicArticle | null
  next: PublicArticle | null
}

export interface PublicCategory {
  id: number
  name: string
  slug: string
  parent_id: number
  article_count: number
}

export interface PublicTag {
  id: number
  name: string
  slug: string
  article_count: number
}

export interface PublicTaxonomyResp {
  categories: PublicCategory[]
  tags: PublicTag[]
}

export interface PublicArticleParams {
  page?: number
  page_size?: number
  keyword?: string
  category?: string
  tag?: string
}

export function getPublicArticles(params: PublicArticleParams) {
  return request({ url: '/api/articles', method: 'get', params }) as Promise<PublicArticleListResp>
}

export function getPublicArticle(slug: string) {
  return request({
    url: `/api/articles/${encodeURIComponent(slug)}`,
    method: 'get',
  }) as Promise<PublicArticleDetail>
}

export function getPublicTaxonomy() {
  return request({ url: '/api/taxonomy', method: 'get' }) as Promise<PublicTaxonomyResp>
}
