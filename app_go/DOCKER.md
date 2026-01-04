# Docker Best Practices - Go Application

## Overview

This document outlines Docker best practices implemented in the Wordle Game application.

## Best Practices Implemented

### 1. Non-Root User

```dockerfile
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser
```

**Why:** Running as root is a security risk. Non-root execution prevents privilege escalation attacks.

### 2. Multi-Stage Build

```dockerfile
# Build stage
FROM golang:1.21-alpine3.19 AS builder
# ... build ...

# Runtime stage
FROM alpine:3.19
# ... only binary ...
```

**Why:** Reduces image from ~300MB to ~15MB by excluding build tools from final image.

### 3. Specific Base Image Versions

```dockerfile
FROM golang:1.21-alpine3.19 AS builder
FROM alpine:3.19
```

**Why:** Reproducible builds and version control.

### 4. Static Binary Compilation

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags '-extldflags "-static" -s -w' \
    -o main .
```

**Flags:**

- `CGO_ENABLED=0`: Fully static binary
- `-s -w`: Strip debug info for smaller size

**Why:** No dynamic dependencies, can run on minimal base images.

### 5. Optimized Layer Caching

```dockerfile
# Dependencies first (change less)
COPY go.mod go.sum ./
RUN go mod download

# Source code later (changes more)
COPY main.go template_funcs.go ./
```

**Why:** Faster builds by leveraging Docker's layer caching.

### 6. Specific File Copying

```dockerfile
COPY main.go template_funcs.go ./
COPY templates ./templates
```

**Why:** Only necessary files, smaller images, better security.

### 7. .dockerignore File

Excludes: `*.test`, `*_test.go`, `.git/`, `.env`

**Why:** Faster builds, prevents accidental secret inclusion.

### 8. Minimal Dependencies

```dockerfile
# hadolint ignore=DL3018
RUN apk add --no-cache git ca-certificates
```

**Note:** We intentionally don't pin Alpine package versions. Base image pinning (`alpine:3.19`) provides stability. Package version pinning can break builds as Alpine packages are ephemeral.

### 9. Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

**Why:** Enables automatic container restart and better orchestration.

### 10. Proper File Ownership

```dockerfile
COPY --from=builder --chown=appuser:appuser /app/main .
```

**Why:** Correct permissions from the start.

## Distroless vs Alpine Comparison

### Alpine Multi-Stage Image

- Builder: `golang:1.21-alpine3.19`
- Runtime: `alpine:3.19`
- Size: ~15MB
- Includes: Shell, wget, ca-certificates
- User: `appuser` (UID 1000)
- **Best for:** Development, debugging

### Distroless Static Image

- Builder: `golang:1.21-alpine3.19`
- Runtime: `gcr.io/distroless/static-debian12:nonroot`
- Size: ~5MB
- Includes: Only static binary runtime
- User: `nonroot` (UID 65532)
- **Best for:** Production, maximum security

### Size Comparison

Run this to compare image sizes:

```bash
make compare-sizes
```

Or manually:

```bash
docker images | grep escalopax
```

### Expected Results

Image Type | Size | What's Included
--- | --- | ---
Full Go | ~300MB | Go toolchain + OS + app
Alpine | ~15MB | Minimal OS + binary
Distroless | ~5MB | Binary + minimal runtime

## Go-Specific Advantages

### Static Compilation

Go compiles to fully static binaries:

- No interpreter needed
- No runtime dependencies
- Perfect for minimal containers
- Can run on `scratch` image

### Performance

- Binary size: ~8-12MB
- Startup: < 100ms
- Memory: 10-20MB
- Concurrent: Built-in goroutines

## Security Benefits

1. **Non-root** - Prevents privilege escalation
2. **Static binary** - No dynamic library vulnerabilities
3. **Minimal base** - Fewer packages = fewer CVEs
4. **No shell** - Cannot execute arbitrary commands
5. **Thread-safe** - sync.RWMutex for concurrent access

## Build and Deploy

```bash
# Build
make build-go
make build-distroless-go

# Test
docker run -d -p 8080:8080 escalopax/wordle-game:latest
curl http://localhost:8080/health

# Push to Docker Hub
make push-go
```

## Distroless Trade-offs

**Advantages:**

- Smallest possible size (5MB)
- No shell = cannot be exploited
- Perfect for static binaries
- Google-maintained

**Disadvantages:**

- Cannot debug with shell
- No health check tools in image
- Requires static compilation
