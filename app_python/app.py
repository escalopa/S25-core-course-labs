"""FastAPI application that displays the current time in Moscow timezone.

This module provides:
- `/` : HTML page showing current Moscow time (increments visit counter)
- `/visits` : JSON endpoint returning total visits
- `/health` : basic healthcheck
- `/metrics` : Prometheus metrics
"""

import os
import time
import threading
from collections.abc import Callable, Coroutine
from datetime import datetime
from typing import Any
from zoneinfo import ZoneInfo

from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, Response
from fastapi.templating import Jinja2Templates
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST

app = FastAPI(title="Moscow Time Display", version="1.0.0")
templates = Jinja2Templates(directory="templates")

# Visit tracking
VISITS_FILE = "/data/visits"
visit_count = 0
visit_lock = threading.Lock()


def load_visits() -> int:
    """Load visit count from file.

    Returns 0 if file is missing or malformed.
    """
    try:
        if os.path.exists(VISITS_FILE):
            with open(VISITS_FILE, "r", encoding="utf-8") as f:
                return int(f.read().strip() or "0")
    except (FileNotFoundError, ValueError):
        pass
    return 0


def save_visits(count: int) -> None:
    """Save visit count to file."""
    try:
        os.makedirs(os.path.dirname(VISITS_FILE), exist_ok=True)
        with open(VISITS_FILE, "w", encoding="utf-8") as f:
            f.write(str(count))
    except IOError as e:
        # Do not crash the application for IO errors; log to stdout
        print(f"Error saving visits: {e}")


def increment_visits() -> int:
    """Increment and persist the visit counter in a thread-safe way."""
    global visit_count
    with visit_lock:
        visit_count += 1
        save_visits(visit_count)
        return visit_count


# Prometheus metrics
REQUEST_COUNT = Counter(
    "http_requests_total",
    "Total HTTP requests",
    ["method", "endpoint", "status"],
)

REQUEST_DURATION = Histogram(
    "http_request_duration_seconds",
    "HTTP request duration in seconds",
    ["method", "endpoint"],
)

ERROR_COUNT = Counter(
    "http_errors_total",
    "Total HTTP errors",
    ["method", "endpoint", "status"],
)


@app.get("/", response_class=HTMLResponse)
async def get_moscow_time(request: Request) -> Response:
    """Render the main page and increment the visit counter."""
    # Increment visit counter (synchronous, quick I/O)
    increment_visits()

    moscow_tz = ZoneInfo("Europe/Moscow")
    current_time = datetime.now(moscow_tz)

    REQUEST_COUNT.labels(method="GET", endpoint="/", status="200").inc()

    return templates.TemplateResponse(
        "index.html",
        {
            "request": request,
            "time": current_time.strftime("%H:%M:%S"),
            "date": current_time.strftime("%B %d, %Y"),
            "timezone": "Moscow (MSK)",
        },
    )


@app.get("/health")
async def health_check() -> dict[str, str]:
    """Health check endpoint for monitoring."""
    REQUEST_COUNT.labels(method="GET", endpoint="/health", status="200").inc()
    return {"status": "healthy", "service": "moscow-time-app"}


@app.get("/visits")
async def get_visits() -> dict[str, int]:
    """Get the number of visits to the main page."""
    REQUEST_COUNT.labels(method="GET", endpoint="/visits", status="200").inc()
    return {"visits": visit_count}


@app.get("/metrics")
async def metrics() -> Response:
    """Prometheus metrics endpoint returning the correct content type."""
    data = generate_latest()
    return Response(content=data, media_type=CONTENT_TYPE_LATEST)


@app.middleware("http")
async def metrics_middleware(
    request: Request,
    call_next: Callable[[Request], Coroutine[Any, Any, Response]],
) -> Response:
    """Middleware to record request count, duration and errors for all endpoints."""
    method = request.method
    endpoint = request.url.path
    start_time = time.time()

    try:
        response = await call_next(request)
        status = str(response.status_code)
        REQUEST_COUNT.labels(method=method, endpoint=endpoint, status=status).inc()
        if response.status_code >= 500:
            ERROR_COUNT.labels(method=method, endpoint=endpoint, status=status).inc()
        return response
    except Exception as exc:  # pylint: disable=broad-except
        # Count unexpected server errors
        ERROR_COUNT.labels(method=method, endpoint=endpoint, status="500").inc()
        raise exc
    finally:
        # Record request duration
        duration = time.time() - start_time
        REQUEST_DURATION.labels(method=method, endpoint=endpoint).observe(duration)


# Load initial visit count on startup
visit_count = load_visits()
