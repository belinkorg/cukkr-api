# ==========================================
# Stage 1: Builder
# ==========================================
FROM golang:1.23-alpine AS builder

# Install dependencies untuk build
RUN apk add --no-cache ca-certificates git make tzdata

WORKDIR /app

# Copy dependency files untuk caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy seluruh source code
COPY . .

# Build binary dengan optimasi
# Sesuai dengan struktur cmd/api/main.go dari taskfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.Version=${VERSION:-dev} -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o bin/api cmd/api/main.go

# ==========================================
# Stage 2: Runtime
# ==========================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl tzdata wget

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/bin/api .

# Copy migrations folder (jika butuh run migration di container)
COPY --from=builder /app/migrations ./migrations

# Copy docs folder (untuk swagger documentation)
COPY --from=builder /app/docs ./docs

# Copy config files jika ada
# COPY --from=builder /app/config ./config

# Set timezone ke WIB (sesuaikan dengan kebutuhan)
ENV TZ=Asia/Jakarta

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (sesuaikan dengan port aplikasi kamu)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run aplikasi
CMD ["./api", "go", "test", "-v", "-cover", "./..."]