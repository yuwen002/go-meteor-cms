// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAdminPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetAdminPasswordLogic {
	return &ResetAdminPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAdminPasswordLogic) ResetAdminPassword(req *types.ResetAdminPasswordReq) (resp *types.CommonResp, err error) {
	id := req.ID
	_, err = l.svcCtx.EntClient.AdminUser.Get(l.ctx, id)
	if err != nil {
		return nil, common.NewBizError(404, "用户不存在")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	_, err = l.svcCtx.EntClient.AdminUser.
		UpdateOneID(id).
		SetPasswordHash(string(hashed)).
		Save(l.ctx)

	if err != nil {
		return nil, common.NewBizError(500, "重置密码失败")
	}

	return &types.CommonResp{Message: "密码已重置"}, nil
}
