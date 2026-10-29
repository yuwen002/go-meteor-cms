# Go Meteor CMS

[![Go Report Card](https://goreportcard.com/badge/github.com/yuwen002/go-meteor-cms)](https://goreportcard.com/report/github.com/yuwen002/go-meteor-cms)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

基于 Go 和 Go-Zero 开发的内容管理系统后端 API 服务。

## 功能特性

- ✅ 管理员登录认证
- 🔑 JWT Token 认证
- 🔐 密码重置功能
- 📝 管理员注册功能
- 📋 管理员列表功能
- 🔍 Token 验证测试接口
- 🚀 高性能 API 服务
- 📦 使用 Ent 作为 ORM
- 🛡️ 统一的错误处理
- 📊 结构化的日志记录
- 🗄️ 数据库自动迁移
- 🌱 初始化数据种子
- 🖼️ 图形验证码功能
- 🔐 管理员密码修改功能
- 🔁 管理员密码重置功能

## 技术栈

- **框架**: [Go-Zero](https://go-zero.dev/) v1.10.3
- **数据库**: MySQL (通过 [Ent ORM](https://entgo.io/) v0.14.6 支持)
- **认证**: JWT (github.com/golang-jwt/jwt/v5 v5.3.1)
- **密码加密**: bcrypt (golang.org/x/crypto v0.58.0)
- **参数验证**: go-playground/validator/v10 v10.30.5
- **图形验证码**: base64Captcha (github.com/mojocn/base64Captcha v1.3.8)
- **API 规范**: RESTful API
- **依赖管理**: Go Modules (Go 1.26.0)
- **代码生成**: goctl 1.10.2 (项目级自定义模板位于 `.goctl/`)

## 快速开始

### 环境要求

- Go 1.26+
- MySQL 5.7+
- Git
- goctl 1.10.2（仅代码生成时需要）

### 安装

1. 克隆仓库

```bash
git clone https://github.com/yuwen002/go-meteor-cms.git
cd go-meteor-cms
```

2. 安装依赖

```bash
go mod tidy
```

3. 数据库配置

确保 MySQL 服务正在运行，并创建数据库：
```sql
CREATE DATABASE go_meteor_cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改数据库连接配置：
```bash
# 数据库连接信息在以下文件中配置
api/cms/v1/etc/cms-api.yaml
cmd/migrate/main.go
```

4. 数据库迁移

```bash
# 执行数据库迁移和初始化数据
go run cmd/migrate/main.go
```

5. 启动服务

```bash
# 启动 API 服务
go run api/cms/v1/cms.go
```

## 项目结构

```
.
├── api/                    # API 定义
│   └── cms/               # CMS 服务
│       └── v1/            # API 版本
├── cmd/                   # 命令行工具
│   └── migrate/           # 数据库迁移工具
├── ent/                   # Ent ORM 实体
│   └── schema/            # 数据库表结构定义
├── internal/              # 内部包
│   ├── common/            # 通用组件
│   ├── seed/              # 初始化数据
│   └── utils/             # 工具函数
└── rpc/                   # RPC 服务定义
```

## API 接口

启动服务后，默认访问地址：`http://localhost:8888`

所有接口统一返回下列结构，HTTP 状态码恒为 200，业务结果由 `code` 字段表示：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

- `code = 0`：成功
- `code > 0`：业务错误，`msg` 为提示信息，完整错误码见 [`internal/common/error_code.go`](internal/common/error_code.go)

### 公共接口

1. **获取验证码**
   - URL: `GET /api/captcha`
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "captcha_id": "验证码ID",
         "captcha_base64": "验证码图片base64编码"
       }
     }
     ```

2. **管理员登录**
   - URL: `POST /admin/login`
   - 参数:
     ```json
     {
       "username": "admin",
       "password": "123456",
       "captcha_id": "验证码ID",
       "captcha": "验证码"
     }
     ```
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
       }
     }
     ```

3. **忘记密码**
   - URL: `POST /admin/forgot-password`
   - 参数:
     ```json
     {
       "username": "admin"
     }
     ```
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "message": "重置密码邮件已发送，请检查邮箱"
       }
     }
     ```

4. **管理员注册**
   - URL: `POST /admin/register`
   - 参数:
     ```json
     {
       "username": "newadmin",
       "password": "password123",
       "email": "newadmin@example.com",
       "nickname": "新管理员"
     }
     ```
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "message": "管理员注册成功"
       }
     }
     ```

### 需要认证的接口

以下接口需要在请求头中添加 Authorization:

```
Authorization: Bearer <token>
```

1. **Token 验证测试**
   - URL: `GET /admin/test-token`
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "message": "Token OK",
         "claims": {
           "user_id": 1,
           "username": "admin"
         }
       }
     }
     ```

2. **获取管理员列表**
   - URL: `GET /admin/admin-users`
   - 查询参数:
     - `page`: 页码（必填）
     - `page_size`: 每页数量（必填）
     - `keyword`: 搜索关键词（用户名/昵称/邮箱，选填）
     - `active`: 启用状态过滤（1-启用，2-禁用，选填）
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "total": 10,
         "page": 1,
         "page_size": 10,
         "list": [
           {
             "id": 1,
             "username": "admin",
             "nickname": "管理员",
             "email": "admin@example.com",
             "phone": "13800138000",
             "avatar": "/uploads/avatars/meteor-default.jpg",
             "isSuper": true,
             "isActive": true,
             "last_login_at": "2023-01-01T12:00:00Z",
             "createdAt": "2023-01-01T12:00:00Z"
           }
         ]
       }
     }
     ```

