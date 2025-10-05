# Generic Makefile for UP Tools
# This file is symlinked as Makefile in each tool directory

# Automatically determine tool name from directory
TOOL := $(notdir $(CURDIR))

# Binary name (main CLI is just "up", everything else is "up-<tool>")
ifeq ($(TOOL),up)
	BINARY := up
else
	BINARY := up-$(TOOL)
endif

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Tool: $(TOOL)'
	@echo 'Binary: $(BINARY)'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run tests
	go test -v -race ./...

.PHONY: lint
lint: ## Run linter
	go tool golangci-lint run ./...

.PHONY: build
build: test ## Build the tool
	go build -v -o $(BINARY) .

.PHONY: clean
clean: ## Clean build artifacts
	go clean
	rm -f $(BINARY)

.PHONY: install
install: ## Install dependencies
	go mod download
	go mod tidy

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: test-ci
test-ci: ## Run CI tests locally using act (requires: brew install act)
	@cd .. && act --container-architecture linux/amd64 -j test

.DEFAULT_GOAL := build
