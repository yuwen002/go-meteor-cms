package upload

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pngHeader 构造一个带真实 PNG 头的 multipart 文件头
func pngHeader(t *testing.T, name string, content []byte) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatal(err)
	}
	w.Close()

	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	f, h, err := req.FormFile("file")
	if err != nil {
		t.Fatal(err)
	}
	// 不提前关闭 f：Save 内部会再次 Open
	_ = f
	return h
}

func pngBytes() []byte {
	// 最小合法 PNG（89 50 4E 47 ...）
	return []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 'I', 'H', 'D', 'R'}
}

func TestSave_OK(t *testing.T) {
	dir := t.TempDir()
	h := pngHeader(t, "avatar.png", pngBytes())

	info, err := Save(h, dir, 5, []string{".png", ".jpg"})
	if err != nil {
		t.Fatalf("Save 失败: %v", err)
	}
	if !strings.HasPrefix(info.URL, "/uploads/") {
		t.Fatalf("URL 前缀错误: %s", info.URL)
	}
	if info.Size != int64(len(pngBytes())) {
		t.Fatalf("Size 不匹配: %d", info.Size)
	}
	// 文件应落在 dir/2006/01/<uuid>.png
	rel := strings.TrimPrefix(info.URL, "/uploads/")
	full := filepath.Join(dir, rel)
	if _, err := os.Stat(full); err != nil {
		t.Fatalf("落盘文件不存在: %v", err)
	}
}

func TestSave_MagicMismatch(t *testing.T) {
	dir := t.TempDir()
	// 内容是纯文本，却用 .png 后缀，应被魔数校验拒绝
	h := pngHeader(t, "fake.png", []byte("this is not an image"))
	_, err := Save(h, dir, 5, []string{".png"})
	if err == nil {
		t.Fatal("期望魔数不匹配错误，但成功了")
	}
	if !strings.Contains(err.Error(), "文件内容与扩展名") {
		t.Fatalf("错误信息不符预期: %v", err)
	}
}

func TestSave_TooLarge(t *testing.T) {
	dir := t.TempDir()
	content := make([]byte, 6*1024*1024) // 6MB
	copy(content, pngBytes())
	h := pngHeader(t, "big.png", content)
	_, err := Save(h, dir, 5, []string{".png"})
	if err == nil {
		t.Fatal("期望大小超限错误，但成功了")
	}
	if !strings.Contains(err.Error(), "文件大小不能超过") {
		t.Fatalf("错误信息不符预期: %v", err)
	}
}

func TestSave_Empty(t *testing.T) {
	dir := t.TempDir()
	h := pngHeader(t, "empty.png", []byte{})
	_, err := Save(h, dir, 5, []string{".png"})
	if err == nil {
		t.Fatal("期望空文件错误，但成功了")
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "2026", "10", "abc.png")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	removed, err := Delete(dir, "/uploads/2026/10/abc.png")
	if err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}
	if !removed {
		t.Fatal("期望删除成功")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("文件应已删除")
	}
}

func TestDelete_DefaultAvatarSkipped(t *testing.T) {
	dir := t.TempDir()
	// 默认头像不应被删除
	removed, err := Delete(dir, "/uploads/avatars/meteor-default.jpg")
	if err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}
	if removed {
		t.Fatal("默认头像不应被删除")
	}
}

func TestDelete_PathTraversalBlocked(t *testing.T) {
	dir := t.TempDir()
	// 尝试穿越到目录外的文件
	_, err := Delete(dir, "/uploads/../../etc/passwd")
	if err == nil {
		t.Fatal("期望路径穿越被拒绝")
	}
}
