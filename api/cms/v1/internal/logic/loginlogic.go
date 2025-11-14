package logic

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
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
	// 查数据库用户
	user, err := l.svcCtx.EntClient.AdminUser.
		Query().
		Where(adminuser.UsernameEQ(req.Username)).
		First(l.ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	// 验证密码
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, ErrUnauthorized
	}

	// 检查用户账号是否激活
	// 如果账号未激活，返回未授权错误
	if user.IsActive == false {
		return nil, ErrUnauthorized
	}

	// 签发 JWT
	tokenStr, err := utils.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, map[string]interface{}{
		"userId":    user.ID,           // 用户ID
		"username":  user.Username,     // 用户名
		"nickname":  user.Nickname,     // 用户昵称
		"email":     user.Email,        // 邮箱
		"phone":     user.Phone,        // 手机号
		"isSuper":   user.IsSuper,      // 是否超级管理员
		"isActive":  user.IsActive,     // 是否激活
		"loginTime": time.Now().Unix(), // 登录时间戳
	})
	if err != nil {
		return nil, err
	}

	go func() {
		// 创建一个新的上下文，设置超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 更新最后登录时间为当前时间
		now := time.Now()
		_, err = l.svcCtx.EntClient.AdminUser.UpdateOneID(user.ID).
			SetLastLoginAt(now).
			Save(ctx)
		if err != nil {
			// 记录更详细的错误信息
			logx.WithContext(ctx).Errorf("更新用户ID: %d， 最后登录时间失败: %v", user.ID, err)
		} else {
			logx.WithContext(ctx).Infof("用户ID: %d， 最后登录时间已更新: %v", user.ID, now)
		}
	}()

	return &types.LoginResp{Token: tokenStr}, nil
}
