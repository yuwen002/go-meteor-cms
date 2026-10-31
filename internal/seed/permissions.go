package seed

import (
	"context"
	"log"

	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminpermission"
	"github.com/yuwen002/go-meteor-cms/internal/common"
)

// 权限类型，定义收敛在 common 包，避免各层各自定义
const (
	// PermissionTypeMenu 菜单
	PermissionTypeMenu = common.PermissionTypeMenu
	// PermissionTypeButton 按钮/功能点
	PermissionTypeButton = common.PermissionTypeButton
	// PermissionTypeAPI 后端 API，权限中间件按此类型匹配路由
	PermissionTypeAPI = common.PermissionTypeAPI
)

// PermissionSpec 权限定义（数据驱动，新增后端接口时只需在此追加一条 Type=API）
type PermissionSpec struct {
	// Name 权限名称
	Name string
	// Permission 权限标识，同一操作在按钮与 API 两种类型上使用相同标识，
	// 便于把「按钮权限」与「接口权限」一次性授权给角色
	Permission string
	// Type 权限类型：1 菜单 2 按钮 3 API
	Type int
	// ParentName 父级权限名称，Type=1 时为空
	ParentName string
	// Path 前端路由路径（Type=1）
	Path string
	// Component 前端组件路径（Type=1）
	Component string
	// Icon 菜单图标（Type=1）
	Icon string
	// Method HTTP 方法（Type=2/3）
	Method string
	// APIPath 后端接口路径（Type=3），需与 cms.api 中的路由保持一致，支持 :id 形式的路径参数
	APIPath string
	// Sort 排序，从小到大
	Sort int
}

