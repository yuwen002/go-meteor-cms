// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionTreeLogic {
	return &PermissionTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionTreeLogic) PermissionTree() (resp []types.PermissionTreeItem, err error) {
	// todo: add your logic here and delete this line

	return
}
