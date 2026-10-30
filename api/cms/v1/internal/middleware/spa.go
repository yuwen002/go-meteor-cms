package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/rest"
)

// NewSPAMiddleware 在配置了 StaticDir 时由本服务直接托管前端 SPA：
//   - 跳过 API 前缀（/admin、/api）与静态资源前缀（/uploads），交由 go-zero 路由处理；
//   - 其余 GET 请求：若路径对应真实文件（如 /assets/index-xxx.js）则直接返回，
//     否则回退到 index.html，以支持前端 history 路由（如 /system/roles 刷新不 404）。
//
// StaticDir 为空时本中间件直接透传，相当于不托管前端。
// 返回 rest.Middleware（即 func(http.HandlerFunc) http.HandlerFunc）。
func NewSPAMiddleware(staticDir string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if staticDir == "" || r.Method != http.MethodGet {
				next(w, r)
				return
			}
			path := r.URL.Path
			if strings.HasPrefix(path, "/admin") ||
				strings.HasPrefix(path, "/api") ||
				strings.HasPrefix(path, "/uploads") {
				next(w, r)
				return
			}

			// 尝试直接返回真实文件（去掉查询串并防目录穿越）
			clean := strings.TrimPrefix(path, "/")
			if clean != "" {
				filePath := filepath.Join(staticDir, filepath.Clean(clean))
				if isInside(staticDir, filePath) {
					if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
						http.ServeFile(w, r, filePath)
						return
					}
				}
			}

			// SPA history 回退
			index := filepath.Join(staticDir, "index.html")
			if _, err := os.Stat(index); err == nil {
				http.ServeFile(w, r, index)
				return
			}

			next(w, r)
		}
	}
}

// isInside 判断 target 是否位于 root 目录内，避免 ../ 穿越读取任意文件。
func isInside(root, target string) bool {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return false
	}
	rootWithSep := absRoot + string(os.PathSeparator)
	return absTarget == absRoot || strings.HasPrefix(absTarget, rootWithSep)
}
