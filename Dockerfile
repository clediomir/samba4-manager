# Build Stage
FROM golang:1.26.0-alpine AS builder

RUN apk add --no-cache upx

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the statically linked binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o samba4-manager .
RUN upx --best --lzma samba4-manager

# Runtime Stage
FROM alpine:latest

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata sqlite-libs

COPY --from=builder /app/samba4-manager .
COPY config.toml /etc/samba4-manager/config.toml

EXPOSE 8080

CMD ["./samba4-manager", "serve", "--config", "/etc/samba4-manager/config.toml"]
