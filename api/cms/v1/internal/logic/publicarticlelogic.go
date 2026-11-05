package logic

import (
	"context"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/article"
	"github.com/yuwen002/go-meteor-cms/ent/category"
	"github.com/yuwen002/go-meteor-cms/ent/tag"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type PublicArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicArticleLogic {
	return &PublicArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// visibleQuery 公开可见条件：状态已发布且发布时间已到。
// 定时发布只需在发布时把 published_at 设为未来时间，到点后本查询自动放行。
func (l *PublicArticleLogic) visibleQuery() *ent.ArticleQuery {
	return l.svcCtx.EntClient.Article.Query().
		Where(
			article.StatusEQ(ArticleStatusPublished),
			article.PublishedAtLTE(time.Now()),
		)
}

// PublicArticleList 前台文章列表，不返回正文
func (l *PublicArticleLogic) PublicArticleList(req *types.PublicArticleListReq) (*types.PublicArticleListResp, error) {
	page, pageSize := utils.NormalizePagination(req.Page, req.PageSize, 10)
	offset, limit := utils.NormalizePage(req.Page, req.PageSize, 10)
	empty := &types.PublicArticleListResp{
		List:     []types.PublicArticleItem{},
		Total:    0,
		Page:     page,
		PageSize: pageSize,
	}

	query := l.visibleQuery()
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		query = query.Where(article.Or(
			article.TitleContainsFold(kw),
			article.SummaryContainsFold(kw),
		))
	}
	if slug := strings.TrimSpace(req.CategorySlug); slug != "" {
		c, err := l.svcCtx.EntClient.Category.Query().
			Where(category.SlugEQ(slug), category.IsActiveEQ(true)).
			Only(l.ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return empty, nil
			}
			logx.Errorf("query public category failed: %v", err)
			return nil, common.NewBizError(common.ErrInternalServer)
		}
		query = query.Where(article.CategoryIDEQ(c.ID))
	}
	if slug := strings.TrimSpace(req.TagSlug); slug != "" {
		query = query.Where(article.HasTagsWith(tag.SlugEQ(slug)))
	}

	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		logx.Errorf("count public articles failed: %v", err)
		return nil, common.NewBizError(common.ErrArticleListFail)
	}
	rows, err := query.
		WithCategory().
		WithTags().
		Order(ent.Desc(article.FieldIsTop), ent.Desc(article.FieldPublishedAt), ent.Desc(article.FieldID)).
		Limit(limit).
		Offset(offset).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list public articles failed: %v", err)
		return nil, common.NewBizError(common.ErrArticleListFail)
	}

	list := make([]types.PublicArticleItem, 0, len(rows))
	for _, a := range rows {
		list = append(list, publicArticleToItem(a))
	}
	return &types.PublicArticleListResp{
		List: list, Total: int64(total), Page: page, PageSize: pageSize,
	}, nil
}

// PublicArticleDetail 前台文章详情，按 slug 访问，返回 Markdown 正文并累加浏览量
func (l *PublicArticleLogic) PublicArticleDetail(req *types.PublicArticleDetailReq) (*types.PublicArticleDetailResp, error) {
	a, err := l.visibleQuery().
		Where(article.SlugEQ(strings.TrimSpace(req.Slug))).
		WithCategory().
		WithTags().
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(common.ErrArticleNotFound)
		}
		logx.Errorf("get public article failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	// 浏览量是弱指标，写入失败不应影响阅读
	if _, err := a.Update().SetViewCount(a.ViewCount + 1).Save(l.ctx); err != nil {
		logx.Errorf("increase article view count failed: %v", err)
	}
	return &types.PublicArticleDetailResp{
		PublicArticleItem: publicArticleToItem(a),
		Content:           a.Content,
		SeoTitle:          a.SeoTitle,
		SeoDescription:    a.SeoDescription,
		SeoKeywords:       a.SeoKeywords,
	}, nil
}

// PublicTaxonomy 前台分类与标签导航，只返回至少有一篇公开文章的条目
func (l *PublicArticleLogic) PublicTaxonomy() (*types.PublicTaxonomyResp, error) {
	cats, err := l.svcCtx.EntClient.Category.Query().
		Where(category.IsActiveEQ(true)).
		Order(ent.Asc(category.FieldSort), ent.Asc(category.FieldID)).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list public categories failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}
	catItems := make([]types.PublicCategoryItem, 0, len(cats))
	for _, c := range cats {
		n, err := l.visibleQuery().Where(article.CategoryIDEQ(c.ID)).Count(l.ctx)
		if err != nil {
			logx.Errorf("count public category articles failed: %v", err)
			return nil, common.NewBizError(common.ErrInternalServer)
		}
		if n == 0 {
			continue
		}
		item := types.PublicCategoryItem{
			ID:           c.ID,
			Name:         c.Name,
			Slug:         c.Slug,
			ArticleCount: int64(n),
		}
		if c.ParentID != nil {
			item.ParentID = *c.ParentID
		}
		catItems = append(catItems, item)
	}

	tags, err := l.svcCtx.EntClient.Tag.Query().
		Where(tag.IsActiveEQ(true)).
		Order(ent.Asc(tag.FieldSort), ent.Asc(tag.FieldID)).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list public tags failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}
	tagItems := make([]types.PublicTagItem, 0, len(tags))
	for _, tg := range tags {
		n, err := l.visibleQuery().
			Where(article.HasTagsWith(tag.IDEQ(tg.ID))).
			Count(l.ctx)
		if err != nil {
			logx.Errorf("count public tag articles failed: %v", err)
			return nil, common.NewBizError(common.ErrInternalServer)
		}
		if n == 0 {
			continue
		}
		tagItems = append(tagItems, types.PublicTagItem{
			ID:           tg.ID,
			Name:         tg.Name,
			Slug:         tg.Slug,
			ArticleCount: int64(n),
		})
	}

	return &types.PublicTaxonomyResp{Categories: catItems, Tags: tagItems}, nil
}

func publicArticleToItem(a *ent.Article) types.PublicArticleItem {
	item := types.PublicArticleItem{
		ID:          a.ID,
		Title:       a.Title,
		Slug:        a.Slug,
		Summary:     a.Summary,
		Cover:       a.Cover,
		TagNames:    make([]string, 0, len(a.Edges.Tags)),
		IsTop:       a.IsTop,
		ViewCount:   a.ViewCount,
		PublishedAt: FormatOptionalTime(a.PublishedAt),
	}
	if a.CategoryID != nil {
		item.CategoryID = *a.CategoryID
	}
	if a.Edges.Category != nil {
		item.CategoryName = a.Edges.Category.Name
		item.CategorySlug = a.Edges.Category.Slug
	}
	for _, tg := range a.Edges.Tags {
		item.TagNames = append(item.TagNames, tg.Name)
	}
	return item
}
