# Build stage
FROM golang:1.21-alpine3.19 AS builder

WORKDIR /app

# Install build dependencies
# hadolint ignore=DL3018
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source files
COPY main.go template_funcs.go ./
COPY templates ./templates

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags '-extldflags "-static" -s -w' \
    -o main .

# Runtime stage - distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Set working directory
WORKDIR /app

# Copy binary and templates from builder
COPY --from=builder --chown=nonroot:nonroot /app/main .
COPY --from=builder --chown=nonroot:nonroot /app/templates ./templates

# Expose port
EXPOSE 8080

# Run the application (distroless uses nonroot user by default, UID 65532)
CMD ["./main"]
