# Go Meteor CMS 接口文档

服务默认监听 `http://localhost:8888`，本文按「模块 + 方法 + 路径」列出全部接口。

## 通用约定

### 响应格式

```json
{ "code": 0, "msg": "success", "data": {} }
```

- `code = 0` 成功
- `code > 0` 业务错误，`msg` 为可直接展示给用户的提示
- 完整错误码见 [`internal/common/error_code.go`](../internal/common/error_code.go)

### 认证

受保护接口需携带：

```
Authorization: Bearer <token>
```

未携带返回 `2001 请先登录`；token 无效或已登出返回 `2009 无效的token`；会话上下文缺失返回 `2004 用户会话已过期，请重新登录`。

### 鉴权

除以下接口外，其余受保护接口均由权限中间件按 `method + api_path` 匹配 `Type=3` 权限记录后校验：

- `POST /admin/logout`
- `GET /admin/test-token`
- `PUT /admin/admin-users/me/change-password`
- `GET /admin/profile`
- `PUT /admin/profile`
- `GET /admin/profile/permissions`

超级管理员跳过全部鉴权与数据范围过滤。

### 分页

`page` 小于 1 按 1 处理；`page_size` 小于 1 取接口默认值，**大于 100 一律截断为 100**，响应中的 `page_size` 为实际生效值。

---

## 公共接口（无需认证）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/captcha` | 获取图形验证码，返回 `captcha_id` 与 base64 图片 |
| POST | `/admin/login` | 登录，需提交验证码 |
| POST | `/admin/register` | 注册管理员 |
| POST | `/admin/forgot-password` | 申请密码重置令牌（账号不存在也返回成功，防账号枚举） |
| POST | `/admin/reset-password` | 用令牌设置新密码，令牌一次性、1 小时有效 |
| POST | `/admin/refresh` | 用 refresh_token 换取新的 access_token |

### 登录

```json
// 请求
{ "username": "admin", "password": "123456", "captcha_id": "xxx", "captcha": "1234" }
// 响应
{
  "code": 0, "msg": "success",
  "data": {
    "access_token": "eyJhbGci...（短期，默认 24h）",
    "refresh_token": "eyJhbGci...（长期，默认 7d，用于刷新）"
  }
}
```

失败场景：`4001` 参数缺失、`2006` 验证码错误、`2010` 用户名或密码错误、`2007` 账号未激活、`2011` 账号已被临时锁定。

> **登录失败锁定**：连续登录失败次数达到 `Login.MaxFailCount`（默认 5）后，账号被锁定 `Login.LockMinutes`（默认 15）分钟；锁定期间登录返回 `2011` 并提示剩余锁定分钟数。登录成功后自动清零失败计数。该机制仅针对已存在的账号，账号不存在时不计数（避免被用于探测账号是否存在）。
>
> **令牌刷新**：access_token 过期后，用登录时拿到的 refresh_token 调用 `POST /admin/refresh` 换取新的 access_token，无需重新登录。refresh_token 与 access_token 使用不同密钥签发，且携带 `token_type=refresh` 标识，不可互换。

### 刷新令牌

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/admin/refresh` | 用 refresh_token 换取新的 access_token |

```json
// 请求
{ "refresh_token": "eyJhbGci..." }
// 响应
{ "code": 0, "msg": "success", "data": { "access_token": "eyJhbGci..." } }
```

> 刷新令牌使用独立密钥（`Auth.RefreshSecret`）与更长有效期（`Auth.RefreshExpire`，默认 7 天）。用户被删除或禁用时刷新立即失效。刷新接口本身**不要求** access_token 鉴权。

> **刷新令牌轮换与并发宽限期**：每次刷新会签发新的 refresh token 并将旧令牌写入黑名单，防止长期复用。为避免并发请求同时用同一个旧令牌刷新时被误判为失效，服务端保留 `Auth.RefreshGraceSeconds`（默认 5 秒）宽限期：刚被轮换的旧令牌在此窗口内仍被接受，超出则判定为盗用重放并拒绝。
> 注意：登出同样会把 refresh token 写入黑名单，因此在登出后的宽限期内理论上存在极短的有效窗口（客户端已即时清除令牌，风险极低）。如需零宽限，可将 `RefreshGraceSeconds` 设为 0 并配合前端的刷新串行化。

> **前端无感刷新**：access_token 过期（后端返回 `code=2001`）时，前端自动用 refresh_token 调用 `/admin/refresh` 换新 token 并重试原请求（并发请求排队共享一次刷新）；刷新接口自身失败或无 refresh_token 时才跳转登录页。前端登出（`POST /admin/logout`）会在请求头携带 `X-Refresh-Token`，后端将其与 access_token 一并加入黑名单（TTL = `Auth.RefreshExpire`），此后该 refresh_token 无法再换发 access_token，实现真正的"登出即失效"。

### 图形验证码

```json
{ "code": 0, "msg": "success", "data": { "captcha_id": "xxx", "image": "data:image/png;base64,..." } }
```

### 重置密码

```json
// 请求
{ "token": "邮件中的令牌", "new_password": "newpassword123" }
// 失败码：3012 令牌无效 / 3013 令牌已过期 / 4001 新密码少于 6 位
```

---

## 个人中心

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/profile` | 免鉴权 | 当前登录管理员的个人信息 |
| PUT | `/admin/profile` | 免鉴权 | 更新昵称/邮箱/手机，未传字段保持不变 |
| GET | `/admin/profile/permissions` | 免鉴权 | 当前用户的菜单树与权限码 |
| PUT | `/admin/admin-users/me/change-password` | 免鉴权 | 修改自己的密码 |
| POST | `/admin/logout` | 免鉴权 | 登出，token 进入黑名单 |
| GET | `/admin/test-token` | 免鉴权 | 校验 token 是否有效 |

