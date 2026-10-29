package utils

import "testing"

func TestNormalizePage(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		defaultSize int
		wantOffset  int
		wantLimit   int
	}{
		{"正常第一页", 1, 10, 10, 0, 10},
		{"正常第三页", 3, 20, 10, 40, 20},
		{"页码为0时纠正", 0, 10, 10, 0, 10},
		{"页码为负时纠正", -5, 10, 10, 0, 10},
		{"每页数量为0时取默认", 1, 0, 15, 0, 15},
		{"每页数量为负时取默认", 1, -3, 15, 0, 15},
		{"每页数量超上限截断", 1, 10000, 10, 0, MaxPageSize},
		{"每页数量正好等于上限", 2, MaxPageSize, 10, MaxPageSize, MaxPageSize},
		{"默认值为0时回落到全局默认", 1, 0, 0, 0, DefaultPageSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, limit := NormalizePage(tt.page, tt.pageSize, tt.defaultSize)
			if offset != tt.wantOffset || limit != tt.wantLimit {
				t.Fatalf("NormalizePage(%d,%d,%d) = (%d,%d), want (%d,%d)",
					tt.page, tt.pageSize, tt.defaultSize, offset, limit, tt.wantOffset, tt.wantLimit)
			}
		})
	}
}

// TestNormalizePagination 回写的页码应与最终生效的每页数量一致
func TestNormalizePagination(t *testing.T) {
	page, pageSize := NormalizePagination(3, 500, 10)
	if pageSize != MaxPageSize {
		t.Fatalf("每页数量应被截断到 %d, got %d", MaxPageSize, pageSize)
	}
	if page != 3 {
		t.Fatalf("页码应保留为 3, got %d", page)
	}

	// 非法入参不应导致除零，limit 恒为正数
	offset, limit := NormalizePage(0, 0, 0)
	if limit <= 0 {
		t.Fatalf("limit 必须为正数, got %d", limit)
	}
	if page := offset/limit + 1; page != 1 {
		t.Fatalf("非法入参应纠正为第 1 页, got %d", page)
	}
}
