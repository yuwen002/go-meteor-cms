package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tag 文章标签表
type Tag struct {
	ent.Schema
}

func (Tag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		SoftDeleteMixin{},
	}
}

func (Tag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tags"},
		entsql.WithComments(true),
		schema.Comment("文章标签表"),
	}
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Immutable().
			Comment("自增主键 ID"),

		field.String("name").
			MaxLen(64).
			NotEmpty().
			Unique().
			Comment("标签名称，全局唯一"),

		field.String("slug").
			MaxLen(191).
			NotEmpty().
			Unique().
			Comment("URL 友好标识，全局唯一"),

		field.String("description").
			MaxLen(512).
			Default("").
			Comment("标签描述"),

		field.Int("sort").
			Default(0).
			Comment("排序，从小到大"),

		field.Bool("is_active").
			Default(true).
			Comment("是否启用"),
	}
}

func (Tag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sort"),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("articles", Article.Type).
			Ref("tags"),
	}
}
