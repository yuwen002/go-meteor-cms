// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLogoutLogic {
	return &AdminLogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminLogoutLogic) AdminLogout() (resp *types.CommonResp, err error) {
	// 从上下文中获取用户信息
	userClaims := utils.GetUserFromCtx(l.ctx)
	if userClaims == nil {
		return nil, common.NewBizError(common.ErrUnauthorized)
	}

	// 从上下文中获取 token
	token, ok := l.ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, common.NewBizError(common.ErrInvalidToken)
	}

	// 获取token过期时间
	exp, ok := userClaims["exp"].(float64)
	if !ok {
		exp = float64(time.Now().Add(24 * time.Hour).Unix()) // 默认24小时后过期
	}
	expTime := time.Unix(int64(exp), 0)

	// 将token加入黑名单
	_, err = l.svcCtx.EntClient.TokenBlacklist.Create().
		SetToken(token).
		SetExpiredAt(expTime).
		Save(l.ctx)

	if err != nil {
		logx.Errorf("Failed to add token to blacklist: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	return &types.CommonResp{
		ID:      0,
		Message: "Logout successful",
	}, nil
}
