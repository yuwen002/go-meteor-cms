package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Category 文章分类表，支持父子层级
type Category struct {
	ent.Schema
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		SoftDeleteMixin{},
	}
}

func (Category) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "categories"},
		entsql.WithComments(true),
		schema.Comment("文章分类表，支持父子层级"),
	}
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Immutable().
			Comment("自增主键 ID"),

		field.String("name").
			MaxLen(128).
			NotEmpty().
			Comment("分类名称"),

		field.String("slug").
			MaxLen(191).
			NotEmpty().
			Unique().
			Comment("URL 友好标识，全局唯一"),

		field.Int64("parent_id").
			Optional().
			Nillable().
			Comment("父级分类 ID，为空表示顶级分类"),

		field.String("description").
			MaxLen(512).
			Default("").
			Comment("分类描述"),

		field.Int("sort").
			Default(0).
			Comment("排序，从小到大"),

		field.Bool("is_active").
			Default(true).
			Comment("是否启用"),
	}
}

func (Category) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("parent_id"),
		index.Fields("sort"),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("parent", Category.Type).
			Ref("children").
			Field("parent_id").
			Unique(),

		edge.To("children", Category.Type),

		edge.To("articles", Article.Type),
	}
}
