# ------------------------------------------------------------------------------
# Project Metadata
# ------------------------------------------------------------------------------
APP_NAME := gochain
API_SPEC := ./api/openapi.yaml
GEN_CONFIG := ./api/oapi-codegen.yaml

# ------------------------------------------------------------------------------
# Go Toolchain Setup
# ------------------------------------------------------------------------------
GO := go

# Determine binary installation directory (GOBIN > GOPATH/bin)
GOBIN ?= $(shell $(GO) env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell $(GO) env GOPATH)/bin
endif

# Tool binary paths (used in recipes)
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
# One-time setup: install tools and generate code
.PHONY: init
init: tools generate

# Daily workflow: tidy, generate, test
.PHONY: all
all: tidy generate test

# ------------------------------------------------------------------------------
# Dependency Management
# ------------------------------------------------------------------------------
.PHONY: tidy
tidy:
	$(GO) mod tidy

# ------------------------------------------------------------------------------
# Tool Installation (idempotent)
# ------------------------------------------------------------------------------
# Installs pinned binaries only if missing; safe to rerun
.PHONY: tools
tools:
	@test -x "$(OAPI_BIN)" || $(GO) install $(OAPI_MOD)
	@test -x "$(LINT_BIN)" || $(GO) install $(LINT_MOD)

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
	cd api && $(GO) run ./cmd/server

# ------------------------------------------------------------------------------
# Quality: Tests, Coverage, Lint
# ------------------------------------------------------------------------------
.PHONY: test
test:
	$(GO) test ./api/... -race -cover -count=1 -v

.PHONY: cover
cover:
	$(GO) test ./api/... -coverprofile=cover.out
	$(GO) tool cover -html=cover.out

.PHONY: lint
lint: tools
	$(LINT_BIN) run ./...

# ------------------------------------------------------------------------------
# Build & Clean
# ------------------------------------------------------------------------------
.PHONY: build
build:
	cd api/cmd/server && $(GO) build -o ../../../bin/$(APP_NAME)

.PHONY: clean
clean:
	rm -rf bin cover.out
