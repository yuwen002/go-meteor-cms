// Package upload 处理后台文件上传。
//
// 安全约束：
//   - 只接受扩展名白名单内的文件，且不依赖客户端声明的扩展名，同时校验文件头魔数；
//   - 文件名统一使用 UUID 重新生成，杜绝路径穿越与同名覆盖；
//   - 限制单文件大小，避免磁盘被写满。
package upload

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// maxDetectBytes 魔数检测读取的最大字节数
const maxDetectBytes = 512

// ErrEmptyFile 空文件
var ErrEmptyFile = errors.New("文件内容为空")

// FileInfo 上传结果
type FileInfo struct {
	// URL 可直接访问的相对路径
	URL string
	// Filename 原始文件名
	Filename string
	// Size 文件大小（字节）
	Size int64
}

// Save 保存上传文件到指定目录，并返回可访问路径
func Save(file *multipart.FileHeader, dir string, maxSizeMB int64, allowedExts []string) (*FileInfo, error) {
	if file == nil {
		return nil, ErrEmptyFile
	}
	if file.Size == 0 {
		return nil, ErrEmptyFile
	}
	if maxSizeMB > 0 && file.Size > maxSizeMB*1024*1024 {
		return nil, fmt.Errorf("文件大小不能超过 %dMB", maxSizeMB)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowed(allowedExts, ext) {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 校验文件头，确认真实类型与扩展名一致
	if err := detectMismatch(src, ext); err != nil {
		return nil, err
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// 按日期分目录，避免单目录文件过多
	subDir := time.Now().Format("2006/01")
	targetDir := filepath.Join(dir, subDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return nil, err
	}

	name := uuid.NewString() + ext
	target := filepath.Join(targetDir, name)

	dst, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	size, err := io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		URL:      filepath.ToSlash(filepath.Join("/uploads", subDir, name)),
		Filename: file.Filename,
		Size:     size,
	}, nil
}

// allowed 判断扩展名是否在白名单内
func allowed(allowedExts []string, ext string) bool {
	for _, allowed := range allowedExts {
		if strings.EqualFold(allowed, ext) {
			return true
		}
	}
	return false
}

// detectMismatch 通过文件头魔数判断真实类型是否与扩展名相符
func detectMismatch(src multipart.File, ext string) error {
	head := make([]byte, maxDetectBytes)
	n, err := src.Read(head)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	head = head[:n]

	matchers := map[string]func([]byte) bool{
		".jpg":  func(b []byte) bool { return bytes.HasPrefix(b, []byte{0xFF, 0xD8, 0xFF}) },
		".jpeg": func(b []byte) bool { return bytes.HasPrefix(b, []byte{0xFF, 0xD8, 0xFF}) },
		".png":  func(b []byte) bool { return bytes.HasPrefix(b, []byte{0x89, 'P', 'N', 'G'}) },
		".gif": func(b []byte) bool {
			return bytes.HasPrefix(b, []byte("GIF87a")) || bytes.HasPrefix(b, []byte("GIF89a"))
		},
		".webp": func(b []byte) bool {
			return len(b) > 12 && string(b[8:12]) == "WEBP"
		},
	}

	check, ok := matchers[ext]
	if !ok {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}
	if !check(head) {
		return fmt.Errorf("文件内容与扩展名 %s 不符", ext)
	}

	return nil
}

// defaultAvatarURL 默认头像，删除时跳过，避免误删系统资源
const defaultAvatarURL = "/uploads/avatars/meteor-default.jpg"

// Delete 删除通过 Save 落盘的文件。
//
// url 形如 /uploads/2006/01/name.ext。出于安全考虑：
//   - 仅接受 /uploads/ 前缀的路径，且不允许包含 ".." 或绝对路径，防止路径穿越误删系统文件；
//   - 默认头像不删除；
//   - 文件不存在时视为已清理，返回 (false, nil)。
func Delete(dir, url string) (bool, error) {
	const prefix = "/uploads/"
	if url == "" || !strings.HasPrefix(url, prefix) {
		return false, nil
	}
	if url == defaultAvatarURL || strings.Contains(url, "meteor-default") {
		return false, nil
	}

	rel := strings.TrimPrefix(url, prefix)
	if strings.Contains(rel, "..") || filepath.IsAbs(rel) {
		return false, fmt.Errorf("非法的文件访问路径: %s", url)
	}

	full := filepath.Clean(filepath.Join(dir, rel))
	base := filepath.Clean(dir)
	if full != base && !strings.HasPrefix(full, base+string(os.PathSeparator)) {
		return false, fmt.Errorf("文件路径越界: %s", url)
	}

	if err := os.Remove(full); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
