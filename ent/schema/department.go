package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Department 部门表
type Department struct {
	ent.Schema
}

func (Department) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		SoftDeleteMixin{},
	}
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "departments",
		},
		entsql.WithComments(true),
		schema.Comment("部门表"),
	}
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Immutable().
			Positive().
			Comment("部门ID"),

		field.String("name").
			NotEmpty().
			Comment("部门名称"),

		field.Int64("parent_id").
			Optional().
			Nillable().
			Comment("父级部门ID"),

		field.Int("level").
			Default(1).
			Comment("部门层级"),

		field.Int("sort").
			Default(0).
			Comment("排序"),

		field.Bool("is_active").
			Default(true).
			Comment("是否启用"),

		field.Int64("leader_id").
			Optional().
			Nillable().
			Comment("部门负责人ID（管理员）"),
	}
}

func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		// 父子部门
		edge.From("parent", Department.Type).
			Ref("children").
			Field("parent_id").
			Unique(),

		edge.To("children", Department.Type),

		// 关联管理员
		edge.To("admin_users", AdminUser.Type),
	}
}
