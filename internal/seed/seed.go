package seed

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
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

	password, err := utils.HashPassword("123456")
	if err != nil {
		panic("生成密码哈希失败: " + err.Error())
	}

	_, err = client.AdminUser.Create().
		SetUsername("admin").
		SetPasswordHash(password).
		SetCreatedAt(time.Now()).
		SetIsSuper(true).
		Save(ctx)
	if err != nil {
		panic("创建管理员用户失败: " + err.Error())
	}
}
