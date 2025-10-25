FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.work ./
COPY api/go.mod api/go.sum ./api/
RUN --mount=type=cache,target=/go/pkg/mod \
    cd api && go mod download

COPY api ./api
WORKDIR /src/api
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/gochaind ./cmd/gochaind

FROM gcr.io/distroless/static-debian12:nonroot
ENV PORT=8080
ENV CHAIN=gochain
WORKDIR /
COPY --from=build /out/gochaind /gochaind
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gochaind"]
