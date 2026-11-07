// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"
	"time"
)

// transactionTimeout 是写事务的独立超时上限，需小于 RestConf.Timeout（见 etc/cms-api.yaml）。
// go-zero 的 Timeout 中间件会给每个请求 ctx 加 deadline；事务若继承这个 ctx，deadline 一到
// database/sql 会自动回滚事务，之后的 Commit/Rollback 全部失败，并产生
// "context deadline exceeded" 与误导性的「事务回滚失败」日志。
const transactionTimeout = 10 * time.Second

// newTxContext 把当前 logic 实例的请求 ctx 切换为写事务专用 ctx：
//   - context.WithoutCancel 保留 trace、日志字段、用户 claims 等 Value，但剥离请求取消信号；
//   - context.WithTimeout 给事务一个独立硬上限，避免客户端断开后事务无限占用数据库连接。
//
// 每个请求都会新建 logic 实例，因此修改 ctx 不会跨请求共享。调用方应在开启事务前调用
// 并 defer cancel()：这样 cancel 的 defer 注册早于事务回滚的 defer，按 LIFO 顺序会先回滚、
// 再释放 ctx。切换后函数体内所有 l.ctx 引用（含事务内 SQL）自动使用新 ctx。
func newTxContext(ctx *context.Context) context.CancelFunc {
	txCtx, cancel := context.WithTimeout(context.WithoutCancel(*ctx), transactionTimeout)
	*ctx = txCtx
	return cancel
}
