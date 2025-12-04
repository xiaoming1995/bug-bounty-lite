# Bug Bounty Lite

一个轻量级的 Web 安全众测平台后端，基于 Golang + Gin + GORM + MySQL 构建。

## ✨ 特性

- ✅ **用户认证系统** - JWT 认证，支持用户注册/登录
- ✅ **漏洞报告管理** - 完整的 CRUD 操作，支持分页查询
- ✅ **用户信息变更** - 信息变更申请流程，支持后台审核
- ✅ **角色权限管理** - 白帽子/厂商/管理员三种角色
- ✅ **数据库迁移** - 基于 GORM 的自动迁移工具
- ✅ **统一响应格式** - 标准化的 API 响应结构
- ✅ **CORS 支持** - 跨域资源共享配置
- ✅ **Clean Architecture** - 清晰的分层架构设计

## 🛠 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| **语言** | Go 1.21+ | 编程语言 |
| **Web框架** | Gin | HTTP Web 框架 |
| **数据库** | MySQL 5.7+ | 关系型数据库 |
| **ORM** | GORM | Go 对象关系映射 |
| **配置管理** | Viper | 配置文件加载 |
| **认证** | JWT | JSON Web Token 认证 |
| **架构** | Clean Architecture | 分层架构设计 |

## 📋 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [配置说明](#配置说明)
- [API 文档](#api-文档)
- [数据库文档](#数据库文档)
- [开发指南](#开发指南)
- [部署](#部署)

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 5.7+ 或 MySQL 8.0+
- Make（可选，用于运行 Makefile 命令）

### 1. 克隆项目

```bash
git clone <repository-url>
cd bug-bounty-lite
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置数据库

#### 方式一：使用本地 MySQL

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

#### 方式二：使用 Docker Compose（推荐）

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8
    container_name: bugbounty-mysql
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: bugbounty
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
```

启动数据库：

```bash
docker-compose up -d
```

### 4. 配置文件

复制配置文件模板：

```bash
cp config/config.yaml.example config/config.yaml
```

编辑 `config/config.yaml`，修改数据库连接信息：

```yaml
database:
  dsn: "root:123456@tcp(localhost:3306)/bugbounty?charset=utf8mb4&parseTime=True&loc=Local"
```

### 5. 运行项目

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

服务启动后，访问: http://localhost:8080

### 6. 验证安装

```bash
# 测试健康检查（如果有的话）
curl http://localhost:8080/api/v1/health

# 测试用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

## 📁 项目结构

```
bug-bounty-lite/
├── cmd/                    # 应用程序入口
│   ├── server/            # HTTP 服务入口
│   │   └── main.go
│   └── migrate/           # 数据库迁移工具
│       └── main.go
├── config/                # 配置文件
│   ├── config.yaml        # 配置文件（需自行创建）
│   └── config.yaml.example # 配置模板
├── internal/              # 内部代码（不对外暴露）
│   ├── domain/            # 领域模型和接口
│   │   ├── user.go
│   │   ├── report.go
│   │   └── user_info_change.go
│   ├── handler/           # HTTP 处理器层
│   │   ├── user_handler.go
│   │   ├── report_handler.go
│   │   └── user_info_change_handler.go
│   ├── middleware/        # 中间件
│   │   ├── auth.go        # JWT 认证中间件
│   │   ├── cors.go        # CORS 中间件
│   │   └── logger.go      # 日志中间件
│   ├── repository/        # 数据访问层
│   │   ├── user_repo.go
│   │   ├── report_repo.go
│   │   └── user_info_change_repo.go
│   ├── router/            # 路由配置
│   │   └── router.go
│   └── service/           # 业务逻辑层
│       ├── user_service.go
│       ├── report_service.go
│       └── user_info_change_service.go
├── pkg/                   # 可复用的公共包
│   ├── config/            # 配置加载
│   ├── database/          # 数据库连接
│   ├── jwt/               # JWT 认证
│   ├── migrate/           # 迁移工具
│   └── response/          # 统一响应格式
├── Dockerfile             # Docker 镜像构建文件
├── Makefile               # 构建脚本
├── go.mod                 # Go 模块依赖
├── go.sum                 # Go 模块校验和
├── README.md              # 项目说明文档
├── API.md                 # API 接口文档
└── DATABASE.md            # 数据库文档
```

### 架构说明

项目采用 **Clean Architecture（清洁架构）** 设计，分为以下层次：

1. **Handler 层** - HTTP 请求处理，参数验证
2. **Service 层** - 业务逻辑处理
3. **Repository 层** - 数据访问，数据库操作
4. **Domain 层** - 领域模型和接口定义

## ⚙️ 配置说明

配置文件位于 `config/config.yaml`：

```yaml
server:
  port: ":8080"      # 服务端口
  mode: "debug"       # 运行模式: debug/release

database:
  dsn: "root:password@tcp(localhost:3306)/bugbounty?charset=utf8mb4&parseTime=True&loc=Local"
  max_idle: 10        # 最大空闲连接数
  max_open: 100       # 最大打开连接数

jwt:
  secret: "your-secret-key-here"  # JWT 密钥（请修改为复杂字符串）
  expire: 7200                     # Token 过期时间（秒，默认2小时）
```

### 环境变量支持

可以通过环境变量覆盖配置（需要修改配置加载代码）：

```bash
export SERVER_PORT=:8080
export DB_DSN="root:password@tcp(localhost:3306)/bugbounty?charset=utf8mb4&parseTime=True&loc=Local"
export JWT_SECRET="your-secret-key"
```

## 📚 API 文档

详细的 API 文档请参考 [API.md](./API.md)

### 主要 API 端点

#### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

#### 漏洞报告相关（需认证）
- `POST /api/v1/reports` - 提交漏洞报告
- `GET /api/v1/reports` - 获取报告列表（支持分页）
- `GET /api/v1/reports/:id` - 获取报告详情
- `PUT /api/v1/reports/:id` - 更新报告

#### 用户信息变更（需认证）
- `POST /api/v1/user/info/change` - 提交信息变更申请
- `GET /api/v1/user/info/changes` - 获取变更申请列表
- `GET /api/v1/user/info/changes/:id` - 获取变更申请详情

## 🗄️ 数据库文档

详细的数据库文档请参考 [DATABASE.md](./DATABASE.md)

### 主要数据表

- `users` - 用户表
- `reports` - 漏洞报告表
- `user_info_change_requests` - 用户信息变更申请表

## 🛠️ 开发指南

### 常用命令

```bash
# 运行项目（不迁移）
make run

# 运行项目（先迁移）
make run-migrate

# 执行数据库迁移
make migrate

# 查看迁移状态
make migrate-status

# 编译项目
make build

# 运行测试
make test

# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 停止运行中的服务
make stop

# 查看所有命令
make help
```

### 数据库迁移

项目使用 GORM 的 AutoMigrate 功能进行数据库迁移：

```bash
# 执行迁移
make migrate

# 查看迁移状态
make migrate-status
```

迁移会自动创建以下表结构：
- `users` - 用户表
- `reports` - 漏洞报告表
- `user_info_change_requests` - 用户信息变更申请表

### 认证流程

1. **用户注册/登录** - 获取 JWT Token
2. **访问受保护接口** - 在请求头中携带 Token：
   ```
   Authorization: Bearer <token>
   ```

### 用户信息变更流程

1. 用户提交信息变更申请（手机号/邮箱/姓名）
2. 申请状态为 `pending`（待审核）
3. 后台管理员审核通过后，状态变为 `approved`，并更新用户信息
4. 审核拒绝后，状态变为 `rejected`

### 角色权限

| 角色 | 说明 | 权限 |
|------|------|------|
| **whitehat** | 白帽子（默认） | 提交报告、查看报告、更新自己的报告 |
| **vendor** | 厂商 | 查看报告、更新报告状态 |
| **admin** | 管理员 | 所有权限 |

### 开发建议

1. **代码规范**: 遵循 Go 官方代码规范
2. **错误处理**: 使用统一的错误响应格式
3. **日志记录**: 使用中间件记录请求日志
4. **测试**: 编写单元测试和集成测试
5. **文档**: 及时更新 API 文档和代码注释

## 🐳 部署

### Docker 部署

#### 1. 构建镜像

```bash
make docker-build
```

或手动构建：

```bash
docker build -t bug-bounty-lite:latest .
```

#### 2. 运行容器

```bash
make docker-run
```

或手动运行：

```bash
docker run -d \
  --name bug-bounty-lite \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  bug-bounty-lite:latest
```

### 生产环境建议

1. **配置管理**
   - 使用环境变量或配置中心管理敏感信息
   - 不要在代码中硬编码密钥

2. **数据库**
   - 使用连接池优化数据库连接
   - 定期备份数据库

3. **安全**
   - 使用 HTTPS
   - 设置强密码策略
   - 定期更新依赖包

4. **监控**
   - 添加日志收集（如 ELK）
   - 添加性能监控（如 Prometheus）
   - 设置告警机制

5. **高可用**
   - 使用负载均衡
   - 数据库主从复制
   - 容器编排（Kubernetes）

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

如有问题或建议，请通过 Issue 反馈。

---

**版本**: 1.0.0  
**最后更新**: 2024
