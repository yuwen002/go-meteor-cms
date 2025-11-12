package utils

import "context"

type ctxKey string

const UserKey ctxKey = "user"

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
