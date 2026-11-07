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
	ArticleStatusPending   = 3 // 待审核
)

// SlugRedirect 的实体类型
const (
	EntityTypeArticle  = "article"
	EntityTypeCategory = "category"
	EntityTypeTag      = "tag"
)

const maxSlugLen = 190

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var markdownImageRe = regexp.MustCompile(`!\[[^\]]*\]\(([^)\s]+)\)`)

var publishTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
}

// IsValidArticleStatus 判断状态值是否合法
func IsValidArticleStatus(status int) bool {
	return status >= ArticleStatusDraft && status <= ArticleStatusPending
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
	case ArticleStatusPending:
		return "待审核"
	default:
		return "未知"
	}
}

// publishTimeTolerance 容忍存储精度带来的 1 秒误差。
// published_at 列是 datetime(0)，写入「立即发布」的 time.Now()（含纳秒）时可能被
// MySQL 四舍五入到下一秒；若严格用 After(now) 判断，刚发布的文章会在 1 秒内被
// 误判成「定时中」且公开不可见。真正的定时发布通常以分钟为单位，不受影响。
const publishTimeTolerance = time.Second

// ArticleDisplayStatus 后台列表的展示状态。已发布但发布时间在未来的文章显示「定时中」，
// 避免运营把未到点的定时文章误认为已经上线。
func ArticleDisplayStatus(status int, publishedAt *time.Time, now time.Time) string {
	if status == ArticleStatusPublished && publishedAt != nil && publishedAt.After(now.Add(publishTimeTolerance)) {
		return "定时中"
	}
	return ArticleStatusText(status)
}

// IsPubliclyVisible 公开可见条件：状态为已发布且发布时间已到。
// 定时发布只需把 published_at 设为未来时间，到点后本条件自动放行，
// 因此不需要额外的后台定时任务或状态扫描。
func IsPubliclyVisible(status int, publishedAt *time.Time, now time.Time) bool {
	return status == ArticleStatusPublished &&
		publishedAt != nil &&
		!publishedAt.After(now.Add(publishTimeTolerance))
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

// TrashSlug 软删除时使用的占位 slug，保证原 slug 的唯一索引被释放
func TrashSlug(id int64, now time.Time) string {
	return fmt.Sprintf("deleted-%d-%d", id, now.UnixNano())
}

// IsTrashSlug 判断是否为软删除占位 slug
func IsTrashSlug(s string) bool {
	return strings.HasPrefix(s, "deleted-")
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

// ExtractMarkdownImageURLs 提取 Markdown 正文里的图片地址，供彻底删除时清理孤儿文件。
// 只接受站内 /uploads/ 路径，外链不参与清理（避免误删他人资源）。
func ExtractMarkdownImageURLs(content string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0)
	for _, m := range markdownImageRe.FindAllStringSubmatch(content, -1) {
		if len(m) < 2 {
			continue
		}
		url := strings.TrimSpace(m[1])
		if !strings.HasPrefix(url, "/uploads/") || seen[url] {
			continue
		}
		seen[url] = true
		out = append(out, url)
	}
	return out
}

// FirstNonEmpty 返回第一个非空字符串，用于 SEO 字段回退
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
