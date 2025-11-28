package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/yuwen002/go-meteor-cms/ent/schema/softdelete"
)

// AdminUser holds the schema definition for the AdminUser entity.
type AdminUser struct {
	ent.Schema
}

func (AdminUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "admin_users"},
		entsql.WithComments(true),
		schema.Comment("后台管理员用户表"),
	}
}

func (AdminUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{}, // created_at / updated_at / deleted_at
		softdelete.SoftDeleteMixin{},
	}
}

// Fields of the AdminUser.
func (AdminUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
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
			Unique().
			Optional(),

		field.String("phone").
			Comment("手机号").
			Optional(),

		field.String("avatar").
			Comment("头像 URL").
			Default("/uploads/avatars/meteor-default.jpg"),

		field.Bool("is_super").
			Default(false).
			Comment("是否超级管理员"),

		field.Bool("is_active").
			Comment("是否启用").
			Default(true),

		field.Time("last_login_at").
			Comment("最后登录时间").
			Optional().
			Nillable(),

		field.String("reset_token").
			Comment("密码重置令牌").
			Optional().
			Nillable(),

		field.Time("reset_expire").
			Comment("密码重置令牌过期时间").
			Optional().
			Nillable(),
	}
}

// Edges of the AdminUser.
func (AdminUser) Edges() []ent.Edge {
	return nil
}
