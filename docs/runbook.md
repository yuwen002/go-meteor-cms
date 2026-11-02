# 运维 Runbook（本地联调 / 生产部署 / 故障排查）

本文档面向开发与运维，覆盖环境搭建、健康检查与常见故障排查。与 `README.md` 的"快速开始"互补：`README` 讲特性与接口，`runbook` 讲"怎么跑起来、出问题怎么查"。

---

## 1. 架构概览

```
浏览器 ──(5173, Vite 代理 /admin,/api,/uploads)──▶ Go-Zero API(:8888) ──▶ MySQL
                                  ▲
                          前端静态资源(web/dist, 生产可反向代理)
```

- **后端**：Go-Zero，默认 `:8888`，JWT 鉴权，RBAC（菜单/按钮/API 三级）。
- **前端**：Vue3 + Element Plus，开发态 `:5173`，生产构建为静态文件。
- **数据库**：MySQL（Ent ORM），启动时由 `cmd/migrate` 自动建表 + 种子数据。

---

## 2. 本地开发

### 2.1 数据库

```sql
CREATE DATABASE go_meteor_cms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

DSN 优先级：`MYSQL_DSN` 环境变量 > `api/cms/v1/etc/cms-api.yaml` 的 `Database.Source` > 内置默认值。

### 2.2 迁移（建表 + 种子，幂等）

```bash
go run cmd/migrate/main.go
# 期望输出：📊 权限统计：菜单 11 / 按钮 52 / API 52
```

> 种子含默认超管 `admin / 123456`，生产环境务必修改密码。

### 2.3 启动后端

```bash
go run ./api/cms/v1 -f api/cms/v1/etc/cms-api.yaml
```

### 2.4 启动前端

```bash
cd web
npm install
npm run dev          # http://localhost:5173
```

`web/vite.config.ts` 已将 `/admin`、`/api`、`/uploads` 代理到 `http://localhost:8888`，因此联调时后端须在 8888 运行。

### 2.5 验证

浏览器打开 `http://localhost:5173`，用 `admin / 123456` 登录（需输入图形验证码）。登录后应看到左侧动态菜单（控制台、管理员管理、部门管理、角色管理、权限管理、操作日志、登录日志、文件管理、系统管理【字典管理 / 系统配置】）。

---

## 3. 生产 / 容器部署

### 3.1 Docker Compose（推荐）

> 若需要分布式登录限流（多副本 API），compose 中需额外加入 `redis` 服务，并在 `cms-api.yaml` 的 `Redis.Host` 指向该服务。

```bash
docker compose up -d --wait        # MySQL(健康检查) → 迁移 → API
docker compose down
```

### 3.2 编译部署

```bash
go build -o bin/cms ./api/cms/v1
go build -o bin/migrate ./cmd/migrate

export MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/go_meteor_cms?parseTime=True&loc=Local"
./bin/migrate
./bin/cms -f api/cms/v1/etc/cms-api.yaml
```

### 3.3 环境变量

| 变量 | 说明 |
| --- | --- |
| `MYSQL_DSN` | 数据库连接串，覆盖配置文件 `Database.Source`，容器化部署推荐 |

### 3.4 后台维护

服务启动后自动清理过期 token 黑名单，按 `Maintenance.TokenBlacklistCleanInterval`（默认 3600s）周期执行。

### 3.5 Redis 与 CORS 配置（鉴权增强）

登录限流默认开启（`Login.RateLimit.Enabled=true`），采用**分布式固定窗口限流**：

- **Redis（多实例/容器部署必需）**：在 `cms-api.yaml` 配置 `Redis` 块：
  ```yaml
  Redis:
    Host: "127.0.0.1:6379"
    Type: "node"
    Pass: ""
    DB: 0
  ```
  - `Host` 为空时，登录限流自动降级为**单实例内存模式**（重启即清空、多实例不共享，仅适合单机开发）。
  - Redis 连接失败/抖动时同样降级为内存模式并打 error 日志，**防护不失效**（不会因 Redis 故障导致全员无法登录）。
  - 生产建议独立 Redis 实例或 Sentinel/Cluster（`Type` 可设 `cluster`），并配置 `Pass` 与网络隔离。
