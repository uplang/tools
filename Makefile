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

.PHONY: bump-patch-version
bump-patch-version: ## Bump patch version (e.g., v0.0.1 -> v0.0.2) and create tag locally
	@echo "Current version: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags found')"
	@NEW_VERSION=$$(go tool svu patch); \
	echo "Creating new patch version: $$NEW_VERSION"; \
	git tag -a $$NEW_VERSION -m "Release version $$NEW_VERSION"
	@echo "Tag created. Push with: git push origin $$(git describe --tags --abbrev=0)"

.PHONY: bump-minor-version
bump-minor-version: ## Bump minor version (e.g., v0.0.1 -> v0.1.0) and create tag locally
	@echo "Current version: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags found')"
	@NEW_VERSION=$$(go tool svu minor); \
	echo "Creating new minor version: $$NEW_VERSION"; \
	git tag -a $$NEW_VERSION -m "Release version $$NEW_VERSION"
	@echo "Tag created. Push with: git push origin $$(git describe --tags --abbrev=0)"

.PHONY: bump-major-version
bump-major-version: ## Bump major version (e.g., v0.0.1 -> v1.0.0) and create tag locally
	@echo "Current version: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags found')"
	@NEW_VERSION=$$(go tool svu major); \
	echo "Creating new major version: $$NEW_VERSION"; \
	git tag -a $$NEW_VERSION -m "Release version $$NEW_VERSION"
	@echo "Tag created. Push with: git push origin $$(git describe --tags --abbrev=0)"

.DEFAULT_GOAL := build
