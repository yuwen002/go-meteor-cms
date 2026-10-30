package logic

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/ent/tokenblacklist"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RefreshToken 使用刷新令牌签发新的访问令牌。
// 刷新令牌使用独立密钥与 access token 区分，且携带 token_type=refresh 标识，防止被当作 access token 使用。
func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	claims, err := utils.ParseToken(l.svcCtx.Config.Auth.RefreshSecret, req.RefreshToken)
	if err != nil {
		return nil, common.NewBizError(common.ErrInvalidToken)
	}

	// 校验令牌类型，防止用 access token 冒充 refresh token
	if typ, _ := claims["token_type"].(string); typ != "refresh" {
		return nil, common.NewBizError(common.ErrInvalidToken)
	}

	// 检查刷新令牌是否已被吊销（登出时写入黑名单）
	rtHash := utils.GenerateTokenHash(req.RefreshToken)
	blacked, berr := l.svcCtx.EntClient.TokenBlacklist.
		Query().
		Where(tokenblacklist.TokenHashEQ(rtHash)).
		Exist(l.ctx)
	if berr != nil {
		return nil, common.NewBizError(common.ErrInternalServer)
	}
	if blacked {
		return nil, common.NewBizError(common.ErrInvalidToken)
	}

	uid, _ := claims["user_id"].(float64)
	userID := int64(uid)
	if userID <= 0 {
		return nil, common.NewBizError(common.ErrInvalidToken)
	}

	// 用户被删除或禁用时，刷新令牌立即失效
	user, err := l.svcCtx.EntClient.AdminUser.
		Query().
		Where(adminuser.IDEQ(userID)).
		Select(
			adminuser.FieldID,
			adminuser.FieldUsername,
			adminuser.FieldNickname,
			adminuser.FieldEmail,
			adminuser.FieldPhone,
			adminuser.FieldIsSuper,
			adminuser.FieldIsActive,
		).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBizError(common.ErrInvalidToken)
		}
		return nil, common.NewBizError(common.ErrInternalServer)
	}
	if !user.IsActive {
		return nil, common.NewBizError(common.ErrAccountInactive)
	}

	accessToken, err := utils.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"is_super":   user.IsSuper,
		"is_active":  user.IsActive,
		"login_time": time.Now().Unix(),
	})
	if err != nil {
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	// 刷新令牌轮换：签发新的刷新令牌，并将本次使用的旧刷新令牌立即吊销，
	// 降低 refresh token 泄露后的重放风险。黑名单 TTL 取 RefreshExpire，可被清理任务回收。
	newRefreshToken, err := utils.GenerateToken(l.svcCtx.Config.Auth.RefreshSecret, l.svcCtx.Config.Auth.RefreshExpire, map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"token_type": "refresh",
	})
	if err != nil {
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	oldHash := utils.GenerateTokenHash(req.RefreshToken)
	if _, e := l.svcCtx.EntClient.TokenBlacklist.Create().
		SetTokenHash(oldHash).
		SetToken(req.RefreshToken).
		SetExpiredAt(time.Now().Add(time.Duration(l.svcCtx.Config.Auth.RefreshExpire) * time.Second)).
		Save(l.ctx); e != nil {
		logx.Errorf("Failed to revoke old refresh token: %v", e)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	return &types.RefreshTokenResp{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}
