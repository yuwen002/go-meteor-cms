package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

// NewPublicRateLimitMiddleware 对免鉴权的 /api/ 前缀接口按客户端 IP 做固定窗口限流。
//
// 公开接口（文章列表/详情/分类标签/RSS 等）不经过 JWT 与权限中间件，是唯一可被
// 匿名高频调用的面，因此需要独立限流；后台接口已有登录限流与鉴权，不受影响。
// 该中间件通过 server.Use 注册，作用于已注册路由（SPA 的 NotFoundHandler 回退
// 不走 server.Use，但那些请求只是返回静态 index.html，不构成数据库压力）。
func NewPublicRateLimitMiddleware(rds *redis.Redis, enabled bool, max, window, blockSeconds int) rest.Middleware {
	if !enabled {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return next
		}
	}
	limiter := &ipRateLimiter{
		rds: rds, prefix: "public_rl:", max: max, window: window,
		blockSeconds: blockSeconds, mem: &sync.Map{},
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				next(w, r)
				return
			}
			blocked, retryAfter := limiter.allow(clientIP(r))
			if blocked {
				writeRateLimited(w, retryAfter)
				return
			}
			next(w, r)
		}
	}
}
