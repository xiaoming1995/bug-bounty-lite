# Bug Bounty Lite - 命令使用手册

> 使用 `make help` 可查看简略帮助，本文档提供更详细的说明。

---

## 📋 目录

- [开发命令](#开发命令)
- [数据库命令](#数据库命令)
- [测试数据填充](#测试数据填充)
- [文章审核](#文章审核)
- [漏洞审核](#漏洞审核)
- [测试命令](#测试命令)
- [Docker 命令](#docker-命令)
- [工具命令](#工具命令)

---

## 🚀 开发命令

| 命令 | 说明 |
|------|------|
| `make run` | 运行服务器（不执行数据库迁移） |
| `make run-migrate` | 运行服务器（先执行数据库迁移） |
| `make build` | 编译项目到 `bin/server` |

```bash
# 日常开发
make run

# 首次运行或有数据库变更时
make run-migrate
```

---

## 🗄️ 数据库命令

| 命令 | 说明 |
|------|------|
| `make migrate` | 执行数据库迁移 |
| `make migrate-status` | 查看迁移状态 |
| `make init` | 初始化系统数据（危害等级等） |
| `make init-force` | 强制初始化系统数据（跳过已存在） |

```bash
# 新环境初始化流程
make migrate
make init
```

---

## 🌱 测试数据填充

### 基础填充命令

| 命令 | 说明 |
|------|------|
| `make seed-all` | 填充所有测试数据（推荐） |
| `make seed-organizations` | 填充组织数据 |
| `make seed-projects` | 填充项目数据 |
| `make seed-users` | 填充用户数据 |
| `make seed-avatars` | 填充头像库数据 |
| `make seed-reports` | 填充漏洞报告数据（需先 seed-users） |

**强制填充**（跳过已存在的数据）：在命令后加 `-force`

```bash
make seed-users-force
make seed-projects-force
```

### 指定用户的项目数据

为特定用户生成项目测试数据：

```bash
# 通过用户 ID
make seed-project-data USER=1

# 通过用户名
make seed-project-data USERNAME=admin

# 清理该用户的数据
make seed-project-data USER=1 CLEAN=1
```

### 学习中心文章数据

为学习中心生成测试文章：

```bash
# 生成所有测试文章（10篇）
make seed-articles

# 清理后重新生成
make seed-articles CLEAN=1

# 生成指定数量
make seed-articles COUNT=5
```

---

## 📝 文章审核

> 在后台管理页面完成前，使用这些命令审核文章。

| 命令 | 说明 |
|------|------|
| `make review-list` | 查看所有待审核文章 |
| `make review-published` | 查看所有已发布文章（含精选标记） |
| `make review-approve ID=<文章ID>` | 审核通过 |
| `make review-reject ID=<文章ID> REASON="原因"` | 驳回文章 |
| `make review-featured ID=<文章ID>` | 设为精选 ⭐ |
| `make review-unfeatured ID=<文章ID>` | 取消精选 |
| `make review-interactive` | 交互式审核模式 |

### 使用示例

```bash
# 1. 查看待审核列表
make review-list

# 2. 审核通过 ID=5 的文章
make review-approve ID=5

# 3. 驳回 ID=3 的文章
make review-reject ID=3 REASON="内容不符合规范，请修改后重新提交"

# 4. 查看已发布文章
make review-published

# 5. 设为精选
make review-featured ID=5

# 6. 取消精选
make review-unfeatured ID=5

# 7. 或使用交互式模式（推荐新手）
make review-interactive
```

---

## 🔒 漏洞审核

> 用于审核用户提交的漏洞报告。

| 命令 | 说明 |
|------|------|
| `make vuln-list` | 查看所有待审核漏洞报告 |
| `make vuln-audited` | 查看所有已审核的报告 |
| `make vuln-all` | 查看所有漏洞报告 |
| `make vuln-approve ID=<ID> SEVERITY=<等级>` | 审核通过 |
| `make vuln-reject ID=<ID>` | 驳回报告 |
| `make vuln-interactive` | 交互式审核模式 |

### 危害等级说明

| 等级 | 英文 | 说明 |
|------|------|------|
| 严重 | Critical | 影响最大，需立即修复 |
| 高危 | High | 影响较大，优先级高 |
| 中危 | Medium | 影响中等，需要关注 |
| 低危 | Low | 影响较小 |
| 无危害 | None | 无实际影响 |

### 使用示例

```bash
# 1. 查看待审核列表
make vuln-list

# 2. 审核通过 ID=5 的报告，评为高危
make vuln-approve ID=5 SEVERITY=High

# 3. 驳回 ID=3 的报告
make vuln-reject ID=3

# 4. 查看已审核报告
make vuln-audited

# 5. 交互式模式（推荐）
make vuln-interactive
```

## 🧪 测试命令

| 命令 | 说明 |
|------|------|
| `make test` | 运行所有测试 |
| `make test-cover` | 运行测试并生成覆盖率报告 |

```bash
# 生成的报告文件
# - coverage.out   (原始数据)
# - coverage.html  (可视化报告)
```

---

## 🐳 Docker 命令

| 命令 | 说明 |
|------|------|
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-run` | 运行 Docker 容器（端口 8080） |
| `make docker-compose-up` | docker-compose 启动 |
| `make docker-compose-down` | docker-compose 停止 |

> 📖 **详细文档**：Docker 环境下执行管理脚本的完整说明请参阅 [DOCKER.md](./DOCKER.md)

---

## 🔧 工具命令

| 命令 | 说明 |
|------|------|
| `make tidy` | 整理 Go 模块依赖 |
| `make lint` | 运行代码检查（需安装 golangci-lint） |
| `make clean` | 清理构建产物 |
| `make help` | 显示帮助信息 |

---

## 💡 常用工作流

### 新环境初始化

```bash
# 1. 数据库迁移
make migrate

# 2. 初始化系统数据
make init

# 3. 填充测试数据
make seed-all

# 4. 启动服务
make run
```

### 日常开发

```bash
# 启动服务
make run

# 代码检查
make lint

# 运行测试
make test
```

### 文章审核工作流

```bash
# 查看待审核
make review-list

# 逐个审核
make review-approve ID=1
make review-approve ID=2
make review-reject ID=3 REASON="需要补充更多细节"
```

### 漏洞审核工作流

```bash
# 查看待审核
make vuln-list

# 逐个审核
make vuln-approve ID=1 SEVERITY=High
make vuln-approve ID=2 SEVERITY=Medium
make vuln-reject ID=3
```
