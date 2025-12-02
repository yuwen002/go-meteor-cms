package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/tokenblacklist"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type JwtMiddleware struct {
	config *config.Config
	db     *ent.Client
}

func NewJwtMiddleware(c *config.Config, db *ent.Client) rest.Middleware {
	m := &JwtMiddleware{
		config: c,
		db:     db,
	}
	return m.Handle
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.Fail(w, common.ErrUnauthorized, common.GetErrorMessage(common.ErrUnauthorized))
			return
		}

		// Bearer xxx
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			common.Fail(w, common.ErrInvalidToken, common.GetErrorMessage(common.ErrInvalidToken))
			return
		}

		tokenStr := parts[1]

		// 生成token的哈希值
		tokenHash := utils.GenerateTokenHash(tokenStr)

		// 检查token是否在黑名单中（使用token_hash进行查询）
		exists, err := m.db.TokenBlacklist.
			Query().
			Where(tokenblacklist.TokenHash(tokenHash)).
			Exist(r.Context())

		if err != nil {
			logx.Errorf("查询token黑名单失败: %v", err)
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}
		if exists {
			common.Fail(w, common.ErrInvalidToken, common.GetErrorMessage(common.ErrInvalidToken))
			return
		}

		// 测试令牌，直接放行
		if tokenStr == "123456" {
			// 创建一个默认的 claims 对象，避免空指针
			claims := map[string]interface{}{
				"user_id":  int64(1), // 测试用户ID
				"username": "admin",
			}
			ctx := utils.WithUserCtx(r.Context(), claims)
			// 存储token到context
			ctx = context.WithValue(ctx, "token", tokenStr)
			r = r.WithContext(ctx)
			next(w, r)
			return
		}

		// 解析 token（使用你自己的 utils.ParseToken）
		claims, err := utils.ParseToken(m.config.Auth.AccessSecret, tokenStr)
		if err != nil {
			common.Fail(w, common.ErrInvalidToken, common.GetErrorMessage(common.ErrInvalidToken))
			return
		}

		// 把用户信息和token塞进 context，后面所有 logic 都能直接拿
		ctx := utils.WithUserCtx(r.Context(), claims)
		// 存储token到context
		ctx = context.WithValue(ctx, "token", tokenStr)
		r = r.WithContext(ctx)

		// 放行
		next(w, r)
	}
}
