package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
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
		field.Int64("role_id").Comment("角色ID"),
		field.Int64("dept_id").Comment("部门ID"),
	}
}