- **CORS（前端独立域名部署时必需）**：`cms-api.yaml` 的 `Cors` 列表配置允许的前端源：
  ```yaml
  Cors:
    - "http://localhost:5173"      # 本地 Vite 联调
    - "https://admin.example.com"
  ```
  - 通过 go-zero 内置 `rest.WithCors` 启用，预检 `OPTIONS` 由框架统一返回 204 并带 CORS 头。
  - 列表含 `"*"` 表示允许任意源（仅适合无 cookie 凭据、纯 token 鉴权场景）。
  - **未配置的源不会收到 `Access-Control-Allow-Origin`，浏览器将拒绝跨域请求**。
  - 同源部署（本服务用 `StaticDir` 直接托管 SPA、或反向代理同源）可不配 `Cors`。

---

## 4. 健康检查

| 检查项 | 方法 | 正常表现 |
| --- | --- | --- |
| 后端进程 | `netstat`/`ss` 看 `:8888` 监听；或 `docker ps` | 端口已监听 / 容器 Up |
| 接口探活 | 带 token `GET /admin/profile/permissions` | `code=0`，返回 `menus` 与 `permissions` |
| 数据库 | 迁移日志 `✅ 数据库迁移完成`、`✅ 数据初始化完成` | 无报错 |
| Redis | 多实例部署时执行 `redis-cli -h <host> ping` 应返回 `PONG`；限流 key 形如 `login_rl:<ip>` | `PONG`；限流生效 |
| 前端 | 浏览器打开页面，能登录并看到菜单 | 菜单正常渲染，无控制台报错 |

> `GET /admin/test-token` 不参与 API 权限校验，但需要有效 token 才能解析 claims，用于验证令牌本身是否有效。

---

## 5. 故障排查矩阵

| 症状 | 可能原因 | 排查 / 修复 |
| --- | --- | --- |
| 仪表盘空白 / 数据不显示 | 该接口响应缺少 `{code,msg,data}` 信封，前端拦截器（`res.code===0`）走错误分支 | 确认 handler 用 `common.Ok` 而非 `httpx.OkJsonCtx`（已修 `d8643a8`）。所有响应必须是统一信封 |
| 非超管访问某接口返回 **403** | 该路由没有对应的 `Type=3` API 权限记录（fail-closed，见 `PermissionMiddleware`） | 在 `internal/seed/permissions.go` 补一条 `Method`+`APIPath` 匹配的 API 权限，重新 `go run cmd/migrate/main.go` |
| 侧边栏某子菜单消失 / 整组丢失 | `buildMenuTree` 父子挂载顺序问题（已修 `e4998ad`）；或菜单 `sort` 与父子关系配置错误 | 确认菜单按 `sort` 升序且父菜单 `parent_id` 正确；复跑迁移 |
| 菜单图标异常（显示默认图标而非 X） | icon 语义键未命中前端 `MenuItem.vue` 的 `iconMap` | 在 `iconMap` 补充语义键；注意未命中只会回落为 `Menu` 图标，**不会**显示 X |
| 前端接口 404 / 网络错误 | Vite 代理未生效或后端未起 | 确认后端 `:8888` 在跑；确认 `web/vite.config.ts` 代理 `/admin`、`/api`、`/uploads` 存在 |
| 登录失败 | 验证码错误/过期、账号被禁用、密码错误 | 重新 `GET /api/captcha` 取新码；检查 `admin_users.is_active` 与密码 |
| 忘记密码邮件收不到 | `Mail.Enabled=false` 时只输出到服务日志，不真发 | 查看 API 服务日志拿重置令牌；或将 `Mail` 节点配置为启用 |
| 数据库迁移失败 | DSN 错误、MySQL 用户权限不足、表锁 | 检查 `MYSQL_DSN` / `cms-api.yaml` 的 `Database.Source`；确认 MySQL 用户有建库建表权限 |
| 文件上传失败 | `uploads` 目录无写权限、超出大小/类型限制 | 检查服务运行目录的 `uploads/` 写权限与文件配置项 |
| 非超管列表数据为空 | 数据权限 fail-closed（当前仅在管理员列表生效）：角色未分配 / 分配部门缺失 / 范围解析失败 → 不返回数据 | 检查用户角色 `data_scope` 与 `dept_id`，以及角色自定义部门（`admin_role_depts`） |
| 登录接口频繁返回 `code=2029`（请求过于频繁） | 同一 IP 在窗口内尝试次数超过 `Login.RateLimit.Max`（默认 10/60s）被限流；或 Redis 未配置导致各实例独立计数 | 等待 `BlockSeconds`(默认 300s) 后重试；多实例务必配置 `Redis` 使计数集中；确认 `Login.RateLimit.Enabled` |
| 多实例部署下登录限流形同虚设 | 未配置 `Redis`，各副本内存独立计数 | 在 `cms-api.yaml` 配置 `Redis.Host` 指向共享 Redis |
| 前端跨域报错（CORS / 预检失败） | 未配置 `Cors`，或前端源不在白名单 | 在 `cms-api.yaml` 的 `Cors` 加入前端实际源（含协议与端口） |
| 带 token 的跨域请求被拦（无 `Access-Control-Allow-Origin`） | 预检/实际请求的 `Origin` 不在 `Cors` 白名单 | 将前端源加入 `Cors`；确认前端实际使用的 `http://x:port` 与配置一致 |

