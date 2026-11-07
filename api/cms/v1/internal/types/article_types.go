package types

// 文章 / 分类 / 标签 / 公开阅读的请求与响应类型。
//
// 本文件独立于 goctl 生成的 types.go：项目里 types.go 与 routes.go 都带有
// "DO NOT EDIT" 头但被手工追加过内容，重跑 make gen-api 会冲掉这些追加段。
// 新模块统一放这里，保证后续重新生成不会丢失定义。

// ---------- 文章管理 ----------

type ArticleItem struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Summary      string   `json:"summary"`
	Cover        string   `json:"cover"`
	Status       int      `json:"status"`
	StatusText   string   `json:"status_text"`
	AuthorName   string   `json:"author_name"`
	CategoryID   int64    `json:"category_id"`
	CategoryName string   `json:"category_name"`
	TagIDs       []int64  `json:"tag_ids"`
	TagNames     []string `json:"tag_names"`
	IsTop        bool     `json:"is_top"`
	ViewCount    int64    `json:"view_count"`
	PublishedAt  string   `json:"published_at"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// ArticleListReq GET 查询串必须使用 form 标签
type ArticleListReq struct {
	Page       int    `form:"page,optional"`
	PageSize   int    `form:"page_size,optional"`
	Keyword    string `form:"keyword,optional"`
	Status     *int   `form:"status,optional"`
	CategoryID int64  `form:"category_id,optional"`
	TagID      int64  `form:"tag_id,optional"`
	// Sort/Order 走白名单，默认「置顶优先 + ID 倒序」
	Sort  string `form:"sort,optional"`
	Order string `form:"order,optional"`
}

// ArticleBatchReq 批量操作的文章 ID 列表，服务端会去重并限制单次数量上限。
type ArticleBatchReq struct {
	IDs []int64 `json:"ids"`
}

type ArticleBatchResp struct {
	Count int `json:"count"`
}

type ArticleListResp struct {
	List     []ArticleItem `json:"list"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type ArticleIDReq struct {
	ID int64 `path:"id"`
}

type ArticleDetailResp struct {
	ArticleItem
	Content        string `json:"content"`
	SeoTitle       string `json:"seo_title"`
	SeoDescription string `json:"seo_description"`
	SeoKeywords    string `json:"seo_keywords"`
}

type ArticleCreateReq struct {
	Title          string  `json:"title"`
	Slug           string  `json:"slug,optional"`
	Summary        string  `json:"summary,optional"`
	Cover          string  `json:"cover,optional"`
	Content        string  `json:"content,optional"`
	CategoryID     int64   `json:"category_id,optional"`
	TagIDs         []int64 `json:"tag_ids,optional"`
	IsTop          bool    `json:"is_top,optional"`
	SeoTitle       string  `json:"seo_title,optional"`
	SeoDescription string  `json:"seo_description,optional"`
	SeoKeywords    string  `json:"seo_keywords,optional"`
}

// ArticleUpdateReq 编辑页整单提交，字段全量覆盖；状态只能通过发布/撤回/审核接口变更。
type ArticleUpdateReq struct {
	ID             int64   `path:"id"`
	Title          string  `json:"title"`
	Slug           string  `json:"slug,optional"`
	Summary        string  `json:"summary,optional"`
	Cover          string  `json:"cover,optional"`
	Content        string  `json:"content,optional"`
	CategoryID     int64   `json:"category_id,optional"`
	TagIDs         []int64 `json:"tag_ids,optional"`
	IsTop          bool    `json:"is_top,optional"`
	SeoTitle       string  `json:"seo_title,optional"`
	SeoDescription string  `json:"seo_description,optional"`
	SeoKeywords    string  `json:"seo_keywords,optional"`
}

// ArticlePublishReq PublishedAt 为空表示立即发布；填未来时间即定时发布，
// 到点后由公开接口的 published_at <= now 条件自动放行，无需后台定时任务。
type ArticlePublishReq struct {
	ID          int64  `path:"id"`
	PublishedAt string `json:"published_at,optional"`
}

// ArticleRejectReq 审核驳回，退回草稿
type ArticleRejectReq struct {
	ID     int64  `path:"id"`
	Remark string `json:"remark,optional"`
}

// ArticleReviewItem 审核动作留痕，action: 1 提交 / 2 通过 / 3 驳回
type ArticleReviewItem struct {
	ID           int64  `json:"id"`
	ArticleID    int64  `json:"article_id"`
	Action       int    `json:"action"`
	ActionText   string `json:"action_text"`
	Remark       string `json:"remark"`
	OperatorID   int64  `json:"operator_id"`
	OperatorName string `json:"operator_name"`
	CreatedAt    string `json:"created_at"`
}

type ArticleReviewListResp struct {
	List []ArticleReviewItem `json:"list"`
}

// ---------- 修订历史 ----------

