package logic

import (
	"context"
	"go-meteor-cms/ent/adminuser"
	"go-meteor-cms/internal/util"
	"time"

	"go-meteor-cms/rpc/v1/admin/admin"
	"go-meteor-cms/rpc/v1/admin/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *admin.LoginReq) (*admin.LoginResp, error) {
	u, err := l.svcCtx.EntClient.AdminUser.
		Query().
		Where(adminuser.UsernameEQ(in.Username)).
		Only(l.ctx)

	if err != nil {
		return &admin.LoginResp{
			Code: 1,
			Msg:  "用户不存在",
		}, nil
	}

	if !util.CheckPassword(u.PasswordHash, in.Password) {
		return &admin.LoginResp{Code: 1, Msg: "密码错误"}, nil
	}

	token, err := util.GenerateToken(u.ID, u.Username)
	if err != nil {
		return &admin.LoginResp{Code: 500, Msg: "生成 Token 失败"}, nil
	}

	// 登录成功后更新最后登录时间
	_ = l.svcCtx.EntClient.AdminUser.UpdateOneID(u.ID).
		SetLastLoginAt(time.Now()).
		Exec(l.ctx)

	return &admin.LoginResp{
		Code:  0,
		Msg:   "登录成功",
		Token: token,
	}, nil
}
