# Build stage - install dependencies
FROM python:3.11-alpine3.19 AS builder

WORKDIR /app

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir --target=/app/packages -r requirements.txt

# Copy application files
COPY app.py .
COPY templates ./templates

# Runtime stage - distroless
FROM gcr.io/distroless/python3-debian12:nonroot

# Set working directory
WORKDIR /app

# Copy installed packages and application from builder
COPY --from=builder --chown=nonroot:nonroot /app/packages /app/packages
COPY --from=builder --chown=nonroot:nonroot /app/app.py .
COPY --from=builder --chown=nonroot:nonroot /app/templates ./templates

# Set Python path to include packages
ENV PYTHONPATH=/app/packages

# Expose port
EXPOSE 5000

# Run the application (distroless uses nonroot user by default, UID 65532)
CMD ["app.py"]
