package utils

import (
	"context"
	"math"
)

type ctxKey string

const (
	// UserKey 用户信息在 context 中的键
	UserKey ctxKey = "user"
	// TokenKey 当前请求 token 在 context 中的键
	TokenKey ctxKey = "token"
)

// WithUserCtx 将用户信息存储到 context 中
// 参数:
//   - ctx: 原始 context
//   - claims: 包含用户信息的 map，通常是从 JWT Token 中解析出的声明
//
// 返回:
//   - context.Context: 包含用户信息的新 context
func WithUserCtx(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, UserKey, claims)
}

// GetUserFromCtx 从 context 中获取用户信息
// 参数:
//   - ctx: 包含用户信息的 context
//
// 返回:
//   - map[string]interface{}: 用户信息，如果不存在则返回 nil
//
// 注意:
//   - 返回的 map 可能包含用户 ID、用户名等认证信息
//   - 使用前应检查返回值是否为 nil
func GetUserFromCtx(ctx context.Context) map[string]interface{} {
	if val := ctx.Value(UserKey); val != nil {
		return val.(map[string]interface{})
	}
	return nil
}

// WithTokenCtx 将当前请求 token 存储到 context 中
func WithTokenCtx(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenKey, token)
}

// GetTokenFromCtx 从 context 中获取 token，不存在时返回空字符串
func GetTokenFromCtx(ctx context.Context) string {
	if val := ctx.Value(TokenKey); val != nil {
		if token, ok := val.(string); ok {
			return token
		}
	}
	return ""
}

// CurrentUserID 从 context 中解析当前登录用户的 ID
// JWT 反序列化后数字类型为 float64，这里统一做类型兼容
func CurrentUserID(ctx context.Context) (int64, bool) {
	claims := GetUserFromCtx(ctx)
	if claims == nil {
		return 0, false
	}
	return ExtractUserID(claims)
}

// IsSuperAdmin 判断 context 中的用户是否为超级管理员，取值缺失时按非超管处理
func IsSuperAdmin(ctx context.Context) bool {
	claims := GetUserFromCtx(ctx)
	if claims == nil {
		return false
	}
	isSuper, ok := claims["is_super"].(bool)
	return ok && isSuper
}

// ExtractUserID 从 claims 中解析用户 ID，兼容多种数字类型
func ExtractUserID(claims map[string]interface{}) (int64, bool) {
	value, ok := claims["user_id"]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		if math.Trunc(v) != v {
			return 0, false
		}
		return int64(v), true
	default:
		return 0, false
	}
}
