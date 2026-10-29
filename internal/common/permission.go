package common

// 权限类型，与 internal/seed 中的权限清单、权限中间件中的路由匹配逻辑保持一致
const (
	// PermissionTypeMenu 菜单（前端路由）
	PermissionTypeMenu = 1
	// PermissionTypeButton 按钮/功能点（前端据此控制按钮显隐）
	PermissionTypeButton = 2
	// PermissionTypeAPI 后端 API（权限中间件按此类型匹配路由进行鉴权）
	PermissionTypeAPI = 3
)