---

## 6. 关键约定（改代码前必读）

1. **统一响应信封**：所有接口返回 `{code,msg,data}`，HTTP 状态码恒 200，业务结果看 `code`（`0` 成功，非 `0` 见 `internal/common/error_code.go`）。handler 统一用 `common.Ok` / `common.Fail`。
2. **分页**：`page` 缺省按 1；`page_size` 缺省取接口默认、上限 100（超出截断为 100）。
3. **权限三级**：菜单(1) / 按钮(2) / API(3)。按钮与 API **共用同一权限标识**（如 `system:role:add`），授权一次即可同时控制按钮显隐与接口可达。
4. **API 权限 `api_path` 必须与 `routes.go` 真实路由一致**（含 method），否则非超管必 403。新增接口务必同步补种子权限。
5. **路由匹配**：按"静态段命中最多者优先"，`/admin/departments/tree` 不会被 `/admin/departments/:id` 抢占。
6. **数据权限**：`data_scope` 1~5（全部/本部门/本部门及子部门/仅自己/自定义部门）；多角色取并集；任一角色为"全部"则不过滤。管理员列表与详情按部门范围收敛（fail-closed：范围为空→空结果 / 越权→按不存在处理）；操作/登录日志对非超管按 `user_id` 收敛为仅本人。角色/权限/部门/字典/文件等系统配置默认全局可见。
7. **超管**：`admin_users.is_super=true` 跳过全部鉴权与数据过滤，便于调试但不应作为常规账户。
8. **动态菜单**：侧边栏由 `GET /admin/profile/permissions` 返回的菜单树驱动（图标语义键 + `sort` 排序 + `parent_id` 层级），不写死在前端。
9. **`is_super` 以数据库为准**：业务层判断超管一律查 DB 的 `admin_users.is_super`，**不信任 token claim**（claim 在权限变更后、旧 token 过期前仍携带旧值）。涉及 `adminupdatelogic` / `resetadminpasswordlogic` / `operationloglistlogic` / `loginloglistlogic`，统一经 `logic.currentUserIsSuper` 取数；权限中间件同理（见 `PermissionMiddleware`）。
10. **密码复杂度**：新密码须 ≥8 位且至少含「大写/小写/数字/特殊字符」中三类（`utils.ValidatePasswordComplexity`）；创建/注册/改密/重置（含他人重置、找回密码）全入口统一校验。种子密码 `Admin123456` 满足该策略。
11. **全局踢人 / 会话批量吊销**：token 携带 `tv`（版本号），Redis `auth:tv:<userID>` 自增即吊销该用户全部旧令牌；触发点：禁用/删除账号、改密（含他人重置、找回密码）。未配置 Redis 时降级为不校验（与限流一致）。
12. **鉴权查库缓存**：权限中间件将「用户→角色→权限」集合与「API 路径→权限」解析按 60s TTL 内存缓存（`go-zero collection.Cache`）；`is_super`/`is_active` 仍每请求实时查库，保证权限变更即时生效。
13. **Redis 依赖**：登录限流、token 版本（全局踢人）依赖 `cms-api.yaml` 的 `Redis`；`Host` 为空则相应能力降级（限流→内存、踢人→不校验），多实例部署务必配置以保证防护一致生效。

---

## 7. 回滚与定位

- 提交历史：`e4998ad`（动态菜单多级层级 + 子菜单丢失修复）、`d8643a8`（统一响应信封修复）。
- 接口一致性审计方法见 `docs/smoke-test.md` 第 6 节，建议改动 `routes.go` 或种子权限后复跑。
