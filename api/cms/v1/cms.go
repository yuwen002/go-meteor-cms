// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/handler"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/middleware"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/internal/pkg/cleanup"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "api/cms/v1/etc/cms-api.yaml", "the config file")

// uploadURLPrefix 上传文件的访问路径前缀
const uploadURLPrefix = "/uploads"

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 环境变量优先：便于容器化部署时注入数据库地址，避免把连接信息写进镜像
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		c.Database.Source = dsn
	}
	// 前端构建产物目录：容器化部署时通过 WEB_DIST 注入，本地开发留空（仅暴露 API）
	if dir := os.Getenv("WEB_DIST"); dir != "" {
		c.StaticDir = dir
	}

	// 静态资源：让 public/uploads 下的文件可通过 /uploads/xxx 访问。
	// WithFileServer 是 RunOption，必须在创建 server 时传入；它仅拦截 GET 且
	// 路径以 /uploads/ 前缀开头且文件存在的请求，其余请求仍走 API 路由。
	uploadDir := c.Upload.Defaults().Dir
	server := rest.MustNewServer(c.RestConf,
		rest.WithFileServer(uploadURLPrefix, http.Dir(uploadDir)))
	defer server.Stop()

	// 若配置了前端产物目录，则直接托管 SPA（含 history 路由回退）
	server.Use(middleware.NewSPAMiddleware(c.StaticDir))
	// 将登出请求携带的 X-Refresh-Token 透传进 context，供登出接口一并吊销
	server.Use(middleware.NewRefreshTokenHeaderMiddleware())

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 启动后台维护任务，随服务停止而退出
	taskCtx, cancelTasks := context.WithCancel(context.Background())
	defer cancelTasks()
	cleanup.StartTokenBlacklistCleaner(
		taskCtx,
		ctx.EntClient,
		time.Duration(c.Maintenance.TokenBlacklistCleanInterval)*time.Second,
	)
	// 操作/登录日志定期清理（retentionDays<=0 时内部不启动，永久保留）
	cleanup.StartAuditLogCleaner(taskCtx, ctx.EntClient, c.Maintenance.AuditLogRetentionDays)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
