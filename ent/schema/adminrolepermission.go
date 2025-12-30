package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.Int64("id").
			Positive().
			Immutable().
			Comment("主键ID"),

		field.Int64("role_id").
			Comment("角色ID，关联 admin_roles.id").
			Positive(),

		field.Int64("permission_id").
			Comment("权限ID，关联 admin_permissions.id").
			Positive(),
	}
}

func (AdminRolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "permission_id").
			Unique(),
	}
}

func (AdminRolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		// 指向角色
		edge.To("role", AdminRole.Type).
			Field("role_id").
			Required().
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),

		// 指向权限
		edge.To("permission", AdminPermission.Type).
			Field("permission_id").
			Required().
			Unique().
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}
