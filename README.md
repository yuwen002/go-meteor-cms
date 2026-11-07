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
- ✍️ 文章发布：Markdown 编辑与预览、草稿 / 待审核 / 发布 / 下线状态机、定时发布、置顶、封面、SEO
- 🗃️ 内容分类与标签：分类树、标签管理，均带 slug 与文章数统计
- 🔎 内容治理：审核流与数据权限、修订历史与版本回滚、回收站（恢复 / 彻底删除）、旧 slug 301
- 🔒 安全与体验：服务端 SEO 注入、浏览量 IP 去重、公开接口限流、登录验证码开关、安全 Markdown 渲染
- 📡 内容分发：相关文章 / 上下篇、作者展示、RSS 2.0、sitemap.xml、robots.txt
- 🌐 公开阅读站点：`/articles` 前台列表与详情，无需登录，正文安全渲染
- ⚡ 性能与运维：组合索引、列表排序与批量操作、公开导航聚合计数、`/healthz` 探活、前端依赖分包、写事务独立超时、权限种子按接口幂等去重
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
└── docker-compose.yml     # 全栈一键编排（MySQL + Redis + 迁移 + API，前端打包进镜像）
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
| 文章 | `GET/POST /admin/articles`、`GET/PUT/DELETE /admin/articles/:id`、`POST /admin/articles/:id/publish`、`/withdraw`、`/toggle-top` | 草稿 / 发布 / 下线状态机、定时发布、置顶、Markdown 正文、封面与 SEO |
| 分类 | `GET/POST /admin/categories`、`PUT/DELETE /admin/categories/:id` | 父子层级、slug、删除前置校验 |
| 标签 | `GET/POST /admin/tags`、`PUT/DELETE /admin/tags/:id` | 名称 / slug 唯一、文章数统计 |
| 审核与修订 | `POST /admin/articles/:id/submit` `/approve` `/reject`、`GET/POST /admin/articles/:id/revisions...` | 待审核状态机、编辑前存档、版本回滚 |
| 回收站 | `GET /admin/trash`、`POST /admin/trash/:type/:id/restore`、`DELETE /admin/trash/:type/:id/purge` | 软删除内容恢复 / 物理删除与孤儿文件清理 |

另有**无需认证**的公开接口：`GET /api/articles`、`GET /api/articles/:slug`、`GET /api/taxonomy`、
`GET /api/rss`、`GET /sitemap.xml`、`GET /robots.txt`，
对应前台阅读页 `/articles` 与 `/articles/:slug`。字段定义与状态机说明见 `docs/api.md`「文章发布」章节。

OpenAPI 文档由 `make gen-doc` 生成：合并 `api/cms/v1/cms.api` 与
`api/cms/v1/content.doc.api`（内容模块文档源，不参与 goctl），输出根目录 `openapi.json`，
镜像内同时发布到 `/openapi.json` 供 `/docs.html` 使用。

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

