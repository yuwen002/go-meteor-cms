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
		return nil, common.NewBizError(401, "未登录")
	}
	userID, ok := claims["user_id"].(int64)
	if !ok {
		return nil, common.NewBizError(400, "用户ID格式错误")
	}

	// 查询当前用户
	user, err := l.svcCtx.EntClient.AdminUser.Get(l.ctx, userID)
	if err != nil {
		return nil, common.NewBizError(404, "用户不存在")
	}

	// 验证旧密码
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) != nil {
		return nil, common.NewBizError(400, "旧密码错误")
	}

	// 设置新密码
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	_, err = l.svcCtx.EntClient.AdminUser.
		UpdateOneID(userID).
		SetPasswordHash(string(hashed)).
		Save(l.ctx)

	if err != nil {
		return nil, common.NewBizError(500, "密码更新失败")
	}

	return &types.CommonResp{Message: "密码修改成功"}, nil
}
