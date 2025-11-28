# Ent 代码生成命令

```bash
go run -mod=mod entgo.io/ent/cmd/ent generate --feature intercept,schema/snapshot ./ent/schema
```

## 功能说明

### 1. 拦截器 (Intercept)
- 允许添加钩子来拦截和修改查询和变更操作
- 用于实现软删除等高级功能

### 2. 模式快照 (Schema Snapshot)
- 跟踪数据库模式变更
- 管理数据库迁移

## 伪删除 (Soft Delete)
- 通过设置 `delete_at` 字段标记记录为已删除
- 查询时自动过滤已删除的记录

## 使用方法
1. 在项目根目录下运行生成命令
2. 确保 `ent/schema` 目录下已有定义好的 schema 文件
3. 生成的代码位于 `ent` 包中

## 注意事项
- 首次运行或 schema 有变更时需要执行此命令
- 确保所有 schema 文件已保存
- 如果遇到问题，可以尝试删除 `ent/generate` 目录后重新生成
