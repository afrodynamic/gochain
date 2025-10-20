# api/Dockerfile
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12
ENV ADDR=0.0.0.0:8080
ENV CHAIN=gochain
COPY --from=build /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
