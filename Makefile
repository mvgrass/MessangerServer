SHELL := /bin/bash
GO := go
PROJECT_NAME := MessangerServer
VERSION := $(shell git describe --tags --always --dirty)
MODULE_DIRS := $(shell find ./services ./libs -type f -name go.mod -exec dirname {} \;)
MODULE_NAMES := $(notdir $(MODULE_DIRS))


OUTPUT_DIR := build
COVERAGE_DIR := coverage


LDFLAGS := -ldflags "-X main.Version=$(VERSION)"
TEST_FLAGS := -race -coverprofile=$(COVERAGE_DIR)/coverage.out

BUILD_TARGETS := $(addprefix build-,$(MODULE_NAMES))
RUN_TARGETS := $(addprefix run-,$(MODULE_NAMES))
TEST_TARGETS := $(addprefix test-,$(MODULE_NAMES))

.PHONY: all build test lint clean mod-tidy work-sync help $(BUILD_TARGETS) $(RUN_TARGETS) $(TEST_TARGETS)


all: build

build: $(BUILD_TARGETS)

test: $(TEST_TARGETS)

define RUN_RULE
run-$(1): build-$(1)
	@echo "Running $(1)..."
	@cd $(OUTPUT_DIR)/$(1) && ./$(1)
endef

#need to do it with foreach because of dependency on build targets
$(foreach module,$(MODULE_NAMES),$(eval $(call RUN_RULE,$(module))))

$(BUILD_TARGETS):
	$(eval MODULE_DIR := $(filter %/$(@:build-%=%), $(MODULE_DIRS)))
	$(eval MODULE_PATH := $(filter %/$*,$(MODULES)))
	@echo "Building $(@:build-%=%)..."
	@mkdir -p $(OUTPUT_DIR)
	cd $(MODULE_DIR) && $(GO) build $(LDFLAGS) -o ../../$(OUTPUT_DIR)/$(@:build-%=%)/$(@:build-%=%) ./cmd

	cp -r global_config/.env $(OUTPUT_DIR)/$(@:build-%=%)/;
	@if [ -d "$(MODULE_DIR)/cmd/config" ]; then \
		echo "Copying config for $(@:build-%=%)..."; \
		cp -r $(MODULE_DIR)/cmd/config $(OUTPUT_DIR)/$(@:build-%=%)/; \
	fi

$(TEST_TARGETS):
	$(eval MODULE_PATH := $(filter %/$*,$(MODULES)))
	@mkdir -p $(COVERAGE_DIR)
	@echo "Testing $(@:test-%=%)..."
	cd $(MODULE_DIR) && $(GO) test $(TEST_FLAGS) ./...

clean:
	@rm -rf $(OUTPUT_DIR) $(COVERAGE_DIR)

mod-tidy:
	@for dir in $(MODULE_DIRS); do \
		echo "Tidying $$(basename $$dir)..."; \
		pushd $$dir && $(GO) mod tidy || exit 1; \
		popd; \
	done

work-sync:
	@echo "Syncing workspace..."
	$(GO) work sync

help:
	@echo "Available commands:"
	@echo "  all              build all modules"
	@echo "  build            build all modules"
	@echo "  test             test all modules"
	@echo "  run-<name>       run one module by it's name (example run-auth)"
	@echo "  build-<name>     build one module by it's name (example build-auth)"
	@echo "  test-<name>      test one module by it's name (example test-auth)"
	@echo "  clean            clean build and test artifacts"
	@echo "  mod-tidy         update dependencies"
	@echo "  work-sync        sync workspace"