# syntax=docker/dockerfile:1

## Build Stage
FROM golang:1.26.0-alpine AS builder

RUN apk add --no-cache gcc musl-dev make upx

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build and compress
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o samba4-manager . \
    && upx --best --lzma samba4-manager

## Runtime Stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata sqlite-libs \
    && addgroup -g 1000 app \
    && adduser -u 1000 -G app -D app

WORKDIR /app

COPY --from=builder --chown=app:app /app/samba4-manager .
COPY --chown=app:app config.toml.example /etc/samba4-manager/config.toml

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

ENTRYPOINT ["./samba4-manager"]
CMD ["serve", "--config", "/etc/samba4-manager/config.toml"]
