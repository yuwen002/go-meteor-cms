package cleanup

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/loginlog"
	"github.com/yuwen002/go-meteor-cms/ent/operationlog"
	"github.com/zeromicro/go-zero/core/logx"
)

// auditLogCleanInterval 审计日志清理周期
const auditLogCleanInterval = 24 * time.Hour

// PurgeOldOperationLogs 删除早于 retentionDays 天的操作日志。
// retentionDays <= 0 时直接返回，表示永久保留。
func PurgeOldOperationLogs(ctx context.Context, client *ent.Client, retentionDays int) (int, error) {
	if retentionDays <= 0 {
		return 0, nil
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return client.OperationLog.
		Delete().
		Where(operationlog.CreatedAtLT(cutoff)).
		Exec(timeoutCtx)
}

// PurgeOldLoginLogs 删除早于 retentionDays 天的登录日志。
// retentionDays <= 0 时直接返回，表示永久保留。
func PurgeOldLoginLogs(ctx context.Context, client *ent.Client, retentionDays int) (int, error) {
	if retentionDays <= 0 {
		return 0, nil
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return client.LoginLog.
		Delete().
		Where(loginlog.LoginAtLT(cutoff)).
		Exec(timeoutCtx)
}

// StartAuditLogCleaner 启动审计日志定期清理任务。
// 启动时立即执行一次，之后每 24 小时执行；retentionDays <= 0 时不启动（永久保留）。
func StartAuditLogCleaner(ctx context.Context, client *ent.Client, retentionDays int) {
	if retentionDays <= 0 {
		return
	}

	go func() {
		run := func() {
			opDeleted, err := PurgeOldOperationLogs(ctx, client, retentionDays)
			if err != nil {
				logx.Errorf("清理过期操作日志失败: %v", err)
				return
			}
			loginDeleted, err := PurgeOldLoginLogs(ctx, client, retentionDays)
			if err != nil {
				logx.Errorf("清理过期登录日志失败: %v", err)
				return
			}
			if opDeleted+loginDeleted > 0 {
				logx.Infof("已清理过期审计日志 操作日志 %d 条 / 登录日志 %d 条", opDeleted, loginDeleted)
			}
		}

		run()

		ticker := time.NewTicker(auditLogCleanInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}
