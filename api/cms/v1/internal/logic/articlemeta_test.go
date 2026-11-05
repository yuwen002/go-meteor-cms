package logic

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeSlug(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"英文标题", "Hello World", "hello-world"},
		{"大小写与下划线", "  Go_Meteor  CMS ", "go-meteor-cms"},
		{"压缩连续分隔符", "a---b   c", "a-b-c"},
		{"纯中文返回空", "文章发布", ""},
		{"中文混排保留英文", "Go 语言入门", "go"},
		{"去首尾连字符", "--abc--", "abc"},
		{"空输入", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NormalizeSlug(c.raw); got != c.want {
				t.Fatalf("NormalizeSlug(%q) = %q, want %q", c.raw, got, c.want)
			}
		})
	}
}

func TestNormalizeSlugTruncates(t *testing.T) {
	got := NormalizeSlug(strings.Repeat("a", 250))
	if len(got) != maxSlugLen {
		t.Fatalf("slug length = %d, want %d", len(got), maxSlugLen)
	}
}

func TestGenerateSlug(t *testing.T) {
	now := time.Unix(1700000000, 123456789)
	if got := GenerateSlug("My Custom Slug", "标题", now); got != "my-custom-slug" {
		t.Fatalf("user slug = %q", got)
	}
	if got := GenerateSlug("", "Hello World", now); got != "hello-world" {
		t.Fatalf("title derived slug = %q", got)
	}
	if got := GenerateSlug("", "纯中文标题", now); got != "post-1700000000123456789" {
		t.Fatalf("fallback slug = %q", got)
	}
}

func TestIsValidSlug(t *testing.T) {
	for _, s := range []string{"hello-world", "post-123", "a"} {
		if !IsValidSlug(s) {
			t.Errorf("expect valid slug: %q", s)
		}
	}
	for _, s := range []string{"", "Hello", "a_b", "-abc", "abc-", "中文", "a--b"} {
		if IsValidSlug(s) {
			t.Errorf("expect invalid slug: %q", s)
		}
	}
}

func TestIsPubliclyVisible(t *testing.T) {
	now := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	if IsPubliclyVisible(ArticleStatusPublished, &past, now) != true {
		t.Error("已发布且发布时间已到，应可见")
	}
	if IsPubliclyVisible(ArticleStatusPublished, &future, now) != false {
		t.Error("定时发布未到点，应不可见")
	}
	if IsPubliclyVisible(ArticleStatusPublished, nil, now) != false {
		t.Error("已发布但没有发布时间，应不可见")
	}
	if IsPubliclyVisible(ArticleStatusDraft, &past, now) != false {
		t.Error("草稿不应可见")
	}
	if IsPubliclyVisible(ArticleStatusOffline, &past, now) != false {
		t.Error("已下线不应可见")
	}
}

func TestParsePublishTime(t *testing.T) {
	now := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)

	if got, err := ParsePublishTime("", now); err != nil || !got.Equal(now) {
		t.Fatalf("empty should mean publish now, got %v err %v", got, err)
	}
	got, err := ParsePublishTime("2026-11-06 09:30:00", now)
	if err != nil || got.Year() != 2026 || got.Hour() != 9 || got.Minute() != 30 {
		t.Fatalf("parse datetime failed: %v err %v", got, err)
	}
	got, err = ParsePublishTime("2026-11-06T09:30", now)
	if err != nil || got.Hour() != 9 || got.Minute() != 30 {
		t.Fatalf("parse datetime-local failed: %v err %v", got, err)
	}
	if _, err = ParsePublishTime("2026/11/06", now); err == nil {
		t.Fatal("unsupported format should return error")
	}
}

func TestArticleStatus(t *testing.T) {
	if !IsValidArticleStatus(ArticleStatusDraft) ||
		!IsValidArticleStatus(ArticleStatusPublished) ||
		!IsValidArticleStatus(ArticleStatusOffline) {
		t.Fatal("all three statuses should be valid")
	}
	if IsValidArticleStatus(-1) || IsValidArticleStatus(3) {
		t.Fatal("out of range status should be invalid")
	}
	if ArticleStatusText(ArticleStatusPublished) != "已发布" || ArticleStatusText(99) != "未知" {
		t.Fatal("unexpected status text")
	}
}
