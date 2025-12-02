# Bug Bounty Lite 数据库表结构

## 概述

- **数据库**: MySQL
- **ORM**: GORM
- **版本**: v1.1.0
- **最后更新**: 2024

---

## 表清单

| 序号 | 表名 | 说明 | 源文件 |
|------|------|------|--------|
| 1 | users | 用户表 | `internal/domain/user.go` |
| 2 | reports | 漏洞报告表 | `internal/domain/report.go` |

---

## 表结构详情

### 1. users - 用户表

存储平台用户信息，包括白帽子、厂商、管理员。

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | 自增 | 用户ID |
| created_at | DATETIME(3) | - | 自动 | 创建时间 |
| updated_at | DATETIME(3) | - | 自动 | 更新时间 |
| username | VARCHAR(64) | UNIQUE, NOT NULL | - | 用户名 |
| password | VARCHAR(255) | NOT NULL | - | 密码(bcrypt加密) |
| role | VARCHAR(20) | - | 'whitehat' | 用户角色 |

**索引**:
| 索引名 | 字段 | 类型 |
|--------|------|------|
| idx_users_username | username | UNIQUE |

**角色枚举值**:
| 值 | 说明 |
|-----|------|
| whitehat | 白帽子（默认） |
| vendor | 厂商 |
| admin | 管理员 |

---

### 2. reports - 漏洞报告表

存储白帽子提交的漏洞报告。

| 字段名 | 类型 | 约束 | 默认值 | 说明 |
|--------|------|------|--------|------|
| id | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | 自增 | 报告ID |
| created_at | DATETIME(3) | - | 自动 | 创建时间 |
| updated_at | DATETIME(3) | - | 自动 | 更新时间 |
| title | VARCHAR(255) | NOT NULL | - | 漏洞标题 |
| description | TEXT | - | NULL | 漏洞描述 |
| type | VARCHAR(50) | - | NULL | 漏洞类型 |
| severity | VARCHAR(20) | - | 'Low' | 危害等级 |
| status | VARCHAR(20) | - | 'Pending' | 报告状态 |
| author_id | BIGINT UNSIGNED | FOREIGN KEY | - | 提交者ID |

**索引**:
| 索引名 | 字段 | 类型 |
|--------|------|------|
| idx_reports_status | status | INDEX |

**外键**:
| 外键名 | 字段 | 引用表 | 引用字段 |
|--------|------|--------|----------|
| fk_reports_author | author_id | users | id |

**危害等级枚举值**:
| 值 | 说明 |
|-----|------|
| Low | 低危 |
| Medium | 中危 |
| High | 高危 |
| Critical | 严重 |

**状态枚举值**:
| 值 | 说明 | 可流转到 |
|-----|------|----------|
| Pending | 待审核 | Triaged, Closed |
| Triaged | 已确认 | Resolved, Closed |
| Resolved | 已修复 | Closed |
| Closed | 已关闭 | - |

**状态流转图**:
```
Pending --> Triaged --> Resolved --> Closed
   |           |            |
   +-----------+------------+
               |
               v
            Closed
```

---

## ER 图

```
+-------------------+          +-------------------+
|      users        |          |     reports       |
+-------------------+          +-------------------+
| PK | id           |<---------| FK | author_id    |
|    | created_at   |          |    | id           |
|    | updated_at   |          |    | created_at   |
|    | username     |          |    | updated_at   |
|    | password     |          |    | title        |
|    | role         |          |    | description  |
+-------------------+          |    | type         |
                               |    | severity     |
                               |    | status       |
                               +-------------------+

关系: users (1) <-----> (N) reports
```

---

## 迭代记录

### v1.1.0 (数据库迁移)

**变更**:
- 数据库从 PostgreSQL 迁移到 MySQL
- DSN 格式变更为 MySQL 格式
- 时间字段类型从 TIMESTAMPTZ 改为 DATETIME(3)
- ID 字段类型从 BIGSERIAL 改为 BIGINT UNSIGNED AUTO_INCREMENT

**说明**:
- 兼容 MySQL 5.7+
- 需要重新执行迁移创建表

---

### v1.0.0 (初始版本)

**新增表**:
- users: 用户表
- reports: 漏洞报告表

**说明**:
- 基础用户认证功能
- 漏洞报告 CRUD 功能
- 状态流转机制

---

### 迭代模板

```markdown
### vX.X.X (YYYY-MM-DD)

**新增表**:
- 表名: 说明

**修改表**:
- 表名:
  - 新增字段: 字段名 (类型) - 说明
  - 修改字段: 字段名 - 修改内容
  - 删除字段: 字段名

**删除表**:
- 表名: 原因

**说明**:
- 变更原因和影响
```

---

## 常用 SQL

### 查看表结构

```sql
-- 查看所有表
SHOW TABLES;

-- 查看表字段
DESCRIBE users;
DESCRIBE reports;

-- 查看表详细结构
SHOW CREATE TABLE users;
SHOW CREATE TABLE reports;

-- 查看索引
SHOW INDEX FROM users;
SHOW INDEX FROM reports;
```

### 数据统计

```sql
-- 用户统计
SELECT role, COUNT(*) FROM users GROUP BY role;

-- 报告状态统计
SELECT status, COUNT(*) FROM reports GROUP BY status;

-- 报告危害等级统计
SELECT severity, COUNT(*) FROM reports GROUP BY severity;
```

---

## 注意事项

1. **密码存储**: password 字段使用 bcrypt 加密，永远不要明文存储
2. **时间字段**: created_at/updated_at 由 GORM 自动维护
3. **软删除**: 当前未启用，如需要可添加 deleted_at 字段
4. **外键约束**: author_id 关联 users 表，删除用户前需处理关联报告
5. **字符集**: 建议使用 utf8mb4 字符集以支持完整的 Unicode

---

## 相关文件

| 文件 | 说明 |
|------|------|
| `internal/domain/user.go` | User 实体定义 |
| `internal/domain/report.go` | Report 实体定义 |
| `pkg/database/mysql.go` | 数据库连接 |
| `pkg/migrate/migrate.go` | 数据库迁移工具 |
| `cmd/migrate/main.go` | 迁移命令入口 |