type ArticleRevisionItem struct {
	ID        int64  `json:"id"`
	ArticleID int64  `json:"article_id"`
	Title     string `json:"title"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}

type ArticleRevisionListResp struct {
	List []ArticleRevisionItem `json:"list"`
}

type ArticleRevisionDetailResp struct {
	ArticleRevisionItem
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	Content string `json:"content"`
}

type ArticleRevisionIDReq struct {
	ID         int64 `path:"id"`
	RevisionID int64 `path:"rid"`
}

// ---------- 回收站 ----------

type TrashItem struct {
	ID         int64  `json:"id"`
	EntityType string `json:"entity_type"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	DeletedAt  string `json:"deleted_at"`
}

type TrashListReq struct {
	EntityType string `form:"type,optional"`
	Page       int    `form:"page,optional"`
	PageSize   int    `form:"page_size,optional"`
}

type TrashListResp struct {
	List     []TrashItem `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type TrashEntityReq struct {
	EntityType string `path:"type"`
	ID         int64  `path:"id"`
}

// ---------- 分类管理 ----------

type CategoryItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     int64  `json:"parent_id"`
	ParentName   string `json:"parent_name"`
	Description  string `json:"description"`
	Sort         int    `json:"sort"`
	IsActive     bool   `json:"is_active"`
	ArticleCount int    `json:"article_count"`
	CreatedAt    string `json:"created_at"`
}

type CategoryListReq struct {
	Keyword  string `form:"keyword,optional"`
	IsActive *bool  `form:"is_active,optional"`
}

type CategoryListResp struct {
	List []CategoryItem `json:"list"`
}

type CategoryIDReq struct {
	ID int64 `path:"id"`
}

type CategoryCreateReq struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,optional"`
	ParentID    int64  `json:"parent_id,optional"`
	Description string `json:"description,optional"`
	Sort        int    `json:"sort,optional"`
	IsActive    *bool  `json:"is_active,optional"`
}

type CategoryUpdateReq struct {
	ID          int64  `path:"id"`
	Name        string `json:"name,optional"`
	Slug        string `json:"slug,optional"`
	ParentID    *int64 `json:"parent_id,optional"`
	Description string `json:"description,optional"`
	Sort        *int   `json:"sort,optional"`
	IsActive    *bool  `json:"is_active,optional"`
}

// ---------- 标签管理 ----------

type TagItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	Sort         int    `json:"sort"`
	IsActive     bool   `json:"is_active"`
	ArticleCount int    `json:"article_count"`
	CreatedAt    string `json:"created_at"`
}

type TagListReq struct {
	Keyword  string `form:"keyword,optional"`
	IsActive *bool  `form:"is_active,optional"`
}

type TagListResp struct {
	List []TagItem `json:"list"`
}

type TagIDReq struct {
	ID int64 `path:"id"`
}

type TagCreateReq struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,optional"`
	Description string `json:"description,optional"`
	Sort        int    `json:"sort,optional"`
	IsActive    *bool  `json:"is_active,optional"`
}

type TagUpdateReq struct {
	ID          int64  `path:"id"`
	Name        string `json:"name,optional"`
	Slug        string `json:"slug,optional"`
	Description string `json:"description,optional"`
	Sort        *int   `json:"sort,optional"`
	IsActive    *bool  `json:"is_active,optional"`
}

// ---------- 公开阅读（无需登录） ----------

type PublicArticleItem struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Summary      string   `json:"summary"`
	Cover        string   `json:"cover"`
	AuthorName   string   `json:"author_name"`
	CategoryID   int64    `json:"category_id"`
	CategoryName string   `json:"category_name"`
	CategorySlug string   `json:"category_slug"`
	TagNames     []string `json:"tag_names"`
	IsTop        bool     `json:"is_top"`
	ViewCount    int64    `json:"view_count"`
	PublishedAt  string   `json:"published_at"`
}

type PublicArticleListReq struct {
	Page         int    `form:"page,optional"`
	PageSize     int    `form:"page_size,optional"`
	Keyword      string `form:"keyword,optional"`
	CategorySlug string `form:"category,optional"`
	TagSlug      string `form:"tag,optional"`
}

type PublicArticleListResp struct {
	List     []PublicArticleItem `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

type PublicArticleDetailReq struct {
	Slug string `path:"slug"`
}

// PublicArticleDetailResp.Content 是 Markdown 源码，前端渲染前必须做转义。
type PublicArticleDetailResp struct {
	PublicArticleItem
	Content        string              `json:"content"`
	SeoTitle       string              `json:"seo_title"`
	SeoDescription string              `json:"seo_description"`
	SeoKeywords    string              `json:"seo_keywords"`
	CanonicalURL   string              `json:"canonical_url"`
	Related        []PublicArticleItem `json:"related"`
	Prev           *PublicArticleItem  `json:"prev"`
	Next           *PublicArticleItem  `json:"next"`
}

type PublicCategoryItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ParentID     int64  `json:"parent_id"`
	ArticleCount int64  `json:"article_count"`
}

type PublicTagItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int64  `json:"article_count"`
}

type PublicTaxonomyResp struct {
	Categories []PublicCategoryItem `json:"categories"`
	Tags       []PublicTagItem      `json:"tags"`
}
