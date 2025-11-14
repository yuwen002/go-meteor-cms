package middleware

import (
	"net/http"
	"strings"

	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func JwtMiddleware(secret string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.Error(w, http.ErrNoCookie)
				return
			}

			// 格式: Bearer <token>
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				httpx.Error(w, http.ErrNoCookie)
				return
			}

			tokenStr := parts[1]
			claims, err := utils.ParseToken(secret, tokenStr)
			if err != nil {
				httpx.Error(w, err)
				return
			}

			// 将用户信息存入 context
			ctx := r.Context()
			ctx = utils.WithUserCtx(ctx, claims)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}
