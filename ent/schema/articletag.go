package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ArticleTag 文章与标签的关联表。
// 与 AdminUserRole 保持一致：刻意不挂 SoftDeleteMixin，避免软删除行占用
// (article_id, tag_id) 复合唯一索引导致重新关联失败。
type ArticleTag struct {
	ent.Schema
}

func (ArticleTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (ArticleTag) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "article_tags"},
		entsql.WithComments(true),
		schema.Comment("文章标签关联表"),
	}
}

func (ArticleTag) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Immutable().
			Comment("主键 ID"),

		field.Int64("article_id").
			Positive().
			Comment("文章 ID，关联 articles.id"),

		field.Int64("tag_id").
			Positive().
			Comment("标签 ID，关联 tags.id"),
	}
}

func (ArticleTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("article_id", "tag_id").
			Unique(),
	}
}

func (ArticleTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("article", Article.Type).
			Field("article_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),

		edge.To("tag", Tag.Type).
			Field("tag_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
