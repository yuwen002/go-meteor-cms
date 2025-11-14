// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterAdminLogic {
	return &RegisterAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterAdminLogic) RegisterAdmin(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. 检查用户名是否存在
	exists, err := l.svcCtx.EntClient.AdminUser.Query().
		Where(adminuser.UsernameEQ(req.Username)).
		Exist(l.ctx)

	if err != nil {
		return nil, common.NewBizError(50001, "检查用户失败")
	}
	if exists {
		return nil, common.NewBizError(40001, "用户名已存在")
	}

	// 2. 检查邮箱是否存在
	emailExists, err := l.svcCtx.EntClient.AdminUser.Query().
		Where(adminuser.EmailEQ(req.Email)).
		Exist(l.ctx)

	if err != nil {
		return nil, common.NewBizError(50002, "检查邮箱失败")
	}
	if emailExists {
		return nil, common.NewBizError(40002, "邮箱已注册")
	}

	// 3. 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.NewBizError(50003, "密码加密失败")
	}

	// 4. 创建用户
	_, err = l.svcCtx.EntClient.AdminUser.
		Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPasswordHash(string(hashed)).
		SetNickname(req.Nickname).
		SetIsSuper(false).
		SetIsActive(false).
		Save(l.ctx)

	if err != nil {
		return nil, common.NewBizError(50004, "用户创建失败")
	}

	return &types.RegisterResp{
		Message: "管理员注册成功",
	}, nil
}
