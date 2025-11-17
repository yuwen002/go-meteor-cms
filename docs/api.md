# Go Meteor CMS API 接口文档

## 概述

本文档详细描述了 Go Meteor CMS 系统的所有 API 接口，包括请求参数、响应格式和使用示例。

## 基础信息

- **服务地址**: `http://localhost:8888`
- **认证方式**: JWT Token
- **数据格式**: JSON

## 公共响应格式

所有接口返回统一的 JSON 格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

### 响应码说明

| 状态码 | 说明 |
|-------|------|
| 0 | 请求成功 |
| 40000 | 参数错误 |
| 401 | 未授权/Token无效 |
| 50000 | 系统错误 |

## 接口列表

### 1. 管理员接口

#### 1.1 管理员登录

**接口地址**: `POST /admin/login`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:

```json
{
  "username": "admin",
  "password": "123456"
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

#### 1.2 忘记密码

**接口地址**: `POST /admin/forgot-password`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| username | string | 是 | 用户名 |

**请求示例**:

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
    "status": 1,
    "message": "密码重置邮件已发送"
  }
}
```

#### 1.3 管理员注册

**接口地址**: `POST /admin/register`

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码（最少6位） |
| email | string | 是 | 邮箱 |
| nickname | string | 否 | 昵称 |

**请求示例**:

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

### 2. 需要认证的接口

以下接口需要在请求头中添加 Authorization:

```
Authorization: Bearer <token>
```

#### 2.1 Token 验证测试

**接口地址**: `GET /admin/test-token`

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

#### 2.2 获取管理员列表

**接口地址**: `GET /admin/users`

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 |
|-------|------|------|------|
| page | int | 是 | 页码 |
| page_size | int | 是 | 每页数量 |
| keyword | string | 否 | 搜索关键词（用户名/昵称/邮箱） |
| active | int | 否 | 启用状态过滤（1-启用，2-禁用） |

**请求示例**:

```
GET /admin/users?page=1&page_size=10&keyword=admin&active=1
```

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

**管理员列表项字段说明**:

| 字段名 | 类型 | 说明 |
|-------|------|------|
| id | int64 | 用户ID |
| username | string | 用户名 |
| nickname | string | 昵称 |
| email | string | 邮箱 |
| phone | string | 手机号 |
| avatar | string | 头像URL |
| isSuper | bool | 是否超级管理员 |
| isActive | bool | 是否启用 |
| last_login_at | string | 最后登录时间 |
| createdAt | string | 创建时间 |

## 错误码说明

| 错误码 | 说明 |
|-------|------|
| 40000 | 参数错误 |
| 40001 | 用户名已存在 |
| 40002 | 邮箱已注册 |
| 401 | 请先登录/Token无效 |
| 50001 | 检查用户失败 |
| 50002 | 检查邮箱失败 |
| 50003 | 密码加密失败 |
| 50004 | 用户创建失败 |

## 数据模型

### 管理员用户模型

| 字段名 | 类型 | 说明 |
|-------|------|------|
| id | int64 | 自增主键ID |
| username | string | 用户名（唯一） |
| password_hash | string | 密码哈希值 |
| nickname | string | 昵称 |
| email | string | 邮箱（唯一） |
| phone | string | 手机号 |
| avatar | string | 头像URL |
| is_super | bool | 是否超级管理员 |
| is_active | bool | 是否启用 |
| last_login_at | time | 最后登录时间 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |
| reset_token | string | 密码重置令牌 |
| reset_expire | time | 密码重置令牌过期时间 |