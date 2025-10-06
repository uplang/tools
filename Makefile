# Makefile for UP Tools Repository

# Find all tool directories with go.mod
TOOLS := $(patsubst %/,%,$(dir $(wildcard */go.mod)))

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Tool targets (generated):'
	@printf "  %-15s %s\n" "test" "Run tests for all tools"
	@printf "  %-15s %s\n" "build" "Build all tools"
	@printf "  %-15s %s\n" "clean" "Clean all tools"
	@printf "  %-15s %s\n" "install" "Install all tools"

# Template for creating tool targets
# Usage: $(eval $(call TOOL_TARGET,target-name,dependencies))
# Creates both 'target' and 'target-<tool>' phony targets
define TOOL_TARGET
$(1)_TOOLS := $$(patsubst %,$(1)-%,$$(TOOLS))

.PHONY: $(1)
$(1): $(2) $$($(1)_TOOLS)

.PHONY: $$($(1)_TOOLS)
$$($(1)_TOOLS):
	$$(MAKE) -C $$(patsubst $(1)-%,%,$$@) $(1)
endef

# Generate targets for each tool operation
$(eval $(call TOOL_TARGET,test))
$(eval $(call TOOL_TARGET,build,test))
$(eval $(call TOOL_TARGET,clean))
$(eval $(call TOOL_TARGET,install,build))

.PHONY: test-ci
test-ci: ## Run CI tests locally using act (requires: brew install act)
	act --container-architecture linux/amd64 -j test
	act --container-architecture linux/amd64 -j build

.PHONY: bump-patch-version
bump-patch-version: ## Bump patch version (e.g., v0.0.1 -> v0.0.2) and create tag locally
	@$(MAKE) tag NEW_VERSION=$(shell go tool svu patch)

.PHONY: bump-minor-version
bump-minor-version: ## Bump minor version (e.g., v0.0.1 -> v0.1.0) and create tag locally
	@$(MAKE) tag NEW_VERSION=$(shell go tool svu minor)

.PHONY: tag
tag: NEW_VERSION ?= $(shell go tool svu patch)
tag:
	@echo "Creating tag: $(NEW_VERSION)"
	git tag -a $(NEW_VERSION) -m "Release version $(NEW_VERSION)"
	@echo "Tag created. Push with: git push origin $(NEW_VERSION)"

.DEFAULT_GOAL := build
