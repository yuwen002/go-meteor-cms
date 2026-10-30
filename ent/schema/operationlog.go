package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OperationLog 操作日志表，记录后台接口调用情况，用于安全审计与问题排查
type OperationLog struct {
	ent.Schema
}

// Annotations 表和字段的额外注解
func (OperationLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		schema.Comment("操作日志表，记录受保护接口的调用情况"),
	}
}

// Fields 定义 OperationLog 实体字段
func (OperationLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Comment("自增主键 ID").
			Unique().
			Positive().
			Immutable(),
		field.Int64("user_id").
			Optional().
			Comment("操作人ID，未登录时为空"),
		field.String("username").
			MaxLen(50).
			Optional().
			Comment("操作人账号"),
		field.String("module").
			MaxLen(50).
			Optional().
			Comment("业务模块，如 管理员管理"),
		field.String("action").
			MaxLen(100).
			Optional().
			Comment("操作描述，如 新增管理员"),
		field.String("method").
			MaxLen(10).
			Comment("HTTP 方法"),
		field.String("path").
			MaxLen(255).
			Comment("请求路径"),
		field.String("ip").
			MaxLen(45).
			Optional().
			Comment("来源IP，兼容 IPv6"),
		field.String("user_agent").
			MaxLen(512).
			Optional().
			Comment("客户端标识"),
		field.Int("status").
			Default(0).
			Comment("操作结果：0 成功，1 失败"),
		field.String("error_msg").
			MaxLen(512).
			Optional().
			Comment("失败原因"),
		field.Int64("duration_ms").
			Default(0).
			Comment("耗时（毫秒）"),
		field.Text("params").
			Optional().
			Comment("请求参数快照，敏感字段已脱敏并截断"),
	}
}

// Mixin 复用创建时间、更新时间等通用字段。
// 注意：审计日志刻意不继承 SoftDeleteMixin，保证记录不可被软删除/篡改，满足审计完整性。
func (OperationLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Indexes 定义索引
func (OperationLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("status"),
	}
}

// Edges 暂不定义关联，避免级联删除把日志一起删掉
func (OperationLog) Edges() []ent.Edge {
	return nil
}
