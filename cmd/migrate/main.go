package main

import (
	"context"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/migrate"
	_ "github.com/yuwen002/go-meteor-cms/ent/runtime"
	"github.com/yuwen002/go-meteor-cms/internal/seed"
	"github.com/zeromicro/go-zero/core/conf"
)

// defaultDSN 兜底的数据库连接串
const defaultDSN = "root:123456@tcp(127.0.0.1:3306)/go_meteor_cms?parseTime=True&loc=Local"

// configPath 与 API 服务共用的配置文件
const configPath = "api/cms/v1/etc/cms-api.yaml"

// resolveDSN 按优先级获取 DSN：环境变量 MYSQL_DSN > cms-api.yaml > 内置默认值
func resolveDSN() string {
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		return dsn
	}

	var c struct {
		Database struct {
			Driver string
			Source string
		}
	}
	if err := conf.Load(configPath, &c); err == nil && c.Database.Source != "" {
		return c.Database.Source
	}

	return defaultDSN
}

func main() {
	log.Println("🚀 开始执行数据库迁移...")

	// 数据库连接信息
	dsn := resolveDSN()

	// 打开数据库连接
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ 无法连接数据库：%v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("❌ 关闭数据库连接时出错: %v", err)
		}
	}()

	log.Println("✅ 数据库连接成功，开始同步表结构...")

	ctx := context.Background()
	// 执行迁移。WithDropColumn 会在 schema 中已不存在的列（如审计日志去掉软删除后残留的
	// deleted_at）从数据库删除，使表结构与代码定义保持一致。
	if err := client.Schema.Create(ctx, migrate.WithDropColumn(true)); err != nil {
		log.Fatalf("❌ 数据库迁移失败：%v", err)
	}

	log.Println("🎉 数据库迁移完成！所有表结构已同步。")

	err = seed.InitSeed(ctx, client)
	if err != nil {
		log.Printf("❌ 初始化数据失败: %v\n", err)
		os.Exit(1)
	}
	log.Println("✅ 数据初始化完成！")

	// 输出权限统计：Type=3 的 API 权限数量直接决定非超级管理员能否正常访问接口
	perms, err := client.AdminPermission.Query().All(ctx)
	if err != nil {
		log.Printf("⚠️  查询权限统计失败: %v\n", err)
		return
	}
	var menu, button, api int
	for _, p := range perms {
		switch p.Type {
		case seed.PermissionTypeMenu:
			menu++
		case seed.PermissionTypeButton:
			button++
		case seed.PermissionTypeAPI:
			api++
		}
	}
	log.Printf("📊 权限统计：菜单 %d / 按钮 %d / API %d\n", menu, button, api)
}
