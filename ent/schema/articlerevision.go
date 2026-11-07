package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ArticleRevision 文章修订快照。每次编辑前把上一版内容存一条，支持回溯与恢复。
// 刻意不挂 SoftDeleteMixin：版本是审计数据，只能通过明确操作清理。
type ArticleRevision struct {
	ent.Schema
}

func (ArticleRevision) Mixin() []ent.Mixin {
	return []ent.Mixin{BaseMixin{}}
}

func (ArticleRevision) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "article_revisions"},
		entsql.WithComments(true),
		schema.Comment("文章修订历史表"),
	}
}

func (ArticleRevision) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive().Immutable().Comment("主键 ID"),
		field.Int64("article_id").Positive().Comment("文章 ID"),
		field.String("title").MaxLen(255).NotEmpty().Comment("该版本的标题"),
		field.String("summary").MaxLen(512).Default("").Comment("该版本的摘要"),
		field.String("cover").MaxLen(512).Default("").Comment("该版本的封面 URL"),
		field.Text("content").Optional().Comment("该版本的 Markdown 正文"),
		field.Int64("editor_id").Optional().Nillable().Comment("编辑人管理员 ID"),
		field.String("remark").MaxLen(255).Default("").Comment("版本备注"),
	}
}

func (ArticleRevision) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("article_id", "id"),
	}
}
