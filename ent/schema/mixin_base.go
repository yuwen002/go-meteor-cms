package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// BaseMixin 提供所有表的公共字段
type BaseMixin struct {
	ent.Mixin
}

// 公共字段
func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Comment("创建时间"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),

		field.Time("deleted_at").
			Optional().
			Nillable().
			Comment("删除时间，用于软删除"),
	}
}

// 以下所有方法必须实现，确保 Ent 不调用 nil 指针

func (BaseMixin) Edges() []ent.Edge {
	return nil
}

func (BaseMixin) Indexes() []ent.Index {
	return nil
}

func (BaseMixin) Annotations() []schema.Annotation {
	return nil
}

func (BaseMixin) Hooks() []ent.Hook {
	return nil
}

func (BaseMixin) Policies() []ent.Policy {
	return nil
}

// 某些 Ent 版本（特别是 2024+）使用 Policy()
func (BaseMixin) Policy() ent.Policy {
	return nil
}

// 2024+ 新版本：拦截器
func (BaseMixin) Interceptors() []ent.Interceptor {
	return nil
}
