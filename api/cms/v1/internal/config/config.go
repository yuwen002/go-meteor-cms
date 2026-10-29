// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import "github.com/zeromicro/go-zero/rest"

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
	}
}
