package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SlugRedirect 记录内容实体的历史 slug。
//
// 三类来源：
//   - 编辑时改了 slug：旧地址 301 到新地址，避免外链与搜索引擎收录失效；
//   - 软删除前的原 slug 快照：回收站恢复时用它还原；
//   - 恢复后该条记录会被移除。
//
// 刻意不挂 SoftDeleteMixin：历史别名本身需要参与唯一约束，且物理清理是明确操作。
type SlugRedirect struct {
	ent.Schema
}

func (SlugRedirect) Mixin() []ent.Mixin {
	return []ent.Mixin{BaseMixin{}}
}

func (SlugRedirect) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "slug_redirects"},
		entsql.WithComments(true),
		schema.Comment("内容实体历史 slug 表，用于 301 与回收站恢复"),
	}
}

func (SlugRedirect) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Positive().Immutable().Comment("主键 ID"),
		field.String("entity_type").MaxLen(32).NotEmpty().
			Comment("实体类型：article / category / tag"),
		field.Int64("entity_id").Positive().Comment("实体 ID"),
		field.String("old_slug").MaxLen(191).NotEmpty().Comment("历史 slug"),
	}
}

func (SlugRedirect) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("entity_type", "old_slug").Unique(),
		index.Fields("entity_type", "entity_id"),
	}
}
