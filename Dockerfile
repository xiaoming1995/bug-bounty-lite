# ==============================
# 阶段 1: 构建 (Builder)
# ==============================
FROM golang:1.24-alpine AS builder

# 优化：设置国内代理加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 安装构建所需工具
RUN apk add --no-cache git

# 1. 先复制依赖文件（利用 Docker 缓存层优化构建速度）
COPY go.mod go.sum ./
RUN go mod download

# 2. 复制源代码
COPY . .

# 3. 构建所有二进制文件
# 主服务
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# 数据库迁移工具
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o migrate_tool ./cmd/migrate

# 系统数据初始化工具
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o init_tool ./cmd/init

# ==============================
# 阶段 2: 运行 (Runner)
# ==============================
FROM alpine:latest

# 安装基础证书和时区
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 创建非 root 用户（安全最佳实践）
RUN adduser -D -g '' appuser

WORKDIR /app

# 复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/migrate_tool .
COPY --from=builder /app/init_tool .

# 复制配置文件
COPY --from=builder /app/config ./config

# 创建启动脚本
RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'set -e' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo '# 自动执行数据库迁移' >> /app/entrypoint.sh && \
    echo 'echo "[Docker] Running database migration..."' >> /app/entrypoint.sh && \
    echo './migrate_tool' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo '# 自动初始化系统数据' >> /app/entrypoint.sh && \
    echo 'echo "[Docker] Initializing system data..."' >> /app/entrypoint.sh && \
    echo './init_tool' >> /app/entrypoint.sh && \
    echo '' >> /app/entrypoint.sh && \
    echo '# 启动主服务' >> /app/entrypoint.sh && \
    echo 'echo "[Docker] Starting server..."' >> /app/entrypoint.sh && \
    echo 'exec ./server' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

# 创建上传目录
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser

# 暴露服务端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/auth/login || exit 1

# 默认入口：自动执行迁移+初始化+启动服务
ENTRYPOINT ["/app/entrypoint.sh"]

# 也可以通过以下方式单独运行各工具：
# docker run --rm <image> ./migrate_tool       # 仅迁移
# docker run --rm <image> ./init_tool          # 仅初始化
# docker run --rm <image> ./server             # 仅启动服务（跳过迁移）