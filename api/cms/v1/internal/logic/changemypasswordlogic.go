// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeMyPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeMyPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeMyPasswordLogic {
	return &ChangeMyPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeMyPasswordLogic) ChangeMyPassword(req *types.ChangeMyPasswordReq) (resp *types.CommonResp, err error) {
	claims := utils.GetUserFromCtx(l.ctx)
	if claims == nil {
		return nil, common.NewBizError(common.ErrUnauthorized)
	}

	userID, ok := claims["user_id"].(int64)
	if !ok {
		return nil, common.NewBizError(common.ErrUserIDFormat)
	}

	// 查询当前用户
	user, err := l.svcCtx.EntClient.AdminUser.Get(l.ctx, userID)
	if err != nil {
		l.Logger.Errorf("获取用户信息失败: %v", err)
		return nil, common.NewBizError(common.ErrAdminUserNotFound)
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return nil, common.NewBizError(common.ErrAdminPasswordIncorrect)
	}

	// 设置新密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Errorf("密码加密失败: %v", err)
		return nil, common.NewBizError(common.ErrAdminPasswordHashFail)
	}

	_, err = l.svcCtx.EntClient.AdminUser.
		UpdateOneID(userID).
		SetPasswordHash(string(hashed)).
		Save(l.ctx)

	if err != nil {
		l.Logger.Errorf("更新密码失败: %v", err)
		return nil, common.NewBizError(common.ErrPasswordUpdateFailed)
	}

	return &types.CommonResp{Message: "密码修改成功"}, nil
}
