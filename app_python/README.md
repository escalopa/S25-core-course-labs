# Moscow Time Display - Python Web Application

![Python CI](https://github.com/escalopa/S25-core-course-labs/workflows/Python%20CI/badge.svg)
[![codecov](https://codecov.io/gh/escalopa/S25-core-course-labs/branch/main/graph/badge.svg?flag=python)](https://codecov.io/gh/escalopa/S25-core-course-labs)

## Overview

FastAPI web application displaying current Moscow (MSK) timezone with modern Bootstrap UI.

## Features

- Real-time Moscow timezone display
- Responsive Bootstrap 5 UI
- Health check endpoint
- Docker containerization with non-root user
- Visit counter with persistence to file
- `/visits` endpoint to track application usage

## Technology Stack

- FastAPI 0.109.0, Uvicorn, Jinja2
- Bootstrap 5, Python 3.11+

## Quick Start

### Using Docker Compose

**With production images from Docker Hub:**

```bash
make compose-prod-up
```

**Or build from source:**

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

- **GET /** - Main page displaying Moscow time (increments visit counter)
- **GET /health** - Health check endpoint
- **GET /visits** - Returns visit count in JSON format

## Persistence

The application tracks visits to the main page and persists the count to a file:

- **Visits File**: `/data/visits` (mounted as volume in Docker)
- **Host Path**: `./app_python/data/visits` (local development)

The visit counter persists across container restarts.

### Example

```bash
# Check visits using Docker
docker-compose exec python-app cat /data/visits

# Check on host machine (after running with volumes)
cat app_python/data/visits
```

## Unit Tests

```bash
make test-python
```

## Testing

```bash
make test
```
