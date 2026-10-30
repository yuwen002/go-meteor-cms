package middleware

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSanitizeParams(t *testing.T) {
	in := `{"username":"alice","password":"secret123","email":"a@b.com","captcha":"X9Z"}`
	out := sanitizeParams("application/json", in)
	if strings.Contains(out, "secret123") {
		t.Fatal("密码未脱敏")
	}
	if strings.Contains(out, "X9Z") {
		t.Fatal("验证码未脱敏")
	}
	if !strings.Contains(out, "******") {
		t.Fatal("敏感字段未替换为 ******")
	}
	if !strings.Contains(out, "a@b.com") {
		t.Fatal("非敏感字段被误改")
	}
}

func TestSanitizeParams_NonJSONDiscarded(t *testing.T) {
	// 非 JSON 表单无法逐字段脱敏，应直接丢弃避免泄露明文
	if got := sanitizeParams("application/x-www-form-urlencoded", "password=secret"); got != "" {
		t.Fatalf("非 JSON 应返回空，实际: %s", got)
	}
}

func TestResolveModule(t *testing.T) {
	cases := map[string]string{
		"/admin/admin-users/1":       "管理员管理",
		"/admin/departments":         "部门管理",
		"/admin/roles/1/departments": "角色管理",
		"/admin/permissions":         "权限管理",
		"/admin/profile":             "个人中心",
		"/admin/logs/operations":     "日志管理",
		"/unknown/path":              "其他",
	}
	for path, want := range cases {
		if got := resolveModule(path); got != want {
			t.Fatalf("resolveModule(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestBusinessResult(t *testing.T) {
	rec := &responseRecorder{status: 200, body: bytes.NewBufferString(`{"code":4001,"msg":"bad"}`)}
	code, msg := rec.businessResult()
	if code != 4001 || msg != "bad" {
		t.Fatalf("businessResult 解析错误: %d %s", code, msg)
	}

	// 非 JSON 响应按 HTTP 状态码判定
	rec2 := &responseRecorder{status: 200, body: bytes.NewBufferString("not json")}
	if code2, _ := rec2.businessResult(); code2 != 0 {
		t.Fatalf("非 JSON 应判成功, 实际 %d", code2)
	}
	rec3 := &responseRecorder{status: 500, body: bytes.NewBufferString("not json")}
	if code3, _ := rec3.businessResult(); code3 != 500 {
		t.Fatalf("非 JSON 且 5xx 应返回状态码, 实际 %d", code3)
	}
}

func TestReadBodySnapshot_MultipartEmpty(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if _, err := w.CreateFormFile("file", "a.png"); err != nil {
		t.Fatal(err)
	}
	w.Close()

	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if got := readBodySnapshot(req); got != "" {
		t.Fatalf("multipart 应返回空快照，实际: %s", got)
	}
}

func TestReadBodySnapshot_JSONMasked(t *testing.T) {
	body := `{"username":"bob","password":"p"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	got := readBodySnapshot(req)
	if got == "" {
		t.Fatal("JSON 请求应返回快照")
	}
	if strings.Contains(got, `"p"`) {
		t.Fatal("快照中密码未脱敏")
	}
	// 原始 body 应被回填，下游仍可读取
	b, _ := io.ReadAll(req.Body)
	if !strings.Contains(string(b), `"username":"bob"`) {
		t.Fatal("请求体未回填，下游无法读取")
	}
}
