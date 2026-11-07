package handler

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/logic"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/ent"
)

// 站点级公开输出：RSS 2.0 / sitemap.xml / robots.txt。
// 三者都是纯文本协议，不走统一 JSON 响应封装。

func xmlEscape(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func siteFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func rssHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, base, ok := loadFeed(w, r, svcCtx, 50)
		if !ok {
			return
		}
		var b strings.Builder
		b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		b.WriteString(`<rss version="2.0"><channel>`)
		b.WriteString("<title>" + xmlEscape("Go Meteor CMS") + "</title>")
		b.WriteString("<link>" + xmlEscape(base+"/articles") + "</link>")
		b.WriteString("<description>" + xmlEscape("最新文章") + "</description>")
		b.WriteString("<language>zh-CN</language>")
		for _, a := range rows {
			link := base + "/articles/" + a.Slug
			b.WriteString("<item>")
			b.WriteString("<title>" + xmlEscape(a.Title) + "</title>")
			b.WriteString("<link>" + xmlEscape(link) + "</link>")
			b.WriteString(`<guid isPermaLink="true">` + xmlEscape(link) + "</guid>")
			if a.PublishedAt != nil {
				b.WriteString("<pubDate>" + a.PublishedAt.UTC().Format(time.RFC1123Z) + "</pubDate>")
			}
			b.WriteString("<description>" + xmlEscape(siteFirstNonEmpty(a.Summary, a.Title)) + "</description>")
			b.WriteString("</item>")
		}
		b.WriteString("</channel></rss>")
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		_, _ = w.Write([]byte(b.String()))
	}
}

func sitemapHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, base, ok := loadFeed(w, r, svcCtx, 200)
		if !ok {
			return
		}
		var b strings.Builder
		b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
		b.WriteString("<url><loc>" + xmlEscape(base+"/articles") + "</loc><changefreq>daily</changefreq><priority>1.0</priority></url>")
		for _, a := range rows {
			b.WriteString("<url><loc>" + xmlEscape(base+"/articles/"+a.Slug) + "</loc>")
			t := a.UpdatedAt
			if a.PublishedAt != nil {
				t = *a.PublishedAt
			}
			b.WriteString("<lastmod>" + t.UTC().Format("2006-01-02") + "</lastmod>")
			b.WriteString("<changefreq>weekly</changefreq><priority>0.8</priority></url>")
		}
		b.WriteString("</urlset>")
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		_, _ = w.Write([]byte(b.String()))
	}
}

func robotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base := strings.TrimRight(svcCtx.Config.WebBaseURL, "/")
		if base == "" {
			base = "http://" + r.Host
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("User-agent: *\nAllow: /\nDisallow: /admin\nSitemap: " + base + "/sitemap.xml\n"))
	}
}

func loadFeed(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext, limit int) ([]*ent.Article, string, bool) {
	rows, base, err := logic.NewPublicArticleLogic(r.Context(), svcCtx).PublicFeed(limit)
	if err != nil {
		http.Error(w, "feed unavailable", http.StatusInternalServerError)
		return nil, "", false
	}
	if base == "" {
		base = "http://" + r.Host
	}
	return rows, base, true
}
