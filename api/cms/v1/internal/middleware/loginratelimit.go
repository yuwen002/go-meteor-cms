package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/rest"
)

// loginRateLimitCode 与统一响应约定一致的业务码（请求过于频繁）
const loginRateLimitCode = 2029

type ipCounter struct {
	mu           sync.Mutex
	count        int
	windowStart  int64
	blockedUntil int64
}

// NewLoginRateLimitMiddleware 按客户端 IP 对登录接口做固定窗口限流，用于缓解爆破。
// 仅对 POST /admin/login 生效；未启用（enabled=false）时直接透传，零额外开销。
func NewLoginRateLimitMiddleware(enabled bool, max, window, blockSeconds int) rest.Middleware {
	if !enabled {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return next
		}
	}
	rl := struct {
		Max          int
		Window       int
		BlockSeconds int
	}{max, window, blockSeconds}
	store := &sync.Map{}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/admin/login" || r.Method != http.MethodPost {
				next(w, r)
				return
			}
			ip := clientIP(r)
			now := time.Now().Unix()
			val, _ := store.LoadOrStore(ip, &ipCounter{})
			c := val.(*ipCounter)
			c.mu.Lock()
			if c.blockedUntil > now {
				c.mu.Unlock()
				writeRateLimited(w)
				return
			}
			if now-c.windowStart > int64(rl.Window) {
				c.count = 0
				c.windowStart = now
			}
			c.count++
			if c.count > rl.Max {
				c.blockedUntil = now + int64(rl.BlockSeconds)
				c.mu.Unlock()
				writeRateLimited(w)
				return
			}
			c.mu.Unlock()
			next(w, r)
		}
	}
}

func writeRateLimited(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code": loginRateLimitCode,
		"msg":  "请求过于频繁，请稍后再试",
		"data": nil,
	})
}

// clientIP 从 X-Forwarded-For 取首个 IP，否则取 RemoteAddr（去除端口）
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.IndexByte(xff, ','); idx >= 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
