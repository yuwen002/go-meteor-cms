// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"os"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FileDelete 删除上传目录下指定文件；目录仅允许在为空时删除，防止误删整棵目录树。
func (l *FileDeleteLogic) FileDelete(req *types.FileDeleteReq) error {
	root := l.svcCtx.Config.Upload.Defaults().Dir
	target, err := utils.SafeResolve(root, req.Path)
	if err != nil {
		return common.NewBizError(common.ErrInvalidParams)
	}
	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return common.NewBizError(common.ErrFileNotFound)
		}
		return common.NewBizError(common.ErrInternalServer)
	}
	if info.IsDir() {
		entries, rerr := os.ReadDir(target)
		if rerr != nil {
			return common.NewBizError(common.ErrInternalServer)
		}
		if len(entries) > 0 {
			return common.NewBizError(common.ErrInvalidParams)
		}
	}
	if err := os.Remove(target); err != nil {
		logx.Errorf("delete %s failed: %v", target, err)
		return common.NewBizError(common.ErrInternalServer)
	}
	return nil
}
