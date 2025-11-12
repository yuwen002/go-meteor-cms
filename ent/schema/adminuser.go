package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// AdminUser holds the schema definition for the AdminUser entity.
type AdminUser struct {
	ent.Schema
}

func (AdminUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true), // 启用表注释
	}
}

// Fields of the AdminUser.
func (AdminUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Comment("自增主键ID").
			Unique().
			Positive().
			Immutable(),

		field.String("username").
			Comment("用户名").
			Unique().
			NotEmpty(),

		field.String("password_hash").
			Comment("密码哈希值").
			NotEmpty(),

		field.String("nickname").
			Comment("昵称").
			Optional(),

		field.String("email").
			Comment("邮箱").
			Optional(),

		field.String("phone").
			Comment("手机号").
			Optional(),

		field.Bool("is_active").
			Comment("是否启用").
			Default(true),

		field.Time("last_login_at").
			Comment("最后登录时间").
			Optional().
			Nillable(),

		field.Time("created_at").
			Comment("创建时间").
			Default(time.Now),

		field.Time("updated_at").
			Comment("更新时间").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the AdminUser.
func (AdminUser) Edges() []ent.Edge {
	return nil
}
