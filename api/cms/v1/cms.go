// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/handler"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/middleware"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/svc"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/internal/pkg/cleanup"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "api/cms/v1/etc/cms-api.yaml", "the config file")

// uploadURLPrefix 上传文件的访问路径前缀
const uploadURLPrefix = "/uploads"

// envBool / envInt 提供环境变量覆盖配置的能力；未设置或值非法时保留配置文件的值。
// 配置不会自动映射环境变量，新增开关必须在这里显式读取，否则容器里传入的 env 不生效。
func envBool(name string, current bool) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(name))) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return current
}

func envInt(name string, current int) int {
	v := strings.TrimSpace(os.Getenv(name))
	if v == "" {
		return current
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return current
	}
	return n
}

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
	// Redis 地址：容器化部署时通过 REDIS_HOST 注入（如 "redis:6379"），本地开发留空（默认 127.0.0.1:6379）
	if rh := os.Getenv("REDIS_HOST"); rh != "" {
		c.Redis.Host = rh
	}
	// 登录验证码开关：默认开启；仅本地开发 / 自动化 E2E 可关闭。
	c.Login.CaptchaEnabled = envBool("LOGIN_CAPTCHA_ENABLED", c.Login.CaptchaEnabled)
	// 公开接口限流：默认开启，防止免鉴权接口被刷。
	c.Public.RateLimit.Enabled = envBool("PUBLIC_RATE_LIMIT_ENABLED", c.Public.RateLimit.Enabled)
	c.Public.RateLimit.Max = envInt("PUBLIC_RATE_LIMIT_MAX", c.Public.RateLimit.Max)
	c.Public.RateLimit.Window = envInt("PUBLIC_RATE_LIMIT_WINDOW", c.Public.RateLimit.Window)
	c.Public.RateLimit.BlockSeconds = envInt("PUBLIC_RATE_LIMIT_BLOCK", c.Public.RateLimit.BlockSeconds)

	// 静态资源：让 public/uploads 下的文件可通过 /uploads/xxx 访问。
	// WithFileServer 是 RunOption，必须在创建 server 时传入；它仅拦截 GET 且
	// 路径以 /uploads/ 前缀开头且文件存在的请求，其余请求仍走 API 路由。
	uploadDir := c.Upload.Defaults().Dir

	// ent client 在 ServiceContext 里创建，但 SPAHandler 必须在 MustNewServer 时构造。
	// 用闭包延迟绑定：resolver 只在请求到达时调用，此时 entClient 已经赋值。
	var entClient *ent.Client
	opts := []rest.RunOption{
		rest.WithFileServer(uploadURLPrefix, http.Dir(uploadDir)),
		// 若配置了前端产物目录，则通过自定义 NotFoundHandler 直接托管 SPA（含 history 路由回退）。
		// 不能用 server.Use 注册：go-zero 的全局中间件不会作用于 404（未匹配）路由，
		// 而 notFoundHandler 在 MustNewServer 时即绑定中间件链，早于 server.Use 调用。
		rest.WithNotFoundHandler(middleware.SPAHandler(
			c.StaticDir,
			middleware.WithPageResolver(func(slug string) (middleware.PageData, bool) {
				if entClient == nil {
					return middleware.PageData{}, false
				}
				return newArticlePageResolver(entClient)(slug)
			}),
		)),
	}
	if len(c.Cors) > 0 {
		opts = append(opts, rest.WithCors(c.Cors...))
	}
	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()

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
	// 公开接口限流：只拦截 /api/ 前缀，后台接口与静态资源不受影响。
	server.Use(middleware.NewPublicRateLimitMiddleware(
		rds,
		c.Public.RateLimit.Enabled,
		c.Public.RateLimit.Max,
		c.Public.RateLimit.Window,
		c.Public.RateLimit.BlockSeconds,
	))

	ctx := svc.NewServiceContext(c, rds)
	entClient = ctx.EntClient
	handler.RegisterHandlers(server, ctx)
	// 文章发布模块独立注册：routes.go 是 goctl 生成文件且被手工追加过，
	// 单独注册可避免重跑 make gen-api 时丢失路由。
	handler.RegisterContentRoutes(server, ctx)

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
