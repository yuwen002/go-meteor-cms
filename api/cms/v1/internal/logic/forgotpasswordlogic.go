// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgotPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgotPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgotPasswordLogic {
	return &ForgotPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForgotPasswordLogic) ForgotPassword(req *types.ForgotPasswordReq) (resp *types.ForgotPasswordResp, err error) {
	// 1. 查找用户
	u, err := l.svcCtx.EntClient.AdminUser.Query().
		Where(adminuser.UsernameEQ(req.Username)).
		First(l.ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(40000, "用户不存在")
		}
		return nil, common.NewBizError(50000, "系统错误，请稍后重试")
	}

	// 2. 生成一个随机 Token，用于重置密码
	resetToken := uuid.NewString()
	expire := time.Now().Add(time.Hour * 1) // 1 小时有效期

	// 3. 存入数据库
	_, err = l.svcCtx.EntClient.AdminUser.
		UpdateOneID(u.ID).
		SetResetToken(resetToken).
		SetResetExpire(expire).
		Save(l.ctx)
	if err != nil {
		return nil, err
	}

	// 4. 发送邮件（这里先模拟）
	logx.Infof("模拟发送邮件: 用户 %s 的重置 Token = %s", u.Username, resetToken)

	// 5. 返回成功信息
	return &types.ForgotPasswordResp{
		Message: "重置密码邮件已发送，请检查邮箱",
	}, nil
}
