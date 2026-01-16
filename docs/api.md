# Go Meteor CMS API 文档

## 概述

Go Meteor CMS 是一个基于 Go 和 Go-Zero 开发的内容管理系统后端 API 服务。

## 基础URL

```
http://localhost:8888
```

## 状态码说明

| 状态码 | 说明         |
| ------ | ------------ |
| 0      | 请求成功     |
| 401    | 未授权       |
| 40000  | 参数错误     |
| 40001  | 用户名已存在 |
| 40002  | 邮箱已注册   |
| 50000  | 系统错误     |

## 公共接口

### 获取验证码

获取图形验证码，用于登录时验证。

**URL**: `GET /api/captcha`

**请求参数**: 无

**响应示例**:
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

### 管理员登录

管理员登录接口，需要提供用户名、密码和验证码。

**URL**: `POST /admin/login`

**请求参数**:
```json
{
  "username": "admin",
  "password": "123456",
  "captcha_id": "验证码ID",
  "captcha": "验证码"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 忘记密码

当管理员忘记密码时，可以通过此接口发送密码重置邮件。

**URL**: `POST /admin/forgot-password`

**请求参数**:
```json
{
  "username": "admin"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "重置密码邮件已发送，请检查邮箱"
  }
}
```

### 管理员注册

新管理员注册接口。

**URL**: `POST /admin/register`

**请求参数**:
```json
{
  "username": "newadmin",
  "password": "password123",
  "email": "newadmin@example.com",
  "nickname": "新管理员"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "管理员注册成功"
  }
}
```

## 需要认证的接口

以下接口需要在请求头中添加 Authorization:

```
Authorization: Bearer <token>
```

### Token 验证测试

测试 Token 是否有效。

**URL**: `GET /admin/test-token`

**请求参数**: 无

**响应示例**:
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

### 获取管理员列表

分页获取管理员列表，支持关键词搜索和状态过滤。

**URL**: `GET /admin/admin-users`

**查询参数**:
- `page`: 页码（必填）
- `page_size`: 每页数量（必填）
- `keyword`: 搜索关键词（用户名/昵称/邮箱，选填）
- `active`: 启用状态过滤（1-启用，2-禁用，选填）

**响应示例**:
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

### 修改自己的密码

管理员修改自己的密码。

**URL**: `PUT /admin/admin-users/me/change-password`

**请求参数**:
```json
{
  "old_password": "oldpassword",
  "new_password": "newpassword"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "密码修改成功"
  }
}
```

### 禁用管理员账号

超级管理员禁用其他管理员的账号。

**URL**: `PUT /admin/admin-users/:id/disable`

**路径参数**:
- `id`: 管理员ID

**请求体**: 无

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "禁用管理员成功"
  }
}
```

**错误码**:
- `3002`: 管理员不存在
- `1000`: 服务器内部错误

### 启用管理员账号

超级管理员启用已被禁用的管理员账号。

**URL**: `PUT /admin/admin-users/:id/enable`

**路径参数**:
- `id`: 管理员ID

**请求体**: 无

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "启用管理员成功"
  }
}
```

**错误码**:
- `3002`: 管理员不存在
- `1000`: 服务器内部错误

### 重置其他管理员的密码

超级管理员重置其他管理员的密码。

**URL**: `PUT /admin/admin-users/reset-password`

**请求参数**:
```json
{
  "id": 2,
  "new_password": "newpassword"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "密码已重置"
  }
}
```

### 管理员登出

管理员登出接口。

**URL**: `POST /admin/logout`

**请求参数**: 无

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "登出成功"
  }
}
```

### 管理员详情

获取指定管理员的详细信息。

**URL**: `GET /admin/admin-users/:id`

**路径参数**:
- `id`: 管理员ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "管理员",
    "email": "admin@example.com",
    "phone": "13800138000",
    "is_active": true
  }
}
```

### 创建管理员

创建新的管理员账号（需要超级管理员权限）。

**URL**: `POST /admin/admin-users/create`

**请求参数**:
```json
{
  "username": "newadmin",
  "password": "password123",
  "nickname": "新管理员",
  "email": "newadmin@example.com",
  "phone": "13800000000",
  "is_active": true,
  "is_super": false
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "创建管理员成功"
  }
}
```

### 更新管理员

更新指定管理员的信息（需要超级管理员权限）。

**URL**: `PUT /admin/admin-users/:id`

**路径参数**:
- `id`: 管理员ID

**请求参数**:
```json
{
  "nickname": "更新的昵称",
  "email": "updated@example.com",
  "phone": "13800000000",
  "is_active": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "更新管理员成功"
  }
}
```

### 删除管理员

删除指定管理员账号（需要超级管理员权限）。

**URL**: `DELETE /admin/admin-users/:id`

**路径参数**:
- `id`: 管理员ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "删除管理员成功"
  }
}
```

### 部门管理

#### 创建部门

创建新的部门（需要超级管理员权限）。

