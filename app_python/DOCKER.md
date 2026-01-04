# Docker Best Practices - Python Application

## Overview

This document outlines Docker best practices implemented in the Moscow Time Display application.

## Best Practices Implemented

### 1. Non-Root User

```dockerfile
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser
USER appuser
```

**Why:** Running as root is a security risk. Non-root execution prevents privilege escalation attacks.

### 2. Specific Base Image Version

```dockerfile
FROM python:3.11-alpine3.19
```

**Why:** Version pinning ensures reproducible builds and prevents breaking changes.

### 3. Optimized Layer Caching

```dockerfile
# Copy requirements first (changes less frequently)
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code later (changes more frequently)
COPY app.py .
COPY templates ./templates
```

**Why:** Faster builds by leveraging Docker's layer caching.

### 4. Specific File Copying

```dockerfile
COPY app.py .
COPY templates ./templates
```

**Why:** Only include necessary files, reduces image size and attack surface.

### 5. .dockerignore File

Excludes: `__pycache__/`, `*.py[cod]`, `venv/`, `.git/`, `.env`

**Why:** Faster builds, smaller images, prevents accidental secret inclusion.

### 6. Environment Variables

```dockerfile
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1
```

**Why:** Prevents `.pyc` files, immediate log output, reduced image size.

### 7. No Pip Cache

```dockerfile
RUN pip install --no-cache-dir -r requirements.txt
```

**Why:** Reduces image size significantly.

### 8. Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:5000/health || exit 1
```

**Why:** Enables automatic container restart and better orchestration.

### 9. Proper File Ownership

```dockerfile
COPY --chown=appuser:appuser app.py .
```

**Why:** Correct permissions from the start, no runtime issues.

## Distroless vs Alpine Comparison

### Alpine Image

- Base: `python:3.11-alpine3.19`
- Size: ~60MB
- Includes: Shell, package manager, basic utilities
- User: `appuser` (UID 1000)
- **Best for:** Development, debugging

### Distroless Image

- Base: `gcr.io/distroless/python3-debian12:nonroot`
- Size: ~45MB
- Includes: Only Python runtime
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

Image Type | Size | Reduction
--- | --- | ---
Alpine | ~60MB | 70% vs Debian
Distroless | ~45MB | 78% vs Debian

## Security Benefits

1. **Non-root** - Prevents privilege escalation
2. **Minimal attack surface** - Fewer binaries = fewer vulnerabilities
3. **Immutable** - No package manager in distroless
4. **Version pinning** - Reproducible and auditable builds

## Build and Deploy

```bash
# Build
make build-python
make build-distroless-python

# Test
docker run -d -p 5000:5000 escalopax/moscow-time-app:latest
curl http://localhost:5000/health

# Push to Docker Hub
make push-python
```

## Distroless Trade-offs

**Advantages:**

- Smaller size
- No shell = cannot be exploited
- Fewer CVEs
- Forces better practices

**Disadvantages:**

- Cannot debug with shell
- Cannot install packages at runtime
- Requires multi-stage build
