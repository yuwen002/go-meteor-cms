package main

import (
	"context"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/internal/seed"
)

func main() {
	log.Println("🚀 开始执行数据库迁移...")

	// 数据库连接信息
	dsn := "root:123456@tcp(127.0.0.1:3306)/go_meteor_cms?parseTime=True&loc=Local"

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
	// 执行迁移
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("❌ 数据库迁移失败：%v", err)
	}

	log.Println("🎉 数据库迁移完成！所有表结构已同步。")

	err = seed.InitSeed(ctx, client)
	if err != nil {
		log.Printf("❌ 初始化数据失败: %v\n", err)
		os.Exit(1)
	}
	log.Println("✅ 数据初始化完成！")
}
