.PHONY: run run-migrate build test clean docker-build docker-run tidy lint migrate migrate-status init init-force seed-projects seed-projects-force seed-users seed-users-force seed-reports seed-reports-force seed-all help

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

## init: 初始化系统必需数据（危害等级等）
init:
	go run cmd/init/main.go

## init-force: 强制初始化系统数据（跳过已存在的数据）
init-force:
	go run cmd/init/main.go -force

## seed-projects: 填充项目测试数据
seed-projects:
	go run cmd/seed-projects/main.go

## seed-projects-force: 强制填充项目测试数据（跳过已存在的数据）
seed-projects-force:
	go run cmd/seed-projects/main.go -force

## seed-users: 填充测试用户数据
seed-users:
	go run cmd/seed-users/main.go

## seed-users-force: 强制填充测试用户数据（跳过已存在的数据）
seed-users-force:
	go run cmd/seed-users/main.go -force

## seed-reports: 填充漏洞报告测试数据（需要先运行 seed-users）
seed-reports:
	go run cmd/seed-reports/main.go

## seed-reports-force: 强制填充漏洞报告测试数据（跳过已存在的数据）
seed-reports-force:
	go run cmd/seed-reports/main.go -force

## seed-all: 填充所有测试数据（项目、用户、报告）
seed-all: seed-projects seed-users seed-reports
	@echo "[OK] All test data seeded successfully!"

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
	@echo "  migrate              Run database migrations"
	@echo "  migrate-status       Show migration status"
	@echo "  init                 Initialize system data (severity levels, etc.)"
	@echo "  init-force           Force init system data (skip existing)"
	@echo "  seed-projects        Seed projects test data"
	@echo "  seed-projects-force  Force seed projects test data (skip existing)"
	@echo "  seed-users           Seed users test data"
	@echo "  seed-users-force     Force seed users test data (skip existing)"
	@echo "  seed-reports         Seed reports test data (requires seed-users)"
	@echo "  seed-reports-force   Force seed reports test data (skip existing)"
	@echo "  seed-all             Seed all test data (projects + users + reports)"
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
