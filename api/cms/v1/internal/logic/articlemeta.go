package logic

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// 文章状态。发布状态是流程状态，刻意不放在数据字典里，避免运营改字典破坏状态机。
const (
	ArticleStatusDraft     = 0
	ArticleStatusPublished = 1
	ArticleStatusOffline   = 2
)

const maxSlugLen = 190

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

var publishTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
}

// IsValidArticleStatus 判断状态值是否合法
func IsValidArticleStatus(status int) bool {
	return status >= ArticleStatusDraft && status <= ArticleStatusOffline
}

// ArticleStatusText 状态中文名
func ArticleStatusText(status int) string {
	switch status {
	case ArticleStatusDraft:
		return "草稿"
	case ArticleStatusPublished:
		return "已发布"
	case ArticleStatusOffline:
		return "已下线"
	default:
		return "未知"
	}
}

// IsPubliclyVisible 公开可见条件：状态为已发布且发布时间已到。
// 定时发布只需把 published_at 设为未来时间，到点后本条件自动放行，
// 因此不需要额外的后台定时任务或状态扫描。
func IsPubliclyVisible(status int, publishedAt *time.Time, now time.Time) bool {
	return status == ArticleStatusPublished &&
		publishedAt != nil &&
		!publishedAt.After(now)
}

// NormalizeSlug 把用户输入或标题规范化为 URL 友好标识。
// 非 ASCII 字符（中文等）会被丢弃，纯中文输入返回空串，由 GenerateSlug 兜底。
func NormalizeSlug(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if b.Len() > 0 && !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > maxSlugLen {
		out = strings.Trim(out[:maxSlugLen], "-")
	}
	return out
}

// IsValidSlug 校验别名格式
func IsValidSlug(s string) bool {
	return s != "" && len(s) <= maxSlugLen && slugPattern.MatchString(s)
}

// GenerateSlug 优先使用用户填写的别名，否则从标题派生；两者都无法生成合法别名时
// 用 post-<纳秒时间戳> 兜底，保证 slug 非空。
func GenerateSlug(raw, title string, now time.Time) string {
	if s := NormalizeSlug(raw); s != "" {
		return s
	}
	if s := NormalizeSlug(title); s != "" {
		return s
	}
	return fmt.Sprintf("post-%d", now.UnixNano())
}

// ParsePublishTime 解析发布时间；空串表示立即发布。
// 支持 RFC3339、YYYY-MM-DD HH:mm:ss 以及 HTML datetime-local 的 YYYY-MM-DDTHH:mm。
func ParsePublishTime(raw string, now time.Time) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return now, nil
	}
	for _, layout := range publishTimeLayouts {
		if t, err := time.ParseInLocation(layout, raw, now.Location()); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported publish time: %s", raw)
}

// FormatOptionalTime 可空时间的统一输出格式，空值返回空串
func FormatOptionalTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
