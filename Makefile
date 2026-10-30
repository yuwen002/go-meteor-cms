.PHONY: help build run test cover fmt vet lint tidy migrate gen-api gen-ent gen-doc docker-build docker-up docker-down clean

BINARY_DIR := bin
API_PKG    := ./api/cms/v1

help: ## 显示可用命令
	@echo "可用命令:"
	@echo "  build         编译 API 服务与迁移工具到 bin/"
	@echo "  run           启动 API 服务"
	@echo "  test          运行全部测试"
	@echo "  cover         生成覆盖率报告"
	@echo "  fmt           格式化代码"
	@echo "  vet           静态检查"
	@echo "  lint          格式化 + 静态检查"
	@echo "  tidy          整理依赖"
	@echo "  migrate       执行数据库迁移（幂等）"
	@echo "  gen-api       按 cms.api 重新生成 API 代码"
	@echo "  gen-ent       重新生成 Ent 代码"
	@echo "  docker-build  构建镜像"
	@echo "  docker-up     启动 MySQL + 迁移 + API"
	@echo "  docker-down   停止容器"
	@echo "  clean         清理构建产物"

build: ## 编译 API 服务与迁移工具到 bin/
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/cms $(API_PKG)
	go build -o $(BINARY_DIR)/migrate ./cmd/migrate

run: ## 启动 API 服务
	go run $(API_PKG)/cms.go

test: ## 运行全部测试
	go test ./...

cover: ## 生成覆盖率报告并打开
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fmt: ## 格式化代码
	gofmt -w api cmd internal ent/schema

vet: ## 静态检查
	go vet ./...

lint: fmt vet ## 格式化 + 静态检查

tidy: ## 整理依赖
	go mod tidy

migrate: ## 执行数据库迁移（幂等）
	go run ./cmd/migrate

gen-api: ## 按 cms.api 重新生成 API 代码（需 goctl 1.10.2，必须指定 --home .goctl）
	goctl api go -api $(API_PKG)/cms.api -dir $(API_PKG) --home .goctl

gen-ent: ## 重新生成 Ent 代码
	go run -mod=mod entgo.io/ent/cmd/ent generate --feature intercept,schema/snapshot ./ent/schema

gen-doc: ## 由 cms.api 生成 OpenAPI 文档（输出 openapi.json）
	go run ./cmd/apidoc

docker-build: ## 构建镜像
	docker build -t go-meteor-cms .

docker-up: ## 启动 MySQL + 迁移 + API（迁移成功后再起服务）
	docker compose up -d --wait

docker-down: ## 停止并清理容器
	docker compose down

clean: ## 清理构建产物
	rm -rf $(BINARY_DIR) coverage.out
