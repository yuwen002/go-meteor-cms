// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileRenameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileRenameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileRenameLogic {
	return &FileRenameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FileRename 在同级目录下重命名文件或目录，新名称不允许包含路径分隔符，且目标不得已存在。
func (l *FileRenameLogic) FileRename(req *types.FileRenameReq) error {
	if strings.TrimSpace(req.NewName) == "" || strings.ContainsAny(req.NewName, "/\\") {
		return common.NewBizError(common.ErrInvalidParams)
	}
	root := l.svcCtx.Config.Upload.Defaults().Dir
	src, err := utils.SafeResolve(root, req.Path)
	if err != nil {
		return common.NewBizError(common.ErrInvalidParams)
	}
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return common.NewBizError(common.ErrFileNotFound)
		}
		return common.NewBizError(common.ErrInternalServer)
	}
	dst := filepath.Join(filepath.Dir(src), req.NewName)
	if _, err := utils.SafeResolve(root, dst); err != nil {
		return common.NewBizError(common.ErrInvalidParams)
	}
	if _, err := os.Stat(dst); err == nil {
		return common.NewBizError(common.ErrInvalidParams)
	}
	if err := os.Rename(src, dst); err != nil {
		logx.Errorf("rename %s -> %s failed: %v", src, dst, err)
		return common.NewBizError(common.ErrInternalServer)
	}
	return nil
}
