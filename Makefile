.PHONY: run run-migrate build test clean docker-build docker-run tidy lint migrate migrate-status init init-force seed-organizations seed-organizations-force seed-avatars seed-avatars-force seed-projects seed-projects-force seed-users seed-users-force seed-reports seed-reports-force seed-all seed-project-data seed-articles review-list review-approve review-reject review-interactive help

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

## seed-organizations: 填充组织测试数据
seed-organizations:
	go run cmd/seed-organizations/main.go

## seed-organizations-force: 强制填充组织测试数据（跳过已存在的数据）
seed-organizations-force:
	go run cmd/seed-organizations/main.go -force

## seed-avatars: 填充头像测试数据
seed-avatars:
	go run cmd/seed-avatars/main.go

## seed-avatars-force: 强制填充头像测试数据（清空并重新填充）
seed-avatars-force:
	go run cmd/seed-avatars/main.go -force

## seed-all: 填充所有测试数据（项目、用户、头像、报告）
seed-all:
	go run cmd/seed-organizations/main.go
	go run cmd/seed-projects/main.go
	go run cmd/seed-users/main.go
	go run cmd/seed-avatars/main.go
	go run cmd/seed-reports/main.go

## seed-project-data: 生成项目测试数据并指派给指定用户
## 用法: make seed-project-data USER=1  (通过用户ID)
##       make seed-project-data USERNAME=admin  (通过用户名)
##       make seed-project-data USER=1 CLEAN=1  (清理数据)
seed-project-data:
	@if [ -n "$(USER)" ]; then \
		if [ -n "$(CLEAN)" ]; then \
			go run cmd/seed-project-data/main.go -user $(USER) -clean; \
		else \
			go run cmd/seed-project-data/main.go -user $(USER); \
		fi \
	elif [ -n "$(USERNAME)" ]; then \
		if [ -n "$(CLEAN)" ]; then \
			go run cmd/seed-project-data/main.go -username $(USERNAME) -clean; \
		else \
			go run cmd/seed-project-data/main.go -username $(USERNAME); \
		fi \
	else \
		echo "请指定用户: make seed-project-data USER=<用户ID> 或 USERNAME=<用户名>"; \
		echo "清理数据: make seed-project-data USER=<用户ID> CLEAN=1"; \
		exit 1; \
	fi

# ===========================
# 学习中心文章数据
# ===========================

## seed-articles: 生成测试文章数据
##       make seed-articles
##       make seed-articles CLEAN=1  (清理后重新生成)
##       make seed-articles COUNT=5  (生成指定数量)
seed-articles:
	@if [ -n "$(CLEAN)" ]; then \
		go run cmd/seed-articles/main.go -clean; \
	elif [ -n "$(COUNT)" ]; then \
		go run cmd/seed-articles/main.go -count $(COUNT); \
	else \
		go run cmd/seed-articles/main.go; \
	fi

# ===========================
# 文章审核命令
# ===========================

## review-list: 查看所有待审核的文章
review-list:
	go run cmd/review-articles/main.go -list

## review-approve: 审核通过文章
## 用法: make review-approve ID=5
review-approve:
	@if [ -z "$(ID)" ]; then \
		echo "请指定文章ID: make review-approve ID=<文章ID>"; \
		exit 1; \
	fi
	go run cmd/review-articles/main.go -approve $(ID)

## review-reject: 驳回文章
## 用法: make review-reject ID=5 REASON="内容不符合规范"
review-reject:
	@if [ -z "$(ID)" ]; then \
		echo "请指定文章ID: make review-reject ID=<文章ID> REASON=\"驳回原因\""; \
		exit 1; \
	fi
	@if [ -z "$(REASON)" ]; then \
		echo "请指定驳回原因: make review-reject ID=<文章ID> REASON=\"驳回原因\""; \
		exit 1; \
	fi
	go run cmd/review-articles/main.go -reject $(ID) -reason "$(REASON)"

## review-interactive: 交互式审核模式
review-interactive:
	go run cmd/review-articles/main.go -i

## review-published: 查看所有已发布的文章
review-published:
	go run cmd/review-articles/main.go -published

## review-featured: 设为精选
## 用法: make review-featured ID=5
review-featured:
	@if [ -z "$(ID)" ]; then \
		echo "请指定文章ID: make review-featured ID=<文章ID>"; \
		exit 1; \
	fi
	go run cmd/review-articles/main.go -featured $(ID)

## review-unfeatured: 取消精选
## 用法: make review-unfeatured ID=5
review-unfeatured:
	@if [ -z "$(ID)" ]; then \
		echo "请指定文章ID: make review-unfeatured ID=<文章ID>"; \
		exit 1; \
	fi
	go run cmd/review-articles/main.go -unfeatured $(ID)


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
	@echo "  seed-avatars         Seed avatar test data (platform avatar library)"
	@echo "  seed-avatars-force   Force seed avatar test data (clear and reseed)"
	@echo "  seed-all             Seed all test data (organizations + projects + users + avatars + reports)"
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
