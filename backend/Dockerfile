# ── Build stage ──────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# 複製 module 定義（go.sum 由 service.sh 預先產生）
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼（cache bust 只在程式碼變動時觸發）
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w -extldflags '-static'" \
    -o trbb-server ./cmd/server

# ── Runtime stage ─────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl && \
    cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime && \
    echo "Asia/Taipei" > /etc/timezone

WORKDIR /app
COPY --from=builder /build/trbb-server .
RUN mkdir -p /app/logs

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./trbb-server"]
