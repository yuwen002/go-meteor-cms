// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	baseconfig "github.com/yuwen002/go-meteor-cms/internal/config"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database struct {
		Driver string
		Source string
	}
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string // 刷新令牌专用密钥，与 AccessSecret 区分，避免 access 泄露后 refresh 被冒用
		RefreshExpire int64  `json:",default=604800"` // 刷新令牌有效期，单位秒（默认 7 天）
		// RefreshGraceSeconds 刷新令牌轮换后的并发宽限期（秒），默认 5。
		// 轮换会把旧令牌签入黑名单，同一旧令牌在宽限期内重复提交仍被接受，避免并发请求误登出；
		// 超出宽限期判定为盗用重放并拒绝。
		RefreshGraceSeconds int `json:",default=5"`
	}
	Mail struct {
		Enabled  bool
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}
	Maintenance struct {
		// TokenBlacklistCleanInterval 过期 token 黑名单的清理间隔（秒），<=0 时按 3600 秒处理
		TokenBlacklistCleanInterval int
		// AuditLogRetentionDays 操作/登录日志保留天数，<=0 表示永久保留（不清理）
		AuditLogRetentionDays int `json:",default=90"`
	}
	Login struct {
		// MaxFailCount 连续登录失败达到该次数后锁定账号，<=0 时不启用锁定
		MaxFailCount int `json:",default=5"`
		// LockMinutes 锁定持续时间（分钟），<=0 时按 15 处理
		LockMinutes int `json:",default=15"`
		// RateLimit 登录接口按客户端 IP 的限流配置（防爆破）；默认关闭，需显式开启
		RateLimit struct {
			Enabled      bool `json:",default=false"`
			Max          int  `json:",default=10"`  // 窗口内允许的最大尝试次数
			Window       int  `json:",default=60"`  // 统计窗口（秒）
			BlockSeconds int  `json:",default=300"` // 触发限流后拒绝时长（秒）
		}
	}
	// WebBaseURL 前端站点地址，用于拼接找回密码等邮件中的跳转链接；为空时仅下发裸令牌
	WebBaseURL string
	// StaticDir 前端构建产物（dist）目录；非空时由本服务直接托管 SPA（含 history 路由回退）。
	// 容器化部署通常通过环境变量 WEB_DIST 注入，本地开发留空即可（仅暴露 API）。
	StaticDir string
	Upload    baseconfig.Upload
}
