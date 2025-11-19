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
	// 获取当前登录用户信息
	claims, ok := l.ctx.Value("user").(map[string]interface{})
	if !ok {
		return nil, common.NewBizError(401, "未登录")
	}
	currentUserID, ok := claims["user_id"].(int64)
	if !ok {
		return nil, common.NewBizError(400, "用户ID格式错误")
	}

	targetID := req.ID

	// 检查是否尝试修改自己的密码
	if currentUserID == targetID {
		return nil, common.NewBizError(400, "不能重置自己的密码，请使用修改密码功能")
	}

	// 检查目标用户是否存在
	_, err = l.svcCtx.EntClient.AdminUser.Get(l.ctx, targetID)
	if err != nil {
		return nil, common.NewBizError(404, "用户不存在")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	_, err = l.svcCtx.EntClient.AdminUser.
		UpdateOneID(targetID).
		SetPasswordHash(string(hashed)).
		Save(l.ctx)

	if err != nil {
		return nil, common.NewBizError(500, "重置密码失败: "+err.Error())
	}

	return &types.CommonResp{Message: "密码已重置"}, nil
}
