// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/types"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/operationlog"
	"github.com/yuwen002/go-meteor-cms/ent/predicate"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogListLogic {
	return &OperationLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationLogListLogic) OperationLogList(req *types.OperationLogListReq) (resp *types.OperationLogListResp, err error) {
	offset, pageSize := utils.NormalizePage(req.Page, req.PageSize, 20)
	page := offset/pageSize + 1

	query := l.svcCtx.EntClient.OperationLog.Query()

	conds := make([]predicate.OperationLog, 0, 8)
	// 数据权限：非超管只能查看自己的操作日志，忽略外部传入的 user_id 参数，防止越权检索他人日志
	if !utils.IsSuperAdmin(l.ctx) {
		uid, ok := utils.CurrentUserID(l.ctx)
		if !ok || uid <= 0 {
			return nil, common.NewBizError(common.ErrUnauthorized)
		}
		conds = append(conds, operationlog.UserIDEQ(uid))
	} else if req.UserID > 0 {
		conds = append(conds, operationlog.UserIDEQ(req.UserID))
	}
	if req.Username != "" {
		conds = append(conds, operationlog.UsernameContainsFold(req.Username))
	}
	if req.Module != "" {
		conds = append(conds, operationlog.ModuleEQ(req.Module))
	}
	if req.Path != "" {
		conds = append(conds, operationlog.PathContainsFold(req.Path))
	}
	if req.Status > 0 {
		conds = append(conds, operationlog.StatusEQ(req.Status))
	}
	if start, ok := parseLogTime(req.StartTime); ok {
		conds = append(conds, operationlog.CreatedAtGTE(start))
	}
	if end, ok := parseLogTime(req.EndTime); ok {
		conds = append(conds, operationlog.CreatedAtLTE(end))
	}
	if len(conds) > 0 {
		query = query.Where(operationlog.And(conds...))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		l.Logger.Errorf("查询操作日志总数失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	list, err := query.
		Order(ent.Desc(operationlog.FieldCreatedAt)).
		Offset(offset).
		Limit(pageSize).
		All(l.ctx)
	if err != nil {
		l.Logger.Errorf("查询操作日志失败: %v", err)
		return nil, common.NewBizError(common.ErrInternalServer)
	}

	items := make([]types.OperationLogItem, 0, len(list))
	for _, item := range list {
		items = append(items, types.OperationLogItem{
			ID:         item.ID,
			UserID:     item.UserID,
			Username:   item.Username,
			Module:     item.Module,
			Action:     item.Action,
			Method:     item.Method,
			Path:       item.Path,
			IP:         item.IP,
			UserAgent:  item.UserAgent,
			Status:     item.Status,
			ErrorMsg:   item.ErrorMsg,
			DurationMs: item.DurationMs,
			Params:     item.Params,
			CreatedAt:  item.CreatedAt.Format(time.RFC3339),
		})
	}

	return &types.OperationLogListResp{
		Total:    int64(total),
		Page:     page,
		PageSize: pageSize,
		List:     items,
	}, nil
}

// parseLogTime 解析可选的起止时间，支持多种常见格式
func parseLogTime(input string) (time.Time, bool) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, false
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, input, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
