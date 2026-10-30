package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SysDictItem 数据字典项表，隶属于某个字典类型
type SysDictItem struct {
	ent.Schema
}

func (SysDictItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("数据字典项表，存放字典类型下的具体枚举值"),
	}
}

func (SysDictItem) Mixin() []ent.Mixin {
	return []ent.Mixin{BaseMixin{}}
}

func (SysDictItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Comment("自增主键 ID").
			Unique().
			Positive().
			Immutable(),

		field.String("dict_code").
			MaxLen(64).
			NotEmpty().
			Comment("所属字典编码，冗余自 sys_dicts.code，便于按编码直接查询"),

		field.String("label").
			MaxLen(128).
			NotEmpty().
			Comment("显示名称"),

		field.String("value").
			MaxLen(128).
			NotEmpty().
			Comment("实际取值"),

		field.Int("status").
			Default(1).
			Comment("状态：1 启用 0 禁用"),

		field.Int("sort").
			Default(0).
			Comment("排序，从小到大"),
	}
}

func (SysDictItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("dict_code"),
		index.Fields("dict_code", "value"),
	}
}

func (SysDictItem) Edges() []ent.Edge {
	return nil
}
