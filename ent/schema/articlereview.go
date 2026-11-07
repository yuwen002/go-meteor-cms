package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ArticleReview 文章审核动作留痕：提交 / 通过 / 驳回各记一条，含操作人与备注。
// 与修订历史一样不挂 SoftDeleteMixin：审核记录属于审计数据，
// 只在文章被彻底删除时由回收站一并物理清理。
type ArticleReview struct {
	ent.Schema
}

func (ArticleReview) Mixin() []ent.Mixin {
	return []ent.Mixin{BaseMixin{}}
}

func (ArticleReview) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "article_reviews"},
		entsql.WithComments(true),
		schema.Comment("文章审核记录表"),
	}
}

func (ArticleReview) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive().Immutable().Comment("主键 ID"),
		field.Int64("article_id").Positive().Comment("文章 ID"),
		field.Int("action").Comment("1 提交审核 / 2 审核通过 / 3 驳回"),
		field.String("remark").MaxLen(500).Default("").Comment("审核备注（驳回原因等）"),
		field.Int64("operator_id").Optional().Default(0).Comment("操作人管理员 ID"),
		// 冗余保存当时的用户名：即便管理员之后改名或离职，审计记录仍可读
		field.String("operator_name").MaxLen(64).Default("").Comment("操作人名称快照"),
	}
}

func (ArticleReview) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("article_id", "id"),
	}
}
