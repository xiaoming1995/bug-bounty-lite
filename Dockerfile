# ==============================
# 阶段 1: 构建 (Builder)
# ==============================
FROM golang:1.24-alpine AS builder

# 优化：设置国内代理
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 安装 git (下载依赖可能需要)
RUN apk add --no-cache git

# 1. 先只复制依赖描述文件 (利用 Docker 缓存层)
COPY go.mod go.sum ./
RUN go mod download

# 2. 再复制其余所有源代码
COPY . .

# 3. 检查一下 migrations 到底在不在 (构建时打印目录结构，方便调试报错)
# 如果构建失败，看日志输出就能知道文件夹名字到底叫什么
RUN ls -F /app

# 4. 构建主程序
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# 5. 构建迁移工具
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o migrate_tool ./cmd/migrate

# ==============================
# 阶段 2: 运行 (Runner)
# ==============================
FROM alpine:latest

# 安装基础证书和时区 (只写一次)
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -g '' appuser

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/migrate_tool .

# 复制配置文件 (前提：你本地根目录下真的有 config 文件夹)
COPY --from=builder /app/config ./config

# 🔴 关键修复：请根据你的实际路径修改这里！
# 如果你确定本地根目录下有 migrations 文件夹，这行就没问题。
# 如果你的 SQL 文件在其他地方，请修改 /app/migrations 为真实路径。
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./server"]