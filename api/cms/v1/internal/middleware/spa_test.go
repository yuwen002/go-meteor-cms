package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newStaticDir 构造一个最小 SPA 产物目录，返回其路径。
func newStaticDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "dist")
	if err := os.MkdirAll(filepath.Join(dir, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("SPA_INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("APP_JS_BUNDLE"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func serveSPA(staticDir, method, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	SPAHandler(staticDir).ServeHTTP(rec, httptest.NewRequest(method, path, nil))
	return rec
}

// TestSPAHandler_HistoryFallback 前端使用 createWebHistory，/、/login、/dashboard、/system/*
// 等未匹配路径都必须回退到 index.html（否则刷新即 404）。
func TestSPAHandler_HistoryFallback(t *testing.T) {
	dir := newStaticDir(t)
	for _, path := range []string{"/", "/login", "/dashboard", "/system/roles", "/profile"} {
		rec := serveSPA(dir, http.MethodGet, path)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s = %d, want 200", path, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), "SPA_INDEX") {
			t.Fatalf("GET %s 未回退到 index.html, body=%q", path, rec.Body.String())
		}
	}
}

// TestSPAHandler_StaticFile 真实存在的产物文件应直接返回文件内容，而不是被 index.html 覆盖。
func TestSPAHandler_StaticFile(t *testing.T) {
	dir := newStaticDir(t)
	rec := serveSPA(dir, http.MethodGet, "/assets/app.js")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /assets/app.js = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "APP_JS_BUNDLE") {
		t.Fatalf("静态文件内容错误, body=%q", rec.Body.String())
	}
}

// TestSPAHandler_ApiPrefixesNotFallback /admin 是纯 API 前缀，/api、/uploads 同理：
// 这些路径不存在对应路由时必须保持 404，不能被 SPA 兜底成 200（否则前端会拿到 HTML 当接口响应）。
func TestSPAHandler_ApiPrefixesNotFallback(t *testing.T) {
	dir := newStaticDir(t)
	for _, path := range []string{"/admin/profile", "/admin/not-exist", "/api/foo", "/uploads/a.png"} {
		rec := serveSPA(dir, http.MethodGet, path)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("GET %s = %d, want 404（API 前缀不应回退）", path, rec.Code)
		}
	}
}

// TestSPAHandler_NonGET 非 GET 请求不做 SPA 兜底。
func TestSPAHandler_NonGET(t *testing.T) {
	dir := newStaticDir(t)
	if rec := serveSPA(dir, http.MethodPost, "/"); rec.Code != http.StatusNotFound {
		t.Fatalf("POST / = %d, want 404", rec.Code)
	}
}

// TestSPAHandler_EmptyStaticDir 本地开发（仅暴露 API）未配置 WEB_DIST 时全部走默认 404。
func TestSPAHandler_EmptyStaticDir(t *testing.T) {
	if rec := serveSPA("", http.MethodGet, "/"); rec.Code != http.StatusNotFound {
		t.Fatalf("staticDir 为空时 GET / = %d, want 404", rec.Code)
	}
}

// TestSPAHandler_NoIndex 产物目录缺失 index.html 时返回 404，不能用空文件冒充成功。
func TestSPAHandler_NoIndex(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "dist")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if rec := serveSPA(dir, http.MethodGet, "/dashboard"); rec.Code != http.StatusNotFound {
		t.Fatalf("缺 index.html 时 GET /dashboard = %d, want 404", rec.Code)
	}
}

// TestSPAHandler_PathTraversal ../ 不能读取产物目录之外的文件。
func TestSPAHandler_PathTraversal(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "dist")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("SPA_INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	const secret = "TOP_SECRET_OUTSIDE_DIST"
	if err := os.WriteFile(filepath.Join(root, "secret.txt"), []byte(secret), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/../secret.txt", "/..%2fsecret.txt"} {
		rec := serveSPA(dir, http.MethodGet, path)
		if strings.Contains(rec.Body.String(), secret) {
			t.Fatalf("GET %s 泄露了产物目录之外的文件, body=%q", path, rec.Body.String())
		}
	}
}

// TestIsInside 路径归属判定的边界：目录内允许、目录外与前缀同名目录拒绝。
func TestIsInside(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		name   string
		target string
		want   bool
	}{
		{"根本身", root, true},
		{"根下文件", filepath.Join(root, "index.html"), true},
		{"根下多级", filepath.Join(root, "assets", "app.js"), true},
		{"父级文件", filepath.Join(root, "..", "secret.txt"), false},
		{"前缀同名目录", root + "-evil", false},
		{"绝对路径", "/etc/passwd", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isInside(root, c.target); got != c.want {
				t.Fatalf("isInside(%q, %q) = %v, want %v", root, c.target, got, c.want)
			}
		})
	}
}
