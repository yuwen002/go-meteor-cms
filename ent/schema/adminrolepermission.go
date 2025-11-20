package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type AdminRolePermission struct {
	ent.Schema
}

func (AdminRolePermission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (AdminRolePermission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "admin_role_permissions",
		},
		entsql.WithComments(true),
		schema.Comment("角色权限关联表"),
	}
}

func (AdminRolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("role_id").
			Comment("角色ID，关联 admin_roles.id").
			Positive(),

		field.Int64("permission_id").
			Comment("权限ID，关联 admin_permissions.id").
			Positive(),

		// ent 不允许只有 mixin 字段，所以加个占位字段
		field.Uint8("dummy").
			Default(0).
			Comment("占位字段，无实际意义"),
	}
}

func (AdminRolePermission) Edges() []ent.Edge {
	return nil
}
