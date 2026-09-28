# ═══════════════════════════════════════════════
#  Build Stage
# ═══════════════════════════════════════════════
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/server

# ═══════════════════════════════════════════════
#  Runtime Stage
# ═══════════════════════════════════════════════
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env.example (default config)
COPY --from=builder /app/.env.example .

EXPOSE 8080

CMD ["./server"]
