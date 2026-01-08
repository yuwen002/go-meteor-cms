// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleCreateLogic {
	return &RoleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleCreateLogic) RoleCreate(req *types.CreateRoleReq) (resp *types.CommonResp, err error) {
	// todo: add your logic here and delete this line

	return
}
