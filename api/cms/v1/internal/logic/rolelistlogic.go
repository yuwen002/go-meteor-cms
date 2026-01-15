// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminrole"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleListLogic) RoleList(req *types.RoleListReq) (resp *types.RoleListResp, err error) {
	// 1. 构建查询条件
	query := l.svcCtx.EntClient.AdminRole.Query()

	// 2. 添加搜索条件
	if req.Keyword != "" {
		query = query.Where(
			adminrole.Or(
				adminrole.NameContains(req.Keyword),
				adminrole.CodeContains(req.Keyword),
			),
		)
	}

	// 3. 添加状态过滤
	if req.IsActive != nil {
		query = query.Where(adminrole.IsActive(*req.IsActive))
	}

	// 4. 获取总数
	total, err := query.Count(l.ctx)
	if err != nil {
		l.Logger.Errorf("获取角色总数失败: %v", err)
		return nil, common.NewBizError(common.ErrRoleListFail)
	}

	// 5. 分页查询
	roles, err := query.
		Order(ent.Asc(adminrole.FieldSort), ent.Desc(adminrole.FieldID)). // 按sort升序，ID降序
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(l.ctx)
	if err != nil {
		l.Logger.Errorf("查询角色列表失败: %v", err)
		return nil, common.NewBizError(common.ErrRoleListFail)
	}

	// 6. 构建响应数据
	list := make([]types.RoleItem, 0, len(roles))
	for _, role := range roles {
		list = append(list, types.RoleItem{
			ID:        role.ID,
			Name:      role.Name,
			Code:      role.Code,
			Desc:      role.Desc,
			DataScope: role.DataScope,
			IsSystem:  role.IsSystem,
			IsActive:  role.IsActive,
			Sort:      role.Sort,
			CreatedAt: role.CreatedAt.Format(time.RFC3339),
			UpdatedAt: role.UpdatedAt.Format(time.RFC3339),
		})
	}

	// 7. 返回分页结果
	return &types.RoleListResp{
		List:     list,
		Total:    int64(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
