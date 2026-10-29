// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package main

import (
	"context"
	"flag"
	"fmt"
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

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 环境变量优先：便于容器化部署时注入数据库地址，避免把连接信息写进镜像
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		c.Database.Source = dsn
	}

	server := rest.MustNewServer(c.RestConf)
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
