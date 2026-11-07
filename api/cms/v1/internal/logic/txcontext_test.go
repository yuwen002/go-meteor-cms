package logic

import (
	"context"
	"testing"
	"time"
)

type txCtxKey struct{}

func TestNewTxContextDetachesRequestCancellation(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.WithValue(context.Background(), txCtxKey{}, "trace-1"))
	defer parentCancel()

	ctx := parent
	cancel := newTxContext(&ctx)
	defer cancel()

	parentCancel()
	select {
	case <-ctx.Done():
		t.Fatalf("写事务 ctx 不应继承请求取消信号，收到: %v", ctx.Err())
	default:
	}

	if got := ctx.Value(txCtxKey{}); got != "trace-1" {
		t.Fatalf("写事务 ctx 应保留 trace/claims 等 Value，得到: %v", got)
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("写事务 ctx 必须有独立超时上限")
	}
	if remaining := time.Until(deadline); remaining <= 0 || remaining > transactionTimeout {
		t.Fatalf("写事务超时不合理: %v", remaining)
	}
}

func TestNewTxContextCancelReleases(t *testing.T) {
	ctx := context.Background()
	cancel := newTxContext(&ctx)
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatalf("defer cancel() 应释放事务 ctx，得到: %v", ctx.Err())
	}
}
