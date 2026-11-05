package tokenversion

import (
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// keyPrefix 用户 token 版本号在 Redis 中的 key 前缀。
// token 版本号用于「全局踢人 / 批量吊销会话」：版本自增后，所有携带旧版本号的令牌立即失效。
const keyPrefix = "auth:tv:"

func key(userID int64) string {
	return keyPrefix + strconv.FormatInt(userID, 10)
}

// Get 读取用户当前的 token 版本号；key 不存在时返回 0（视为初始版本）。
// rdb 为 nil（未配置 Redis）时直接返回 0，由调用方降级为「不校验版本」。
func Get(ctx context.Context, rdb *redis.Redis, userID int64) (int64, error) {
	if rdb == nil {
		return 0, nil
	}
	val, err := rdb.Get(key(userID))
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	// go-zero 的 Redis.Get 在 key 不存在时返回 ("", nil)（内部已把 redis.Nil 转成空串），
	// 必须显式视为初始版本 0；否则 ParseInt("") 会误报错误，被 JWT 中间件当成
	// 「Redis 异常」走降级放行并刷错误日志，跳过版本比对。
	if val == "" {
		return 0, nil
	}
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Bump 将用户 token 版本号自增 1 并返回新值。
// 触发场景：禁用/删除账号、修改密码等——使该用户所有旧令牌立即失效（全局踢人）。
// rdb 为 nil（未配置 Redis）时为 no-op，返回 0。
func Bump(ctx context.Context, rdb *redis.Redis, userID int64) (int64, error) {
	if rdb == nil {
		return 0, nil
	}
	return rdb.Incr(key(userID))
}
