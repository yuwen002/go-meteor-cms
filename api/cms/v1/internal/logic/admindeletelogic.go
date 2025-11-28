// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteLogic {
	return &AdminDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteLogic) AdminDelete(req *types.DeleteAdminReq) (resp *types.CommonResp, err error) {
	// 检查用户是否存在
	_, err = l.svcCtx.EntClient.AdminUser.Get(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(common.ErrAdminUserNotFound)
		}
		l.Logger.Errorf("查询管理员失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	// 执行软删除：更新 deleted_at 字段
	_, err = l.svcCtx.EntClient.AdminUser.UpdateOneID(req.Id).
		SetDeletedAt(time.Now()).
		Save(l.ctx)

	if err != nil {
		l.Logger.Errorf("删除管理员失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	return &types.CommonResp{
		ID:      req.Id,
		Message: "删除成功",
	}, nil
}
