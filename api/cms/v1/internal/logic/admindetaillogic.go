// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDetailLogic {
	return &AdminDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDetailLogic) AdminDetail() (resp *types.AdminDetailResp, err error) {
	id := l.ctx.Value("id").(int64)
	fmt.Println("id:", id)

	admin, err := l.svcCtx.EntClient.AdminUser.Get(l.ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(common.ErrAdminUserNotFound)
		}
		l.Logger.Errorf("获取管理员详情失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	return &types.AdminDetailResp{
		Id:       admin.ID,
		Username: admin.Username,
		Nickname: admin.Nickname,
		Email:    admin.Email,
		Phone:    admin.Phone,
		IsActive: admin.IsActive,
	}, nil
}
