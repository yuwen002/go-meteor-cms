// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent/department"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDepartmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentCreateLogic {
	return &DepartmentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentCreateLogic) DepartmentCreate(req *types.DepartmentCreateReq) (resp *types.CommonResp, err error) {
	// 父级校验（如果 parent_id != 0）
	if req.ParentId != 0 {
		exist, err := l.svcCtx.EntClient.Department.
			Query().
			Where(department.IDEQ(req.ParentId)).
			Exist(l.ctx)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, common.NewBizError(common.ErrDepartmentParentNotExist)
		}
	}

	dept, err := l.svcCtx.EntClient.Department.
		Create().
		SetName(req.Name).
		SetParentID(req.ParentId).
		SetSort(req.Sort).
		SetIsActive(req.IsActive).
		SetNillableLeaderID(&req.LeaderId).
		Save(l.ctx)

	if err != nil {
		return nil, common.NewBizError(common.ErrDepartmentCreateFail)
	}

	return &types.CommonResp{
		ID:      dept.ID,
		Message: "创建成功",
	}, nil
}
