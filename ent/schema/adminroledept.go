package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type AdminRoleDept struct {
	ent.Schema
}

func (AdminRoleDept) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		SoftDeleteMixin{},
	}
}

func (AdminRoleDept) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "admin_role_depts",
		},
		entsql.WithComments(true),
		schema.Comment("角色部门关联表"),
	}
}

func (AdminRoleDept) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Immutable().
			Comment("主键ID"),
		field.Int64("role_id").
			Positive().
			Comment("角色ID，关联 admin_roles.id"),

		field.Int64("dept_id").
			Positive().
			Comment("部门ID，关联 departments.id"),
	}
}

func (AdminRoleDept) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "dept_id").
			Unique(),
	}
}

func (AdminRoleDept) Edges() []ent.Edge {
	return nil
}
