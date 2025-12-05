# Bug Bounty Lite API 文档

## 目录

- [概述](#概述)
- [基础信息](#基础信息)
- [认证说明](#认证说明)
- [API 端点](#api-端点)
  - [认证相关](#认证相关)
  - [漏洞报告相关](#漏洞报告相关)
  - [用户信息变更](#用户信息变更)
  - [项目管理](#项目管理)
  - [系统配置](#系统配置)
  - [文件上传](#文件上传)
- [数据模型](#数据模型)
- [错误处理](#错误处理)
- [快速开始](#快速开始)
  - [curl 命令示例](#curl-命令示例)
  - [代码示例](#代码示例)

---

## 概述

Bug Bounty Lite 是一个轻量级的 Web 安全众测平台后端 API，提供用户认证、漏洞报告管理和用户信息变更等功能。

**API 版本**: v1  
**Base URL**: `http://localhost:8080/api/v1`  
**协议**: HTTP/HTTPS  
**数据格式**: JSON  
**字符编码**: UTF-8

---

## 基础信息

### 请求规范

#### 通用请求头

| 请求头 | 值 | 必填 | 说明 |
|--------|-----|------|------|
| Content-Type | `application/json` | 是（POST/PUT） | 请求体格式 |
| Accept | `application/json` | 否 | 期望的响应格式 |
| Authorization | `Bearer <token>` | 视接口而定 | JWT 认证令牌 |

#### 请求体格式

- 使用 JSON 格式
- 字段名使用 `snake_case`（下划线命名）
- 字符串值使用双引号
- 布尔值使用 `true`/`false`
- 空值使用 `null`

#### 响应格式

**成功响应**:
```json
{
  "message": "操作成功",
  "data": { ... }
}
```

**错误响应**:
```json
{
  "error": "错误描述信息"
}
```

---

## 认证说明

### JWT Token 认证

登录成功后会返回 JWT Token，访问需要认证的接口时需要在请求头中携带：

```
Authorization: Bearer <token>
```

### Token 结构

JWT Token 包含以下信息（Payload）：
```json
{
  "user_id": 1,
  "username": "test",
  "role": "whitehat",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Token 有效期

- 默认有效期：7200 秒（2小时）
- 过期后需要重新登录获取新 Token

### 接口权限一览

| 接口 | 方法 | 认证 | 说明 |
|------|------|------|------|
| `/api/v1/auth/register` | POST | 否 | 用户注册 |
| `/api/v1/auth/login` | POST | 否 | 用户登录 |
| `/api/v1/reports` | POST | 是 | 提交报告 |
| `/api/v1/reports` | GET | 是 | 获取报告列表 |
| `/api/v1/reports/:id` | GET | 是 | 获取报告详情 |
| `/api/v1/reports/:id` | PUT | 是 | 更新报告 |
| `/api/v1/user/info/change` | POST | 是 | 提交信息变更申请 |
| `/api/v1/user/info/changes` | GET | 是 | 获取变更申请列表 |
| `/api/v1/user/info/changes/:id` | GET | 是 | 获取变更申请详情 |
| `/api/v1/projects` | POST | 是 | 创建项目（仅admin） |
| `/api/v1/projects` | GET | 是 | 获取项目列表 |
| `/api/v1/projects/:id` | GET | 是 | 获取项目详情 |
| `/api/v1/projects/:id` | PUT | 是 | 更新项目（仅admin） |
| `/api/v1/projects/:id` | DELETE | 是 | 删除项目（仅admin） |
| `/api/v1/configs/:type` | GET | 是 | 获取配置列表 |
| `/api/v1/configs/:type/:id` | GET | 是 | 获取配置详情 |
| `/api/v1/configs/:type` | POST | 是 | 创建配置（仅admin） |
| `/api/v1/configs/:type/:id` | PUT | 是 | 更新配置（仅admin） |
| `/api/v1/configs/:type/:id` | DELETE | 是 | 删除配置（仅admin） |
| `/api/v1/upload` | POST | 是 | 上传文件 |

---

## API 端点

### 认证相关

#### 1. 用户注册

创建新用户账号。

**接口**: `POST /api/v1/auth/register`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |

**请求体参数**:
| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| username | string | 是 | 1-64字符 | 用户名，唯一 |
| password | string | 是 | 最少6位 | 密码 |

**请求示例**:
```json
{
  "username": "whitehat_user",
  "password": "secure123"
}
```

**响应示例**:

成功 (201 Created):
```json
{
  "message": "User registered successfully"
}
```

失败 (400 Bad Request):
```json
{
  "error": "username already exists"
}
```

---

#### 2. 用户登录

用户登录获取 JWT Token。

**接口**: `POST /api/v1/auth/login`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |

**请求体参数**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**:
```json
{
  "username": "whitehat_user",
  "password": "secure123"
}
```

**响应示例**:

成功 (200 OK):
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "whitehat_user",
    "role": "whitehat",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

失败 (401 Unauthorized):
```json
{
  "error": "Invalid credentials"
}
```

> **重要**: 请保存返回的 `token`，后续请求需要在 Header 中携带。

---

### 漏洞报告相关

#### 3. 提交漏洞报告

提交新的漏洞报告。

**接口**: `POST /api/v1/reports`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**请求体参数**:
| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| project_id | integer | 是 | - | 项目ID（关联项目） |
| vulnerability_name | string | 是 | 最大255字符 | 漏洞名称 |
| vulnerability_type_id | integer | 是 | - | 漏洞类型配置ID（从系统配置获取） |
| vulnerability_impact | string | 否 | 无限制 | 漏洞的危害 |
| self_assessment | string | 否 | 无限制 | 危害自评 |
| vulnerability_url | string | 否 | URL格式 | 漏洞链接 |
| vulnerability_detail | string | 否 | 无限制 | 漏洞详情 |
| attachment_url | string | 否 | URL格式 | 附件地址（文件上传后的URL） |
| severity | string | 否 | 枚举值 | 危害等级，默认 `Low` |
| title | string | 否 | 最大255字符 | 漏洞标题（保留字段，向后兼容） |
| description | string | 否 | 无限制 | 漏洞描述（保留字段，向后兼容） |
| type | string | 否 | 最大50字符 | 漏洞类型（保留字段，向后兼容） |

**severity 可选值**: `Low`, `Medium`, `High`, `Critical`

**请求示例**:
```json
{
  "project_id": 1,
  "vulnerability_name": "SQL注入漏洞",
  "vulnerability_type_id": 1,
  "vulnerability_impact": "可能导致数据泄露",
  "self_assessment": "高危漏洞",
  "vulnerability_url": "https://example.com/vuln",
  "vulnerability_detail": "详细描述漏洞情况...",
  "attachment_url": "https://example.com/uploads/reports/2024/01/abc123.pdf",
  "severity": "High"
}
```

**响应示例**:

成功 (200 OK):
```json
{
  "code": 200,
  "message": "漏洞报告提交成功"
}
```

失败 (400 Bad Request):
```json
{
  "code": 400,
  "message": "项目ID不能为空"
}
```

或

```json
{
  "code": 400,
  "message": "漏洞名称不能为空"
}
```

或

```json
{
  "code": 400,
  "message": "漏洞类型不能为空"
}
```

或

```json
{
  "code": 400,
  "message": "请求参数错误: ..."
}
```

失败 (401 Unauthorized):
```json
{
  "code": 401,
  "message": "用户未认证"
}
```

---

#### 4. 获取报告列表

获取漏洞报告列表，支持分页。

**接口**: `GET /api/v1/reports`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**查询参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | integer | 否 | 1 | 页码（从1开始） |
| page_size | integer | 否 | 10 | 每页数量（最大100） |

**请求示例**:
```
GET /api/v1/reports?page=1&page_size=10
```

**响应示例**:

成功 (200 OK):
```json
{
  "data": [
    {
      "id": 2,
      "title": "XSS Vulnerability in Comment Section",
      "description": "The comment section allows...",
      "type": "XSS",
      "severity": "Medium",
      "status": "Pending",
      "author_id": 1,
      "author": {
        "id": 1,
        "username": "whitehat_user",
        "role": "whitehat",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      },
      "created_at": "2024-01-02T00:00:00Z",
      "updated_at": "2024-01-02T00:00:00Z"
    }
  ],
  "total": 2,
  "page": 1
}
```

---

#### 5. 获取报告详情

根据 ID 获取单个报告的详细信息。

**接口**: `GET /api/v1/reports/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 报告 ID |

**请求示例**:
```
GET /api/v1/reports/1
```

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "id": 1,
    "title": "SQL Injection in Login Form",
    "description": "The login form is vulnerable to SQL injection attacks...",
    "type": "SQL Injection",
    "severity": "High",
    "status": "Pending",
    "author_id": 1,
    "author": {
      "id": 1,
      "username": "whitehat_user",
      "role": "whitehat",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

失败 (404 Not Found):
```json
{
  "error": "Report not found"
}
```

---

#### 6. 更新报告

更新报告信息或状态。

**接口**: `PUT /api/v1/reports/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 报告 ID |

**请求体参数**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| project_id | integer | 否 | 项目ID |
| vulnerability_name | string | 否 | 漏洞名称 |
| vulnerability_type_id | integer | 否 | 漏洞类型配置ID |
| vulnerability_impact | string | 否 | 漏洞的危害 |
| self_assessment | string | 否 | 危害自评 |
| vulnerability_url | string | 否 | 漏洞链接（URL格式） |
| vulnerability_detail | string | 否 | 漏洞详情 |
| attachment_url | string | 否 | 附件地址（URL格式） |
| severity | string | 否 | 危害等级 |
| status | string | 否 | 状态（仅 admin/vendor） |
| title | string | 否 | 漏洞标题（保留字段） |
| description | string | 否 | 漏洞描述（保留字段） |
| type | string | 否 | 漏洞类型（保留字段） |

**severity 可选值**: `Low`, `Medium`, `High`, `Critical`

**status 可选值**: `Pending`, `Triaged`, `Resolved`, `Closed`

**请求示例**:
```json
{
  "status": "Triaged",
  "severity": "Critical"
}
```

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "id": 1,
    "title": "SQL Injection in Login Form",
    "description": "The login form is vulnerable to SQL injection attacks...",
    "type": "SQL Injection",
    "severity": "Critical",
    "status": "Triaged",
    "author_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

失败 (400 Bad Request):
```json
{
  "error": "permission denied"
}
```

**权限说明**:
- 报告作者可以更新 `title`、`description`、`type`、`severity`
- 只有 `admin` 或 `vendor` 角色可以更新 `status`
- 状态流转规则: `Pending` -> `Triaged` -> `Resolved` -> `Closed`

---

### 用户信息变更

#### 7. 提交信息变更申请

提交用户信息变更申请（手机号、邮箱、姓名），需要后台审核。

**接口**: `POST /api/v1/user/info/change`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**请求体参数**:
| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| phone | string | 否 | 最大20字符 | 手机号 |
| email | string | 否 | 邮箱格式 | 邮箱 |
| name | string | 否 | 最大50字符 | 姓名 |

> **注意**: 至少需要提供一个要变更的字段（phone、email 或 name）

**请求示例**:
```json
{
  "phone": "13800138000",
  "email": "newemail@example.com",
  "name": "张三"
}
```

**响应示例**:

成功 (201 Created):
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "phone": "13800138000",
    "email": "newemail@example.com",
    "name": "张三",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

失败 (400 Bad Request):
```json
{
  "error": "至少需要提供一个要变更的字段（手机号、邮箱或姓名）"
}
```

---

#### 8. 获取变更申请列表

获取当前用户的所有信息变更申请。

**接口**: `GET /api/v1/user/info/changes`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**响应示例**:

成功 (200 OK):
```json
{
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "phone": "13800138000",
      "email": "newemail@example.com",
      "name": "张三",
      "status": "pending",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

#### 9. 获取变更申请详情

根据 ID 获取单个变更申请的详细信息。

**接口**: `GET /api/v1/user/info/changes/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 变更申请 ID |

**请求示例**:
```
GET /api/v1/user/info/changes/1
```

**响应示例**:

成功 (200 OK):
```json
{
  "message": "获取成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "phone": "13800138000",
    "email": "newemail@example.com",
    "name": "张三",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

失败 (400 Bad Request):
```json
{
  "error": "变更申请不存在或无权限访问"
}
```

**状态说明**:
- `pending`: 待审核
- `approved`: 已通过（审核通过后，用户信息会被更新）
- `rejected`: 已拒绝

---

### 项目管理

#### 10. 创建项目

创建新项目（仅管理员）。

**接口**: `POST /api/v1/projects`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**请求体参数**:
| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| name | string | 是 | 最大255字符 | 项目名称 |
| description | string | 否 | 无限制 | 项目描述 |
| note | string | 否 | 无限制 | 备注 |

**请求示例**:
```json
{
  "name": "某公司官网",
  "description": "公司官方网站项目",
  "note": "重要项目，需要重点关注"
}
```

**响应示例**:

成功 (201 Created):
```json
{
  "data": {
    "id": 1,
    "name": "某公司官网",
    "description": "公司官方网站项目",
    "note": "重要项目，需要重点关注",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 11. 获取项目列表

获取项目列表，支持分页。

**接口**: `GET /api/v1/projects`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**查询参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | integer | 否 | 1 | 页码（从1开始） |
| page_size | integer | 否 | 10 | 每页数量（最大100） |

**权限说明**:
- `whitehat`/`vendor`: 只能查看 `status='active'` 的项目
- `admin`: 可以查看所有项目（包括 `inactive`）

**响应示例**:

成功 (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "name": "某公司官网",
      "description": "公司官方网站项目",
      "note": "重要项目",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1
}
```

---

#### 12. 获取项目详情

根据 ID 获取单个项目的详细信息。

**接口**: `GET /api/v1/projects/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 项目 ID |

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "id": 1,
    "name": "某公司官网",
    "description": "公司官方网站项目",
    "note": "重要项目，需要重点关注",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 13. 更新项目

更新项目信息（仅管理员）。

**接口**: `PUT /api/v1/projects/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 项目 ID |

**请求体参数**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 项目名称 |
| description | string | 否 | 项目描述 |
| note | string | 否 | 备注 |
| status | string | 否 | 项目状态（active/inactive） |

**请求示例**:
```json
{
  "name": "某公司官网（更新）",
  "status": "inactive"
}
```

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "id": 1,
    "name": "某公司官网（更新）",
    "description": "公司官方网站项目",
    "note": "重要项目",
    "status": "inactive",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

#### 14. 删除项目

删除项目（仅管理员）。

**接口**: `DELETE /api/v1/projects/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 项目 ID |

**响应示例**:

成功 (200 OK):
```json
{
  "message": "项目删除成功"
}
```

---

### 系统配置

#### 15. 获取配置列表

根据配置类型获取配置列表。

**接口**: `GET /api/v1/configs/:type`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 配置类型（如：vulnerability_type） |

**权限说明**:
- `whitehat`/`vendor`: 只能查看 `status='active'` 的配置
- `admin`: 可以查看所有配置（包括 `inactive`）

**请求示例**:
```
GET /api/v1/configs/vulnerability_type
```

**响应示例**:

成功 (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "config_type": "vulnerability_type",
      "config_key": "SQL_INJECTION",
      "config_value": "SQL注入",
      "description": "SQL注入漏洞",
      "sort_order": 1,
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "config_type": "vulnerability_type",
      "config_key": "XSS",
      "config_value": "XSS跨站脚本",
      "description": "跨站脚本攻击",
      "sort_order": 2,
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

#### 16. 获取配置详情

根据 ID 获取单个配置的详细信息。

**接口**: `GET /api/v1/configs/:type/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 配置类型 |
| id | integer | 是 | 配置 ID |

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "id": 1,
    "config_type": "vulnerability_type",
    "config_key": "SQL_INJECTION",
    "config_value": "SQL注入",
    "description": "SQL注入漏洞",
    "sort_order": 1,
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 17. 创建配置

创建新的配置项（仅管理员）。

**接口**: `POST /api/v1/configs/:type`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 配置类型 |

**请求体参数**:
| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| config_key | string | 是 | 最大100字符 | 配置键（如：SQL_INJECTION） |
| config_value | string | 是 | 最大255字符 | 配置值（显示名称） |
| description | string | 否 | 无限制 | 配置描述 |
| sort_order | integer | 否 | - | 排序顺序 |
| status | string | 否 | 枚举值 | 状态，默认 `active` |

**请求示例**:
```json
{
  "config_key": "NEW_VULN_TYPE",
  "config_value": "新漏洞类型",
  "description": "新发现的漏洞类型",
  "sort_order": 10,
  "status": "active"
}
```

**响应示例**:

成功 (201 Created):
```json
{
  "data": {
    "id": 9,
    "config_type": "vulnerability_type",
    "config_key": "NEW_VULN_TYPE",
    "config_value": "新漏洞类型",
    "description": "新发现的漏洞类型",
    "sort_order": 10,
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 18. 更新配置

更新配置信息（仅管理员）。

**接口**: `PUT /api/v1/configs/:type/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 配置类型 |
| id | integer | 是 | 配置 ID |

**请求体参数**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| config_key | string | 否 | 配置键 |
| config_value | string | 否 | 配置值 |
| description | string | 否 | 配置描述 |
| sort_order | integer | 否 | 排序顺序 |
| status | string | 否 | 状态（active/inactive） |

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "id": 1,
    "config_type": "vulnerability_type",
    "config_key": "SQL_INJECTION",
    "config_value": "SQL注入（已更新）",
    "description": "SQL注入漏洞",
    "sort_order": 1,
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

---

#### 19. 删除配置

删除配置（仅管理员）。

**接口**: `DELETE /api/v1/configs/:type/:id`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| type | string | 是 | 配置类型 |
| id | integer | 是 | 配置 ID |

**响应示例**:

成功 (200 OK):
```json
{
  "message": "配置删除成功"
}
```

---

### 文件上传

#### 20. 上传文件

上传单个文件（用于报告附件等）。

**接口**: `POST /api/v1/upload`

**请求头**:
| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**请求体**:
- 使用 `multipart/form-data` 格式
- 字段名：`file`

**支持的文件类型**:
- PDF: `application/pdf`
- 图片: `image/jpeg`, `image/png`, `image/gif`
- 文档: `application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- 文本: `text/plain`

**文件大小限制**: 最大 10MB

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer <TOKEN>" \
  -F "file=@/path/to/file.pdf"
```

**响应示例**:

成功 (200 OK):
```json
{
  "data": {
    "url": "http://localhost:8080/uploads/reports/2024/01/1234567890.pdf",
    "filename": "vulnerability_report.pdf",
    "size": 1024000,
    "mime_type": "application/pdf"
  }
}
```

失败 (400 Bad Request):
```json
{
  "error": "文件大小超过限制（最大10MB）"
}
```

或

```json
{
  "error": "不支持的文件类型: image/bmp"
}
```

**文件访问**:
上传后的文件可以通过以下 URL 访问：
```
http://localhost:8080/uploads/reports/{year}/{month}/{filename}
```

---

## 数据模型

### User（用户对象）

```typescript
interface User {
  id: number;
  username: string;
  role: 'whitehat' | 'vendor' | 'admin';
  phone?: string;
  email?: string;
  name?: string;
  created_at: string;  // ISO 8601 格式
  updated_at: string;  // ISO 8601 格式
}
```

### Report（报告对象）

```typescript
interface Report {
  id: number;
  project_id: number;                    // 必填，关联项目ID
  project?: Project;                     // 关联的项目信息
  vulnerability_name: string;            // 必填，漏洞名称
  vulnerability_type_id: number;         // 必填，关联漏洞类型配置ID
  vulnerability_type?: SystemConfig;     // 关联的漏洞类型配置
  vulnerability_impact?: string;         // 漏洞的危害
  self_assessment?: string;              // 危害自评
  vulnerability_url?: string;            // 漏洞链接
  vulnerability_detail?: string;         // 漏洞详情
  attachment_url?: string;               // 附件地址
  title: string;                         // 保留字段，与vulnerability_name同步
  description?: string;                  // 保留字段，与vulnerability_detail同步
  type?: string;                         // 保留字段，从vulnerability_type同步
  severity: 'Low' | 'Medium' | 'High' | 'Critical';
  status: 'Pending' | 'Triaged' | 'Resolved' | 'Closed';
  author_id: number;
  author?: User;                         // 列表和详情接口会返回
  created_at: string;                    // ISO 8601 格式
  updated_at: string;                    // ISO 8601 格式
}
```

### Project（项目对象）

```typescript
interface Project {
  id: number;
  name: string;
  description?: string;
  note?: string;
  status: 'active' | 'inactive';
  created_at: string;  // ISO 8601 格式
  updated_at: string;  // ISO 8601 格式
}
```

### SystemConfig（系统配置对象）

```typescript
interface SystemConfig {
  id: number;
  config_type: string;                  // 配置类型（如：vulnerability_type）
  config_key: string;                    // 配置键（如：SQL_INJECTION）
  config_value: string;                  // 配置值（显示名称）
  description?: string;                  // 配置描述
  sort_order: number;                    // 排序顺序
  status: 'active' | 'inactive';         // 配置状态
  extra_data?: any;                      // 扩展数据（JSON格式）
  created_at: string;                    // ISO 8601 格式
  updated_at: string;                    // ISO 8601 格式
}
```

### UserInfoChangeRequest（用户信息变更申请对象）

```typescript
interface UserInfoChangeRequest {
  id: number;
  user_id: number;
  phone?: string;
  email?: string;
  name?: string;
  status: 'pending' | 'approved' | 'rejected';
  reviewed_at?: string;  // ISO 8601 格式
  reviewer_id?: number;
  review_note?: string;
  created_at: string;  // ISO 8601 格式
  updated_at: string;  // ISO 8601 格式
}
```

### 枚举值说明

**报告状态流转**:
```
Pending (待审) -> Triaged (已确认) -> Resolved (已修复) -> Closed (关闭)
```

**报告危害等级**:
- `Low`: 低危
- `Medium`: 中危
- `High`: 高危
- `Critical`: 严重

**用户角色**:
| 角色 | 说明 | 权限 |
|------|------|------|
| whitehat | 白帽子（默认） | 提交报告、查看报告、更新自己的报告、查看活跃项目和配置 |
| vendor | 厂商 | 查看报告、更新报告状态、查看活跃项目和配置 |
| admin | 管理员 | 所有权限（包括项目管理、配置管理） |

**变更申请状态**:
- `pending`: 待审核
- `approved`: 已通过
- `rejected`: 已拒绝

---

## 错误处理

### HTTP 状态码

| 状态码 | 说明 | 常见场景 |
|--------|------|----------|
| 200 | 成功 | GET/PUT 请求成功 |
| 201 | 创建成功 | POST 请求创建资源成功 |
| 400 | 请求错误 | 参数验证失败、业务逻辑错误 |
| 401 | 未授权 | Token 缺失、无效或过期 |
| 403 | 禁止访问 | 权限不足 |
| 404 | 未找到 | 资源不存在 |
| 500 | 服务器错误 | 服务端异常 |

### 错误响应格式

所有错误响应统一格式：
```json
{
  "error": "错误描述信息"
}
```

### 常见错误

| 错误信息 | HTTP 状态码 | 说明 |
|----------|-------------|------|
| `username already exists` | 400 | 用户名已存在 |
| `Invalid credentials` | 401 | 用户名或密码错误 |
| `Authorization header is required` | 401 | 缺少 Authorization 请求头 |
| `Invalid authorization format` | 401 | Authorization 格式错误 |
| `invalid token` | 401 | Token 无效 |
| `token has expired` | 401 | Token 已过期 |
| `permission denied` | 400 | 权限不足 |
| `title is required` | 400 | 缺少必填字段 |
| `Report not found` | 404 | 报告不存在 |
| `invalid status transition` | 400 | 非法的状态流转 |
| `only admin or vendor can change status` | 400 | 无权修改状态 |
| `至少需要提供一个要变更的字段（手机号、邮箱或姓名）` | 400 | 变更申请至少需要一个字段 |

---

## 快速开始

### curl 命令示例

#### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

#### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

#### 3. 获取漏洞类型配置列表

```bash
curl -X GET http://localhost:8080/api/v1/configs/vulnerability_type \
  -H "Authorization: Bearer <TOKEN>"
```

#### 4. 上传文件

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer <TOKEN>" \
  -F "file=@/path/to/vulnerability_report.pdf"
```

#### 5. 提交漏洞报告

```bash
# 替换 <TOKEN> 为登录返回的 token
curl -X POST http://localhost:8080/api/v1/reports \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "project_id": 1,
    "vulnerability_name": "SQL注入漏洞",
    "vulnerability_type_id": 1,
    "vulnerability_impact": "可能导致数据泄露",
    "self_assessment": "高危漏洞",
    "vulnerability_url": "https://example.com/vuln",
    "vulnerability_detail": "详细描述漏洞情况...",
    "attachment_url": "http://localhost:8080/uploads/reports/2024/01/abc123.pdf",
    "severity": "High"
  }'
```

**响应示例**:
```json
{
  "code": 200,
  "message": "漏洞报告提交成功"
}
```

#### 6. 获取报告列表

```bash
curl -X GET "http://localhost:8080/api/v1/reports?page=1&page_size=10" \
  -H "Authorization: Bearer <TOKEN>"
```

#### 7. 获取报告详情

```bash
curl -X GET http://localhost:8080/api/v1/reports/1 \
  -H "Authorization: Bearer <TOKEN>"
```

#### 8. 更新报告

```bash
curl -X PUT http://localhost:8080/api/v1/reports/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "vulnerability_name": "SQL注入漏洞（更新）",
    "vulnerability_impact": "更新后的危害描述",
    "severity": "Critical",
    "status": "Triaged"
  }'
```

#### 9. 提交信息变更申请

```bash
curl -X POST http://localhost:8080/api/v1/user/info/change \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"phone":"13800138000","email":"newemail@example.com","name":"张三"}'
```

#### 10. 获取变更申请列表

```bash
curl -X GET http://localhost:8080/api/v1/user/info/changes \
  -H "Authorization: Bearer <TOKEN>"
```

#### 11. 获取变更申请详情

```bash
curl -X GET http://localhost:8080/api/v1/user/info/changes/1 \
  -H "Authorization: Bearer <TOKEN>"
```

#### 12. 创建项目（仅admin）

```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -d '{"name":"某公司官网","description":"公司官方网站项目","note":"重要项目"}'
```

#### 13. 获取项目列表

```bash
curl -X GET "http://localhost:8080/api/v1/projects?page=1&page_size=10" \
  -H "Authorization: Bearer <TOKEN>"
```

#### 14. 获取项目详情

```bash
curl -X GET http://localhost:8080/api/v1/projects/1 \
  -H "Authorization: Bearer <TOKEN>"
```

#### 15. 更新项目（仅admin）

```bash
curl -X PUT http://localhost:8080/api/v1/projects/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -d '{"name":"某公司官网（更新）","status":"inactive"}'
```

#### 16. 删除项目（仅admin）

```bash
curl -X DELETE http://localhost:8080/api/v1/projects/1 \
  -H "Authorization: Bearer <ADMIN_TOKEN>"
```

#### 17. 获取配置列表

```bash
curl -X GET http://localhost:8080/api/v1/configs/vulnerability_type \
  -H "Authorization: Bearer <TOKEN>"
```

#### 18. 创建配置（仅admin）

```bash
curl -X POST http://localhost:8080/api/v1/configs/vulnerability_type \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -d '{"config_key":"NEW_TYPE","config_value":"新类型","description":"新漏洞类型","sort_order":10}'
```

### curl 注意事项

1. **JSON 格式**: 使用单引号包裹 JSON 字符串，避免 shell 解析问题
2. **Content-Type**: POST/PUT 请求必须设置 `Content-Type: application/json`（文件上传除外）
3. **文件上传**: 使用 `-F` 参数上传文件，不要设置 `Content-Type` 头（curl 会自动设置）
4. **Authorization**: 需认证的接口必须携带 `Authorization: Bearer <token>`
5. **Windows 用户**: 使用双引号时需要转义内部双引号，或使用 PowerShell

---

### 代码示例

#### JavaScript / TypeScript (Fetch API)

```typescript
const BASE_URL = 'http://localhost:8080/api/v1';

// 获取保存的 token
function getToken(): string | null {
  return localStorage.getItem('token');
}

// 创建带认证的请求头
function getAuthHeaders(): HeadersInit {
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  const token = getToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

// 用户注册
async function register(username: string, password: string) {
  const response = await fetch(`${BASE_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  return await response.json();
}

// 用户登录
async function login(username: string, password: string) {
  const response = await fetch(`${BASE_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  const data = await response.json();
  if (data.token) {
    localStorage.setItem('token', data.token);
  }
  return data;
}

// 获取漏洞类型配置列表
async function getVulnerabilityTypes() {
  const response = await fetch(`${BASE_URL}/configs/vulnerability_type`, {
    headers: getAuthHeaders(),
  });
  return await response.json();
}

// 上传文件
async function uploadFile(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  
  const token = getToken();
  const headers: HeadersInit = {};
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  const response = await fetch(`${BASE_URL}/upload`, {
    method: 'POST',
    headers,
    body: formData,
  });
  return await response.json();
}

// 提交漏洞报告
// 成功时返回 (200 OK): { code: 200, message: "漏洞报告提交成功" }
async function submitReport(report: {
  project_id: number;
  vulnerability_name: string;
  vulnerability_type_id: number;
  vulnerability_impact?: string;
  self_assessment?: string;
  vulnerability_url?: string;
  vulnerability_detail?: string;
  attachment_url?: string;
  severity?: string;
}) {
  const response = await fetch(`${BASE_URL}/reports`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(report),
  });
  return await response.json();
}

// 获取报告列表
async function getReports(page: number = 1, pageSize: number = 10) {
  const response = await fetch(
    `${BASE_URL}/reports?page=${page}&page_size=${pageSize}`,
    { headers: getAuthHeaders() }
  );
  return await response.json();
}

// 获取报告详情
async function getReport(id: number) {
  const response = await fetch(`${BASE_URL}/reports/${id}`, {
    headers: getAuthHeaders(),
  });
  return await response.json();
}

// 更新报告
async function updateReport(id: number, data: {
  project_id?: number;
  vulnerability_name?: string;
  vulnerability_type_id?: number;
  vulnerability_impact?: string;
  self_assessment?: string;
  vulnerability_url?: string;
  vulnerability_detail?: string;
  attachment_url?: string;
  severity?: string;
  status?: string;
}) {
  const response = await fetch(`${BASE_URL}/reports/${id}`, {
    method: 'PUT',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  return await response.json();
}

// 获取项目列表
async function getProjects(page: number = 1, pageSize: number = 10) {
  const response = await fetch(
    `${BASE_URL}/projects?page=${page}&page_size=${pageSize}`,
    { headers: getAuthHeaders() }
  );
  return await response.json();
}

// 获取项目详情
async function getProject(id: number) {
  const response = await fetch(`${BASE_URL}/projects/${id}`, {
    headers: getAuthHeaders(),
  });
  return await response.json();
}

// 获取配置列表
async function getConfigs(configType: string) {
  const response = await fetch(`${BASE_URL}/configs/${configType}`, {
    headers: getAuthHeaders(),
  });
  return await response.json();
}

// 提交信息变更申请
async function submitInfoChange(data: {
  phone?: string;
  email?: string;
  name?: string;
}) {
  const response = await fetch(`${BASE_URL}/user/info/change`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  return await response.json();
}

// 获取变更申请列表
async function getInfoChanges() {
  const response = await fetch(`${BASE_URL}/user/info/changes`, {
    headers: getAuthHeaders(),
  });
  return await response.json();
}

// 获取变更申请详情
async function getInfoChange(id: number) {
  const response = await fetch(`${BASE_URL}/user/info/changes/${id}`, {
    headers: getAuthHeaders(),
  });
  return await response.json();
}
```

#### 使用 Axios

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

// 请求拦截器：自动添加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器：处理认证错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      // 可以在这里跳转到登录页
    }
    return Promise.reject(error);
  }
);

// 用户注册
export const register = (username: string, password: string) =>
  api.post('/auth/register', { username, password });

// 用户登录
export const login = async (username: string, password: string) => {
  const response = await api.post('/auth/login', { username, password });
  if (response.data.token) {
    localStorage.setItem('token', response.data.token);
  }
  return response;
};

// 获取漏洞类型配置
export const getVulnerabilityTypes = () => api.get('/configs/vulnerability_type');

// 上传文件
export const uploadFile = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  return api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
};

// 提交报告
// 成功时返回 (200 OK): { code: 200, message: "漏洞报告提交成功" }
export const submitReport = (report: {
  project_id: number;
  vulnerability_name: string;
  vulnerability_type_id: number;
  vulnerability_impact?: string;
  self_assessment?: string;
  vulnerability_url?: string;
  vulnerability_detail?: string;
  attachment_url?: string;
  severity?: string;
}) => api.post('/reports', report);

// 获取报告列表
export const getReports = (page: number = 1, pageSize: number = 10) =>
  api.get('/reports', { params: { page, page_size: pageSize } });

// 获取报告详情
export const getReport = (id: number) => api.get(`/reports/${id}`);

// 更新报告
export const updateReport = (id: number, data: {
  project_id?: number;
  vulnerability_name?: string;
  vulnerability_type_id?: number;
  vulnerability_impact?: string;
  self_assessment?: string;
  vulnerability_url?: string;
  vulnerability_detail?: string;
  attachment_url?: string;
  severity?: string;
  status?: string;
}) => api.put(`/reports/${id}`, data);

// 获取项目列表
export const getProjects = (page: number = 1, pageSize: number = 10) =>
  api.get('/projects', { params: { page, page_size: pageSize } });

// 获取项目详情
export const getProject = (id: number) => api.get(`/projects/${id}`);

// 获取配置列表
export const getConfigs = (configType: string) => api.get(`/configs/${configType}`);

// 提交信息变更申请
export const submitInfoChange = (data: {
  phone?: string;
  email?: string;
  name?: string;
}) => api.post('/user/info/change', data);

// 获取变更申请列表
export const getInfoChanges = () => api.get('/user/info/changes');

// 获取变更申请详情
export const getInfoChange = (id: number) => api.get(`/user/info/changes/${id}`);
```

---

**文档版本**: 3.0.0  
**最后更新**: 2024
