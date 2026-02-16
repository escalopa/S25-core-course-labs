"""
FastAPI application that displays the current time in Moscow timezone.
"""
from datetime import datetime
from zoneinfo import ZoneInfo

from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, PlainTextResponse
from fastapi.templating import Jinja2Templates
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST

app = FastAPI(title="Moscow Time Display", version="1.0.0")
templates = Jinja2Templates(directory="templates")

# Prometheus metrics
REQUEST_COUNT = Counter(
    'http_requests_total',
    'Total HTTP requests',
    ['method', 'endpoint', 'status']
)

REQUEST_DURATION = Histogram(
    'http_request_duration_seconds',
    'HTTP request duration in seconds',
    ['method', 'endpoint']
)

ERROR_COUNT = Counter(
    'http_errors_total',
    'Total HTTP errors',
    ['method', 'endpoint', 'status']
)


@app.get("/", response_class=HTMLResponse)
async def get_moscow_time(request: Request):
    """
    Display the current time in Moscow timezone.

    Returns:
        HTMLResponse: Rendered HTML page with Moscow time
    """
    with REQUEST_DURATION.labels(method='GET', endpoint='/').time():
        moscow_tz = ZoneInfo("Europe/Moscow")
        current_time = datetime.now(moscow_tz)

        REQUEST_COUNT.labels(method='GET', endpoint='/', status='200').inc()

        return templates.TemplateResponse(
            "index.html",
            {
                "request": request,
                "time": current_time.strftime("%H:%M:%S"),
                "date": current_time.strftime("%B %d, %Y"),
                "timezone": "Moscow (MSK)"
            }
        )


@app.get("/health")
async def health_check():
    """
    Health check endpoint for monitoring.

    Returns:
        dict: Status of the application
    """
    REQUEST_COUNT.labels(method='GET', endpoint='/health', status='200').inc()
    return {"status": "healthy", "service": "moscow-time-app"}


@app.get("/metrics", response_class=PlainTextResponse)
async def metrics():
    """
    Prometheus metrics endpoint.

    Returns:
        PlainTextResponse: Prometheus metrics in text format
    """
    return generate_latest()
