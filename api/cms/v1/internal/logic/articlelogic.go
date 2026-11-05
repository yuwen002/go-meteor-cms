package logic

import (
	"context"
	"fmt"
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

type ArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ArticleList 后台文章列表：支持关键字、状态、分类、标签筛选与分页
func (l *ArticleLogic) ArticleList(req *types.ArticleListReq) (*types.ArticleListResp, error) {
	query := l.svcCtx.EntClient.Article.Query()
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		query = query.Where(article.Or(
			article.TitleContainsFold(kw),
			article.SummaryContainsFold(kw),
		))
	}
	if req.Status != nil {
		query = query.Where(article.StatusEQ(*req.Status))
	}
	if req.CategoryID > 0 {
		query = query.Where(article.CategoryIDEQ(req.CategoryID))
	}
	if req.TagID > 0 {
		query = query.Where(article.HasTagsWith(tag.IDEQ(req.TagID)))
	}

	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		logx.Errorf("count articles failed: %v", err)
		return nil, common.NewBizError(common.ErrArticleListFail)
	}

	page, pageSize := utils.NormalizePagination(req.Page, req.PageSize, 20)
	offset, limit := utils.NormalizePage(req.Page, req.PageSize, 20)
	rows, err := query.
		WithCategory().
		WithTags().
		Order(ent.Desc(article.FieldIsTop), ent.Desc(article.FieldID)).
		Limit(limit).
		Offset(offset).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list articles failed: %v", err)
		return nil, common.NewBizError(common.ErrArticleListFail)
	}

	list := make([]types.ArticleItem, 0, len(rows))
	for _, a := range rows {
		list = append(list, articleToItem(a))
	}
	return &types.ArticleListResp{
		List: list, Total: int64(total), Page: page, PageSize: pageSize,
	}, nil
}

// ArticleDetail 文章详情，含 Markdown 正文与 SEO 字段
func (l *ArticleLogic) ArticleDetail(req *types.ArticleIDReq) (*types.ArticleDetailResp, error) {
	a, err := l.svcCtx.EntClient.Article.Query().
		Where(article.IDEQ(req.ID)).
		WithCategory().
		WithTags().
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(common.ErrArticleNotFound)
		}
		logx.Errorf("get article failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}
	return &types.ArticleDetailResp{
		ArticleItem:    articleToItem(a),
		Content:        a.Content,
		SeoTitle:       a.SeoTitle,
		SeoDescription: a.SeoDescription,
		SeoKeywords:    a.SeoKeywords,
	}, nil
}

// ArticleCreate 新建文章，一律以草稿状态落库，发布必须走 ArticlePublish
func (l *ArticleLogic) ArticleCreate(req *types.ArticleCreateReq) (*types.CommonResp, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, common.NewBizErrorWithMsg(common.ErrInvalidParams, "文章标题不能为空")
	}
	rawSlug := strings.TrimSpace(req.Slug)
	if rawSlug != "" && !IsValidSlug(rawSlug) {
		return nil, common.NewBizError(common.ErrArticleInvalidSlug)
	}
	slug := GenerateSlug(rawSlug, title, time.Now())
	if err := l.ensureArticleSlugFree(slug, 0); err != nil {
		return nil, err
	}
	tagIDs, err := l.validateArticleTaxonomy(req.CategoryID, req.TagIDs)
	if err != nil {
		return nil, err
	}

	create := l.svcCtx.EntClient.Article.Create().
		SetTitle(title).
		SetSlug(slug).
		SetSummary(strings.TrimSpace(req.Summary)).
		SetCover(strings.TrimSpace(req.Cover)).
		SetContent(req.Content).
		SetStatus(ArticleStatusDraft).
		SetIsTop(req.IsTop).
		SetSeoTitle(strings.TrimSpace(req.SeoTitle)).
		SetSeoDescription(strings.TrimSpace(req.SeoDescription)).
		SetSeoKeywords(strings.TrimSpace(req.SeoKeywords))
	if req.CategoryID > 0 {
		create = create.SetCategoryID(req.CategoryID)
	}
	if userID, ok := utils.CurrentUserID(l.ctx); ok {
		create = create.SetAuthorID(userID)
	}
	if len(tagIDs) > 0 {
		create = create.AddTagIDs(tagIDs...)
	}

	a, err := create.Save(l.ctx)
	if err != nil {
		logx.Errorf("create article failed: %v", err)
		return nil, common.NewBizError(common.ErrArticleCreateFail)
	}
	return &types.CommonResp{ID: a.ID, Message: "文章已保存为草稿"}, nil
}

