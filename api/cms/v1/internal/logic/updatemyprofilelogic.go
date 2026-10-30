// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyProfileLogic {
	return &UpdateMyProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateMyProfile 更新当前登录管理员的个人资料。
//
// 仅允许维护昵称、邮箱、手机号；角色、启用状态、是否超管等由管理员管理接口负责，
// 避免普通账号通过个人资料接口自我提权。
func (l *UpdateMyProfileLogic) UpdateMyProfile(req *types.UpdateProfileReq) (resp *types.AdminProfileResp, err error) {
	userID, ok := utils.CurrentUserID(l.ctx)
	if !ok {
		return nil, common.NewBizError(common.ErrMissingUserContext)
	}

	// 邮箱查重（排除自身）
	if req.Email != nil && *req.Email != "" {
		exists, err := l.svcCtx.EntClient.AdminUser.
			Query().
			Where(adminuser.EmailEQ(*req.Email), adminuser.IDNEQ(userID)).
			Exist(l.ctx)
		if err != nil {
			l.Logger.Errorf("检查邮箱失败: %v", err)
			return nil, common.NewBizError(common.ErrInternalServer)
		}
		if exists {
			return nil, common.NewBizError(common.ErrAdminEmailAlreadyExists)
		}
	}

	// 按传入字段做局部更新
	update := l.svcCtx.EntClient.AdminUser.UpdateOneID(userID)
	changed := false
	if req.Nickname != nil {
		update = update.SetNickname(*req.Nickname)
		changed = true
	}
	if req.Email != nil {
		update = update.SetEmail(*req.Email)
		changed = true
	}
	if req.Phone != nil {
		update = update.SetPhone(*req.Phone)
		changed = true
	}
	if req.Avatar != nil {
		update = update.SetAvatar(utils.Truncate(*req.Avatar, 512))
		changed = true
	}

	if !changed {
		user, err := l.svcCtx.EntClient.AdminUser.Get(l.ctx, userID)
		if err != nil {
			l.Logger.Errorf("查询当前用户信息失败: %v", err)
			return nil, common.NewBizError(common.ErrInternalServer)
		}
		return toProfileResp(user), nil
	}

	user, err := update.Save(l.ctx)
	if err != nil {
		l.Logger.Errorf("更新个人资料失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	return toProfileResp(user), nil
}
