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

**URL**: `PUT /admin/change-password`

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

**URL**: `PUT /admin/users/:id/reset-password`

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