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

type AdminCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateLogic {
	return &AdminCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCreateLogic) AdminCreate(req *types.CreateAdminReq) (resp *types.CommonResp, err error) {
	// 检查用户是否存在
	exists, err := l.svcCtx.EntClient.AdminUser.
		Query().
		Where(adminuser.Username(req.Username)).
		Exist(l.ctx)

	if err != nil {
		logx.Errorf("查询管理员用户失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	if exists {
		return nil, common.NewBizError(common.ErrAdminUserAlreadyExists)
	}

	// 密码加密
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logx.Errorf("密码加密失败: %v", err)
		return nil, common.NewBizError(common.ErrAdminPasswordHashFail)
	}

	// 创建用户
	user, err := l.svcCtx.EntClient.AdminUser.
		Create().
		SetUsername(req.Username).
		SetPasswordHash(string(hashed)).
		SetNickname(req.Nickname).
		SetIsSuper(req.IsSuper).
		SetIsActive(req.IsActive).
		Save(l.ctx)

	if err != nil {
		logx.Errorf("创建管理员失败: %v", err)
		return nil, common.NewBizError(common.ErrAdminCreateFailed)
	}

	return &types.CommonResp{ID: user.ID, Message: "创建成功"}, nil
}
