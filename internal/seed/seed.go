package seed

import (
	"context"
	"time"

	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminrole"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
)

func InitSeed(ctx context.Context, client *ent.Client) error {
	// 1️⃣ 初始化超级管理员用户
	user, err := initAdminUser(ctx, client)
	if err != nil {
		return err
	}

	// 2️⃣ 初始化超级管理员角色
	role, err := initSuperAdminRole(ctx, client)
	if err != nil {
		return err
	}

	// 3️⃣ 初始化默认权限
	perms, err := initDefaultPermissions(ctx, client)
	if err != nil {
		return err
	}

	// 4️⃣ 用户角色关联
	_, err = client.AdminUserRole.Create().
		SetUserID(user.ID).
		SetRoleID(role.ID).
		Save(ctx)
	if err != nil {
		return err
	}

	// 5️⃣ 角色权限关联
	for _, perm := range perms {
		_, err = client.AdminRolePermission.Create().
			SetRoleID(role.ID).
			SetPermissionID(perm.ID).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return nil

}

// initAdminUser 初始化管理员用户
// 如果用户名 "admin" 已存在，则返回已存在的用户
// 如果不存在，则创建新的超级管理员用户，默认密码为 "123456"
// 返回创建或获取的管理员用户信息和可能的错误
func initAdminUser(ctx context.Context, client *ent.Client) (*ent.AdminUser, error) {
	exist, err := client.AdminUser.Query().
		Where(adminuser.UsernameEQ("admin")).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exist {
		return client.AdminUser.Query().Where(adminuser.UsernameEQ("admin")).Only(ctx)
	}
	password, err := utils.HashPassword("123456")
	if err != nil {
		return nil, err
	}

	return client.AdminUser.Create().
		SetUsername("admin").
		SetPasswordHash(password).
		SetIsSuper(true).
		SetCreatedAt(time.Now()).
		Save(ctx)
}

// initSuperAdminRole 初始化超级管理员角色
// 检查是否已存在代码为 "SUPER_ADMIN" 的角色
// 如果存在，直接返回该角色
// 如果不存在，则创建新的超级管理员角色并返回
// 返回创建或获取的角色信息和可能的错误
func initSuperAdminRole(ctx context.Context, client *ent.Client) (*ent.AdminRole, error) {
	// 1️⃣ 检查超级管理员角色是否存在
	role, err := client.AdminRole.Query().
		Where(adminrole.CodeEQ("SUPER_ADMIN")).
		Only(ctx)
	if err == nil {
		// 已存在，直接返回
		return role, nil
	}
	// 2️⃣ 创建超级管理员角色
	role, err = client.AdminRole.Create().
		SetName("超级管理员").
		SetCode("SUPER_ADMIN").
		SetDataScope(1).
		SetIsSystem(true). // 系统内置，禁止删除
		SetIsActive(true). // 默认启用
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func initDefaultPermissions(ctx context.Context, client *ent.Client) ([]*ent.AdminPermission, error) {
	// 检查是否已有权限
	count, err := client.AdminPermission.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return client.AdminPermission.Query().All(ctx)
	}

	// 定义权限列表（先建菜单，再建按钮/API）
	menus := []struct {
		Name       string
		Permission string
		Type       int // 1 菜单 2 按钮 3 API
		Path       string
		Component  string
		Method     string
		APIPath    string
	}{
		{"控制台", "dashboard:view", 1, "/dashboard", "views/dashboard/index.vue", "", ""},
		{"管理员管理", "system:admin_user:list", 1, "/admin/admin-users", "views/user/list.vue", "", ""},
	}

	buttons := []struct {
		Name       string
		Permission string
		Type       int
		ParentName string
		Method     string
		APIPath    string
	}{
		{"新增用户", "system:admin_user:add", 2, "用户管理", "POST", "/admin/users"},
		{"编辑用户", "system:admin_user:edit", 2, "用户管理", "PUT", "/admin/users/:id"},
		{"删除用户", "system:admin_user:delete", 2, "用户管理", "DELETE", "/admin/users/:id"},
	}

	// 1️⃣ 创建菜单
	menuMap := make(map[string]*ent.AdminPermission)
	for _, m := range menus {
		p, err := client.AdminPermission.Create().
			SetName(m.Name).
			SetPermission(m.Permission).
			SetType(m.Type).
			SetPath(m.Path).
			SetComponent(m.Component).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		menuMap[m.Name] = p
	}

	// 2️⃣ 创建按钮/API，设置 parent_id
	var result []*ent.AdminPermission
	for _, b := range buttons {
		parent := menuMap[b.ParentName]
		p, err := client.AdminPermission.Create().
			SetName(b.Name).
			SetPermission(b.Permission).
			SetType(b.Type).
			SetParentID(parent.ID).
			SetMethod(b.Method).
			SetAPIPath(b.APIPath).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	// 3️⃣ 合并菜单和按钮，返回全部权限
	for _, m := range menuMap {
		result = append(result, m)
	}

	return result, nil

}
