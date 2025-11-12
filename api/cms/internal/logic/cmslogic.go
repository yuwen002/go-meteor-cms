// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/yuwen002/go-meteor-cms/api/cms/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CmsLogic {
	return &CmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CmsLogic) Cms(req *types.HelloReq) (resp *types.HelloResp, err error) {
	// todo: add your logic here and delete this line
	return &types.HelloResp{
		Message: "Hello " + req.Name,
	}, nil
}