# 重新生成全量 OpenAPI 文档并校验内容模块覆盖
make gen-doc && go test ./cmd/apidoc
```

当前已覆盖的纯逻辑单测（均不依赖外部 MySQL / Redis / Docker，可直接进 CI）：

- `internal/utils`：分页参数规范化（默认值、负值纠正、上限截断）
- `internal/pkg/datascope`：部门树子孙收集，含脏数据成环时不死循环
- `internal/pkg/mail`：邮件报文构造（UTF-8 主题与正文编码）、SMTP 参数校验
- `internal/pkg/tokenversion`：token 版本号读取 / 自增、按用户隔离、未配置 Redis 时降级（进程内 miniredis）
- `api/cms/v1/internal/middleware`：路由匹配规则、静态路由优先于动态路由、免鉴权路由判定；
  SPA 托管的 history 回退与目录穿越防护；JWT 中间件的 token 版本比对与全局踢人
- `api/cms/v1/internal/logic`：文章状态机与公开可见性（含定时发布到点放行）、slug 规范化与中文兜底、发布时间格式解析、写事务 ctx 剥离请求取消信号后保留 trace/claims
- `internal/seed`：权限种子的 `(method, api_path)` 唯一性、文章模块 API / 按钮 / 菜单齐备、菜单父级声明顺序、contentroutes 全路由登记覆盖、同权限码多接口不被覆盖
- `cmd/apidoc`：OpenAPI 输出覆盖内容模块路径、schema、query / path 参数与 RSS / sitemap 纯文本响应
- JWT 中间件测试用进程内 Redis + 空结果集 ent driver 替代数据库，无需 MySQL / Docker

只跑文章相关的回归用例：

```bash
go test ./api/cms/v1/internal/logic -run 'Test(NormalizeSlug|GenerateSlug|IsValidSlug|IsPubliclyVisible|ParsePublishTime|ArticleStatus)'
go test ./internal/seed
```

## 部署

### Docker（全栈一键）

`docker-compose.yml` 编排完整栈：MySQL（数据） + Redis（限流 / 全局踢人） + 迁移（表结构 + 种子） + API（含前端 SPA）。
镜像由仓库根目录 `Dockerfile` 多阶段构建：前端 `web/` 的构建产物被打进同一镜像，由 Go 服务在 8888 端口统一托管。

前端由 Go 服务通过 `rest.WithNotFoundHandler` 托管（见 `api/cms/v1/internal/middleware/spa.go`）：
`/`、`/login`、`/dashboard`、`/system/*` 等未匹配路由回退到 `index.html`（保证 `createWebHistory` 刷新不 404），
`/assets/*` 等真实文件直接返回；`/admin`、`/api`、`/uploads` 前缀保留给 API 与上传资源，不参与回退。
注意不要用 `server.Use` 注册 SPA 兜底：go-zero 的全局中间件不作用于 404（未匹配）路由。

> Windows / macOS 需先启动 Docker Desktop 并等引擎就绪（`docker info` 能返回 ServerVersion），
> 否则 `docker compose` 会报 `dockerDesktopLinuxEngine` 管道不存在。

```bash
cp .env.example .env          # 可选：按需修改端口、口令、构建镜像源
# 一键启动：mysql 就绪 → migrate 迁移+种子 → api 启动并托管前端
docker compose up -d

# 跟踪日志（看迁移进度 / API 启动）
docker compose logs -f api

# 停止并保留数据
docker compose down
# 停止并清空数据（含库）
docker compose down -v
```

应用通过环境变量连接依赖，无需改动 `cms-api.yaml`：`MYSQL_DSN` 指向 `mysql:3306`、`REDIS_HOST` 指向 `redis:6379`、
`WEB_DIST=/app/web/dist`。这些值由 compose 从 `.env` 注入容器，**部署到别的机器时改 `.env` 即可，不要改 `docker-compose.yml`**。

配置项集中在 `.env`（复制 `.env.example` 得到，含开发默认值）：可覆盖 root 密码、库名、各服务端口等；
不提供 `.env` 时 compose 的 `${VAR:-默认值}` 自动生效，开箱即用。注意改密码 / 库名时要同步改 `MYSQL_DSN` 中的对应部分。

受限网络下容器构建可能拉不到默认官方源，可在 `.env` 里启用镜像源（它们是 **build args**，只影响构建、不进入运行容器）：

```bash
NPM_REGISTRY=https://registry.npmmirror.com/
GOPROXY=https://goproxy.cn,direct
```

修改后需重新构建：`docker compose build && docker compose up -d`。

> 端口说明：mysql 默认**不映射**宿主机端口（避免与本地 MySQL 的 3306 冲突），
> 如需从宿主机直连排查，取消 `docker-compose.yml` 中 mysql 的 `ports` 注释并设 `MYSQL_PORT`（如 `3307`）。
> redis 默认映射 `${REDIS_PORT:-6379}:6379` 供本机调试；若 6379 已被本地 Redis 占用，
> 在 `.env` 里改 `REDIS_PORT`（如 `6380`）即可——容器内服务间仍走 `redis:6379`，不受影响。

仅想用容器里的 MySQL / Redis、应用仍在宿主机跑（需先按上方说明打开 mysql 的 `ports` 映射；
redis 默认已映射 `${REDIS_PORT:-6379}:6379`）：

```bash
docker compose up -d mysql redis
make migrate && make run                  # 或 go run ./api/cms/v1 -f api/cms/v1/etc/cms-api.yaml
```

或单独构建运行：

```bash
docker build -t go-meteor-cms .
docker run -d -p 8888:8888 \
  -e MYSQL_DSN="root:123456@tcp(host.docker.internal:3306)/go_meteor_cms?parseTime=True&loc=Local" \
  -e REDIS_HOST="host.docker.internal:6379" \
  -e WEB_DIST="/app/web/dist" \
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

应用进程读取的环境变量（本地直跑 / 单独 `docker run` 时同样适用）：

| 变量 | 说明 |
| --- | --- |
| `MYSQL_DSN` | 数据库连接串，设置后覆盖配置文件中的 `Database.Source`（容器化部署推荐；未设置时默认 `127.0.0.1:3306`） |
| `REDIS_HOST` | Redis 地址，容器内置为 `redis:6379`；未设置时按配置文件默认 `127.0.0.1:6379` |
| `WEB_DIST` | 前端构建产物目录；容器内置为 `/app/web/dist`，本地开发留空则仅暴露 API |

compose 专用配置（写在 `.env`，由 `${VAR:-默认值}` 插值；不提供 `.env` 时用默认值）：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `CMS_IMAGE` | `go-meteor-cms:local` | migrate / api 使用的镜像；生产固定为已发布版本，改回旧版本即可回滚 |
| `MYSQL_ROOT_PASSWORD` | `123456` | MySQL root 口令，需与 `MYSQL_DSN` 保持一致 |
| `MYSQL_DATABASE` | `go_meteor_cms` | 库名，需与 `MYSQL_DSN` 保持一致 |
| `API_PORT` | `8888` | API / 前端对外端口 |
| `REDIS_PORT` | `6379` | Redis 映射到宿主机的端口；被占用时可改（如 `6380`） |
| `MYSQL_PORT` | `3306` | 仅在取消 compose 中 mysql `ports` 注释后生效 |
| `NPM_REGISTRY` | `https://registry.npmjs.org/` | 构建参数：前端 `npm ci` 源，受限网络可换 npmmirror |
| `GOPROXY` | `https://proxy.golang.org,direct` | 构建参数：Go 模块代理，受限网络可换 goproxy.cn |

### 版本发布与回滚

发布由 `.github/workflows/release.yml` 完成：推送 `vX.Y.Z` 标签即自动验证、构建并推送镜像到 GHCR，同时创建 GitHub Release。
使用仓库内置 `GITHUB_TOKEN`（`packages: write`），**无需配置额外密钥**；首次使用需在 GitHub 把该 Package 设为 Public，或在目标机器先执行 `docker login ghcr.io`。

```bash
# 在已验证的 commit 上打标签触发发版
git tag -a v1.2.3 -m "release v1.2.3"
git push origin v1.2.3
# 产出：ghcr.io/yuwen002/go-meteor-cms:1.2.3 / :1.2 / :latest
```

生产部署在 `.env` 中固定版本（不要用 `latest`）：

```bash
CMS_IMAGE=ghcr.io/yuwen002/go-meteor-cms:1.2.3
docker compose pull && docker compose up -d
```

> compose 同时保留了 `build:`，因此**必须先 `pull` 且不要加 `--build`**，否则会在本机重新构建并覆盖已发布镜像。

回滚只需把 `CMS_IMAGE` 改回上一个版本号并重复上述命令：

```bash
CMS_IMAGE=ghcr.io/yuwen002/go-meteor-cms:1.2.2
docker compose pull && docker compose up -d
```

> `migrate` 服务在每次 `up -d` 时幂等执行迁移（只增不减）。回滚镜像前请确认旧版本与当前表结构兼容；若发版包含破坏性 schema 变更，需另行评估数据回滚方案。

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