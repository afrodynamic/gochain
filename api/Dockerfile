FROM golang:1.25-alpine AS build
WORKDIR /src

COPY api/go.mod api/go.sum ./api/
RUN cd api && go mod download

COPY api ./api
COPY proto ./proto

WORKDIR /src/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o /out/gochaind ./cmd/gochaind

FROM gcr.io/distroless/static-debian12:nonroot
ENV PORT=8080
ENV CHAIN=gochain
COPY --from=build /out/gochaind /gochaind
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gochaind"]
