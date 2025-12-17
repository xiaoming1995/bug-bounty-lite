# 构建阶段
FROM golang:1.24-alpine AS builder

# 优化：设置国内代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# 任务2：编译迁移工具 (Migrate Tool)
# 这样我们在最终的精简镜像里也能执行数据库迁移，而不需要 'go run'
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o migrate_tool ./cmd/migrate

# 运行阶段
FROM alpine:latest

# 安装 CA 证书（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates tzdata

# 安装基础证书和时区
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -g '' appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

#  从 builder 阶段复制编译好的迁移工具
COPY --from=builder /app/migrate_tool .

# 复制配置文件
COPY --from=builder /app/config ./config

#  关键：将 migrations SQL 文件夹复制过去，否则迁移工具找不到 SQL 文件
COPY --from=builder /app/migrations ./migrations

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./server"]

