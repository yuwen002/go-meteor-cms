package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

import "errors"

var (
	ErrUnauthorized = errors.New("用户名或密码错误")
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	if req.Username != "admin" || req.Password != "123456" {
		return nil, ErrUnauthorized
	}

	tokenStr, err := utils.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, map[string]interface{}{
		"userId": 1,
		"name":   req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{Token: tokenStr}, nil
}
