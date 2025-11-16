package middleware

import (
	"net/http"
	"strings"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/rest"
)

//
//type JwtMiddleware struct {
//	Secret string
//}
//
//func NewJwtMiddleware(secret string) *JwtMiddleware {
//	return &JwtMiddleware{Secret: secret}
//}
//
//func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		authHeader := r.Header.Get("Authorization")
//		if authHeader == "" {
//			httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
//				"message": "missing Authorization header",
//			})
//			return
//		}
//
//		parts := strings.SplitN(authHeader, " ", 2)
//		if len(parts) != 2 || parts[0] != "Bearer" {
//			httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
//				"message": "invalid Authorization format",
//			})
//			return
//		}
//
//		tokenStr := parts[1]
//		claims, err := utils.ParseToken(m.Secret, tokenStr)
//		if err != nil {
//			httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{
//				"message": "invalid or expired token",
//			})
//			return
//		}
//
//		// 存入 context
//		ctx := utils.WithUserCtx(r.Context(), claims)
//		r = r.WithContext(ctx)
//
//		next(w, r)
//	}
//}

func JwtMiddleware(c *config.Config) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				common.Fail(w, 401, "请先登录")
				return
			}

			// Bearer xxx
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				common.Fail(w, 401, "Authorization 格式错误，需为 Bearer xxx")
				return
			}

			tokenStr := parts[1]

			// 解析 token（使用你自己的 utils.ParseToken）
			claims, err := utils.ParseToken(c.Auth.AccessSecret, tokenStr)
			if err != nil {
				common.Fail(w, 401, "登录已过期或 token 无效")
				return
			}

			// 把用户信息塞进 context，后面所有 logic 都能直接拿
			ctx := utils.WithUserCtx(r.Context(), claims)
			r = r.WithContext(ctx)

			// 放行
			next(w, r)
		}
	}
}
