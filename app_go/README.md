# Wordle Game - Go Web Application

## Overview

Interactive Wordle game with 6 attempts to guess a 5-letter word. Features Bootstrap UI with animations and thread-safe session management.

## Features

- Full Wordle game implementation
- Bootstrap 5 UI with CSS animations
- Color-coded feedback (green, yellow, gray)
- Thread-safe concurrent access
- Health check endpoint

## Technology Stack

- Go 1.21, net/http, html/template
- Bootstrap 5, UUID

## Quick Start

### Using Docker Compose

```bash
make compose-up
```

Access at: <http://localhost:8080>

### Local Development

```bash
cd app_go
go mod download
go run .
```

## How to Play

1. Visit <http://localhost:8080> to start a new game
2. Enter a 5-letter word
3. Observe feedback:
   - 🟩 Green: Correct position
   - 🟨 Yellow: Wrong position
   - ⬜ Gray: Not in word
4. Win in 6 attempts or less

## Docker

### Build

```bash
# Standard Alpine multi-stage image
make build-go

# Distroless image
make build-distroless-go
```

### Pull from Docker Hub

```bash
make pull-go
```

### Run

```bash
# Standard image
make run-go

# Distroless image
make run-distroless-go
```

### Compare Sizes

```bash
make compare-sizes
```

This shows both regular and distroless image sizes.

### Image Comparison

Image Type | Size | Build Type
--- | --- | ---
Full Go | ~300MB | Single-stage
Alpine | ~15MB | Multi-stage
Distroless | ~5MB | Multi-stage + minimal base

See [DOCKER.md](./DOCKER.md) for best practices details.

## API Endpoints

- **GET /** - Redirects to new game
- **GET /game/{game_id}** - Game page
- **POST /guess** - Submit guess
- **GET /health** - Health check

## Testing

```bash
make test
```
