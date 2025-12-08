// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminBindDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminBindDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminBindDepartmentLogic {
	return &AdminBindDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminBindDepartmentLogic) AdminBindDepartment(req *types.AdminBindDeptReq) (resp *types.CommonResp, err error) {
	// 1️⃣ 校验管理员是否存在
	user, err := l.svcCtx.EntClient.AdminUser.
		Query().
		Where(adminuser.IDEQ(req.Id)).
		Only(l.ctx)

	return
}
