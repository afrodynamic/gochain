APP := gochaind
API_DIR := api
PROTO_DIR := proto
BIN_DIR := bin

GO := go

GOBIN := $(shell cd $(API_DIR) && $(GO) env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell cd $(API_DIR) && $(GO) env GOPATH)/bin
endif
export PATH := $(GOBIN):$(PATH)

PROTOC_GEN_GO := google.golang.org/protobuf/cmd/protoc-gen-go@latest
PROTOC_GEN_GO_GRPC := google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
GRPC_GATEWAY := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
OPENAPI_V2 := github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
BUF_CLI := github.com/bufbuild/buf/cmd/buf@latest
GOLANGCI := github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: init
init: tools buf-gen tidy

.PHONY: tools
tools:
	@command -v buf >/dev/null 2>&1 || (cd $(API_DIR) && $(GO) install $(BUF_CLI))
	@command -v protoc-gen-go >/dev/null 2>&1 || (cd $(API_DIR) && $(GO) install $(PROTOC_GEN_GO))
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || (cd $(API_DIR) && $(GO) install $(PROTOC_GEN_GO_GRPC))
	@command -v protoc-gen-grpc-gateway >/dev/null 2>&1 || (cd $(API_DIR) && $(GO) install $(GRPC_GATEWAY))
	@command -v protoc-gen-openapiv2 >/dev/null 2>&1 || (cd $(API_DIR) && $(GO) install $(OPENAPI_V2))
	@command -v golangci-lint >/dev/null 2>&1 || (cd $(API_DIR) && $(GO) install $(GOLANGCI))

.PHONY: buf-update
buf-update:
	cd $(PROTO_DIR) && buf dep update

.PHONY: generate
generate: tools
	cd $(PROTO_DIR) && buf generate

.PHONY: generate-go
generate-go:
	cd $(PROTO_DIR) && buf generate --template buf.gen.go.yaml

.PHONY: tidy
tidy:
	cd $(API_DIR) && $(GO) mod tidy

.PHONY: build
build:
	cd $(API_DIR)/cmd/$(APP) && $(GO) build -o ../../$(BIN_DIR)/$(APP)

.PHONY: run
run:
	cd $(API_DIR) && $(GO) run ./cmd/$(APP)

.PHONY: test
test:
	cd $(API_DIR) && $(GO) test ./... -race -cover -count=1

.PHONY: lint
lint: tools
	cd $(API_DIR) && golangci-lint run ./...

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: clean-gen
clean-gen:
	rm -rf $(API_DIR)/proto $(PROTO_DIR)/openapi
