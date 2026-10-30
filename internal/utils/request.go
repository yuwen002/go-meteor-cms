package utils

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type reqCtxKey string

const (
	// ClientIPKey 客户端 IP 在 context 中的键
	ClientIPKey reqCtxKey = "client_ip"
	// UserAgentKey 客户端标识在 context 中的键
	UserAgentKey reqCtxKey = "user_agent"
	// RequestKey 原始 HTTP 请求在 context 中的键
	RequestKey reqCtxKey = "http_request"
)

// WithClientInfo 将客户端信息注入 context，供 logic 层读取
func WithClientInfo(ctx context.Context, ip, userAgent string) context.Context {
	ctx = context.WithValue(ctx, ClientIPKey, ip)
	ctx = context.WithValue(ctx, UserAgentKey, userAgent)
	return ctx
}

// WithRequest 将原始 HTTP 请求注入 context，供需要读取 multipart 等内容的 logic 使用
func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, RequestKey, r)
}

// GetRequest 从 context 中获取原始 HTTP 请求
func GetRequest(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(RequestKey).(*http.Request); ok {
		return r
	}
	return nil
}

// GetClientIP 从 context 中获取客户端 IP
func GetClientIP(ctx context.Context) string {
	if v, ok := ctx.Value(ClientIPKey).(string); ok {
		return v
	}
	return ""
}

// GetUserAgent 从 context 中获取客户端标识
func GetUserAgent(ctx context.Context) string {
	if v, ok := ctx.Value(UserAgentKey).(string); ok {
		return v
	}
	return ""
}

// ClientIPFromRequest 提取客户端 IP，优先信任反向代理注入的请求头
func ClientIPFromRequest(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if first := strings.Split(forwarded, ",")[0]; strings.TrimSpace(first) != "" {
			return strings.TrimSpace(first)
		}
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// Truncate 按字符数截断字符串，避免截断多字节字符产生乱码
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
