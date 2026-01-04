# Moscow Time Display - Python Web Application

## Overview

FastAPI web application displaying current Moscow (MSK) timezone with modern Bootstrap UI.

## Features

- Real-time Moscow timezone display
- Responsive Bootstrap 5 UI
- Health check endpoint
- Docker containerization with non-root user

## Technology Stack

- FastAPI 0.109.0, Uvicorn, Jinja2
- Bootstrap 5, Python 3.11+

## Quick Start

### Using Docker Compose

```bash
make compose-up
```

Access at: <http://localhost:5000>

### Local Development

```bash
cd app_python
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt
uvicorn app:app --host 0.0.0.0 --port 5000
```

## Docker

### Build

```bash
# Standard Alpine image
make build-python

# Distroless image
make build-distroless-python
```

### Pull from Docker Hub

```bash
make pull-python
```

### Run

```bash
# Standard image
make run-python

# Distroless image
make run-distroless-python
```

### Compare Sizes

```bash
make compare-sizes
```

This shows both regular and distroless image sizes.

### Image Comparison

Image Type | Size | Security | Use Case
--- | --- | --- | ---
Alpine | ~60MB | Good | Development, debugging
Distroless | ~45MB | Excellent | Production

See [DOCKER.md](./DOCKER.md) for best practices details.

## API Endpoints

- **GET /** - Main page displaying Moscow time
- **GET /health** - Health check

## Testing

```bash
make test
```
