# Bug Bounty Lite (Go)

这是一个轻量级的 Web 安全众测平台后端，基于 Golang + Gin + Gorm + MySQL 构建。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL 5.7+
- **ORM**: Gorm
- **配置**: Viper
- **认证**: JWT
- **架构**: Modular Monolith (Clean Architecture)

## 功能特性

- ✅ 用户注册/登录（JWT 认证）
- ✅ 漏洞报告管理（CRUD）
- ✅ 用户信息变更申请（需后台审核）
- ✅ 角色权限管理（白帽子/厂商/管理员）
- ✅ 数据库迁移工具
- ✅ 统一响应格式
- ✅ CORS 支持

## 快速开始

### 1. 环境准备

确保本地已安装 MySQL，并创建数据库：

```bash
# macOS (Homebrew)
brew install mysql
brew services start mysql

# 或使用 Docker
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -p 3306:3306 \
  mysql:8

# 创建数据库
mysql -u root -p123456 -e "CREATE DATABASE bugbounty CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 2. 配置

复制配置文件模板：

```bash
cp config/config.yaml.example config/config.yaml
```

修改 `config/config.yaml` 中的数据库连接信息：

```yaml
database:
  dsn: "root:YOUR_PASSWORD@tcp(localhost:3306)/bugbounty?charset=utf8mb4&parseTime=True&loc=Local"
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 运行

**方式一：直接运行（不执行迁移）**

```bash
make run
```

**方式二：运行并执行数据库迁移**

```bash
make run-migrate
```

**方式三：先迁移再运行**

```bash
make migrate    # 执行迁移
make run        # 运行服务
```

服务启动在: http://localhost:8080

## 常用命令

```bash
make run            # 运行项目（不迁移）
make run-migrate    # 运行项目（先迁移）
make migrate        # 执行数据库迁移
make migrate-status # 查看迁移状态
make build          # 编译项目
make test           # 运行测试
make docker-build   # 构建 Docker 镜像
make docker-run     # 运行 Docker 容器
make stop           # 停止运行中的服务
make help           # 查看所有命令
```

## 项目结构

```
bug-bounty-lite/
├── cmd/
│   ├── server/main.go      # HTTP 服务入口
│   └── migrate/main.go     # 数据库迁移工具
├── config/
│   ├── config.yaml         # 配置文件
│   └── config.yaml.example # 配置模板
├── internal/
│   ├── domain/             # 领域模型和接口
│   │   ├── user.go         # 用户实体
│   │   ├── report.go       # 漏洞报告实体
│   │   └── user_info_change.go # 用户信息变更申请实体
│   ├── handler/            # HTTP 处理器
│   │   ├── user_handler.go
│   │   ├── report_handler.go
│   │   └── user_info_change_handler.go
│   ├── middleware/         # 中间件
│   │   ├── auth.go         # JWT 认证
│   │   ├── cors.go         # CORS
│   │   └── logger.go       # 日志
│   ├── repository/         # 数据访问层
│   │   ├── user_repo.go
│   │   ├── report_repo.go
│   │   └── user_info_change_repo.go
│   ├── router/             # 路由配置
│   │   └── router.go
│   └── service/            # 业务逻辑层
│       ├── user_service.go
│       ├── report_service.go
│       └── user_info_change_service.go
├── pkg/
│   ├── config/             # 配置加载
│   ├── database/           # 数据库连接
│   ├── jwt/                # JWT 认证
│   ├── migrate/            # 迁移工具
│   └── response/           # 统一响应
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── README.md
├── API.md                  # API 文档
└── DATABASE.md             # 数据库文档
```

## API 文档

详见 [API.md](./API.md)

主要 API 端点：

- **认证相关**
  - `POST /api/v1/auth/register` - 用户注册
  - `POST /api/v1/auth/login` - 用户登录

- **漏洞报告相关**（需认证）
  - `POST /api/v1/reports` - 提交漏洞报告
  - `GET /api/v1/reports` - 获取报告列表
  - `GET /api/v1/reports/:id` - 获取报告详情
  - `PUT /api/v1/reports/:id` - 更新报告

- **用户信息变更**（需认证）
  - `POST /api/v1/user/info/change` - 提交信息变更申请
  - `GET /api/v1/user/info/changes` - 获取变更申请列表
  - `GET /api/v1/user/info/changes/:id` - 获取变更申请详情

## 数据库文档

详见 [DATABASE.md](./DATABASE.md)

主要数据表：

- `users` - 用户表
- `reports` - 漏洞报告表
- `user_info_change_requests` - 用户信息变更申请表

## 开发说明

### 数据库迁移

项目使用 GORM 的 AutoMigrate 功能进行数据库迁移：

```bash
# 执行迁移
make migrate

# 查看迁移状态
make migrate-status
```

### 认证流程

1. 用户注册/登录获取 JWT Token
2. 访问需要认证的接口时，在请求头中携带 Token：
   ```
   Authorization: Bearer <token>
   ```

### 用户信息变更流程

1. 用户提交信息变更申请（手机号/邮箱/姓名）
2. 申请状态为 `pending`（待审核）
3. 后台管理员审核通过后，状态变为 `approved`，并更新用户信息
4. 审核拒绝后，状态变为 `rejected`

## 配置说明

配置文件位于 `config/config.yaml`：

```yaml
server:
  port: ":8080"      # 服务端口
  mode: "debug"       # 运行模式: debug/release

database:
  dsn: "..."         # MySQL 连接字符串
  max_idle: 10       # 最大空闲连接数
  max_open: 100      # 最大打开连接数

jwt:
  secret: "..."      # JWT 密钥（请修改为复杂字符串）
  expire: 7200       # Token 过期时间（秒）
```

## 许可证

MIT License
