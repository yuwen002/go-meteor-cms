// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/sysdict"
	"github.com/yuwen002/go-meteor-cms/ent/sysdictitem"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictLogic {
	return &DictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DictList 分页查询字典类型，支持按编码/名称模糊搜索与状态过滤
func (l *DictLogic) DictList(req *types.DictListReq) (resp *types.DictListResp, err error) {
	query := l.svcCtx.EntClient.SysDict.Query()
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		query = query.Where(sysdict.Or(
			sysdict.CodeContainsFold(kw),
			sysdict.NameContainsFold(kw),
		))
	}
	if req.Status != nil {
		query = query.Where(sysdict.StatusEQ(*req.Status))
	}

	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		logx.Errorf("count dict failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	rows, err := query.
		Order(ent.Asc(sysdict.FieldSort), ent.Asc(sysdict.FieldID)).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list dict failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	list := make([]types.DictItem, 0, len(rows))
	for _, d := range rows {
		list = append(list, types.DictItem{
			ID:          d.ID,
			Code:        d.Code,
			Name:        d.Name,
			Description: d.Description,
			Status:      d.Status,
			Sort:        d.Sort,
			CreatedAt:   d.CreatedAt.Format(time.RFC3339),
		})
	}
	return &types.DictListResp{List: list, Total: int64(total)}, nil
}

// DictCreate 新增字典类型，编码全局唯一
func (l *DictLogic) DictCreate(req *types.DictCreateReq) error {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return common.NewBizError(common.ErrInvalidParams)
	}
	exists, err := l.svcCtx.EntClient.SysDict.
		Query().
		Where(sysdict.CodeEQ(code)).
		Exist(l.ctx)
	if err != nil {
		logx.Errorf("check dict code failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}
	if exists {
		return common.NewBizError(common.ErrDictCodeExists)
	}

	status := req.Status
	if status == 0 {
		status = 1
	}
	if _, err := l.svcCtx.EntClient.SysDict.
		Create().
		SetCode(code).
		SetName(name).
		SetDescription(req.Description).
		SetStatus(status).
		SetSort(req.Sort).
		Save(l.ctx); err != nil {
		logx.Errorf("create dict failed: %v", err)
		return common.NewBizError(common.ErrDictCreateFail)
	}
	return nil
}

// DictUpdate 更新字典类型，仅覆盖传入的字段
func (l *DictLogic) DictUpdate(req *types.DictUpdateReq) error {
	d, err := l.svcCtx.EntClient.SysDict.Get(l.ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return common.NewBizError(common.ErrDictNotFound)
		}
		logx.Errorf("get dict failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}

	update := d.Update()
	if name := strings.TrimSpace(req.Name); name != "" {
		update = update.SetName(name)
	}
	if req.Description != "" {
		update = update.SetDescription(req.Description)
	}
	if req.Status != nil {
		update = update.SetStatus(*req.Status)
	}
	if req.Sort != nil {
		update = update.SetSort(*req.Sort)
	}
	if _, err := update.Save(l.ctx); err != nil {
		logx.Errorf("update dict failed: %v", err)
		return common.NewBizError(common.ErrDictUpdateFail)
	}
	return nil
}

// DictDelete 删除字典类型，并级联清理其下的字典项
func (l *DictLogic) DictDelete(req *types.DictDetailReq) error {
	d, err := l.svcCtx.EntClient.SysDict.Get(l.ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return common.NewBizError(common.ErrDictNotFound)
		}
		logx.Errorf("get dict failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}

	if _, err := l.svcCtx.EntClient.SysDictItem.
		Delete().
		Where(sysdictitem.DictCodeEQ(d.Code)).
		Exec(l.ctx); err != nil {
		logx.Errorf("delete dict items failed: %v", err)
		return common.NewBizError(common.ErrDictDeleteFail)
	}
	if err := l.svcCtx.EntClient.SysDict.DeleteOneID(req.ID).Exec(l.ctx); err != nil {
		logx.Errorf("delete dict failed: %v", err)
		return common.NewBizError(common.ErrDictDeleteFail)
	}
	return nil
}
