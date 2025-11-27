// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateLogic {
	return &AdminUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateLogic) AdminUpdate(req *types.UpdateAdminReq) (resp *types.CommonResp, err error) {
	id := req.ID

	_, err = l.svcCtx.EntClient.AdminUser.UpdateOneID(id).
		SetNickname(req.Nickname).
		SetEmail(req.Email).
		SetPhone(req.Phone).
		SetIsActive(req.IsActive).
		Save(l.ctx)

	if err != nil {
		return nil, err
	}

	return &types.CommonResp{
		ID:      id,
		Message: "更新成功",
	}, nil
}
