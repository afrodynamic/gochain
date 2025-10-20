# ------------------------------------------------------------------------------
# Project Metadata
# ------------------------------------------------------------------------------
APP_NAME := gochain
API_DIR := api
API_SPEC := $(API_DIR)/openapi.yaml
GEN_CONFIG := $(API_DIR)/oapi-codegen.yaml

# ------------------------------------------------------------------------------
# Go Toolchain Setup
# ------------------------------------------------------------------------------
GO := go
GOBIN ?= $(shell cd $(API_DIR) && $(GO) env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell cd $(API_DIR) && $(GO) env GOPATH)/bin
endif

OAPI_BIN := $(GOBIN)/oapi-codegen
LINT_BIN := $(GOBIN)/golangci-lint

# ------------------------------------------------------------------------------
# Tool Versions
# ------------------------------------------------------------------------------
OAPI_MOD := github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
LINT_MOD := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0

# ------------------------------------------------------------------------------
# Primary Targets
# ------------------------------------------------------------------------------
.PHONY: init
init: tools generate

.PHONY: all
all: tidy generate test

# ------------------------------------------------------------------------------
# Dependency Management
# ------------------------------------------------------------------------------
.PHONY: tidy
tidy:
	cd $(API_DIR) && $(GO) mod tidy

# ------------------------------------------------------------------------------
# Tool Installation
# ------------------------------------------------------------------------------
.PHONY: tools
tools:
	@test -x "$(OAPI_BIN)" || (cd $(API_DIR) && $(GO) install $(OAPI_MOD))
	@test -x "$(LINT_BIN)" || (cd $(API_DIR) && $(GO) install $(LINT_MOD))

# ------------------------------------------------------------------------------
# Code Generation
# ------------------------------------------------------------------------------
.PHONY: generate
generate: tools
	$(OAPI_BIN) -config $(GEN_CONFIG) $(API_SPEC)

# ------------------------------------------------------------------------------
# Run Command
# ------------------------------------------------------------------------------
.PHONY: run
run:
	cd $(API_DIR) && $(GO) run ./cmd/server

# ------------------------------------------------------------------------------
# Quality: Tests, Coverage, Lint
# ------------------------------------------------------------------------------
.PHONY: test
test:
	cd $(API_DIR) && $(GO) test ./... -race -cover -count=1 -v

.PHONY: cover
cover:
	cd $(API_DIR) && $(GO) test ./... -coverprofile=cover.out
	cd $(API_DIR) && $(GO) tool cover -html=cover.out

.PHONY: lint
lint: tools
	cd $(API_DIR) && $(LINT_BIN) run ./...

# ------------------------------------------------------------------------------
# Build & Clean
# ------------------------------------------------------------------------------
.PHONY: build
build:
	cd $(API_DIR)/cmd/server && $(GO) build -o ../../../bin/$(APP_NAME)

.PHONY: clean
clean:
	rm -rf bin $(API_DIR)/cover.out
