# Build stage
FROM golang:1.21-alpine3.19 AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
# hadolint ignore=DL3018
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy only necessary source files
COPY main.go template_funcs.go ./
COPY templates ./templates

# Build the application as static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags '-extldflags "-static" -s -w' \
    -o main .

# Runtime stage
FROM alpine:3.19

# Install ca-certificates and wget for health checks
# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates wget && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy the binary and templates from builder
COPY --from=builder --chown=appuser:appuser /app/main .
COPY --from=builder --chown=appuser:appuser /app/templates ./templates

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
