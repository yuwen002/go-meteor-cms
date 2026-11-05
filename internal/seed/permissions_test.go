package seed

import (
	"strings"
	"testing"
)

// 权限中间件按 method + api_path 查 Type=3 记录，重复配置会判定为歧义并直接拒绝请求。
func TestAPIPermissionsAreUnique(t *testing.T) {
	seen := make(map[string]string)
	for _, p := range DefaultPermissions {
		if p.Type != PermissionTypeAPI {
			continue
		}
		if p.Method == "" || p.APIPath == "" {
			t.Errorf("权限 %s 缺少 Method 或 APIPath", p.Permission)
			continue
		}
		key := strings.ToUpper(p.Method) + " " + p.APIPath
		if prev, ok := seen[key]; ok {
			t.Errorf("API 权限路径重复：%s（%s 与 %s）", key, prev, p.Permission)
		}
		seen[key] = p.Permission
	}
}

// 文章模块的每个受保护接口都必须在种子里有一条 Type=3 记录，
// 否则非超级管理员会被 PermissionMiddleware 以「未登记权限」拒绝。
func TestContentAPIPermissionsComplete(t *testing.T) {
	want := []string{
		"cms:article:list", "cms:article:detail", "cms:article:add", "cms:article:edit",
		"cms:article:delete", "cms:article:publish", "cms:article:withdraw", "cms:article:toggle_top",
		"cms:category:list", "cms:category:add", "cms:category:edit", "cms:category:delete",
		"cms:tag:list", "cms:tag:add", "cms:tag:edit", "cms:tag:delete",
	}
	apis := make(map[string]bool)
	buttons := make(map[string]bool)
	menus := make(map[string]bool)
	for _, p := range DefaultPermissions {
		switch p.Type {
		case PermissionTypeAPI:
			apis[p.Permission] = true
		case PermissionTypeButton:
			buttons[p.Permission] = true
		case PermissionTypeMenu:
			menus[p.Permission] = true
		}
	}
	for _, code := range want {
		if !apis[code] {
			t.Errorf("缺少 Type=3 API 权限：%s", code)
		}
		if !buttons[code] {
			t.Errorf("缺少同标识的 Type=2 按钮权限：%s", code)
		}
	}
	for _, code := range []string{"cms:group:content", "cms:article:list", "cms:category:list", "cms:tag:list"} {
		if !menus[code] {
			t.Errorf("缺少菜单：%s", code)
		}
	}
}

// 菜单的 ParentName 必须能在菜单列表里找到，且父级排在子级之前，
// 否则 UpsertPermissions 解析不到 parentID，子菜单会挂在根节点。
func TestMenuParentsDeclaredBeforeChildren(t *testing.T) {
	declared := make(map[string]bool)
	for _, p := range DefaultPermissions {
		if p.Type != PermissionTypeMenu {
			continue
		}
		if p.ParentName != "" && !declared[p.ParentName] {
			t.Errorf("菜单 %s 的父级 %s 尚未声明", p.Name, p.ParentName)
		}
		declared[p.Name] = true
	}
}