// ArticleUpdate 编辑文章：整单覆盖除状态与发布时间以外的字段
func (l *ArticleLogic) ArticleUpdate(req *types.ArticleUpdateReq) error {
	a, err := l.svcCtx.EntClient.Article.Get(l.ctx, req.ID)
	if err != nil {
		return articleGetError(err)
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return common.NewBizErrorWithMsg(common.ErrInvalidParams, "文章标题不能为空")
	}

	slug := a.Slug
	if raw := strings.TrimSpace(req.Slug); raw != "" {
		if !IsValidSlug(raw) {
			return common.NewBizError(common.ErrArticleInvalidSlug)
		}
		slug = raw
		if err := l.ensureArticleSlugFree(slug, req.ID); err != nil {
			return err
		}
	}
	tagIDs, err := l.validateArticleTaxonomy(req.CategoryID, req.TagIDs)
	if err != nil {
		return err
	}

	update := a.Update().
		SetTitle(title).
		SetSlug(slug).
		SetSummary(strings.TrimSpace(req.Summary)).
		SetCover(strings.TrimSpace(req.Cover)).
		SetContent(req.Content).
		SetIsTop(req.IsTop).
		SetSeoTitle(strings.TrimSpace(req.SeoTitle)).
		SetSeoDescription(strings.TrimSpace(req.SeoDescription)).
		SetSeoKeywords(strings.TrimSpace(req.SeoKeywords))
	if req.CategoryID > 0 {
		update = update.SetCategoryID(req.CategoryID)
	} else {
		update = update.ClearCategory()
	}
	update = update.ClearTags()
	if len(tagIDs) > 0 {
		update = update.AddTagIDs(tagIDs...)
	}
	if _, err := update.Save(l.ctx); err != nil {
		logx.Errorf("update article failed: %v", err)
		return common.NewBizError(common.ErrArticleUpdateFail)
	}
	return nil
}

// ArticleDelete 软删除文章。软删除不会触发数据库唯一索引释放，
// 因此先把 slug 改成一个一次性的值，避免同名文章永远无法再创建。
func (l *ArticleLogic) ArticleDelete(req *types.ArticleIDReq) error {
	a, err := l.svcCtx.EntClient.Article.Get(l.ctx, req.ID)
	if err != nil {
		return articleGetError(err)
	}
	if _, err := a.Update().
		SetSlug(fmt.Sprintf("deleted-%d-%d", a.ID, time.Now().UnixNano())).
		Save(l.ctx); err != nil {
		logx.Errorf("release article slug failed: %v", err)
		return common.NewBizError(common.ErrArticleDeleteFail)
	}
	if err := l.svcCtx.EntClient.Article.DeleteOneID(req.ID).Exec(l.ctx); err != nil {
		logx.Errorf("delete article failed: %v", err)
		return common.NewBizError(common.ErrArticleDeleteFail)
	}
	return nil
}

// ArticlePublish 发布/定时发布。传未来时间即定时发布，到点自动对前台可见。
func (l *ArticleLogic) ArticlePublish(req *types.ArticlePublishReq) error {
	a, err := l.svcCtx.EntClient.Article.Get(l.ctx, req.ID)
	if err != nil {
		return articleGetError(err)
	}
	if strings.TrimSpace(a.Content) == "" {
		return common.NewBizErrorWithMsg(common.ErrInvalidParams, "正文为空，不能发布")
	}
	publishAt, err := ParsePublishTime(req.PublishedAt, time.Now())
	if err != nil {
		return common.NewBizError(common.ErrArticleInvalidTime)
	}
	if _, err := a.Update().
		SetStatus(ArticleStatusPublished).
		SetPublishedAt(publishAt).
		Save(l.ctx); err != nil {
		logx.Errorf("publish article failed: %v", err)
		return common.NewBizError(common.ErrArticlePublishFail)
	}
	return nil
}

