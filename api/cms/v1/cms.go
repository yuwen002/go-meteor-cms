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

	// 静态资源：让 public/uploads 下的文件可通过 /uploads/xxx 访问。
	// WithFileServer 是 RunOption，必须在创建 server 时传入；它仅拦截 GET 且
	// 路径以 /uploads/ 前缀开头且文件存在的请求，其余请求仍走 API 路由。
	uploadDir := c.Upload.Defaults().Dir
	server := rest.MustNewServer(c.RestConf,
		rest.WithFileServer(uploadURLPrefix, http.Dir(uploadDir)))
	defer server.Stop()

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

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
