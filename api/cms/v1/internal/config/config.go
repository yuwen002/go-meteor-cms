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
		AccessSecret string
		AccessExpire int64
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
	Upload baseconfig.Upload
}