### GET `/admin/profile`

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1, "username": "admin", "nickname": "超级管理员",
    "email": "admin@example.com", "phone": "", "avatar": "",
    "is_super": true, "dept_id": 0,
    "last_login_at": "2026-10-30T10:00:00+08:00",
    "created_at": "2026-10-01T10:00:00+08:00"
  }
}
```

### PUT `/admin/profile`

```json
{ "nickname": "新昵称", "email": "new@example.com", "phone": "13800138000" }
```

邮箱已被他人使用时返回 `3008 邮箱已被注册`。

### GET `/admin/profile/permissions`

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "is_super": false,
    "menus": [
      { "id": 2, "name": "管理员管理", "path": "/admin/admin-users",
        "component": "views/admin/user/list.vue", "icon": "user",
        "permission": "system:admin_user:list", "children": [] }
    ],
    "permissions": ["system:admin_user:list", "system:admin_user:add"]
  }
}
```

---

## 管理员管理

列表接口按当前用户的数据权限收敛结果，详见 [数据权限说明](#数据权限)。

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/admin-users` | `system:admin_user:list` | 管理员列表（分页、按用户名/手机/状态筛选） |
| GET | `/admin/admin-users/:id` | `system:admin_user:detail` | 管理员详情 |
| POST | `/admin/admin-users/create` | `system:admin_user:add` | 新增管理员 |
| PUT | `/admin/admin-users/:id` | `system:admin_user:edit` | 更新资料（局部更新，未传字段不变） |
| DELETE | `/admin/admin-users/:id` | `system:admin_user:delete` | 删除管理员（软删除），不可删自己与超管 |
| PUT | `/admin/admin-users/:id/disable` | `system:admin_user:disable` | 禁用，不可禁用自己与超管 |
| PUT | `/admin/admin-users/:id/enable` | `system:admin_user:enable` | 启用 |
| PUT | `/admin/admin-users/:id/bind-department` | `system:admin_user:bind_dept` | 绑定部门 |
| PUT | `/admin/admin-users/reset-password` | `system:admin_user:reset_pwd` | 重置他人密码，仅超管，目标不能是超管 |
| GET | `/admin/admin-users/:id/roles` | `system:admin_user:role_list` | 查询已分配角色 |
| PUT | `/admin/admin-users/:id/roles` | `system:admin_user:role_assign` | 分配角色（覆盖式，事务保证） |
| DELETE | `/admin/admin-users/:id/roles/:role_id` | `system:admin_user:role_remove` | 移除单个角色 |

### 新增管理员

```json
{
  "username": "editor", "password": "123456", "nickname": "编辑",
  "email": "editor@example.com", "phone": "13800138000",
  "is_active": true, "is_super": false, "dept_id": 1
}
```

用户名重复返回 `3001 管理员已存在`，邮箱重复返回 `3008 邮箱已被注册`。

### 分配角色

```json
{ "role_ids": [1, 2] }
```

角色 ID 去重后覆盖写入，角色不存在返回 `5201 角色不存在`。

---

## 部门管理

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/departments/tree` | `system:dept:list` | 部门树（仅返回启用部门） |
| GET | `/admin/departments/:id` | `system:dept:detail` | 部门详情 |
| POST | `/admin/departments` | `system:dept:add` | 新增部门 |
| PUT | `/admin/departments/:id` | `system:dept:edit` | 更新部门 |
| DELETE | `/admin/departments/:id` | `system:dept:delete` | 删除部门，存在子部门或仍有关联管理员时拒绝 |

部门树节点含 `children`，父部门未启用时其下的启用子部门不会出现在树中。

---

## 角色管理

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/roles` | `system:role:list` | 角色列表（分页、按名称/编码/状态筛选） |
| GET | `/admin/roles/:id` | `system:role:detail` | 角色详情 |
| POST | `/admin/roles` | `system:role:add` | 新增角色，含 `data_scope` 数据范围 |
| PUT | `/admin/roles/:id` | `system:role:edit` | 更新角色，系统内置角色不可改 |
| DELETE | `/admin/roles/:id` | `system:role:delete` | 删除角色，仍有关联用户时拒绝 |
| GET | `/admin/roles/:id/permissions` | `system:role:perm_list` | 角色已分配的权限 ID 列表 |
| GET | `/admin/roles/:id/permissions-with-names` | `system:role:perm_list_names` | 同上，附带权限名称 |
| PUT | `/admin/roles/:id/permissions` | `system:role:perm_assign` | 分配权限（覆盖式，事务保证） |
| GET | `/admin/roles/:id/departments` | `system:role:dept_list` | 自定义数据范围的部门列表 |
| PUT | `/admin/roles/:id/departments` | `system:role:dept_assign` | 设置部门范围，仅 `data_scope=5` 可用 |

### 分配权限

```json
{ "permission_ids": [1, 2, 3] }
```

系统内置角色返回 `5208 系统内置角色，禁止修改`；权限不存在返回 `5401 权限不存在`。

### 设置部门范围

```json
{ "dept_ids": [1, 2] }
```

非 `data_scope=5` 的角色返回 `4001 参数错误`；部门不存在或未启用返回 `5001 部门不存在`。

---

## 权限管理

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/permissions` | `system:permission:list` | 权限列表（分页，可按类型/名称筛选） |
| GET | `/admin/permissions/tree` | `system:permission:tree` | 权限树 |
| GET | `/admin/permissions/:id` | `system:permission:detail` | 权限详情 |
| POST | `/admin/permissions` | `system:permission:add` | 新增权限 |
| PUT | `/admin/permissions/:id` | `system:permission:edit` | 更新权限 |
| DELETE | `/admin/permissions/:id` | `system:permission:delete` | 删除权限，存在子权限时拒绝 |

### 权限类型

| type | 含义 | 用途 |
| --- | --- | --- |
| 1 | 菜单 | 前端路由，`/admin/profile/permissions` 据此下发菜单树 |
| 2 | 按钮 | 前端功能点，控制按钮显隐 |
| 3 | API | 后端接口鉴权，按 `method + api_path` 匹配 |

新增受保护接口后，请同步在 `internal/seed/permissions.go` 追加一条 `Type=3` 权限，并执行 `go run ./cmd/migrate` 写入，否则非超管访问会返回 `2002 权限不足`。

---

## 数据权限

角色通过 `data_scope` 声明可见数据范围：

| data_scope | 含义 |
| --- | --- |
| 1 | 全部数据 |
| 2 | 本部门 |
| 3 | 本部门及子部门 |
| 4 | 仅自己 |
| 5 | 自定义部门（读取 `admin_role_depts`） |

- 多角色取并集，任一角色为「全部数据」时不过滤
- 范围解析失败、未分配启用角色、或要求按部门过滤但未分配部门时，按最严格处理（返回空集）
- 当前仅在管理员列表接口落地

---

## 审计日志

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/logs/operations` | `system:log:operation_list` | 操作日志列表 |
| GET | `/admin/logs/logins` | `system:log:login_list` | 登录日志列表 |

### GET `/admin/logs/operations`

查询后台接口调用记录，用于安全审计与问题排查。支持按操作人、模块、路径、结果、时间范围筛选。

**请求参数（query）**：

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| page | int | 页码，默认 1 |
| page_size | int | 每页条数，默认 20，上限 100 |
| user_id | int64 | 操作人 ID（仅超管可指定；非超管强制只看自己） |
| username | string | 操作人账号（模糊匹配） |
| module | string | 业务模块，如「管理员管理」 |
| path | string | 请求路径（模糊匹配） |
| status | int | 0 成功 / 1 失败 |
| start_time | string | 起始时间，支持 `2006-01-02 15:04:05` / `2006-01-02` / RFC3339 |
| end_time | string | 结束时间 |

**响应 data**：

| 字段 | 说明 |
| --- | --- |
| total | 总条数 |
| list[].module | 业务模块 |
| list[].action | 操作描述（取自权限树，如「新增管理员」） |
| list[].status | 0 成功 / 1 失败 |
| list[].params | 请求参数快照（密码等敏感字段已脱敏、超长截断） |
| list[].duration_ms | 耗时（毫秒） |

> 数据权限：非超管只能查询自己的操作日志，外部传入的 `user_id` 会被忽略。按 `created_at` 倒序返回。
> 日志保留策略由 `Maintenance.AuditLogRetentionDays` 配置（默认 90 天，<=0 表示永久保留），由后台任务定期清理。

### GET `/admin/logs/logins`

查询登录成功与失败记录，用于账号安全审计。参数与操作日志类似（`user_id`/`username`/`status`/`start_time`/`end_time`），按 `login_at` 倒序。

**响应 data.list[]** 字段：`id, user_id, username, ip, user_agent, status, fail_reason, login_at`。

---

## 文件上传

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| POST | `/admin/upload` | `system:upload:file` | 上传文件（multipart/form-data，字段名 `file`） |

### POST `/admin/upload`

接收单个文件，落盘到 `public/uploads/YYYY/MM/<uuid>.<ext>`，返回可直接访问的 URL（前缀 `/uploads/`）。静态资源通过 `rest.WithFileServer` 对外暴露，无需鉴权即可 GET。

**校验规则**：
- 扩展名白名单：`.jpg .jpeg .png .gif .webp`（配置于 `etc/cms-api.yaml` 的 `Upload.AllowedExts`）
- 文件头魔数校验：真实类型必须与扩展名一致，伪造后缀会被拒绝（返回 `4001`）
- 大小上限：默认 5MB（`Upload.MaxSizeMB`），超限返回 `4001`
- 文件名统一用 UUID 重命名，杜绝路径穿越与同名覆盖

**请求**：`multipart/form-data`，字段名 `file`

**响应 data**：

| 字段 | 说明 |
| --- | --- |
| url | 可访问的相对路径，如 `/uploads/2026/10/xxxx.png` |
| filename | 原始文件名 |
| size | 文件大小（字节） |

> 回填头像：把返回的 `url` 通过 `PUT /admin/profile` 的 `avatar` 字段更新即可；更新头像时会自动清理旧头像文件（默认头像除外）。
> 上传目录（`public/uploads/*`，默认头像除外）已在 `.gitignore` 忽略，避免运行时文件污染仓库。

---

## 文件管理

| 方法 | 路径 | 权限码 | 说明 |
| --- | --- | --- | --- |
| GET | `/admin/files` | `system:file:list` | 列出上传目录下的文件（支持目录、关键字、分页） |
| POST | `/admin/files/delete` | `system:file:delete` | 删除指定文件（空目录也可删） |
| PUT | `/admin/files/rename` | `system:file:rename` | 同级目录内重命名 |

> 以上接口需登录（JWT 鉴权）并具备对应权限码，超级管理员天然放行；权限清单由迁移时的种子数据自动写入（`internal/seed/permissions.go`）。路径均相对上传根目录，服务端做防穿越校验（`../` 会被拒绝）。

### GET `/admin/files`

**请求参数（query）**：

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| dir | string | 子目录，相对上传根；为空表示根目录 |
| keyword | string | 按文件名模糊匹配 |
| page | int | 页码，默认 1 |
| page_size | int | 每页条数，默认 20，上限 100 |

**响应 data**：`list[]`（字段 `name / path / url / size / modified / is_dir`）、`total`、`current_dir`。

### POST `/admin/files/delete`

请求体：`{ "path": "相对上传根目录的路径" }`。目录仅允许在为空时删除。

### PUT `/admin/files/rename`

请求体：`{ "path": "...", "new_name": "新名称" }`。`new_name` 不含路径分隔符，且目标名称不得已存在。

---

## 仪表盘统计

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/admin/dashboard/stats` | 后台首页统计指标，需登录（JWT 鉴权） |

### GET `/admin/dashboard/stats`

聚合常用运营指标，单项统计失败不影响其余字段返回（失败项记日志、值为 0）。

**响应 data**：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| admin_count | int64 | 管理员总数（admin_users） |
| dept_count | int64 | 部门总数（departments） |
| role_count | int64 | 角色总数（admin_roles） |
| permission_count | int64 | 权限总数（admin_permissions） |
| today_login_count | int64 | 今日登录成功次数（login_logs.status=0 且 login_at 为今日） |
| today_fail_login | int64 | 今日登录失败次数（login_logs.status=1 且 login_at 为今日） |
| week_op_count | int64 | 近 7 天操作日志数（operation_logs.created_at 在 7 天内） |

**示例**：

```json
{
  "code": 0, "msg": "success",
  "data": {
    "admin_count": 3, "dept_count": 5, "role_count": 4, "permission_count": 42,
    "today_login_count": 12, "today_fail_login": 1, "week_op_count": 318
  }
}
```

