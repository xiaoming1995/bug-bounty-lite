# Bug Bounty Lite 数据库表结构

## 概述

- **数据库**: MySQL
- **ORM**: GORM
- **版本**: v2.1.0
- **最后更新**: 2024-12-10

---

## 表清单

| 序号 | 表名 | 说明 | 源文件 |
|------|------|------|--------|
| 1 | users | 用户表 | `internal/domain/user.go` |
| 2 | reports | 漏洞报告表 | `internal/domain/report.go` |
| 3 | projects | 项目表 | `internal/domain/project.go` |
| 4 | system_configs | 系统配置表 | `internal/domain/system_config.go` |
| 5 | user_info_change_requests | 用户信息变更申请表 | `internal/domain/user_info_change.go` |

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
| project_id | BIGINT UNSIGNED | NOT NULL, FOREIGN KEY | - | 关联项目ID |
| vulnerability_name | VARCHAR(255) | NOT NULL | - | 漏洞名称 |
| vulnerability_type_id | BIGINT UNSIGNED | NOT NULL, FOREIGN KEY | - | 漏洞类型配置ID |
| vulnerability_impact | TEXT | - | NULL | 漏洞危害 |
| self_assessment_id | BIGINT UNSIGNED | FOREIGN KEY | NULL | 危害自评配置ID（关联system_configs表，config_type='severity_level'） |
| vulnerability_url | VARCHAR(500) | - | NULL | 漏洞链接 |
| vulnerability_detail | TEXT | - | NULL | 漏洞详情 |
| attachment_url | VARCHAR(500) | - | NULL | 附件地址 |
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
| fk_reports_project | project_id | projects | id |
| fk_reports_vuln_type | vulnerability_type_id | system_configs | id |
| fk_reports_self_assessment | self_assessment_id | system_configs | id |

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
+-------------------+          +-------------------------+          +-------------------+
|      users        |          |        reports          |          |     projects      |
+-------------------+          +-------------------------+          +-------------------+
| PK | id           |<---------| FK | author_id          |--------->| PK | id           |
|    | created_at   |          |    | id                 |          |    | name         |
|    | updated_at   |          |    | created_at         |          |    | description  |
|    | username     |          |    | updated_at         |          |    | note         |
|    | password     |          | FK | project_id         |          |    | status       |
|    | role         |          | FK | vulnerability_type_id         +-------------------+
+-------------------+          |    | vulnerability_name |
                               |    | vulnerability_impact|          +-------------------+
                               |    | self_assessment    |          |  system_configs   |
                               |    | vulnerability_url  |          +-------------------+
                               |    | vulnerability_detail--------->| PK | id           |
                               |    | attachment_url     |          |    | config_type  |
                               |    | severity           |          |    | config_key   |
                               |    | status             |          |    | config_value |
                               +-------------------------+          +-------------------+

关系: 
  - users (1) <-----> (N) reports
  - projects (1) <-----> (N) reports
  - system_configs (1) <-----> (N) reports (漏洞类型)
```

---

## 迭代记录

### v2.1.0 (2024-12-10)

**修改表**:
- reports:
  - 修改字段: self_assessment (TEXT) -> self_assessment_id (BIGINT UNSIGNED, 可为NULL)
  - 新增外键: fk_reports_self_assessment (self_assessment_id -> system_configs.id)

**说明**:
- 危害自评字段改为关联配置表，使用 `self_assessment_id` 关联到 `system_configs` 表
- `self_assessment_id` 必须对应 `config_type='severity_level'` 的配置
- 字段可为 NULL，表示未设置危害自评

---

### v2.0.0 (表结构优化)

**修改表**:
- reports:
  - 新增字段: project_id (BIGINT UNSIGNED) - 关联项目ID
  - 新增字段: vulnerability_name (VARCHAR(255)) - 漏洞名称
  - 新增字段: vulnerability_type_id (BIGINT UNSIGNED) - 漏洞类型配置ID
  - 新增字段: vulnerability_impact (TEXT) - 漏洞危害
  - 新增字段: self_assessment_id (BIGINT UNSIGNED) - 危害自评配置ID（关联system_configs表）
  - 新增字段: vulnerability_url (VARCHAR(500)) - 漏洞链接
  - 新增字段: vulnerability_detail (TEXT) - 漏洞详情
  - 新增字段: attachment_url (VARCHAR(500)) - 附件地址
  - 删除字段: title - 已被 vulnerability_name 替代
  - 删除字段: description - 已被 vulnerability_detail 替代
  - 删除字段: type - 已被 vulnerability_type_id 关联替代

**新增表**:
- projects: 项目表
- system_configs: 系统配置表
- user_info_change_requests: 用户信息变更申请表

**说明**:
- 报告表结构重构，使用关联表替代冗余字段
- 新增项目管理功能
- 新增系统配置管理功能

---

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

