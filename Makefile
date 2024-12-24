.DEFAULT_GOAL := help

SHELL := bash
PATH := $(CURDIR)/.dev/go-tools/bin:$(PATH)
COMMIT_HASH := $(shell git rev-parse HEAD)

# Release version: You need to update this version when you release the project.
VERSION := 0.2.0
BUILD_LDFLAGS = "-s -w -X github.com/kohkimakimoto/actions-gateway/version.CommitHash=$(COMMIT_HASH) -X github.com/kohkimakimoto/actions-gateway/version.Version=$(VERSION)"

# Load .env file if it exists.
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

.PHONY: help
help: ## Show help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'


# --------------------------------------------------------------------------------------
# Development environment
# --------------------------------------------------------------------------------------
.PHONY: setup
setup: ## Setup development environment
	@echo "==> Setting up development environment..."
	@mkdir -p $(CURDIR)/.dev/go-tools
	@export GOPATH=$(CURDIR)/.dev/go-tools && \
		go install honnef.co/go/tools/cmd/staticcheck@latest && \
		go install github.com/golang/mock/mockgen@latest && \
		go install github.com/Songmu/goxz/cmd/goxz@latest && \
		go install github.com/tcnksm/ghr@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html/cmd/gocov-html@latest
	@export GOPATH=$(CURDIR)/.dev/go-tools && go clean -modcache && rm -rf $(CURDIR)/.dev/go-tools/pkg

.PHONY: clean
clean: ## Clean up development environment
	@rm -rf .dev

.PHONY: clean/build
clean/build: ## Clean up build directory
	@rm -rf .dev/build


# --------------------------------------------------------------------------------------
# Build
# --------------------------------------------------------------------------------------
.PHONY: build
build: ## Build dev binary
	@mkdir -p .dev/build/dev
	@CGO_ENABLED=0 go build -ldflags=$(BUILD_LDFLAGS) -o .dev/build/dev/actions-gateway ./cmd/actions-gateway

build/release: ## Build release binary
	@mkdir -p .dev/build/release
	@CGO_ENABLED=0 go build -ldflags=$(BUILD_LDFLAGS) -trimpath -o .dev/build/release/actions-gateway ./cmd/actions-gateway

.PHONY: build/dist
build/dist: ## Build cross-platform binaries for distribution
	@mkdir -p .dev/build/dist
	@CGO_ENABLED=0 goxz -n actions-gateway -os=linux,darwin -static -build-ldflags=$(BUILD_LDFLAGS) -trimpath -d=.dev/build/dist ./cmd/actions-gateway

.PHONY: build/server-frontend ## Build server frontend
build/server-frontend:
	@if [[ ! -d "server/node_modules" ]]; then \
		cd server && npm install; \
	fi
	@cd server && npm run build


# --------------------------------------------------------------------------------------
# Testing, Formatting and etc.
# --------------------------------------------------------------------------------------
.PHONY: format
format: ## Format source code
	@go fmt ./...

.PHONY: lint
lint: ## Lint source code
	@staticcheck ./...

.PHONY: test
test: ## Run tests
	@go test -race -timeout 30m ./...

.PHONY: test/short
test/short: ## Run short tests
	@go test -short -race -timeout 30m ./...

.PHONY: test/verbos
test/verbose: ## Run tests with verbose outputting
	@go test -race -timeout 30m -v ./...

.PHONY: test/cover
test/cover: ## Run tests with coverage report
	@mkdir -p $(CURDIR)/.dev/test
	@go test -race -coverpkg=./... -coverprofile=$(CURDIR)/.dev/test/coverage.out ./...
	@gocov convert $(CURDIR)/.dev/test/coverage.out | gocov-html > $(CURDIR)/.dev/test/coverage.html

.PHONY: open/coverage
open/coverage: ## Open coverage report
	@open $(CURDIR)/.dev/test/coverage.html


# --------------------------------------------------------------------------------------
# Go commands
# --------------------------------------------------------------------------------------
.PHONY: go-generate
go-generate: ## Run go generate
	@go generate ./...

.PHONY: go-mod-tidy
go-mod-tidy: ## Run go mod tidy
	@go mod tidy


# --------------------------------------------------------------------------------------
# Release
# --------------------------------------------------------------------------------------
.PHONY: release
release: guard-GITHUB_TOKEN ## Release the project with the specified version and tags it
	@$(MAKE) clean/build
	@$(MAKE) build/server-frontend
	@$(MAKE) build/dist
	@ghr -b "Release v$(VERSION)" v$(VERSION) .dev/build/dist


# --------------------------------------------------------------------------------------
# Utilities
# --------------------------------------------------------------------------------------
# This is a utility for checking variable definition
guard-%:
	@if [[ -z '${${*}}' ]]; then echo 'ERROR: variable $* not set' && exit 1; fi
