"""
FastAPI application that displays the current time in Moscow timezone.
"""

import time
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
    """Display the current time in Moscow timezone.

    Args:
        request: FastAPI request object

    Returns:
        HTMLResponse: Rendered HTML page with Moscow time
    """
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
    """Health check endpoint for monitoring.

    Returns:
        dict: Status of the application
    """
    REQUEST_COUNT.labels(method="GET", endpoint="/health", status="200").inc()
    return {"status": "healthy", "service": "moscow-time-app"}


@app.get("/metrics")
async def metrics() -> Response:
    """Prometheus metrics endpoint returning the correct content type.

    Returns:
        Response: Prometheus metrics in text format
    """
    data = generate_latest()
    return Response(content=data, media_type=CONTENT_TYPE_LATEST)


@app.middleware("http")
async def metrics_middleware(
    request: Request,
    call_next: Callable[[Request], Coroutine[Any, Any, Response]],
) -> Response:
    """Middleware to record request count, duration and errors for all endpoints.

    Args:
        request: FastAPI request object
        call_next: Next middleware function

    Returns:
        Response: The response from the next middleware
    """
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
