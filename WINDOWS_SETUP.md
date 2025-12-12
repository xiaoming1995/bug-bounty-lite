# Windows 环境快速开始指南

## 前置要求

- ✅ Go 1.21 或更高版本
- ✅ MySQL 5.7+ 或 MySQL 8.0+
- ✅ Git（可选）

## 一、安装 MySQL

### 方式1：使用 Docker（推荐）

```powershell
# 启动 MySQL 容器
docker run -d --name mysql `
  -e MYSQL_ROOT_PASSWORD=123456 `
  -e MYSQL_DATABASE=bugbounty `
  -p 3306:3306 `
  mysql:8

# 验证 MySQL 是否运行
docker ps
```

### 方式2：安装本地 MySQL

1. 下载 MySQL 安装包：https://dev.mysql.com/downloads/mysql/
2. 安装并启动 MySQL 服务
3. 创建数据库：

```sql
CREATE DATABASE bugbounty CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 二、配置项目

### 1. 检查配置文件

确保 `config/config.yaml` 文件存在。如果不存在，复制示例文件：

```powershell
# PowerShell
Copy-Item config\config.yaml.example config\config.yaml

# 或者使用 CMD
copy config\config.yaml.example config\config.yaml
```

### 2. 修改数据库配置

编辑 `config/config.yaml`，修改数据库密码：

```yaml
database:
  dsn: "root:123456@tcp(localhost:3306)/bugbounty?charset=utf8mb4&parseTime=True&loc=Local"
```

**注意**：将 `123456` 替换为你的 MySQL 密码。

## 三、初始化数据库

### 使用批处理脚本（推荐）

```powershell
# 1. 执行数据库迁移（创建表结构）
.\run.bat migrate

# 2. 初始化系统数据（危害等级等）
.\run.bat init

# 3. 填充所有测试数据
.\run.bat seed-all
```

### 或使用 Go 命令

```powershell
# 1. 执行数据库迁移
go run cmd/migrate/main.go

# 2. 初始化系统数据
go run cmd/init/main.go

# 3. 填充测试数据
go run cmd/seed-projects/main.go
go run cmd/seed-users/main.go
go run cmd/seed-reports/main.go
```

## 四、运行项目

### 使用批处理脚本

```powershell
.\run.bat run
```

### 或使用 Go 命令

```powershell
go run cmd/server/main.go
```

服务启动后，访问：http://localhost:8080

## 五、验证安装

### 测试用户登录

```powershell
# 使用 PowerShell
$body = @{
    username = "whitehat_zhang"
    password = "password123"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body
```

### 或使用 curl（如果已安装）

```powershell
curl -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"whitehat_zhang\",\"password\":\"password123\"}'
```

## 六、常用命令

### 批处理脚本命令

```powershell
.\run.bat help              # 显示帮助信息
.\run.bat run               # 运行服务器
.\run.bat migrate           # 执行数据库迁移
.\run.bat init              # 初始化系统数据
.\run.bat seed-projects     # 填充项目测试数据
.\run.bat seed-users        # 填充用户测试数据
.\run.bat seed-reports      # 填充报告测试数据
.\run.bat seed-all          # 填充所有测试数据
.\run.bat build             # 编译项目
.\run.bat test              # 运行测试
```

### 编译和运行

```powershell
# 编译项目
.\run.bat build

# 运行编译后的程序
.\bin\server.exe
```

## 七、测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| whitehat_zhang | password123 | 白帽子 |
| whitehat_li | password123 | 白帽子 |
| vendor_test | password123 | 厂商 |
| admin_test | admin123 | 管理员 |

## 八、常见问题

### 1. 配置文件找不到

**错误**：`Fatal error config file: Config File "config" Not Found`

**解决**：
- 确保从项目根目录运行命令
- 确保 `config/config.yaml` 文件存在
- 使用提供的 `run.bat` 脚本

### 2. 数据库连接失败

**错误**：`Error 1045: Access denied for user 'root'@'localhost'`

**解决**：
- 检查 `config/config.yaml` 中的数据库密码是否正确
- 确保 MySQL 服务正在运行
- 确保数据库 `bugbounty` 已创建

### 3. 端口被占用

**错误**：`bind: address already in use`

**解决**：
- 修改 `config/config.yaml` 中的端口号
- 或关闭占用 8080 端口的程序

```powershell
# 查看占用 8080 端口的进程
netstat -ano | findstr :8080

# 结束进程（替换 PID 为实际进程ID）
taskkill /PID <PID> /F
```

## 九、开发工具推荐

- **IDE**: GoLand / VS Code
- **API 测试**: Postman / Insomnia
- **数据库管理**: MySQL Workbench / DBeaver
- **Git 客户端**: Git for Windows / GitHub Desktop

## 十、下一步

- 📖 阅读 [API.md](./API.md) 了解 API 接口
- 📖 阅读 [DATABASE.md](./DATABASE.md) 了解数据库结构
- 🔧 开始开发你的功能
- 📝 编写测试用例

---

**提示**：如果遇到问题，请查看项目根目录的 `README.md` 获取更多信息。