// DefaultPermissions 系统内置权限清单。
//
// 约定：
//   - Type=1 菜单使用「列表」权限标识（如 system:role:list），菜单可见即拥有该标识；
//   - Type=2 按钮与 Type=3 API 使用完全相同的权限标识，保证授权一次即可同时生效；
//   - Type=3 的 APIPath 必须与 cms.api 的真实路由一致，否则非超级管理员会被权限中间件拒绝。
var DefaultPermissions = []PermissionSpec{
	// ---------- 菜单（Type=1） ----------
	{Name: "控制台", Permission: "system:dashboard:view", Type: PermissionTypeMenu,
		Path: "/dashboard", Component: "views/dashboard/index.vue", Icon: "dashboard", Sort: 1},
	{Name: "管理员管理", Permission: "system:admin_user:list", Type: PermissionTypeMenu,
		Path: "/system/users", Component: "views/admin/user/index.vue", Icon: "user", Sort: 10},
	{Name: "部门管理", Permission: "system:dept:list", Type: PermissionTypeMenu,
		Path: "/system/depts", Component: "views/admin/dept/index.vue", Icon: "dept", Sort: 20},
	{Name: "角色管理", Permission: "system:role:list", Type: PermissionTypeMenu,
		Path: "/system/roles", Component: "views/admin/role/index.vue", Icon: "role", Sort: 30},
	{Name: "权限管理", Permission: "system:permission:list", Type: PermissionTypeMenu,
		Path: "/system/permissions", Component: "views/admin/permission/index.vue", Icon: "permission", Sort: 40},
	{Name: "操作日志", Permission: "system:log:operation_list", Type: PermissionTypeMenu,
		Path: "/system/operation-log", Component: "views/system/log/operation.vue", Icon: "document", Sort: 50},
	{Name: "登录日志", Permission: "system:log:login_list", Type: PermissionTypeMenu,
		Path: "/system/login-log", Component: "views/system/log/login.vue", Icon: "clock", Sort: 55},

	// ---------- 控制台（菜单与统计接口同用一个标识） ----------
	{Name: "查看控制台", Permission: "system:dashboard:view", Type: PermissionTypeButton, ParentName: "控制台", Sort: 2},
	{Name: "查看控制台", Permission: "system:dashboard:view", Type: PermissionTypeAPI, ParentName: "控制台", Method: "GET", APIPath: "/admin/dashboard/stats", Sort: 2},

	// ---------- 管理员管理（Type=2 按钮 / Type=3 API） ----------
	{Name: "查看管理员", Permission: "system:admin_user:list", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 11},
	{Name: "查看管理员", Permission: "system:admin_user:list", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "GET", APIPath: "/admin/admin-users", Sort: 11},
	{Name: "管理员详情", Permission: "system:admin_user:detail", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 12},
	{Name: "管理员详情", Permission: "system:admin_user:detail", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "GET", APIPath: "/admin/admin-users/:id", Sort: 12},
	{Name: "新增管理员", Permission: "system:admin_user:add", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 13},
	{Name: "新增管理员", Permission: "system:admin_user:add", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "POST", APIPath: "/admin/admin-users/create", Sort: 13},
	{Name: "编辑管理员", Permission: "system:admin_user:edit", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 14},
	{Name: "编辑管理员", Permission: "system:admin_user:edit", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "PUT", APIPath: "/admin/admin-users/:id", Sort: 14},
	{Name: "删除管理员", Permission: "system:admin_user:delete", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 15},
	{Name: "删除管理员", Permission: "system:admin_user:delete", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "DELETE", APIPath: "/admin/admin-users/:id", Sort: 15},
	{Name: "禁用管理员", Permission: "system:admin_user:disable", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 16},
	{Name: "禁用管理员", Permission: "system:admin_user:disable", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "PUT", APIPath: "/admin/admin-users/:id/disable", Sort: 16},
	{Name: "启用管理员", Permission: "system:admin_user:enable", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 17},
	{Name: "启用管理员", Permission: "system:admin_user:enable", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "PUT", APIPath: "/admin/admin-users/:id/enable", Sort: 17},
	{Name: "重置管理员密码", Permission: "system:admin_user:reset_pwd", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 18},
	{Name: "重置管理员密码", Permission: "system:admin_user:reset_pwd", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "PUT", APIPath: "/admin/admin-users/reset-password", Sort: 18},
	{Name: "绑定部门", Permission: "system:admin_user:bind_dept", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 19},
	{Name: "绑定部门", Permission: "system:admin_user:bind_dept", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "PUT", APIPath: "/admin/admin-users/:id/bind-department", Sort: 19},
	{Name: "查看管理员角色", Permission: "system:admin_user:role_list", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 21},
	{Name: "查看管理员角色", Permission: "system:admin_user:role_list", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "GET", APIPath: "/admin/admin-users/:id/roles", Sort: 21},
	{Name: "分配管理员角色", Permission: "system:admin_user:role_assign", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 22},
	{Name: "分配管理员角色", Permission: "system:admin_user:role_assign", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "PUT", APIPath: "/admin/admin-users/:id/roles", Sort: 22},
	{Name: "移除管理员角色", Permission: "system:admin_user:role_remove", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 23},
	{Name: "移除管理员角色", Permission: "system:admin_user:role_remove", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "DELETE", APIPath: "/admin/admin-users/:id/roles/:role_id", Sort: 23},

	// ---------- 部门管理 ----------
	{Name: "查看部门", Permission: "system:dept:list", Type: PermissionTypeButton, ParentName: "部门管理", Sort: 24},
	{Name: "查看部门", Permission: "system:dept:list", Type: PermissionTypeAPI, ParentName: "部门管理", Method: "GET", APIPath: "/admin/departments/tree", Sort: 24},
	{Name: "部门详情", Permission: "system:dept:detail", Type: PermissionTypeButton, ParentName: "部门管理", Sort: 25},
	{Name: "部门详情", Permission: "system:dept:detail", Type: PermissionTypeAPI, ParentName: "部门管理", Method: "GET", APIPath: "/admin/departments/:id", Sort: 25},
	{Name: "新增部门", Permission: "system:dept:add", Type: PermissionTypeButton, ParentName: "部门管理", Sort: 26},
	{Name: "新增部门", Permission: "system:dept:add", Type: PermissionTypeAPI, ParentName: "部门管理", Method: "POST", APIPath: "/admin/departments", Sort: 26},
	{Name: "编辑部门", Permission: "system:dept:edit", Type: PermissionTypeButton, ParentName: "部门管理", Sort: 27},
	{Name: "编辑部门", Permission: "system:dept:edit", Type: PermissionTypeAPI, ParentName: "部门管理", Method: "PUT", APIPath: "/admin/departments/:id", Sort: 27},
	{Name: "删除部门", Permission: "system:dept:delete", Type: PermissionTypeButton, ParentName: "部门管理", Sort: 28},
	{Name: "删除部门", Permission: "system:dept:delete", Type: PermissionTypeAPI, ParentName: "部门管理", Method: "DELETE", APIPath: "/admin/departments/:id", Sort: 28},

	// ---------- 角色管理 ----------
	{Name: "查看角色", Permission: "system:role:list", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 31},
	{Name: "查看角色", Permission: "system:role:list", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "GET", APIPath: "/admin/roles", Sort: 31},
	{Name: "角色详情", Permission: "system:role:detail", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 32},
	{Name: "角色详情", Permission: "system:role:detail", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "GET", APIPath: "/admin/roles/:id", Sort: 32},
	{Name: "新增角色", Permission: "system:role:add", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 33},
	{Name: "新增角色", Permission: "system:role:add", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "POST", APIPath: "/admin/roles", Sort: 33},
	{Name: "编辑角色", Permission: "system:role:edit", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 34},
	{Name: "编辑角色", Permission: "system:role:edit", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "PUT", APIPath: "/admin/roles/:id", Sort: 34},
	{Name: "删除角色", Permission: "system:role:delete", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 35},
	{Name: "删除角色", Permission: "system:role:delete", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "DELETE", APIPath: "/admin/roles/:id", Sort: 35},
	{Name: "查看角色权限", Permission: "system:role:perm_list", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 36},
	{Name: "查看角色权限", Permission: "system:role:perm_list", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "GET", APIPath: "/admin/roles/:id/permissions", Sort: 36},
	{Name: "查看角色权限(含名称)", Permission: "system:role:perm_list_names", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 36},
	{Name: "查看角色权限(含名称)", Permission: "system:role:perm_list_names", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "GET", APIPath: "/admin/roles/:id/permissions-with-names", Sort: 36},
	{Name: "分配角色权限", Permission: "system:role:perm_assign", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 37},
	{Name: "分配角色权限", Permission: "system:role:perm_assign", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "PUT", APIPath: "/admin/roles/:id/permissions", Sort: 37},
	{Name: "查看角色部门", Permission: "system:role:dept_list", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 38},
	{Name: "查看角色部门", Permission: "system:role:dept_list", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "GET", APIPath: "/admin/roles/:id/departments", Sort: 38},
	{Name: "设置角色部门", Permission: "system:role:dept_assign", Type: PermissionTypeButton, ParentName: "角色管理", Sort: 39},
	{Name: "设置角色部门", Permission: "system:role:dept_assign", Type: PermissionTypeAPI, ParentName: "角色管理", Method: "PUT", APIPath: "/admin/roles/:id/departments", Sort: 39},

	// ---------- 权限管理 ----------
	{Name: "查看权限", Permission: "system:permission:list", Type: PermissionTypeButton, ParentName: "权限管理", Sort: 41},
	{Name: "查看权限", Permission: "system:permission:list", Type: PermissionTypeAPI, ParentName: "权限管理", Method: "GET", APIPath: "/admin/permissions", Sort: 41},
	{Name: "权限树", Permission: "system:permission:tree", Type: PermissionTypeButton, ParentName: "权限管理", Sort: 42},
	{Name: "权限树", Permission: "system:permission:tree", Type: PermissionTypeAPI, ParentName: "权限管理", Method: "GET", APIPath: "/admin/permissions/tree", Sort: 42},
	{Name: "权限详情", Permission: "system:permission:detail", Type: PermissionTypeButton, ParentName: "权限管理", Sort: 43},
	{Name: "权限详情", Permission: "system:permission:detail", Type: PermissionTypeAPI, ParentName: "权限管理", Method: "GET", APIPath: "/admin/permissions/:id", Sort: 43},
	{Name: "新增权限", Permission: "system:permission:add", Type: PermissionTypeButton, ParentName: "权限管理", Sort: 44},
	{Name: "新增权限", Permission: "system:permission:add", Type: PermissionTypeAPI, ParentName: "权限管理", Method: "POST", APIPath: "/admin/permissions", Sort: 44},
	{Name: "编辑权限", Permission: "system:permission:edit", Type: PermissionTypeButton, ParentName: "权限管理", Sort: 45},
	{Name: "编辑权限", Permission: "system:permission:edit", Type: PermissionTypeAPI, ParentName: "权限管理", Method: "PUT", APIPath: "/admin/permissions/:id", Sort: 45},
	{Name: "删除权限", Permission: "system:permission:delete", Type: PermissionTypeButton, ParentName: "权限管理", Sort: 46},
	{Name: "删除权限", Permission: "system:permission:delete", Type: PermissionTypeAPI, ParentName: "权限管理", Method: "DELETE", APIPath: "/admin/permissions/:id", Sort: 46},

	// ---------- 日志管理 ----------
	{Name: "查看操作日志", Permission: "system:log:operation_list", Type: PermissionTypeButton, ParentName: "操作日志", Sort: 51},
	{Name: "查看操作日志", Permission: "system:log:operation_list", Type: PermissionTypeAPI, ParentName: "操作日志", Method: "GET", APIPath: "/admin/logs/operations", Sort: 51},
	{Name: "查看登录日志", Permission: "system:log:login_list", Type: PermissionTypeButton, ParentName: "登录日志", Sort: 52},
	{Name: "查看登录日志", Permission: "system:log:login_list", Type: PermissionTypeAPI, ParentName: "登录日志", Method: "GET", APIPath: "/admin/logs/logins", Sort: 52},
	{Name: "上传文件", Permission: "system:upload:file", Type: PermissionTypeButton, ParentName: "管理员管理", Sort: 25},
	{Name: "上传文件", Permission: "system:upload:file", Type: PermissionTypeAPI, ParentName: "管理员管理", Method: "POST", APIPath: "/admin/upload", Sort: 25},

	// ---------- 文件管理 ----------
	{Name: "文件管理", Permission: "system:file:list", Type: PermissionTypeMenu,
		Path: "/system/files", Component: "views/system/file/index.vue", Icon: "file", Sort: 60},
	{Name: "查看文件", Permission: "system:file:list", Type: PermissionTypeButton, ParentName: "文件管理", Sort: 61},
	{Name: "查看文件", Permission: "system:file:list", Type: PermissionTypeAPI, ParentName: "文件管理", Method: "GET", APIPath: "/admin/files", Sort: 61},
	{Name: "重命名文件", Permission: "system:file:rename", Type: PermissionTypeButton, ParentName: "文件管理", Sort: 62},
	{Name: "重命名文件", Permission: "system:file:rename", Type: PermissionTypeAPI, ParentName: "文件管理", Method: "PUT", APIPath: "/admin/files/rename", Sort: 62},
	{Name: "删除文件", Permission: "system:file:delete", Type: PermissionTypeButton, ParentName: "文件管理", Sort: 63},
	{Name: "删除文件", Permission: "system:file:delete", Type: PermissionTypeAPI, ParentName: "文件管理", Method: "POST", APIPath: "/admin/files/delete", Sort: 63},

	// ---------- 系统管理（分组，用于演示多级菜单层级） ----------
	{Name: "系统管理", Permission: "system:group:system", Type: PermissionTypeMenu,
		Path: "", Component: "", Icon: "setting", Sort: 65},

	// ---------- 数据字典 ----------
	{Name: "字典管理", Permission: "system:dict:list", Type: PermissionTypeMenu,
		Path: "/system/dicts", Component: "views/system/dict/index.vue", Icon: "dict", Sort: 70, ParentName: "系统管理"},
	{Name: "查看字典", Permission: "system:dict:list", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 71},
	{Name: "查看字典", Permission: "system:dict:list", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "GET", APIPath: "/admin/dicts", Sort: 71},
	{Name: "新增字典", Permission: "system:dict:add", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 72},
	{Name: "新增字典", Permission: "system:dict:add", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "POST", APIPath: "/admin/dicts", Sort: 72},
	{Name: "编辑字典", Permission: "system:dict:edit", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 73},
	{Name: "编辑字典", Permission: "system:dict:edit", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "PUT", APIPath: "/admin/dicts/:id", Sort: 73},
	{Name: "删除字典", Permission: "system:dict:delete", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 74},
	{Name: "删除字典", Permission: "system:dict:delete", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "DELETE", APIPath: "/admin/dicts/:id", Sort: 74},
	{Name: "查看字典项", Permission: "system:dict_item:list", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 75},
	{Name: "查看字典项", Permission: "system:dict_item:list", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "GET", APIPath: "/admin/dicts/items", Sort: 75},
	{Name: "新增字典项", Permission: "system:dict_item:add", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 76},
	{Name: "新增字典项", Permission: "system:dict_item:add", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "POST", APIPath: "/admin/dicts/items", Sort: 76},
	{Name: "编辑字典项", Permission: "system:dict_item:edit", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 77},
	{Name: "编辑字典项", Permission: "system:dict_item:edit", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "PUT", APIPath: "/admin/dicts/items/:id", Sort: 77},
	{Name: "删除字典项", Permission: "system:dict_item:delete", Type: PermissionTypeButton, ParentName: "字典管理", Sort: 78},
	{Name: "删除字典项", Permission: "system:dict_item:delete", Type: PermissionTypeAPI, ParentName: "字典管理", Method: "DELETE", APIPath: "/admin/dicts/items/:id", Sort: 78},

	// ---------- 系统配置 ----------
	{Name: "系统配置", Permission: "system:config:list", Type: PermissionTypeMenu,
		Path: "/system/configs", Component: "views/system/config/index.vue", Icon: "setting", Sort: 80, ParentName: "系统管理"},
	{Name: "查看配置", Permission: "system:config:list", Type: PermissionTypeButton, ParentName: "系统配置", Sort: 81},
	{Name: "查看配置", Permission: "system:config:list", Type: PermissionTypeAPI, ParentName: "系统配置", Method: "GET", APIPath: "/admin/configs", Sort: 81},
	{Name: "新增配置", Permission: "system:config:add", Type: PermissionTypeButton, ParentName: "系统配置", Sort: 82},
	{Name: "新增配置", Permission: "system:config:add", Type: PermissionTypeAPI, ParentName: "系统配置", Method: "POST", APIPath: "/admin/configs", Sort: 82},
	{Name: "编辑配置", Permission: "system:config:edit", Type: PermissionTypeButton, ParentName: "系统配置", Sort: 83},
	{Name: "编辑配置", Permission: "system:config:edit", Type: PermissionTypeAPI, ParentName: "系统配置", Method: "PUT", APIPath: "/admin/configs/:id", Sort: 83},
	{Name: "删除配置", Permission: "system:config:delete", Type: PermissionTypeButton, ParentName: "系统配置", Sort: 84},
	{Name: "删除配置", Permission: "system:config:delete", Type: PermissionTypeAPI, ParentName: "系统配置", Method: "DELETE", APIPath: "/admin/configs/:id", Sort: 84},
}

