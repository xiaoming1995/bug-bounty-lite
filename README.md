# Bug Bounty Lite (Go)

这是一个轻量级的 Web 安全众测平台后端，基于 Golang + Gin + Gorm + PostgreSQL 构建。

## 🛠 技术栈
- **语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: Gorm
- **配置**: Viper
- **架构**: Modular Monolith (Clean Architecture)

## 🚀 快速开始

### 1. 环境准备
确保本地已安装 PostgreSQL，并创建数据库 `bugbounty`。

### 2. 配置
复制配置文件模板：
cp config/config.yaml.example config/config.yaml

修改 `config/config.yaml` 中的数据库账号密码。

### 3. 运行
go run cmd/server/main.go

服务启动在: http://localhost:8080