# Makefile for UP Tools Repository

TOOLS := up examples language-server repl

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run tests for all tools
	@for tool in $(TOOLS); do \
		if [ -d $$tool ] && [ -f $$tool/go.mod ]; then \
			echo "Testing $$tool..."; \
			$(MAKE) -C $$tool test || exit 1; \
		fi \
	done

.PHONY: build
build: test ## Build all tools
	@for tool in $(TOOLS); do \
		if [ -d $$tool ] && [ -f $$tool/go.mod ]; then \
			echo "Building $$tool..."; \
			$(MAKE) -C $$tool build || exit 1; \
		fi \
	done

.PHONY: clean
clean: ## Clean all tools
	@for tool in $(TOOLS); do \
		if [ -d $$tool ] && [ -f $$tool/go.mod ]; then \
			echo "Cleaning $$tool..."; \
			$(MAKE) -C $$tool clean || exit 1; \
		fi \
	done

.PHONY: install
install: build ## Install all tools
	@for tool in $(TOOLS); do \
		if [ -d $$tool ] && [ -f $$tool/go.mod ]; then \
			echo "Installing $$tool..."; \
			$(MAKE) -C $$tool install || exit 1; \
		fi \
	done

.PHONY: test-ci
test-ci: ## Run CI tests locally using act (requires: brew install act)
	act --container-architecture linux/amd64 -j test
	act --container-architecture linux/amd64 -j build

.DEFAULT_GOAL := build
