FROM golang:1.25-alpine AS build

WORKDIR /app
ENV PATH="/go/bin:${PATH}"

RUN apk add --no-cache protobuf git build-base

COPY api/go.mod api/go.sum ./api/
WORKDIR /app/api
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

WORKDIR /app
COPY proto/ ./proto/
COPY api/ ./api/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
    go install github.com/bufbuild/buf/cmd/buf@latest

WORKDIR /app/proto
RUN buf generate --template buf.gen.go.yaml

WORKDIR /app/api

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -trimpath -o /out/app ./cmd/gochaind

FROM gcr.io/distroless/static-debian12:nonroot

ENV PORT=8080
ENV CHAIN=gochain

COPY --from=build /out/app /app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
