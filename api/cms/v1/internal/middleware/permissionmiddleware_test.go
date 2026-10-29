package middleware

import (
	"net/http"
	"testing"
)

// TestMatchPermissionPath 覆盖权限中间件的路由匹配规则：
// 静态段与动态段需逐段比对，且路径参数段不能匹配空值。
func TestMatchPermissionPath(t *testing.T) {
	tests := []struct {
		name     string
		template string
		request  string
		wantHit  bool
	}{
		{"完全匹配", "/admin/roles", "/admin/roles", true},
		{"路径参数匹配", "/admin/roles/:id", "/admin/roles/12", true},
		{"末位斜杠容错", "/admin/roles", "/admin/roles/", true},
		{"缺少前置斜杠容错", "admin/roles", "/admin/roles", true},
		{"段数不同不匹配", "/admin/roles", "/admin/roles/1/permissions", false},
		{"静态段不同不匹配", "/admin/roles", "/admin/users", false},
		{"路径参数段为空不匹配", "/admin/roles/:id", "/admin/roles/", false},
		{"多级路径参数匹配", "/admin/users/:id/roles/:role_id", "/admin/users/1/roles/9", true},
		{"嵌套路径参数不匹配", "/admin/users/:id/roles", "/admin/users/1/permissions", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := matchPermissionPath(tt.template, tt.request)
			if ok != tt.wantHit {
				t.Fatalf("matchPermissionPath(%q, %q) = %v, want %v", tt.template, tt.request, ok, tt.wantHit)
			}
		})
	}
}

// TestMatchPermissionPathPriority 静态路由应优先于同名动态路由：
// 例如 GET /admin/departments/tree 必须命中 tree 而不是 /admin/departments/:id
func TestMatchPermissionPathPriority(t *testing.T) {
	// /admin/departments/tree 共 3 段且全部为静态段
	staticHits, ok := matchPermissionPath("/admin/departments/tree", "/admin/departments/tree")
	if !ok || staticHits != 3 {
		t.Fatalf("静态路径命中数异常: hits=%d ok=%v, want 3/true", staticHits, ok)
	}

	dynamicHits, ok := matchPermissionPath("/admin/departments/:id", "/admin/departments/tree")
	if !ok || dynamicHits != 2 {
		t.Fatalf("动态路径命中数异常: hits=%d ok=%v, want 2/true", dynamicHits, ok)
	}

	if !(staticHits > dynamicHits) {
		t.Fatalf("静态路由优先级应高于动态路由: static=%d dynamic=%d", staticHits, dynamicHits)
	}
}

// TestIsPermissionBypassRoute 放行路由的判断必须与方法、路径完全对应
func TestIsPermissionBypassRoute(t *testing.T) {
	if !isPermissionBypassRoute(http.MethodGet, "/admin/profile/permissions") {
		t.Fatal("我的菜单接口必须放行，否则普通管理员无法取回自己的菜单")
	}
	if !isPermissionBypassRoute(http.MethodPut, "/admin/admin-users/me/change-password") {
		t.Fatal("修改自己的密码必须放行")
	}
	if isPermissionBypassRoute(http.MethodDelete, "/admin/admin-users/1") {
		t.Fatal("删除管理员不应被放行")
	}
	if isPermissionBypassRoute(http.MethodGet, "/admin/admin-users") {
		t.Fatal("管理员列表不应被放行")
	}
}

// TestNormalizePath 标准化路径用于消除首尾斜杠差异
func TestNormalizePath(t *testing.T) {
	cases := map[string]string{
		"":                "/",
		"/":               "/",
		"admin/roles":     "/admin/roles",
		"/admin/roles/":   "/admin/roles",
		"/admin/roles//":  "/admin/roles/",
		"  /admin/roles ": "/admin/roles",
	}
	for input, want := range cases {
		if got := normalizePath(input); got != want {
			t.Fatalf("normalizePath(%q) = %q, want %q", input, got, want)
		}
	}
}
