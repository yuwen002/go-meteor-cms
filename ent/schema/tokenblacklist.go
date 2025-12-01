package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TokenBlacklist 表示 JWT Token 黑名单表
type TokenBlacklist struct {
	ent.Schema
}

// Annotations 表和字段的额外注解
func (TokenBlacklist) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true), // 开启字段注释支持
		schema.Comment("JWT Token 黑名单表，用于存储已注销或失效的 JWT token"), // 表注释
	}
}

// Mixin 复用创建时间、更新时间等通用字段
func (TokenBlacklist) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{}, // 假设你已经定义了包含 created_at、updated_at 等字段的 BaseMixin
	}
}

// Fields 定义 TokenBlacklist 实体字段
func (TokenBlacklist) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Comment("自增主键 ID").
			Unique().
			Positive().
			Immutable(),

		field.String("token").
			NotEmpty().
			Annotations(entsql.Prefix(191)).
			Comment("被加入黑名单的完整 JWT token 字符串"),

		field.Time("expired_at").
			Comment("该 token 原本的过期时间，用于定时清理黑名单（过期后可删除）"),
	}
}

// Indexes 定义索引
func (TokenBlacklist) Indexes() []ent.Index {
	return []ent.Index{
		// token 必须唯一，防止同一 token 重复加入黑名单
		index.Fields("token").
			Unique(),

		// 可选：给过期时间加个普通索引，方便定时清理任务快速查找已过期记录
		index.Fields("expired_at"),
	}
}

// Edges 如果以后需要关联用户等实体，可以在这里定义，目前不需要
func (TokenBlacklist) Edges() []ent.Edge {
	return nil
}
