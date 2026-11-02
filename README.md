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
- 🎯 RBAC 权限模型（菜单 / 按钮 / API 三级鉴权）
- 🧭 我的菜单接口（按角色动态下发菜单树与权限码）
- 🔒 数据权限（按部门范围收敛可见数据）
- ✉️ 忘记密码闭环（邮件令牌 + 一次性重置）
- 📄 统一分页参数（每页上限 100）
- 🧭 动态菜单：侧边栏由后端按角色下发的菜单树渲染，含图标 / 排序 / 多级层级
- 📚 数据字典：字典类型 + 字典项两级管理（字典项按 `dict_code` 归属）
- ⚙️ 系统配置：键值对配置项管理
- 🖼️ 文件管理：上传 / 列表 / 重命名 / 删除
- 📈 仪表盘统计：管理员 / 部门 / 角色 / 权限总量与今日登录、本周操作数
- 🗂️ 操作日志与登录日志：分页查询，便于审计
- 🖥️ 配套 Vue3 + Element Plus 前端管理后台（见「前端 Web 管理后台」一节）

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
- **前端管理后台**:
  - **框架**: [Vue 3](https://vuejs.org/) + [Vue Router 4](https://router.vuejs.org/) + [Pinia](https://pinia.vuejs.dev/)
  - **UI 组件库**: [Element Plus](https://element-plus.org/) + [@element-plus/icons-vue](https://element-plus.org/)
  - **构建工具**: [Vite](https://vitejs.dev/) 5 + [TypeScript](https://www.typescriptlang.org/) 5
  - **类型检查**: vue-tsc

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

数据库连接信息按以下优先级获取（API 服务与迁移工具共用同一份配置）：

1. 环境变量 `MYSQL_DSN`（推荐用于 CI / 生产）
2. `api/cms/v1/etc/cms-api.yaml` 中的 `Database.Source`

```bash
export MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/go_meteor_cms?parseTime=True&loc=Local"
```

邮件配置（`etc/cms-api.yaml` 的 `Mail` 节点）用于忘记密码时投递重置令牌：

```yaml
Mail:
  Enabled: false                  # false 时只输出到服务日志，不发真实邮件
  Host: "smtp.example.com"
  Port: 587                       # 587(STARTTLS) / 465(SSL) / 25
  Username: "noreply@example.com"
  Password: ""                    # 密码或授权码
  From: "Go Meteor CMS <noreply@example.com>"
```

4. 数据库迁移

```bash
# 执行数据库迁移和初始化数据（幂等，可重复执行）
go run cmd/migrate/main.go
```

迁移完成后会输出权限统计，形如 `📊 权限统计：菜单 11 / 按钮 52 / API 52`（合计 115 条），其中 API 权限直接决定非超级管理员能否访问各接口。

5. 启动服务

```bash
# 启动 API 服务（需用 -f 指定配置文件）
go run ./api/cms/v1 -f api/cms/v1/etc/cms-api.yaml
```

## 前端 Web 管理后台

仓库内 `web/` 为配套的前端管理后台，技术栈见上方「技术栈 · 前端管理后台」。

### 环境要求

- Node.js 18+（建议 LTS）
- 包管理器 npm

### 本地开发

```bash
cd web
npm install
npm run dev          # 启动开发服务器，默认 http://localhost:5173
```

开发服务器通过 Vite 代理将 `/admin`、`/api`、`/uploads` 转发到后端 `http://localhost:8888`，
因此本地联调时需要先在 8888 端口启动后端服务。

### 构建与预览

```bash
npm run build        # 类型检查（vue-tsc）+ 生产构建，产物在 web/dist
npm run preview      # 本地预览构建产物
```

### 动态菜单说明

侧边栏**不写死在代码里**，而是登录后由 `GET /admin/profile/permissions` 返回的菜单树驱动：

- 后端按角色聚合出 `menus`（字段含 `id / name / path / component / icon / children`），
  前端 `Sidebar.vue` 递归渲染 `MenuItem.vue`（`el-sub-menu` + `el-menu-item`）。
- 图标为语义键（如 `dashboard`、`dict`、`setting`），前端 `MenuItem.vue` 的 `iconMap` 映射到 Element Plus 图标组件；未命中时回落为 `Menu` 图标。
- 菜单按 `sort` 升序，后端在 `buildMenuTree` 中按 `parent_id` 组装多级层级；当父菜单对当前用户无权限时，其子菜单会自动降级为顶级菜单，不会丢失。

> 新增自定义页面时，请同步在种子 `internal/seed/permissions.go` 中维护对应「菜单 / 按钮 / API」三类权限，
> 并在 `MenuItem.vue` 的 `iconMap` 中补充图标键，否则侧边栏图标会回落为默认图标。

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
│   ├── common/            # 通用组件（响应、错误码、权限类型）
│   ├── pkg/               # 可复用子包（mail 邮件、datascope 数据权限、cleanup 维护任务）
│   ├── seed/              # 初始化数据（含权限清单）
│   └── utils/             # 工具函数（jwt、password、captcha、page）
├── web/                   # 前端管理后台（Vue3 + Element Plus，侧边栏由后端菜单树驱动）
├── rpc/                   # RPC 服务定义
├── Makefile               # 常用开发命令
├── Dockerfile             # 多阶段构建镜像
└── docker-compose.yml     # MySQL + 迁移 + API 一键编排
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

> 分页参数通用约定：`page` 小于 1 时按 1 处理，`page_size` 小于 1 时取接口默认值，大于 100 一律截断为 100，响应中的 `page_size` 为实际生效值。
> 完整的接口清单（含角色、权限、部门模块）见 [`docs/api.md`](docs/api.md) 与 `docs/api-docs/go-meteor-cms.apifox.json`。

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

3. **忘记密码（第一步：申请重置令牌）**
   - URL: `POST /admin/forgot-password`
   - 说明：无论账号是否存在均返回成功，避免账号枚举；令牌 1 小时内有效
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

4. **忘记密码（第二步：用令牌重置密码）**
   - URL: `POST /admin/reset-password`
   - 说明：令牌一次性，使用后立即失效；新密码不少于 6 位
   - 参数:
     ```json
     {
       "token": "邮件中的重置令牌",
       "new_password": "newpassword123"
     }
     ```
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "message": "密码重置成功，请使用新密码登录"
       }
     }
     ```

5. **管理员注册**
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

除管理员列表接口会额外按数据权限收敛结果外，认证接口的可见性由 RBAC 决定，详见下方「权限模型」与「数据权限」两节。

1. **获取当前用户的菜单与权限**（前端据此渲染菜单与按钮）
   - URL: `GET /admin/profile/permissions`
   - 响应:
     ```json
     {
       "code": 0,
       "msg": "success",
       "data": {
         "is_super": false,
         "menus": [
           {
             "id": 3,
             "name": "部门管理",
             "path": "/admin/departments",
             "component": "views/admin/dept/list.vue",
             "icon": "dept",
             "permission": "system:dept:list",
             "children": []
           }
         ],
         "permissions": ["system:dept:list", "system:dept:add"]
       }
     }
     ```

2. **角色部门范围**
   - URL: `GET /admin/roles/:id/departments`
   - URL: `PUT /admin/roles/:id/departments`
   - 说明：仅 `data_scope=5`（自定义部门）的角色支持设置，系统内置角色不可修改
   - PUT 参数:
     ```json
     { "dept_ids": [1, 2] }
     ```

3. **Token 验证测试**
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

4. **获取管理员列表**（按当前用户的数据权限收敛结果）
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

5. **修改自己的密码**
   - URL: `PUT /admin/admin-users/me/change-password`
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

6. **重置其他管理员的密码** (需要超级管理员权限，且目标不能是超级管理员)
   - URL: `PUT /admin/admin-users/reset-password`
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

### 业务模块一览

除上方「公共接口」与「需要认证的接口」示例外，其余业务接口按模块归纳如下（均需 `Authorization: Bearer <token>`，完整字段定义见 `docs/api.md`）：

| 模块 | 主要接口 | 说明 |
| --- | --- | --- |
| 仪表盘 | `GET /admin/dashboard/stats` | 管理员 / 部门 / 角色 / 权限总量，今日登录与本周操作数 |
| 管理员 | `GET/POST /admin/admin-users`、`GET/PUT/DELETE /admin/admin-users/:id`、禁用/启用、重置密码、绑定部门、角色分配 | 分页 + 数据权限收敛 |
| 部门 | `GET /admin/departments/tree`、`POST /admin/departments`、`GET/PUT/DELETE /admin/departments/:id` | 树形结构 |
| 角色 | 角色 CRUD、`GET/PUT /admin/roles/:id/permissions`、`GET/PUT /admin/roles/:id/departments` | 权限与数据范围分配 |
| 权限 | `GET /admin/permissions`、`GET /admin/permissions/tree`、CRUD | 菜单 / 按钮 / API 三级定义 |
| 数据字典 | `GET/POST /admin/dicts` + CRUD、`GET/POST /admin/dicts/items` + 字典项 CRUD | 字典类型 + 字典项 |
| 系统配置 | `GET/POST /admin/configs` + CRUD | 键值对配置 |
| 文件 | `GET /admin/files`、`PUT /admin/files/rename`、`POST /admin/files/delete`、`POST /admin/upload` | 上传 / 列表 / 重命名 / 删除 |
| 日志 | `GET /admin/logs/operations`、`GET /admin/logs/logins` | 操作 / 登录审计 |
| 个人中心 | `GET/PUT /admin/profile`、`PUT /admin/admin-users/me/change-password`、`GET /admin/profile/permissions` | 自身资料与菜单权限 |

## 默认管理员账户

系统初始化时会自动创建默认管理员账户：

- 用户名: `admin`
- 密码: `123456`

请在生产环境中及时修改默认密码。

## 权限模型

权限记录存放在 `admin_permissions` 表，由迁移时的种子写入，分三类：

| 类型 | 取值 | 用途 |
| --- | --- | --- |
| 菜单 | 1 | 前端路由菜单，`GET /admin/profile/permissions` 据此下发菜单树 |
| 按钮 | 2 | 前端功能点，同 `permissions` 字段下发，用于控制按钮显隐 |
| API | 3 | 后端接口，权限中间件按 `method + api_path` 匹配后校验角色是否拥有对应标识 |

关键约定：

- **按钮与 API 使用相同的权限标识**（如 `system:role:add`），因此授权一次即可同时让按钮可见、接口可调用
- **API 权限的 `api_path` 必须与 `cms.api` 中的真实路由一致**，否则非超级管理员访问时会被判定为「未配置权限」而拒绝；新增受保护接口后请同步追加一条 Type=3 权限
- 路由匹配按「静态段命中数最多者优先」，`/admin/departments/tree` 不会被 `/admin/departments/:id` 抢占
- 以下路由不参与 API 鉴权：`/admin/logout`、`/admin/test-token`、`/admin/admin-users/me/change-password`、`/admin/profile/permissions`
- 超级管理员（`admin_users.is_super`）跳过全部鉴权与数据范围过滤

## 数据权限

角色通过 `data_scope` 声明可见数据范围：

| data_scope | 含义 | 说明 |
| --- | --- | --- |
| 1 | 全部数据 | 不做行级过滤 |
| 2 | 本部门 | 与当前管理员的 `dept_id` 一致 |
| 3 | 本部门及子部门 | 沿部门树展开全部子孙部门 |
| 4 | 仅自己 | 只能看到自己的记录 |
| 5 | 自定义部门 | 读取 `admin_role_depts`，通过 `PUT /admin/roles/:id/departments` 维护 |

- 多个角色时取各角色可见范围的并集；任一角色为「全部数据」时不过滤
- **安全原则（fail-closed）**：范围解析失败、用户未分配任何启用角色、或要求按部门过滤但未分配部门时，一律按最严格处理（不返回数据），绝不退回全量
- 当前仅在管理员列表接口落地，其他列表接口暂不过滤

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
# 运行全部测试
go test ./...

# 只跑带用例的包（其余包显示 no test files）
go test ./internal/... ./api/cms/v1/internal/middleware/...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

当前已覆盖的纯逻辑单测：

- `internal/utils`：分页参数规范化（默认值、负值纠正、上限截断）
- `internal/pkg/datascope`：部门树子孙收集，含脏数据成环时不死循环
- `internal/pkg/mail`：邮件报文构造（UTF-8 主题与正文编码）、SMTP 参数校验
- `api/cms/v1/internal/middleware`：路由匹配规则、静态路由优先于动态路由、免鉴权路由判定

## 部署

### Docker（推荐）

项目已内置 `Dockerfile`（多阶段构建，最终镜像基于 alpine，以非 root 用户运行）与 `docker-compose.yml`。

```bash
# 一键启动：MySQL（健康检查通过）→ 迁移（执行成功）→ API 服务
docker compose up -d --wait

# 停止
docker compose down
```

或单独构建运行：

```bash
docker build -t go-meteor-cms .
docker run -d -p 8888:8888 \
  -e MYSQL_DSN="root:123456@tcp(host.docker.internal:3306)/go_meteor_cms?parseTime=True&loc=Local" \
  go-meteor-cms
```

容器内的迁移工具位于 `/app/migrate`，可在部署流程中单独执行：

```bash
docker run --rm -e MYSQL_DSN="..." --entrypoint /app/migrate go-meteor-cms
```

### 编译部署

```bash
go build -o bin/cms ./api/cms/v1
go build -o bin/migrate ./cmd/migrate

# 先迁移（幂等），再启动服务
./bin/migrate
./bin/cms -f api/cms/v1/etc/cms-api.yaml
```

### 环境变量

| 变量 | 说明 |
| --- | --- |
| `MYSQL_DSN` | 数据库连接串，设置后覆盖配置文件中的 `Database.Source`（容器化部署推荐） |

### 后台维护任务

服务启动时会自动清理已过期的 token 黑名单记录，之后按 `Maintenance.TokenBlacklistCleanInterval`（默认 3600 秒）周期性执行，避免 `token_blacklists` 表随每次登出无限增长。

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