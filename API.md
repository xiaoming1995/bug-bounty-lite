# Bug Bounty Lite API 文档

## 目录

- [基础信息](#基础信息)
- [请求规范](#请求规范)
- [认证说明](#认证说明)
- [API 端点](#api-端点)
  - [认证相关](#认证相关)
  - [漏洞报告相关](#漏洞报告相关)
- [数据模型](#数据模型)
- [错误处理](#错误处理)
- [curl 命令示例](#curl-命令示例)
- [代码示例](#代码示例)

---

## 基础信息

| 项目 | 值 |
|------|-----|
| Base URL | `http://localhost:8080` |
| API 前缀 | `/api/v1` |
| 协议 | HTTP/HTTPS |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |

---

## 请求规范

### 通用请求头

所有请求必须包含以下请求头：

| 请求头 | 值 | 必填 | 说明 |
|--------|-----|------|------|
| Content-Type | `application/json` | 是（POST/PUT） | 请求体格式 |
| Accept | `application/json` | 否 | 期望的响应格式 |
| Authorization | `Bearer <token>` | 视接口而定 | JWT 认证令牌 |

### 请求头示例

**公开接口（无需认证）**:
```http
POST /api/v1/auth/register HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{"username":"test","password":"123456"}
```

**需认证接口**:
```http
POST /api/v1/reports HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{"title":"SQL Injection","severity":"High"}
```

### 请求体格式

- 使用 JSON 格式
- 字段名使用 snake_case（下划线命名）
- 字符串值使用双引号
- 布尔值使用 true/false
- 空值使用 null

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

---

## API 端点

### 认证相关

#### 1. 用户注册

创建新用户账号。

**请求**

```
POST /api/v1/auth/register
```

**请求头**

| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |

**请求体参数**

| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| username | string | 是 | 1-64字符 | 用户名，唯一 |
| password | string | 是 | 最少6位 | 密码 |

**请求示例**
```json
{
  "username": "whitehat_user",
  "password": "secure123"
}
```

**响应**

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

**请求**

```
POST /api/v1/auth/login
```

**请求头**

| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |

**请求体参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**
```json
{
  "username": "whitehat_user",
  "password": "secure123"
}
```

**响应**

成功 (200 OK):
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJyb2xlIjoid2hpdGVoYXQiLCJleHAiOjE3MzI5NTQwMDB9.xxxxx",
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

**请求**

```
POST /api/v1/reports
```

**请求头**

| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**请求体参数**

| 字段 | 类型 | 必填 | 约束 | 说明 |
|------|------|------|------|------|
| title | string | 是 | 最大255字符 | 漏洞标题 |
| description | string | 否 | 无限制 | 漏洞描述 |
| type | string | 否 | 最大50字符 | 漏洞类型 |
| severity | string | 否 | 枚举值 | 危害等级，默认 `Low` |

**severity 可选值**: `Low`, `Medium`, `High`, `Critical`

**请求示例**
```json
{
  "title": "SQL Injection in Login Form",
  "description": "The login form is vulnerable to SQL injection attacks...",
  "type": "SQL Injection",
  "severity": "High"
}
```

**响应**

成功 (201 Created):
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
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

失败 (400 Bad Request):
```json
{
  "error": "title is required"
}
```

失败 (401 Unauthorized):
```json
{
  "error": "Authorization header is required"
}
```

---

#### 4. 获取报告列表

获取漏洞报告列表，支持分页。

**请求**

```
GET /api/v1/reports?page=1&page_size=10
```

**请求头**

| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**查询参数**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | integer | 否 | 1 | 页码（从1开始） |
| page_size | integer | 否 | 10 | 每页数量（最大100） |

**响应**

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

**请求**

```
GET /api/v1/reports/:id
```

**请求头**

| 名称 | 值 | 必填 |
|------|-----|------|
| Authorization | Bearer {token} | 是 |

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 报告 ID |

**响应**

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

**请求**

```
PUT /api/v1/reports/:id
```

**请求头**

| 名称 | 值 | 必填 |
|------|-----|------|
| Content-Type | application/json | 是 |
| Authorization | Bearer {token} | 是 |

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 报告 ID |

**请求体参数**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 否 | 漏洞标题 |
| description | string | 否 | 漏洞描述 |
| type | string | 否 | 漏洞类型 |
| severity | string | 否 | 危害等级 |
| status | string | 否 | 状态（仅 admin/vendor） |

**severity 可选值**: `Low`, `Medium`, `High`, `Critical`

**status 可选值**: `Pending`, `Triaged`, `Resolved`, `Closed`

**请求示例**
```json
{
  "status": "Triaged",
  "severity": "Critical"
}
```

**响应**

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

## 数据模型

### User（用户对象）

```typescript
interface User {
  id: number;
  username: string;
  role: 'whitehat' | 'vendor' | 'admin';
  created_at: string;  // ISO 8601 格式
  updated_at: string;  // ISO 8601 格式
}
```

### Report（报告对象）

```typescript
interface Report {
  id: number;
  title: string;
  description?: string;
  type?: string;
  severity: 'Low' | 'Medium' | 'High' | 'Critical';
  status: 'Pending' | 'Triaged' | 'Resolved' | 'Closed';
  author_id: number;
  author?: User;  // 列表和详情接口会返回
  created_at: string;  // ISO 8601 格式
  updated_at: string;  // ISO 8601 格式
}
```

### 状态说明

**报告状态流转**:
```
Pending (待审) -> Triaged (已确认) -> Resolved (已修复) -> Closed (关闭)
```

**用户角色**:
| 角色 | 说明 | 权限 |
|------|------|------|
| whitehat | 白帽子（默认） | 提交报告、查看报告、更新自己的报告 |
| vendor | 厂商 | 查看报告、更新报告状态 |
| admin | 管理员 | 所有权限 |

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

---

## curl 命令示例

### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

### 3. 提交漏洞报告

```bash
# 替换 <TOKEN> 为登录返回的 token
curl -X POST http://localhost:8080/api/v1/reports \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"title":"SQL Injection","description":"Found SQL injection in login","type":"SQL Injection","severity":"High"}'
```

### 4. 获取报告列表

```bash
curl -X GET "http://localhost:8080/api/v1/reports?page=1&page_size=10" \
  -H "Authorization: Bearer <TOKEN>"
```

### 5. 获取报告详情

```bash
curl -X GET http://localhost:8080/api/v1/reports/1 \
  -H "Authorization: Bearer <TOKEN>"
```

### 6. 更新报告

```bash
curl -X PUT http://localhost:8080/api/v1/reports/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"severity":"Critical","status":"Triaged"}'
```

### curl 注意事项

1. **JSON 格式**: 使用单引号包裹 JSON 字符串，避免 shell 解析问题
2. **Content-Type**: POST/PUT 请求必须设置 `Content-Type: application/json`
3. **Authorization**: 需认证的接口必须携带 `Authorization: Bearer <token>`
4. **Windows 用户**: 使用双引号时需要转义内部双引号，或使用 PowerShell

---

## 代码示例

### JavaScript / TypeScript

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

// 提交漏洞报告
async function submitReport(report: {
  title: string;
  description?: string;
  type?: string;
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
  title?: string;
  description?: string;
  type?: string;
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
```

### 使用 Axios

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

// 提交报告
export const submitReport = (report: {
  title: string;
  description?: string;
  type?: string;
  severity?: string;
}) => api.post('/reports', report);

// 获取报告列表
export const getReports = (page: number = 1, pageSize: number = 10) =>
  api.get('/reports', { params: { page, page_size: pageSize } });

// 获取报告详情
export const getReport = (id: number) => api.get(`/reports/${id}`);

// 更新报告
export const updateReport = (id: number, data: {
  title?: string;
  description?: string;
  type?: string;
  severity?: string;
  status?: string;
}) => api.put(`/reports/${id}`, data);
```

---

**文档版本**: 2.1.0  
**更新日期**: 2024
