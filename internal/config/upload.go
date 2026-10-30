// Package config 存放跨层共用的配置项。
package config

import "errors"

// Upload 上传相关配置
type Upload struct {
	// Dir 文件保存根目录，相对进程工作目录
	Dir string
	// MaxSizeMB 单个文件大小上限（MB）
	MaxSizeMB int64
	// AllowedExts 允许的扩展名白名单
	AllowedExts []string
}

// ErrFileTooLarge 文件超过大小上限
var ErrFileTooLarge = errors.New("文件超过大小上限")

// ErrFileTypeNotAllowed 文件类型不在白名单内
var ErrFileTypeNotAllowed = errors.New("不支持的文件类型")

// Default 返回配置的默认值，用于未配置时的兜底
func (u Upload) withDefaults() Upload {
	if u.Dir == "" {
		u.Dir = "public/uploads"
	}
	if u.MaxSizeMB <= 0 {
		u.MaxSizeMB = 5
	}
	if len(u.AllowedExts) == 0 {
		u.AllowedExts = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	}
	return u
}

// Defaults 返回补齐默认值后的配置副本
func (u Upload) Defaults() Upload {
	return u.withDefaults()
}
