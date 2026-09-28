# ═══════════════════════════════════════════════
#  Products API — Makefile
# ═══════════════════════════════════════════════

BINARY_NAME := products-api
MODULE      := github.com/abdallah-elngar/products-api
GO          := go
GOTEST      := $(GO) test
GOFLAGS     :=

GREEN  := \033[0;32m
RED    := \033[0;31m
YELLOW := \033[1;33m
CYAN   := \033[0;36m
BOLD   := \033[1m
NC     := \033[0m

.PHONY: help
help: ## Show this help
	@echo ""
	@echo "$(BOLD)$(CYAN)Products API — Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

.DEFAULT_GOAL := help

.PHONY: setup
setup: ## Setup development environment
	@echo "$(BOLD)$(CYAN)🔧 Setting up...$(NC)"
	@$(GO) mod download
	@$(GO) mod verify
	@cp -n .env.example .env || true
	@echo "$(GREEN)✅ Setup complete!$(NC)"

.PHONY: run
run: ## Run the server
	@echo "$(BOLD)$(CYAN)🚀 Starting server...$(NC)"
	@$(GO) run ./cmd/server

.PHONY: build
build: ## Build binary
	@echo "$(BOLD)$(CYAN)🔨 Building...$(NC)"
	@mkdir -p bin
	@$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) ./cmd/server
	@echo "$(GREEN)✅ Build complete!$(NC)"

.PHONY: test
test: ## Run tests
	@echo "$(BOLD)$(CYAN)🧪 Running tests...$(NC)"
	@$(GOTEST) -v -race ./...
	@echo "$(GREEN)✅ Tests passed!$(NC)"

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage: coverage.html$(NC)"

.PHONY: lint
lint: ## Run linter
	@echo "$(BOLD)$(CYAN)🔍 Linting...$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 || \
		{ echo "$(RED)❌ golangci-lint not installed$(NC)"; exit 1; }
	@golangci-lint run --timeout=5m
	@echo "$(GREEN)✅ Lint complete!$(NC)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(BOLD)$(CYAN)🎨 Formatting...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✅ Format complete!$(NC)"

.PHONY: vet
vet: ## Run go vet
	@$(GO) vet ./...

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BOLD)$(CYAN)🧹 Cleaning...$(NC)"
	@rm -rf bin/ tmp/
	@rm -f coverage.out coverage.html
	@rm -f products.db
	@echo "$(GREEN)✅ Clean complete!$(NC)"

.PHONY: seed
seed: ## Seed database with test data
	@echo "$(BOLD)$(CYAN)🌱 Seeding data...$(NC)"
	@chmod +x seed-data.sh
	@./seed-data.sh

.PHONY: test-routes
test-routes: ## Run full routes test
	@echo "$(BOLD)$(CYAN)🧪 Testing routes...$(NC)"
	@chmod +x test-routes.sh
	@./test-routes.sh

.PHONY: ci
ci: vet test build ## CI pipeline

.PHONY: docker-build
docker-build: ## Build Docker image
	@docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@docker run --rm -it -p 8080:8080 $(BINARY_NAME):latest

.PHONY: version
version: ## Show version
	@echo "Version: $$(git describe --tags --always 2>/dev/null || echo 'dev')"
