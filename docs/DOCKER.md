# Bug Bounty Lite - Docker 部署与管理指南

> 本文档介绍如何在 Docker 环境中部署和管理项目。

---

## 📋 目录

- [镜像构建](#镜像构建)
- [服务管理](#服务管理)
- [容器内执行管理脚本](#容器内执行管理脚本)
- [常见问题](#常见问题)

---

## 🏗️ 镜像构建

### 使用 docker-compose 构建

```bash
# 构建并启动 go-api 服务
docker compose up -d --build go-api

# 仅构建不启动
docker compose build go-api

# 查看构建日志
docker compose logs -f go-api
```

### 单独构建镜像

```bash
# 构建镜像
docker build -t bug-bounty-lite .

# 运行容器
docker run -d -p 8080:8080 --env-file .env bug-bounty-lite
```

---

## 🚀 服务管理

```bash
# 启动服务
docker compose up -d go-api

# 停止服务
docker compose stop go-api

# 重启服务
docker compose restart go-api

# 查看日志
docker compose logs -f go-api

# 查看运行状态
docker compose ps
```

---

## 🛠️ 容器内执行管理脚本

镜像内预编译了以下管理工具：

| 工具 | 说明 |
|------|------|
| `./review_reports` | 漏洞审核 |
| `./review_articles` | 文章审核 |
| `./seed_articles` | 填充学习中心文章数据 |
| `./migrate_tool` | 数据库迁移 |
| `./init_tool` | 系统初始化 |

### 漏洞审核

```bash
# 查看待审核列表
docker compose exec go-api ./review_reports -list

# 查看已审核列表
docker compose exec go-api ./review_reports -audited

# 查看所有报告
docker compose exec go-api ./review_reports -all

# 审核通过（需指定危害等级）
docker compose exec go-api ./review_reports -approve 5 -severity High

# 驳回报告
docker compose exec go-api ./review_reports -reject 3

# 交互式审核模式（推荐）
docker compose exec -it go-api ./review_reports -i
```

**危害等级说明：**

| 等级 | 英文 | 说明 |
|------|------|------|
| 严重 | Critical | 影响最大，需立即修复 |
| 高危 | High | 影响较大，优先级高 |
| 中危 | Medium | 影响中等，需要关注 |
| 低危 | Low | 影响较小 |
| 无危害 | None | 无实际影响 |

### 文章审核

```bash
# 查看待审核列表
docker compose exec go-api ./review_articles -list

# 查看已发布文章
docker compose exec go-api ./review_articles -published

# 审核通过
docker compose exec go-api ./review_articles -approve 5

# 驳回文章
docker compose exec go-api ./review_articles -reject 3 -reason "内容不符合规范"

# 设为精选
docker compose exec go-api ./review_articles -featured 5

# 取消精选
docker compose exec go-api ./review_articles -unfeatured 5

# 交互式审核模式
docker compose exec -it go-api ./review_articles -i
```

### 数据填充

```bash
# 填充学习中心测试文章（10篇）
docker compose exec go-api ./seed_articles

# 清理后重新填充
docker compose exec go-api ./seed_articles -clean
```

### 数据库工具

```bash
# 执行数据库迁移
docker compose exec go-api ./migrate_tool

# 初始化系统数据（危害等级、漏洞类型等）
docker compose exec go-api ./init_tool
```

---

## ❓ 常见问题

### 1. 交互式模式无法输入

确保使用 `-it` 参数：

```bash
docker compose exec -it go-api ./review_reports -i
```

### 2. 无法连接数据库

检查环境变量配置：

```bash
# 查看容器环境变量
docker compose exec go-api env | grep DB
```

### 3. 查看容器内文件

```bash
# 进入容器 shell
docker compose exec go-api /bin/sh

# 列出可用工具
ls -la
```

### 4. 重新构建镜像

如果代码有更新，需要重新构建：

```bash
docker compose up -d --build go-api
```
