package utils

// 分页常量
const (
	// DefaultPageSize 默认每页数量
	DefaultPageSize = 10
	// MaxPageSize 每页数量上限，防止一次性拉取全表
	MaxPageSize = 100
)

// NormalizePage 规范化分页参数
// 规则：
//   - page < 1 时取 1
//   - pageSize < 1 时取 defaultSize（<=0 时用 DefaultPageSize）
//   - pageSize > MaxPageSize 时截断到 MaxPageSize
//
// 返回规范后的 (offset, limit)，调用方直接用于 Ent 查询
func NormalizePage(page, pageSize, defaultSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if defaultSize < 1 {
		defaultSize = DefaultPageSize
	}
	if pageSize < 1 {
		pageSize = defaultSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return (page - 1) * pageSize, pageSize
}

// NormalizePagination 规范化分页参数，返回规范后的页码和每页数量
// 用于需要把 page/pageSize 回写给响应体的场景
func NormalizePagination(page, pageSize, defaultSize int) (int, int) {
	offset, limit := NormalizePage(page, pageSize, defaultSize)
	return offset/limit + 1, limit
}
