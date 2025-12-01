// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDisableAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableAdminLogic {
	return &DisableAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DisableAdminLogic) DisableAdmin(req *types.DisableAdminReq) (resp *types.CommonResp, err error) {
	// 1. 检查管理员是否存在
	_, err = l.svcCtx.EntClient.AdminUser.Query().Where(
		adminuser.IDEQ(req.ID),
	).Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(common.ErrAdminUserNotFound)
		}
		l.Logger.Errorf("查询管理员失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	// 2. 更新管理员状态为禁用
	_, err = l.svcCtx.EntClient.AdminUser.UpdateOneID(req.ID).
		SetIsActive(false).
		Save(l.ctx)
	if err != nil {
		l.Logger.Errorf("禁用管理员失败: %v", err)
		return nil, common.NewBizErrorWithMsg(common.ErrInternalServer, "禁用管理员失败")
	}

	return &types.CommonResp{
		ID:      req.ID,
		Message: "禁用管理员成功",
	}, nil
}
