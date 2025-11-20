package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// AdminUserRole 用户角色关联表
type AdminUserRole struct {
	ent.Schema
}

func (AdminUserRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (AdminUserRole) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "admin_user_roles",
		},
		entsql.WithComments(true),
		schema.Comment("用户角色关联表"),
	}
}

func (AdminUserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("user_id").
			Comment("用户ID，关联 admin_users.id").
			Positive(),

		field.Int64("role_id").
			Comment("角色ID，关联 admin_roles.id").
			Positive(),

		// 复合唯一
		field.Uint8("dummy"). // ent 不允许无字段表，所以加个占位字段
					Comment("占位字段，无实际意义").
					Default(0),
	}
}

// Edges 建立与用户、角色的关系
func (AdminUserRole) Edges() []ent.Edge {
	return nil
}