// legacyPermissions 历史版本的权限标识，已被 DefaultPermissions 中的新标识取代。
// 首次运行迁移生成的旧记录若不清理，会出现在权限列表里成为无法解释的孤儿数据。
var legacyPermissions = []struct {
	Type       int
	Permission string
}{
	{PermissionTypeMenu, "dashboard:view"},
}

// CleanupLegacyPermissions 清理已被废弃的历史权限记录（幂等）
func CleanupLegacyPermissions(ctx context.Context, client *ent.Client) error {
	for _, legacy := range legacyPermissions {
		deleted, err := client.AdminPermission.
			Delete().
			Where(
				adminpermission.TypeEQ(legacy.Type),
				adminpermission.PermissionEQ(legacy.Permission),
			).
			Exec(ctx)
		if err != nil {
			return err
		}
		if deleted > 0 {
			log.Printf("🧹 清理废弃权限 %s %d 条", legacy.Permission, deleted)
		}
	}
	return nil
}

// UpsertPermissions 幂等写入权限清单，返回全部权限实体。
//
// 同一 (type, permission) 只会存在一条记录：已存在时更新展示信息与 API 属性，不存在时创建。
// 因此该过程可以随每次迁移重复执行，不会因唯一索引而中断。
func UpsertPermissions(ctx context.Context, client *ent.Client, specs []PermissionSpec) ([]*ent.AdminPermission, error) {
	result := make([]*ent.AdminPermission, 0, len(specs))
	// 菜单名称 -> 实体 ID，用于解析 ParentName
	parentIndex := make(map[string]int64, len(specs))

	writePermission := func(spec PermissionSpec, parentID int64) (*ent.AdminPermission, error) {
		existing, err := client.AdminPermission.
			Query().
			Where(
				adminpermission.TypeEQ(spec.Type),
				adminpermission.PermissionEQ(spec.Permission),
			).
			First(ctx)
		if err == nil {
			// 已存在则同步展示信息与接口属性
			update := existing.Update().
				SetName(spec.Name).
				SetSort(spec.Sort).
				SetIsActive(true)
			if spec.Type == PermissionTypeMenu {
				update = update.SetPath(spec.Path).SetComponent(spec.Component).SetIcon(spec.Icon)
				if parentID > 0 {
					update = update.SetParentID(parentID)
				}
			} else {
				update = update.SetMethod(spec.Method).SetAPIPath(spec.APIPath)
				if parentID > 0 {
					update = update.SetParentID(parentID)
				}
			}
			return update.Save(ctx)
		}
		if !ent.IsNotFound(err) {
			return nil, err
		}

		create := client.AdminPermission.
			Create().
			SetName(spec.Name).
			SetPermission(spec.Permission).
			SetType(spec.Type).
			SetSort(spec.Sort).
			SetIsActive(true)
		if spec.Type == PermissionTypeMenu {
			create = create.SetPath(spec.Path).SetComponent(spec.Component).SetIcon(spec.Icon)
			if parentID > 0 {
				create = create.SetParentID(parentID)
			}
		} else {
			create = create.SetMethod(spec.Method).SetAPIPath(spec.APIPath)
			if parentID > 0 {
				create = create.SetParentID(parentID)
			}
		}

		return create.Save(ctx)
	}

	// 1. 菜单（父级菜单需在切片中位于其子菜单之前，parentIndex 才能正确解析）
	for _, spec := range specs {
		if spec.Type != PermissionTypeMenu {
			continue
		}
		var parentID int64
		if spec.ParentName != "" {
			parentID = parentIndex[spec.ParentName]
		}
		p, err := writePermission(spec, parentID)
		if err != nil {
			return nil, err
		}
		parentIndex[spec.Name] = p.ID
		result = append(result, p)
	}

	// 2. 按钮与 API
	for _, spec := range specs {
		if spec.Type == PermissionTypeMenu {
			continue
		}
		var parentID int64
		if spec.ParentName != "" {
			parentID = parentIndex[spec.ParentName]
		}
		p, err := writePermission(spec, parentID)
		if err != nil {
			return nil, err
		}
		parentIndex[spec.Name] = p.ID
		result = append(result, p)
	}

	return result, nil
}
