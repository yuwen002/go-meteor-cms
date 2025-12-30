package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.Int64("id").
			Positive().
			Immutable().
			Comment("主键ID"),

		field.Int64("user_id").
			Comment("用户ID，关联 admin_users.id").
			Positive(),

		field.Int64("role_id").
			Comment("角色ID，关联 admin_roles.id").
			Positive(),
	}
}

func (AdminUserRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "role_id").
			Unique(),
	}
}

func (AdminUserRole) Edges() []ent.Edge {
	return []ent.Edge{
		// 指向用户
		edge.To("user", AdminUser.Type).
			Field("user_id").
			Required().
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// 指向角色
		edge.To("role", AdminRole.Type).
			Field("role_id").
			Required().
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