**URL**: `POST /admin/departments`

**请求参数**:
```json
{
  "name": "技术部",
  "parent_id": 0,
  "sort": 1,
  "is_active": true,
  "leader_id": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "message": "创建部门成功"
  }
}
```

#### 更新部门

更新指定部门信息（需要超级管理员权限）。

**URL**: `PUT /admin/departments/:id`

**路径参数**:
- `id`: 部门ID

**请求参数**:
```json
{
  "name": "更新的部门名称",
  "parent_id": 0,
  "sort": 2,
  "is_active": true,
  "leader_id": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "message": "更新部门成功"
  }
}
```

#### 删除部门

删除指定部门（需要超级管理员权限）。

**URL**: `DELETE /admin/departments/:id`

**路径参数**:
- `id`: 部门ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "message": "删除部门成功"
  }
}
```

#### 部门详情

获取指定部门的详细信息。

**URL**: `GET /admin/departments/:id`

**路径参数**:
- `id`: 部门ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "技术部",
    "parent_id": 0,
    "sort": 1,
    "is_active": true,
    "leader_id": 1
  }
}
```

#### 部门树

获取部门树结构。

**URL**: `GET /admin/departments/tree`

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "技术部",
      "sort": 1,
      "is_active": true,
      "children": [
        {
          "id": 2,
          "name": "前端组",
          "sort": 1,
          "is_active": true,
          "children": []
        }
      ]
    }
  ]
}
```

### 管理员绑定部门

将管理员绑定到指定部门（需要超级管理员权限）。

**URL**: `PUT /admin/admin-users/:id/bind-department`

**路径参数**:
- `id`: 管理员ID

**请求参数**:
```json
{
  "dept_id": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "message": "绑定部门成功"
  }
}
```

### 角色管理

#### 角色列表

获取角色列表（需要超级管理员权限）。

**URL**: `GET /admin/roles`

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `keyword`: 搜索关键词（角色名称/编码，选填）
- `is_active`: 是否启用（选填）

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 10,
    "page": 1,
    "page_size": 20,
    "list": [
      {
        "id": 1,
        "name": "管理员",
        "code": "ADMIN",
        "desc": "系统管理员",
        "data_scope": 1,
        "is_system": true,
        "is_active": true,
        "sort": 0,
        "created_at": "2023-01-01T12:00:00Z",
        "updated_at": "2023-01-01T12:00:00Z"
      }
    ]
  }
}
```

#### 角色详情

获取指定角色的详细信息（需要超级管理员权限）。

**URL**: `GET /admin/roles/:id`

**路径参数**:
- `id`: 角色ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "管理员",
    "code": "ADMIN",
    "desc": "系统管理员",
    "data_scope": 1,
    "is_system": true,
    "is_active": true,
    "sort": 0,
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  }
}
```

#### 创建角色

创建新的角色（需要超级管理员权限）。

**URL**: `POST /admin/roles`

**请求参数**:
```json
{
  "name": "普通用户",
  "code": "USER",
  "desc": "普通用户角色",
  "data_scope": 1,
  "is_active": true,
  "sort": 0
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "创建角色成功"
  }
}
```

#### 更新角色

更新指定角色的信息（需要超级管理员权限）。

**URL**: `PUT /admin/roles/:id`

**路径参数**:
- `id`: 角色ID

**请求参数**:
```json
{
  "name": "更新的角色名称",
  "code": "UPDATED_CODE",
  "desc": "更新的角色描述",
  "data_scope": 2,
  "is_active": true,
  "sort": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "更新角色成功"
  }
}
```

#### 删除角色

删除指定角色（需要超级管理员权限）。

**URL**: `DELETE /admin/roles/:id`

**路径参数**:
- `id`: 角色ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 2,
    "message": "删除角色成功"
  }
}
```

#### 获取角色权限ID列表

获取指定角色拥有的权限ID列表（需要超级管理员权限）。

**URL**: `GET /admin/roles/:id/permissions`

**路径参数**:
- `id`: 角色ID

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "permission_ids": [1, 2, 3]
  }
}
```

#### 分配角色权限

为指定角色分配权限（需要超级管理员权限）。

**URL**: `PUT /admin/roles/:id/permissions`

**路径参数**:
- `id`: 角色ID

**请求参数**:
```json
{
  "permission_ids": [1, 2, 3]
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "message": "分配权限成功"
  }
}
```

### 修改自己的密码

管理员修改自己的密码。

**URL**: `PUT /admin/admin-users/me/change-password`

**请求参数**:
```json
{
  "old_password": "oldpassword",
  "new_password": "newpassword"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "密码修改成功"
  }
}
```

### 重置其他管理员的密码

超级管理员重置其他管理员的密码。

**URL**: `PUT /admin/admin-users/reset-password`

**请求参数**:
```json
{
  "id": 2,
  "new_password": "newpassword"
}
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message": "密码已重置"
  }
}
```