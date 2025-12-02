.PHONY: run run-migrate build test clean docker-build docker-run tidy lint migrate migrate-status help

# 默认目标
.DEFAULT_GOAL := help

# 变量
APP_NAME := bug-bounty-lite
BINARY := server
BUILD_DIR := bin
DOCKER_IMAGE := $(APP_NAME)

# ===========================
# 开发命令
# ===========================

## run: 运行项目（不执行迁移）
run:
	go run cmd/server/main.go

## run-migrate: 运行项目（先执行迁移）
run-migrate:
	go run cmd/server/main.go --migrate

## build: 编译项目
build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY) ./cmd/server

# ===========================
# 数据库迁移
# ===========================

## migrate: 执行数据库迁移
migrate:
	go run cmd/migrate/main.go

## migrate-status: 查看迁移状态
migrate-status:
	go run cmd/migrate/main.go -status

# ===========================
# 测试命令
# ===========================

## test: 运行测试
test:
	go test -v ./...

## test-cover: 运行测试并生成覆盖率报告
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# ===========================
# 工具命令
# ===========================

## clean: 清理构建产物
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## tidy: 整理依赖
tidy:
	go mod tidy

## lint: 代码检查
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# ===========================
# Docker 命令
# ===========================

## docker-build: 构建 Docker 镜像
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## docker-run: 运行 Docker 容器
docker-run:
	docker run -p 8080:8080 --rm $(DOCKER_IMAGE)

## docker-compose-up: 使用 docker-compose 启动
docker-compose-up:
	docker-compose up -d

## docker-compose-down: 使用 docker-compose 停止
docker-compose-down:
	docker-compose down

# ===========================
# 帮助
# ===========================

## help: 显示帮助信息
help:
	@echo "Bug Bounty Lite - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Development:"
	@echo "  run            Run server (skip migrations)"
	@echo "  run-migrate    Run server with migrations"
	@echo "  build          Build binary"
	@echo ""
	@echo "Database:"
	@echo "  migrate        Run database migrations"
	@echo "  migrate-status Show migration status"
	@echo ""
	@echo "Testing:"
	@echo "  test           Run tests"
	@echo "  test-cover     Run tests with coverage"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container"
	@echo ""
	@echo "Other:"
	@echo "  tidy           Tidy go modules"
	@echo "  lint           Run linter"
	@echo "  clean          Clean build artifacts"