// ArticleWithdraw 撤回（下线）文章，保留内容供重新编辑后再发布
func (l *ArticleLogic) ArticleWithdraw(req *types.ArticleIDReq) error {
	a, err := l.svcCtx.EntClient.Article.Get(l.ctx, req.ID)
	if err != nil {
		return articleGetError(err)
	}
	if _, err := a.Update().SetStatus(ArticleStatusOffline).Save(l.ctx); err != nil {
		logx.Errorf("withdraw article failed: %v", err)
		return common.NewBizError(common.ErrArticleWithdrawFail)
	}
	return nil
}

// ArticleToggleTop 切换置顶状态
func (l *ArticleLogic) ArticleToggleTop(req *types.ArticleIDReq) error {
	a, err := l.svcCtx.EntClient.Article.Get(l.ctx, req.ID)
	if err != nil {
		return articleGetError(err)
	}
	if _, err := a.Update().SetIsTop(!a.IsTop).Save(l.ctx); err != nil {
		logx.Errorf("toggle article top failed: %v", err)
		return common.NewBizError(common.ErrArticleSetTopFail)
	}
	return nil
}

func articleGetError(err error) error {
	if ent.IsNotFound(err) {
		return common.NewBizError(common.ErrArticleNotFound)
	}
	logx.Errorf("get article failed: %v", err)
	return common.NewBizError(common.ErrInternalServer)
}

func articleToItem(a *ent.Article) types.ArticleItem {
	item := types.ArticleItem{
		ID:          a.ID,
		Title:       a.Title,
		Slug:        a.Slug,
		Summary:     a.Summary,
		Cover:       a.Cover,
		Status:      a.Status,
		StatusText:  ArticleStatusText(a.Status),
		TagIDs:      make([]int64, 0, len(a.Edges.Tags)),
		TagNames:    make([]string, 0, len(a.Edges.Tags)),
		IsTop:       a.IsTop,
		ViewCount:   a.ViewCount,
		PublishedAt: FormatOptionalTime(a.PublishedAt),
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}
	if a.CategoryID != nil {
		item.CategoryID = *a.CategoryID
	}
	if a.Edges.Category != nil {
		item.CategoryName = a.Edges.Category.Name
	}
	for _, tg := range a.Edges.Tags {
		item.TagIDs = append(item.TagIDs, tg.ID)
		item.TagNames = append(item.TagNames, tg.Name)
	}
	return item
}

// ensureArticleSlugFree 校验 slug 未被其它文章占用（软删除的文章已被拦截器过滤）
func (l *ArticleLogic) ensureArticleSlugFree(slug string, excludeID int64) error {
	query := l.svcCtx.EntClient.Article.Query().Where(article.SlugEQ(slug))
	if excludeID > 0 {
		query = query.Where(article.IDNEQ(excludeID))
	}
	exists, err := query.Exist(l.ctx)
	if err != nil {
		logx.Errorf("check article slug failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}
	if exists {
		return common.NewBizError(common.ErrArticleSlugExists)
	}
	return nil
}

// validateArticleTaxonomy 校验分类与标签存在且启用，并返回去重后的标签 ID 列表
func (l *ArticleLogic) validateArticleTaxonomy(categoryID int64, tagIDs []int64) ([]int64, error) {
	if categoryID > 0 {
		ok, err := l.svcCtx.EntClient.Category.Query().
			Where(category.IDEQ(categoryID), category.IsActiveEQ(true)).
			Exist(l.ctx)
		if err != nil {
			logx.Errorf("check article category failed: %v", err)
			return nil, common.NewBizError(common.ErrInternalServer)
		}
		if !ok {
			return nil, common.NewBizError(common.ErrArticleCategoryInvalid)
		}
	}

	unique := make([]int64, 0, len(tagIDs))
	seen := make(map[int64]bool, len(tagIDs))
	for _, id := range tagIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}
	if len(unique) == 0 {
		return unique, nil
	}
	count, err := l.svcCtx.EntClient.Tag.Query().
		Where(tag.IDIn(unique...), tag.IsActiveEQ(true)).
		Count(l.ctx)
	if err != nil {
		logx.Errorf("check article tags failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}
	if count != len(unique) {
		return nil, common.NewBizError(common.ErrArticleTagInvalid)
	}
	return unique, nil
}
