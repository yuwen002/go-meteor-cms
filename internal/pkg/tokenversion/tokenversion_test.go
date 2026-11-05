package tokenversion

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// newTestRedis 启动进程内 Redis（miniredis），使本包测试不依赖外部 Docker/Redis 服务。
func newTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Redis) {
	t.Helper()
	s := miniredis.RunT(t)
	rdb, err := redis.NewRedis(redis.RedisConf{
		Host: s.Addr(),
		Type: redis.NodeType,
	})
	if err != nil {
		t.Fatalf("NewRedis failed: %v", err)
	}
	return s, rdb
}

// TestGet_MissingKey 契约：key 不存在时视为初始版本 0，且不应返回错误。
// go-zero 的 Redis.Get 在 key 不存在时返回 ("", nil)（内部已吞掉 redis.Nil），
// 实现必须显式处理空串，否则 strconv.ParseInt("") 会被误报成 Redis 异常。
func TestGet_MissingKey(t *testing.T) {
	_, rdb := newTestRedis(t)

	got, err := Get(context.Background(), rdb, 1)
	if err != nil {
		t.Fatalf("key 不存在时不应返回错误, got err=%v", err)
	}
	if got != 0 {
		t.Fatalf("key 不存在时应返回初始版本 0, got %d", got)
	}
}

// TestBump_IncrementsAndGet 覆盖「全局踢人」核心：Bump 自增并把新版本写入 auth:tv:<uid>，
// 且返回值与随后 Get 的值一致（JWT 中间件正是拿该值比对 token 内的 tv）。
func TestBump_IncrementsAndGet(t *testing.T) {
	s, rdb := newTestRedis(t)
	ctx := context.Background()

	if got := key(1); got != "auth:tv:1" {
		t.Fatalf("版本 key = %q, want %q", got, "auth:tv:1")
	}

	for _, want := range []int64{1, 2, 3} {
		got, err := Bump(ctx, rdb, 1)
		if err != nil {
			t.Fatalf("Bump #%d failed: %v", want, err)
		}
		if got != want {
			t.Fatalf("Bump 返回 %d, want %d", got, want)
		}
		cur, err := Get(ctx, rdb, 1)
		if err != nil || cur != want {
			t.Fatalf("Bump 后 Get = (%d, %v), want (%d, nil)", cur, err, want)
		}
	}

	// 直接读 miniredis，确认真实写入的 key/value（而非仅内存自洽）
	if v, err := s.Get("auth:tv:1"); err != nil || v != "3" {
		t.Fatalf("auth:tv:1 = %q, %v; want \"3\"", v, err)
	}
}

// TestBump_KeyIsolatedByUser 不同用户版本号互不影响（避免踢 A 误伤 B）。
func TestBump_KeyIsolatedByUser(t *testing.T) {
	_, rdb := newTestRedis(t)
	ctx := context.Background()

	if _, err := Bump(ctx, rdb, 1); err != nil {
		t.Fatal(err)
	}
	got2, err := Get(ctx, rdb, 2)
	if err != nil || got2 != 0 {
		t.Fatalf("用户 2 版本 = (%d, %v), want (0, nil)", got2, err)
	}
}

// TestGet_InvalidValue 非数字版本值应返回错误，由调用方按「Redis 异常」降级放行。
func TestGet_InvalidValue(t *testing.T) {
	s, rdb := newTestRedis(t)
	if err := s.Set("auth:tv:1", "not-a-number"); err != nil {
		t.Fatal(err)
	}
	if got, err := Get(context.Background(), rdb, 1); err == nil {
		t.Fatalf("非法版本值应返回错误, got (%d, nil)", got)
	}
}

// TestNilRedis_Degrade 未配置 Redis 时：Get 返回 0、Bump 为 no-op，不 panic 不报错。
func TestNilRedis_Degrade(t *testing.T) {
	ctx := context.Background()
	if got, err := Get(ctx, nil, 1); err != nil || got != 0 {
		t.Fatalf("nil redis Get = (%d, %v), want (0, nil)", got, err)
	}
	if got, err := Bump(ctx, nil, 1); err != nil || got != 0 {
		t.Fatalf("nil redis Bump = (%d, %v), want (0, nil)", got, err)
	}
}
