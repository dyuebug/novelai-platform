FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache ca-certificates && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/gateway ./cmd/gateway

FROM alpine:3.19
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /out/gateway /app/gateway
COPY ./configs/config.yaml.example /app/config.yaml.example
USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/app/gateway"]