3. **修改自己的密码**
   - URL: `PUT /admin/change-password`
   - 参数:
     ```json
     {
       "old_password": "oldpassword",
       "new_password": "newpassword"
     }
     ```
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "message": "密码修改成功"
       }
     }
     ```

4. **重置其他管理员的密码** (需要超级管理员权限)
   - URL: `PUT /admin/users/:id/reset-password`
   - 参数:
     ```json
     {
       "id": 2,
       "new_password": "newpassword"
     }
     ```
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "message": "密码已重置"
       }
     }
     ```

## 默认管理员账户

系统初始化时会自动创建默认管理员账户：

- 用户名: `admin`
- 密码: `123456`

请在生产环境中及时修改默认密码。

## 开发指南

### 代码生成

项目在 `.goctl/` 下维护了自定义的 handler 模板（统一使用 `common.Ok` / `common.Fail` 响应），因此生成时必须用 `--home` 指定该目录，否则会退回内置模板产出风格不一致的代码：

```bash
# 安装指定版本的 goctl
go install github.com/zeromicro/go-zero/tools/goctl@v1.10.2

# 生成 API 代码 (基于 api/cms/v1/cms.api)
# 注意：goctl 不会覆盖已存在的 handler/logic 文件，只会新增缺失的文件
goctl api go -api api/cms/v1/cms.api -dir api/cms/v1 --home .goctl

# 生成 Ent 代码 (详见 docs/ent.md)
go run -mod=mod entgo.io/ent/cmd/ent generate --feature intercept,schema/snapshot ./ent/schema
```

> 不要执行 `goctl template init --home .goctl`，该命令会覆盖 `.goctl/api/handler.tpl` 中的自定义模板。

### 测试

```bash
# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## 部署

项目暂未内置 Dockerfile 与 docker-compose.yml，容器化部署需自行补充，运行参数参考：

```bash
# 构建镜像（需自行提供 Dockerfile）
docker build -t go-meteor-cms .

# 运行容器
docker run -d -p 8888:8888 go-meteor-cms
```

## 贡献指南

欢迎提交 Issue 和 Pull Request。

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

[MIT](LICENSE) © 2023 Your Name

## 致谢

- [Go-Zero](https://go-zero.dev/)
- [Ent](https://entgo.io/)
- [Go](https://golang.org/)