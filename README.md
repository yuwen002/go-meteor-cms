# Go Meteor CMS

[![Go Report Card](https://goreportcard.com/badge/github.com/yuwen002/go-meteor-cms)](https://goreportcard.com/report/github.com/yuwen002/go-meteor-cms)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

基于 Go 和 Go-Zero 开发的内容管理系统后端 API 服务。

## 功能特性

- ✅ 管理员登录认证
- 🔑 JWT Token 认证
- 🔐 密码重置功能
- 🚀 高性能 API 服务
- 📦 使用 Ent 作为 ORM
- 🛡️ 统一的错误处理
- 📊 结构化的日志记录
- 🗄️ 数据库自动迁移
- 🌱 初始化数据种子

## 技术栈

- **框架**: [Go-Zero](https://go-zero.dev/) v1.9.2
- **数据库**: MySQL (通过 [Ent ORM](https://entgo.io/) v0.14.5 支持)
- **认证**: JWT (github.com/golang-jwt/jwt/v5 v5.3.0)
- **密码加密**: bcrypt (golang.org/x/crypto)
- **API 规范**: RESTful API
- **依赖管理**: Go Modules (Go 1.25.3)

## 快速开始

### 环境要求

- Go 1.25+
- MySQL 5.7+
- Git

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

### 管理员接口

1. **管理员登录**
   - URL: `POST /admin/login`
   - 参数:
     ```json
     {
       "username": "admin",
       "password": "123456"
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

2. **忘记密码**
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
         "status": 1,
         "message": "密码重置邮件已发送"
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

```bash
# 生成 API 代码 (基于 api/cms/v1/cms.api)
goctl api go -api api/cms/v1/cms.api -dir .

# 生成 Ent 代码 (基于 ent/schema/*.go)
go generate ./ent
```

### 测试

```bash
# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## 部署

### 使用 Docker

```bash
# 构建镜像
docker build -t go-meteor-cms .

# 运行容器
docker run -d -p 8888:8888 go-meteor-cms
```

### 使用 Docker Compose

```bash
docker-compose up -d
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