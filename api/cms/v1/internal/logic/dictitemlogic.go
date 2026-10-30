// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"strings"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/sysdictitem"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemLogic {
	return &DictItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DictItemList 查询指定字典下的字典项，按 sort、id 升序
func (l *DictItemLogic) DictItemList(req *types.DictItemListReq) (resp *types.DictItemListResp, err error) {
	dictCode := strings.TrimSpace(req.DictCode)
	if dictCode == "" {
		return nil, common.NewBizError(common.ErrInvalidParams)
	}

	query := l.svcCtx.EntClient.SysDictItem.Query().
		Where(sysdictitem.DictCodeEQ(dictCode))
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		query = query.Where(sysdictitem.Or(
			sysdictitem.LabelContainsFold(kw),
			sysdictitem.ValueContainsFold(kw),
		))
	}
	if req.Status != nil {
		query = query.Where(sysdictitem.StatusEQ(*req.Status))
	}

	rows, err := query.
		Order(ent.Asc(sysdictitem.FieldSort), ent.Asc(sysdictitem.FieldID)).
		All(l.ctx)
	if err != nil {
		logx.Errorf("list dict items failed: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	list := make([]types.DictDataItem, 0, len(rows))
	for _, it := range rows {
		list = append(list, types.DictDataItem{
			ID:       it.ID,
			DictCode: it.DictCode,
			Label:    it.Label,
			Value:    it.Value,
			Status:   it.Status,
			Sort:     it.Sort,
		})
	}
	return &types.DictItemListResp{List: list}, nil
}

// DictItemCreate 新增字典项，同一字典下 value 不允许重复
func (l *DictItemLogic) DictItemCreate(req *types.DictItemCreateReq) error {
	dictCode := strings.TrimSpace(req.DictCode)
	label := strings.TrimSpace(req.Label)
	value := strings.TrimSpace(req.Value)
	if dictCode == "" || label == "" || value == "" {
		return common.NewBizError(common.ErrInvalidParams)
	}

	exists, err := l.svcCtx.EntClient.SysDictItem.
		Query().
		Where(
			sysdictitem.DictCodeEQ(dictCode),
			sysdictitem.ValueEQ(value),
		).
		Exist(l.ctx)
	if err != nil {
		logx.Errorf("check dict item failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}
	if exists {
		return common.NewBizError(common.ErrDictItemExists)
	}

	status := req.Status
	if status == 0 {
		status = 1
	}
	if _, err := l.svcCtx.EntClient.SysDictItem.
		Create().
		SetDictCode(dictCode).
		SetLabel(label).
		SetValue(value).
		SetStatus(status).
		SetSort(req.Sort).
		Save(l.ctx); err != nil {
		logx.Errorf("create dict item failed: %v", err)
		return common.NewBizError(common.ErrDictCreateFail)
	}
	return nil
}

// DictItemUpdate 更新字典项，仅覆盖传入的字段
func (l *DictItemLogic) DictItemUpdate(req *types.DictItemUpdateReq) error {
	it, err := l.svcCtx.EntClient.SysDictItem.Get(l.ctx, req.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return common.NewBizError(common.ErrDictItemNotFound)
		}
		logx.Errorf("get dict item failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}

	update := it.Update()
	if label := strings.TrimSpace(req.Label); label != "" {
		update = update.SetLabel(label)
	}
	if value := strings.TrimSpace(req.Value); value != "" {
		update = update.SetValue(value)
	}
	if req.Status != nil {
		update = update.SetStatus(*req.Status)
	}
	if req.Sort != nil {
		update = update.SetSort(*req.Sort)
	}
	if _, err := update.Save(l.ctx); err != nil {
		logx.Errorf("update dict item failed: %v", err)
		return common.NewBizError(common.ErrDictUpdateFail)
	}
	return nil
}

// DictItemDelete 删除字典项
func (l *DictItemLogic) DictItemDelete(req *types.DictItemDetailReq) error {
	exists, err := l.svcCtx.EntClient.SysDictItem.
		Query().
		Where(sysdictitem.IDEQ(req.ID)).
		Exist(l.ctx)
	if err != nil {
		logx.Errorf("check dict item failed: %v", err)
		return common.NewBizError(common.ErrInternalServer)
	}
	if !exists {
		return common.NewBizError(common.ErrDictItemNotFound)
	}
	if err := l.svcCtx.EntClient.SysDictItem.DeleteOneID(req.ID).Exec(l.ctx); err != nil {
		logx.Errorf("delete dict item failed: %v", err)
		return common.NewBizError(common.ErrDictDeleteFail)
	}
	return nil
}
