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
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
	// 跨域：通过 go-zero 内置 WithCors 启用（内部以 NotAllowedHandler 接管 OPTIONS 预检并返回 204）。
	// 仅在显式配置 Cors 源时开启；未配置则保持同源部署（不暴露 CORS 头）。
	opts := []rest.RunOption{
		rest.WithFileServer(uploadURLPrefix, http.Dir(uploadDir)),
	}
	if len(c.Cors) > 0 {
		opts = append(opts, rest.WithCors(c.Cors...))
	}
	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()

	// 若配置了前端产物目录，则直接托管 SPA（含 history 路由回退）
	server.Use(middleware.NewSPAMiddleware(c.StaticDir))
	// 将登出请求携带的 X-Refresh-Token 透传进 context，供登出接口一并吊销
	server.Use(middleware.NewRefreshTokenHeaderMiddleware())

	// 登录限流 / token 版本（全局踢人）：配置了 Redis 则启用分布式能力，
	// 否则（含 Redis 不可用）自动降级为内存限流 / 不校验版本，防护不失效。
	var rds *redis.Redis
	if c.Redis.Host != "" {
		if r, err := redis.NewRedis(c.Redis); err != nil {
			logx.Errorf("Redis 初始化失败，登录限流将使用内存模式: %v", err)
		} else {
			rds = r
		}
	}
	server.Use(middleware.NewLoginRateLimitMiddleware(
		rds,
		c.Login.RateLimit.Enabled,
		c.Login.RateLimit.Max,
		c.Login.RateLimit.Window,
		c.Login.RateLimit.BlockSeconds,
	))

	ctx := svc.NewServiceContext(c, rds)
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
