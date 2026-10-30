package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// SafeResolve 将相对路径 rel 解析到 root 目录下，防止路径穿越（如 ../etc）。
// 返回基于 root 的绝对路径；若 rel 试图逃出 root 则返回错误。
func SafeResolve(root, rel string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	joined := filepath.Join(rootAbs, filepath.Clean(rel))
	if joined != rootAbs && !strings.HasPrefix(joined, rootAbs+string(os.PathSeparator)) {
		return "", os.ErrInvalid
	}
	return joined, nil
}
