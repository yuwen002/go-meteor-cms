package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AdminRole 后台角色表
type AdminRole struct {
	ent.Schema
}

func (AdminRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		SoftDeleteMixin{},
	}
}

func (AdminRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "admin_roles",
		},
		entsql.WithComments(true),
		schema.Comment("后台角色表"),
	}
}

func (AdminRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Comment("主键ID，自增").
			Immutable().
			Positive(),

		field.String("name").
			NotEmpty().
			Unique().
			Comment("角色名称"),

		field.String("code").
			Comment("角色编码，用于系统标识，如 SUPER_ADMIN").
			NotEmpty().
			Unique(),

		field.String("desc").
			Default("").
			Optional().
			Comment("角色描述"),

		field.Int("data_scope").
			Default(1).
			Comment("数据权限范围：1 全公司 / 2 本部门 / 3 部门及子部门 / 4 仅自己 / 5 自定义部门"),

		field.Bool("is_system").
			Comment("是否系统内置角色（禁止删除）").
			Default(false),

		field.Bool("is_active").
			Comment("是否启用").
			Default(true),

		field.Int("sort").
			Comment("排序，从小到大").
			Default(0),
	}
}

func (AdminRole) Edges() []ent.Edge {
	return []ent.Edge{
		// 角色 - 管理员（多对多）
		edge.To("roles", AdminRole.Type).
			Through("admin_user_roles", AdminUserRole.Type),

		// 角色 - 权限（多对多）
		edge.To("permissions", AdminPermission.Type).
			Through("admin_role_permissions", AdminRolePermission.Type),

		// 角色 - 自定义部门（data_scope = 5 时使用）
		edge.To("departments", Department.Type).
			Through("admin_role_departments", AdminRoleDept.Type),
	}
}
