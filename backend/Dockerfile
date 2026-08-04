# FASE A.4: Docker Hardening - Multi-stage build for production
# CIS Docker Benchmark compliant
# Twelve-Factor App compliant

# Build stage
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Create non-root user for build
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o server \
    ./cmd/server

# Production stage
FROM scratch

# Import CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import user and group files from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy binary from builder
COPY --from=builder /build/server /server

# Create non-root user
USER appuser

# Set working directory
WORKDIR /

# Make binary executable
RUN chmod +x /server

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD /server health || exit 1

# Run as non-root user
USER appuser

# Read-only filesystem (except for necessary directories)
# Volumes will be mounted for writable paths

# Expose port
EXPOSE 8080

# Run application
ENTRYPOINT ["/server"]
