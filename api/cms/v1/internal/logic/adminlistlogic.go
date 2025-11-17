// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListLogic {
	return &AdminListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminListLogic) AdminList(req *types.AdminListReq) (resp *types.AdminListResp, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10 // 默认每页 10 条
	}

	// ent 查询器
	query := l.svcCtx.EntClient.AdminUser.Query()

	// 关键字搜索
	if req.Keyword != "" {
		query = query.Where(
			adminuser.Or(
				adminuser.UsernameContains(req.Keyword),
				adminuser.NicknameContains(req.Keyword),
				adminuser.EmailContains(req.Keyword),
			),
		)
	}

	// 状态过滤
	if req.Active == 1 {
		query = query.Where(adminuser.IsActive(true))
	} else if req.Active == 2 {
		query = query.Where(adminuser.IsActive(false))
	}

	// 总数
	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, common.NewBizError(50001, "获取管理员总数失败")
	}

	// 分页数据
	list, err := query.
		Order(ent.Desc(adminuser.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(l.ctx)
	if err != nil {
		return nil, common.NewBizError(50002, "获取管理员列表失败")
	}

	// 转换结构
	items := make([]types.AdminItem, 0, len(list))
	for _, u := range list {
		item := types.AdminItem{
			Id:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			Avatar:    u.Avatar,
			IsSuper:   u.IsSuper,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		}
		if u.LastLoginAt != nil {
			item.LastLoginAt = u.LastLoginAt.Format(time.RFC3339)
		}
		items = append(items, item)
	}

	// 返回带上 Page 和 PageSize
	return &types.AdminListResp{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     items,
	}, nil
}
