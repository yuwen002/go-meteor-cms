package seed

import (
	"context"
	"go-meteor-cms/internal/util"
	"time"

	"go-meteor-cms/ent"
	"go-meteor-cms/ent/adminuser"
)

func InitSeed(ctx context.Context, client *ent.Client) {
	initAdmin(ctx, client)
	// 以后可以在这里初始化角色、权限等
}

func initAdmin(ctx context.Context, client *ent.Client) {
	exist, _ := client.AdminUser.Query().
		Where(adminuser.UsernameEQ("admin")).
		Exist(ctx)
	if exist {
		return
	}

	password, err := util.HashPassword("123456")
	if err != nil {
		panic("生成密码哈希失败: " + err.Error())
	}

	_, err = client.AdminUser.Create().
		SetUsername("admin").
		SetPasswordHash(password).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		panic("创建管理员用户失败: " + err.Error())
	}
}
