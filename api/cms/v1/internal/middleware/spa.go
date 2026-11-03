package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SPAHandler 在配置了 StaticDir 时由本服务直接托管前端 SPA，作为自定义 404(NotFound) 处理器使用：
//   - 跳过 API 前缀（/admin、/api）与静态资源前缀（/uploads），交由 go-zero 路由处理；
//   - 其余 GET 请求：若路径对应真实文件（如 /assets/index-xxx.js）则直接返回，
//     否则回退到 index.html，以支持前端 history 路由（如 /system/roles 刷新不 404）。
//
// 必须以 rest.WithNotFoundHandler 注册，而不是 server.Use：
// go-zero 的 server.Use 注册的全局中间件不会作用于 404（未匹配）路由，
// 而 notFoundHandler 在 MustNewServer 时即绑定中间件链，早于 server.Use 调用，
// 因此用 server.Use 兜底的 SPA 中间件对 /、/admin/users 这类未匹配路径不会生效。
func SPAHandler(staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 未配置产物目录、或非 GET 请求：交回 go-zero 默认 404
		if staticDir == "" || r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		path := r.URL.Path
		if strings.HasPrefix(path, "/admin") ||
			strings.HasPrefix(path, "/api") ||
			strings.HasPrefix(path, "/uploads") {
			http.NotFound(w, r)
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

		// SPA history 路由回退到 index.html
		index := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(index); err == nil {
			http.ServeFile(w, r, index)
			return
		}

		// 既无对应文件、也无 index.html：返回 404
		http.NotFound(w, r)
	})
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
