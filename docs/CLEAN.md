# Bug Bounty Lite - 测试数据清理

> 清理测试数据时会**保护 admin 用户**和**系统配置数据**，不会被删除。

---

## 📋 快速开始

```bash
# 1. 查看当前数据量
run.bat clean-stats

# 2. 交互式清理（会提示确认）
run.bat clean

# 3. 强制清理（自动化脚本使用）
run.bat clean-force
```

---

## 🖥️ Windows 批处理命令

| 命令 | 说明 |
|------|------|
| `run.bat clean-stats` | 查看当前数据统计（不删除） |
| `run.bat clean` | 交互式清理所有测试数据（需确认） |
| `run.bat clean-force` | 强制清理所有测试数据（无需确认） |

---

## 💻 PowerShell 命令

| 命令 | 说明 |
|------|------|
| `.\run.ps1 clean-stats` | 查看当前数据统计 |
| `.\run.ps1 clean-data` | 交互式清理所有测试数据 |
| `.\run.ps1 clean-force` | 强制清理所有测试数据 |

---

## 🎯 按类型清理

可以只清理特定类型的数据：

```bash
# 只清理用户（保留 admin）
go run cmd/clean/main.go -users -confirm

# 只清理项目
go run cmd/clean/main.go -projects -confirm

# 只清理报告
go run cmd/clean/main.go -reports -confirm

# 只清理文章
go run cmd/clean/main.go -articles -confirm
```

### 命令行参数

| 参数 | 说明 |
|------|------|
| `-all` | 清理所有测试数据 |
| `-users` | 只清理用户数据（保留 admin） |
| `-projects` | 只清理项目数据 |
| `-reports` | 只清理报告数据 |
| `-articles` | 只清理文章数据 |
| `-stats` | 只显示统计信息，不删除 |
| `-confirm` | 跳过确认提示（用于自动化） |
| `-help` | 显示帮助信息 |

---

## 🛡️ 安全保护

| 数据类型 | 清理行为 |
|----------|----------|
| admin 用户 | ✅ **保留** |
| 普通用户（whitehat/vendor） | ❌ 删除 |
| 系统配置（漏洞类型、危害等级） | ✅ **保留** |
| 项目 | ❌ 删除 |
| 漏洞报告 | ❌ 删除 |
| 文章 | ❌ 删除 |
| 头像库 | ❌ 删除 |
| 组织 | ❌ 删除 |

---

## 📖 使用示例

### 场景 1：完整清理后重新填充测试数据

```bash
# 1. 查看当前数据
run.bat clean-stats

# 2. 清理所有测试数据
run.bat clean-force

# 3. 重新填充
run.bat seed-all
```

### 场景 2：只清理用户相关数据

```bash
go run cmd/clean/main.go -users -confirm
```

### 场景 3：CI/CD 自动化脚本

```bash
# 无交互式确认，直接清理
run.bat clean-force

# 或使用 Go 命令
go run cmd/clean/main.go -all -confirm
```

---

## 🔗 相关命令

- 测试数据填充：参见 [COMMANDS.md](./COMMANDS.md#测试数据填充)
- 数据库迁移：参见 [COMMANDS.md](./COMMANDS.md#数据库命令)
